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
	fmt.Printf("Feed Updated: %s\n", feed.Updated)

	fmt.Println("\n\nFeed Items:")

	for _, value := range feed.Items[0:5] {
		n := 100 // You can set n to any value you prefer
		fmt.Printf("Title: %s\n", value.Title[:min(len(value.Title), n)])
		fmt.Printf("Description: %s\n", value.Description[:min(len(value.Description), n)])
		fmt.Printf("Link: %s\n", value.Link[:min(len(value.Link), n)])
		fmt.Printf("Links: %v\n", value.Links)
		fmt.Printf("Updated: %s\n", value.Updated[:min(len(value.Updated), n)])
		fmt.Printf("UpdatedParsed: %v\n", value.UpdatedParsed)
		fmt.Printf("Published: %s\n", value.Published[:min(len(value.Published), n)])
		fmt.Printf("PublishedParsed: %v\n", value.PublishedParsed)
		if value.Author != nil {
			fmt.Printf("Author: %s\n", value.Author.Name[:min(len(value.Author.Name), n)])
		}
		fmt.Printf("Authors: %v\n", value.Authors)
		fmt.Printf("GUID: %s\n", value.GUID[:min(len(value.GUID), n)])
		if value.Image != nil {
			fmt.Printf("Image: %s\n", value.Image.URL[:min(len(value.Image.URL), n)])
		}
		fmt.Printf("Categories: %v\n", value.Categories)
		fmt.Printf("Extensions: %v\n", value.Extensions)
		for key, val := range value.Custom {
			fmt.Printf("Custom[%s]: %s\n", key, val[:min(len(val), n)])
		}
	}

	fmt.Println("\n\n ")

}
