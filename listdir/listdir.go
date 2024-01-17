package main

import (
	"io/fs"
	"os"
	"slices"
)

// non recursive
func ListDirectory(directory string, filter func(fs.FileInfo) bool) ([]string, error) {
	items := []string{}

	dir, err := os.Open(directory)
	if err != nil {
		return items, err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return items, err
	}

	for _, fileInfo := range fileInfos {
		if filter == nil || filter(fileInfo) {
			items = append(items, fileInfo.Name())
		}
	}

	slices.Sort(items)
	return items, nil
}
