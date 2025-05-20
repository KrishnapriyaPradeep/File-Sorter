package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name string
	Path string
	Size int64
	Perm os.FileMode
	Ext  string
}

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
	var fileInfos []FileInfo
	for _, file := range entries {
		info, error := os.Stat(file)
		if error != nil {
			fmt.Printf("There's an error!")
			continue
		}
		fileInfos = append(fileInfos, FileInfo{
			Name: info.Name(),
			Path: file,
			Size: info.Size(),
			Perm: info.Mode().Perm(),
			Ext:  filepath.Ext(file),
		})
	}

	for i, file := range fileInfos {
		fmt.Printf("%d: %s — Size: %d bytes — Permissions: %s\n", i+1, file.Path, file.Size, file.Perm)
	}

	jsonData, error := json.MarshalIndent(fileInfos, "", "  ")
	if error != nil {
		fmt.Print("There is an error!")
	}

	jsonFileName := "metadatafile.json"
	write := os.WriteFile(jsonFileName, jsonData, 0644)
	if write != nil {
		fmt.Print("There is an error!")
	}

	fmt.Printf("\nMetadata saved to '%s'\n", jsonFileName)

}
