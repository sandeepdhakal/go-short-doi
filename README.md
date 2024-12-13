A simple Go program to update long DOIs in a bib file with short DOIs using the [shortDOI Service](https://shortdoi.org/).

# Usage
```bash
Usage of ./short-doi:
  -i string
        input file
  -o string
        output file
```

## Example
If we have a *ref.bib* file with the following reference:

<pre>
@article{jones_characterising_2020,
	title = {Characterising the {Digital} {Twin}: {A} {Systematic} {Literature} {Review}},
	volume = {29},
	journal = {CIRP Journal of Manufacturing Science and Technology},
	author = {Jones, David and Snider, Chris and Nassehi, Aydin and Yon, Jason and Hicks, Ben},
	month = may,
	year = {2020},
	pages = {36--52},
	<mark>note = {doi: {10.1016/j.cirpj.2020.02.002}},</mark>
}
</pre>

You can replace all DOIs with their short counterparts as follows:

```bash
./short-doi -i ref.bib -o short-doi-ref.bib
``` 

We'll get *short-doi-ref.bib* with a short DOI.
<pre>
@article{jones_characterising_2020,
	title = {Characterising the {Digital} {Twin}: {A} {Systematic} {Literature} { Review}},
	volume = {29},
	journal = {CIRP Journal of Manufacturing Science and Technology},
	author = {Jones, David and Snider, Chris and Nassehi, Aydin and Yon, Jason and Hicks, Ben},
	month = may,
	year = {2020},
	pages = {36--52},
	<mark>note = {doi: {10/ghg846}},</mark>
}
</pre>

