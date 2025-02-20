package main

import (
	dw "arxiv-filter/download"
)

func main() {
	dw.DownloadUnpack("http://export.arxiv.org/api/query?search_query=all:Marc+Mezard")
	// only keep these extensions
	dw.FilterTex([]string{".tex", ".bib"})
}
