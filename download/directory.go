package download

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	keepSuffix = ".tex"
	srcDir     = "sources/"
	unpackDir  = "unpacked/"
)

func setupDirs() bool {
	return mkdir(srcDir) == nil && mkdir(unpackDir) == nil
}

// FilterTex recursively deletes all files in the 'unpackDir', which are not tex files
func FilterTex() {
	err := filepath.Walk(unpackDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), keepSuffix) {
			fmt.Println("Removing ", path)
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error cleaning up: %v\n", err)
	}
}

func mkdir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// EmptyFolders removes all files in the 'srcDir' and 'unpackDir'
func EmptyFolders() {
	removeFilesInDir(srcDir)
	removeFilesInDir(unpackDir)
}

func removeFilesInDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		fmt.Printf("Error cleaning up directory %s: %v\n", dir, err)
	}
}

func getDownloadPath(src Source) string {
	return fmt.Sprintf("%s%s", srcDir, src.id)
}

func getUnpackDir(src Source) string {
	return fmt.Sprintf("%s%s/", unpackDir, src.id)
}
