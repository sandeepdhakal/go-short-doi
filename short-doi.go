package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

const URL string = "https://shortdoi.org/"

type response struct {
	DOI      string
	ShortDOI string
	IsNew    bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Return the short doi received from shortdoi.org for long `doi`.
func GetShortDOI(doi string) string {
	doi = strings.ReplaceAll(doi, `\`, "")
	resp, err := http.Get(URL + doi + "?format=json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result response
	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println(doi)
		panic(err)
	}
	return result.ShortDOI
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
				shortDoi := GetShortDOI(doi)

				lock.Lock()
				defer lock.Unlock()
				doiMap[doi] = shortDoi
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
		panic("An input file must be provided!")
	}

	// open file for reading
	f, err := os.Open(inFile)
	check(err)

	defer f.Close()

	// get all short DOIs
	doiMap := getShortDOIs(f)

	// open file for writing
	var out *os.File
	if outFile != "" {
		out, err = os.Create(outFile)
		check(err)
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
