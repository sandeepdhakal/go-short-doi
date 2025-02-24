package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

// API endpoint
const URL string = "https://shortdoi.org/"

type response struct {
	DOI      string
	ShortDOI string
	IsNew    bool
}

func handleError(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

// Return the short doi received from shortdoi.org for long `doi`.
func GetShortDOI(doi string) (string, error) {
	doi = strings.ReplaceAll(doi, `\`, "")
	resp, err := http.Get(URL + doi + "?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result response
	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.ShortDOI, nil
}

// Get short DOIs for each DOI found in the file `f`.
// Returns a map of LongDOI -> ShortDOI.
func getShortDOIs(f *os.File) map[string]string {
	r := regexp.MustCompile(`10\.\d{4,9}/[-.\\_;()/:A-Za-z0-9]+`)

	var wg sync.WaitGroup
	lock := sync.RWMutex{}
	doiMap := make(map[string]string)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		m := r.FindString(line)
		if len(m) > 0 {
			wg.Add(1)
			go func(doi string) {
				defer wg.Done()
				shortDOI, e := GetShortDOI(doi)
				if e != nil {
					shortDOI = doi
				}

				lock.Lock()
				defer lock.Unlock()
				doiMap[doi] = shortDOI
			}(m)
		}
	}
	wg.Wait()

	return doiMap
}

// Write to file
func writeToFileOrPrint(f *os.File, line string) {
	if f != nil {
		fmt.Fprintln(f, line)
	} else {
		fmt.Println(line)
	}
}

func main() {
	r := regexp.MustCompile(`10\.\d{4,9}/[-.\\_;()/:A-Za-z0-9]+`)

	var inFile, outFile string
	flag.StringVar(&inFile, "i", "", "input file")
	flag.StringVar(&outFile, "o", "", "output file")
	flag.Parse()

	if inFile == "" {
		// check if we have a DOI input
		longDOIs := flag.Args()
		if len(longDOIs) > 0 {
			longDOI := longDOIs[0] // we take only the first input parameter
			m := r.FindString(longDOI)
			if m == "" {
				handleError(errors.New("Please provide a valid DOI"))
			}

			shortDOI, e := GetShortDOI(longDOI)
			handleError(e)

			fmt.Println(shortDOI)
			os.Exit(0)
		}

		fmt.Println("Please provide an input file or a DOI.")
		flag.Usage()
		os.Exit(1)
	}

	// open file for reading
	f, err := os.Open(inFile)
	handleError(err)

	defer f.Close()

	// get all short DOIs
	doiMap := getShortDOIs(f)

	// open file for writing
	var out *os.File
	if outFile != "" {
		out, err = os.Create(outFile)
		handleError(err)
	}

	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		m := r.FindString(line)
		if len(m) > 0 {
			line = strings.Replace(line, m, doiMap[m], 1)
		}
		writeToFileOrPrint(out, line)
	}
}
