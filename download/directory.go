package download

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	srcDir    = "sources/"
	unpackDir = "unpacked/"
)

func setupDirs() bool {
	return mkdir(srcDir) == nil && mkdir(unpackDir) == nil
}

// FilterTex recursively deletes all files in the 'unpackDir', which are not tex files
func FilterTex(keepSuffix []string) {
	err := filepath.Walk(unpackDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}
		if !info.IsDir() {
			if shouldDelete(info.Name(), keepSuffix) {
				fmt.Println("Removing ", path)
				if err := os.Remove(path); err != nil {
					fmt.Printf("Error removing file %s: %v\n", path, err)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", unpackDir, err)
	}
}

// shouldDelete determines whether a file should be deleted based on its extension.
// It takes a filename and a slice of extensions to keep. If the filename ends with
// any of the specified extensions, the function returns false, indicating the file
// should not be deleted. Otherwise, it returns true, indicating the file should be deleted.
//
// Parameters:
// - filename: The name of the file to check.
// - keepExtensions: A slice of string extensions to keep.
//
// Returns:
// - bool: true if the file should be deleted, false otherwise.
func shouldDelete(filename string, keepExtensions []string) bool {
	for _, suffix := range keepExtensions {
		if strings.HasSuffix(filename, suffix) {
			return false
		}
	}
	return false
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
