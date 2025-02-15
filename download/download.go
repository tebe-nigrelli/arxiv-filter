package download

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"os/exec"
	"time"

	"github.com/mmcdole/gofeed"
)

func DownloadUnpack(query string) {

	if !setupDirs() {
		fmt.Println("Could not mkdirs")
		return
	}

	start := time.Now()
	feed, err := ContextQuery(query, 5)
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
	allottedTime := time.Duration(len(sources)) * time.Second
	downloadContext, cancelDownloadContext := context.WithTimeout(context.Background(), allottedTime)
	defer cancelDownloadContext()

	for i := range sources {
		sources[i].SrcFile = DownloadSource(sources[i], downloadContext)
	}

	elapsed = time.Since(start)
	fmt.Printf("Download took %s\n\n", elapsed)

	for i := range sources {
		fmt.Println("src file: ", sources[i].SrcFile)
		UnpackSource(sources[i])
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

// DownloadSource downloads the source of the article
// and saves it in the sources directory, under the name of the article
func DownloadSource(src Source, ctx context.Context) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	downloadPath := getDownloadPath(src)

	req, _ := http.NewRequestWithContext(ctx, "HEAD", src.SrcLink, nil)
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		if ct := resp.Header.Get("Content-Type"); ct != "" {
			exts, _ := mime.ExtensionsByType(ct)
			if len(exts) > 0 {
				cmd := exec.CommandContext(ctx, "curl", "-o", downloadPath+exts[0], src.SrcLink)
				err = cmd.Run()
				if err != nil {
					fmt.Printf("Error 1 downloading source for %s: %v\n", src.id, err)
				}
				return src.id + exts[0]
			}
		}
	}

	// Fallback to original behavior if extension detection fails
	cmd := exec.CommandContext(ctx, "curl", "-o", downloadPath, src.SrcLink)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error 2 downloading source for %s: %v\n", src.id, err)
		return ""
	}
	return downloadPath
}
