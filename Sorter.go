package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func insidedirectory(path string) []string {
	var filesList []string
	entries, error := os.ReadDir(path)
	if error != nil {
		fmt.Print("There's an error!")
	}
	for _, file := range entries {
		fullPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			subfiles := insidedirectory(fullPath)
			filesList = append(filesList, subfiles...)
		} else {
			filesList = append(filesList, fullPath)
		}
	}
	return filesList
}
func main() {
	path := os.Args[1]
	fmt.Printf("Folder path = %v", path)
	var entries []string
	entries = append(entries, insidedirectory(path)...)
	for index, file := range entries {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Printf("There's an error!")
			continue
		}
		perm := info.Mode().Perm()
		fmt.Printf("%v : %s — Permissions: %s\n", index, file, perm)
	}
}
