package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext("http://export.arxiv.org/api/query?search_query=all:electron+AND+all:proton", ctx)
	if err != nil {
		fmt.Println("Error parsing feed:", err)
		return
	}

	fmt.Println("\n\nFeed Description:")
	fmt.Printf("Feed Title: %s\n", feed.Title)

	sources, err := downloadTopResults(feed.Items)
	if err != nil {
		fmt.Println("Error downloading top results:", err)
		return
	}

	for _, source := range sources {
		fmt.Println("Source Link:", source.SrcLink)
	}

	fmt.Println("\n\n ")

}

type Source struct {
	Link    string
	SrcLink string
}

func downloadTopResults(list []*gofeed.Item) ([]Source, error) {
	n := len(list)
	srcLinks := make([]Source, 0, n)
	for _, value := range list[0:min(len(list), n)] {
		srcString := value.GUID[:17] + "src" + value.GUID[20:]
		srcLinks = append(srcLinks, Source{Link: value.Link, SrcLink: srcString})
	}
	return srcLinks, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
