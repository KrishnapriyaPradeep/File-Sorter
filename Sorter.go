package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func insidedirectory(path string) []string {
	var filesList []string
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("There's an error1!%v\n", err)
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

func organizeFiles(files []OrganizedFile) {
	for _, file := range files {
		newDir := filepath.Dir(file.NewPath)
		os.MkdirAll(newDir, os.ModePerm)

		if _, err := os.Stat(file.CurrentPath); os.IsNotExist(err) {
			fmt.Printf("File does not exist: %s\n", file.CurrentPath)
			continue
		}

		err := os.Rename(file.CurrentPath, file.NewPath)
		if err != nil {
			fmt.Printf("Failed to move %s to %s: %v\n", file.CurrentPath, file.NewPath, err)
		} else {
			fmt.Printf("Moved %s -> %s\n", file.CurrentPath, file.NewPath)
		}

	}
}

func callGPT(jsonFile string) []OrganizedFile {
	apikey := os.Getenv("OPENAI_API_KEY")
	fmt.Println("API Key:", apikey)

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Printf("There's an error2!%v\n", err)
		return nil
	}
	prompt := fmt.Sprintf(`You are a file organizer.Attached is a metadata file of files in a folder(files inside subfolders also included). You have to parse through each file, and using the informations like type, size,permissions, etc organize it based on the best logical reasons. It maybe based on File Types:Documents, Images, Audios, Videos, Others etc; File Size:100-500MB,500-1000MB, etc; Date Modified: Past Month,Past week, etc; or any other like this. The most apt reason must be chosen. Respond Only with a JSON array of objects. Each object must have :"name", "current_path",new_path. No explanations. Here is the data:%s.`, string(data))

	body := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonPayload, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("There's an error3!%v\n", err)
		return nil
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("There's an error4!%v\n", err)
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+apikey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("There's an error5!%v\n", err)
		return nil
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(bodyBytes, &result)

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		fmt.Printf("There's an error6!%v\n", err)
		return nil
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		fmt.Printf("There's an error7!%v\n", err)
		return nil
	}

	messageMap, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		fmt.Printf("There's an error8!%v\n", err)
		return nil
	}

	message, ok := messageMap["content"].(string)
	if !ok {
		fmt.Printf("There's an error9!%v\n", err)
		return nil
	}

	var organizedFiles []OrganizedFile
	err = json.Unmarshal([]byte(message), &organizedFiles)
	if err != nil {
		fmt.Printf("There's an error10!%v\n", err)
		return nil
	}

	fmt.Println("GPT Response:\n", message)

	return organizedFiles
}

type OrganizedFile struct {
	Name        string `json:"name"`
	CurrentPath string `json:"current_path"`
	NewPath     string `json:"new_path"`
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
		info, err := os.Stat(file)
		if err != nil {
			fmt.Printf("There's an error11!%v\n", err)
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

	jsonData, err := json.MarshalIndent(fileInfos, "", "  ")
	if err != nil {
		fmt.Printf("There's an error12!%v\n", err)
	}

	jsonFileName := "metadatafile.json"
	err = os.WriteFile(jsonFileName, jsonData, 0644)
	if err != nil {
		fmt.Printf("There's an error13!%v\n", err)
	}

	fmt.Printf("\nMetadata saved to '%s'\n", jsonFileName)
	organized := callGPT(jsonFileName)
	if organized != nil {
		organizeFiles(organized)
	}

}
