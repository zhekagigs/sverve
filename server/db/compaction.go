package db

import (
	"io/fs"
	"os"
	"path/filepath"
)

func closeSegment() {

}

func compactSegment(segment []string) {

}

func mergeSegments(segmentPaths []string) {

}

func runWorkThread() {

}

func getFileToWrite(root string) ([]string, error) {
    var files []string
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() && filepath.Ext(path) == ".seg" {
            files = append(files, path)
        }
		os.Stat("stat")
        return nil
    })
    return files, err
}
