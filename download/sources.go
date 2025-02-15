package download

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mmcdole/gofeed"
)

type Source struct {
	// Title of the article
	Title string
	// link to the page of the article
	Link string
	// link to the tex source file
	SrcLink string
	// filename of the downloaded package
	SrcFile string
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

func UnpackSource(src Source) error {
	err := os.MkdirAll(getUnpackDir(src), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating unpack directory for %s: %v\n", src.Title, err)
		return err
	}
	cmd := exec.Command("tar", "-xvf", srcDir+src.SrcFile, "-C", getUnpackDir(src))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error unpacking source %s: \nOutput: %s\n", src.id, string(output))
	}
	return err
}
