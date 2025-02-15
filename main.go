package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	srcDir    = "sources"
	unpackDir = "unpacked"
)

func main() {

	start := time.Now()
	feed, err := ContextQuery("http://export.arxiv.org/api/query?search_query=all:marc+mezard", 3)
	if err != nil {
		fmt.Println("Error parsing feed:", err)
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("ContextQuery took %s\n\n", elapsed)

	// Feed
	start = time.Now()
	fmt.Printf("Feed Title: %s\n", feed.Title)

	sources, err := ResultToSources(feed.Items)
	if err != nil {
		fmt.Println("Error downloading top results:", err)
		return
	}
	elapsed = time.Since(start)
	fmt.Printf("Construct sources took: %s\n\n", elapsed)

	// Download
	downloadContext, cancelDownloadContext := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelDownloadContext()

	for _, source := range sources {
		fmt.Println("Source:", source.SrcLink)
		DownloadSource(source, downloadContext)
	}

	elapsed = time.Since(start)
	fmt.Printf("Download took %s\n\n", elapsed)

	for _, source := range sources {
		UnpackSource(source)
	}

	elapsed = time.Since(start)
	fmt.Printf("unpack took %s\n\n", elapsed)

}

// ContextQuery produces a context of length sec seconds and parses the URL result
func ContextQuery(query string, sec int) (*gofeed.Feed, error) {
	if sec <= 0 {
		return nil, fmt.Errorf("timeout must be greater than zero")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sec)*time.Second)
	defer cancel()

	return gofeed.NewParser().ParseURLWithContext(query, ctx)
}

type Source struct {
	// Title of the article
	Title string
	// link to the page of the article
	Link string
	// link to the tex source file
	SrcLink string
	// GUID of the article
	id string
}

func ResultToSources(list []*gofeed.Item) ([]Source, error) {
	n := len(list)
	srcLinks := make([]Source, 0, n)
	for _, value := range list[0:min(len(list), n)] {
		if len(value.GUID) < 21 {
			return nil, fmt.Errorf("GUID too short: %s", value.GUID)
		}

		idString := value.GUID[strings.LastIndex(value.GUID, "/")+1:]
		srcString := value.GUID[:17] + "src" + value.GUID[20:]
		srcLinks = append(srcLinks, Source{Title: value.Title, Link: value.Link, SrcLink: srcString, id: idString})
	}
	return srcLinks, nil
}

func getDownloadDir(src Source) string {
	return fmt.Sprintf("%s/%s", srcDir, src.id)
}

func getUnpackDir(src Source) string {
	return fmt.Sprintf("%s/%s/", unpackDir, src.id)
}

// DownloadSource downloads the source of the article
// and saves it in the sources directory, under the name of the article
func DownloadSource(src Source, ctx context.Context) {
	cmd := exec.CommandContext(ctx, "curl", "-o", getDownloadDir(src), src.SrcLink)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error downloading source for %s: %v\n", src.Title, err)
	}
}

func UnpackSource(src Source) error {
	err := os.MkdirAll(getUnpackDir(src), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating unpack directory for %s: %v\n", src.Title, err)
		return err
	}

	cmd := exec.Command("tar", "-xvf", getDownloadDir(src), "-C", getUnpackDir(src))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error unpacking source for %s: %v\nOutput: %s\n", src.Title, err, string(output))
	}
	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
