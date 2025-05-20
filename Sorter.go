package main

import (
	"encoding/json"
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

type FileInfo struct {
	Name string
	Path string
	Size int64
	Perm os.FileMode
	Ext  string
}

func callGPT(jsonFile string) {
	apikey := os.Getenv("OPENAI_API_KEY")
	data, error := os.ReadFile(jsonFile)
	if error != nil {
		fmt.Println("There is an error")
		return
	}
	prompt := fmt.Sprintf("You are a file organizer.Attached is a metadata file of files in a folder(files inside subfolders also included). You have to parse through each file, and using the informations like type, size,permissions, etc organize it based on the best logical reasons. It maybe based on File Types:Documents, Images, Audios, Videos, Others etc; File Size:100-500MB,500-1000MB, etc; Date Modified: Past Month,Past week, etc; or any other like this. The most apt reason must be chosen. Return Result in a JSON format along with the subfolder names and current path and new path of file.Here is the data:%s", string(data))

	body := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <folder_path>")
		return
	}
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
	error = os.WriteFile(jsonFileName, jsonData, 0644)
	if error != nil {
		fmt.Print("There is an error!")
	}

	fmt.Printf("\nMetadata saved to '%s'\n", jsonFileName)
	callGPT(jsonFileName)
}
