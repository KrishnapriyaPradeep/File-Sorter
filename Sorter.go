package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <folder_path>")
		return
	}
	path := os.Args[1]
	rootDir := path
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

	chunks := chunkFiles(fileInfos, 30)
	var allOrganized []OrganizedFile

	for i, chunk := range chunks {
		jsonData, err := json.MarshalIndent(chunk, "", "  ")
		if err != nil {
			fmt.Printf("There's an error in chunk %d: %v\n", i, err)
			continue
		}

		jsonFileName := fmt.Sprintf("metadata_chunk_%d.json", i)
		err = os.WriteFile(jsonFileName, jsonData, 0644)
		if err != nil {
			fmt.Printf("Can't write chunk %d: %v\n", i, err)
			continue
		}

		fmt.Printf("\nChunk %d metadata saved to '%s'\n", i, jsonFileName)
		organized := callGPT(jsonFileName, rootDir)
		if organized != nil {
			allOrganized = append(allOrganized, organized...)
		}
	}
	organizeFiles(allOrganized, rootDir)
	sortedMap := make(map[string]bool)
	for _, f := range allOrganized {
		sortedMap[f.CurrentPath] = true
	}
	for _, f := range fileInfos {
		if !sortedMap[f.Path] {
			unsortedPath := filepath.Join(rootDir, "Unsorted", filepath.Base(f.Path))
			os.MkdirAll(filepath.Dir(unsortedPath), os.ModePerm)
			os.Rename(f.Path, unsortedPath)
			fmt.Printf("Unsorted -> %s\n", unsortedPath)
		}
	}
	deleteEmptyDirs(rootDir)

}

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

func chunkFiles(files []FileInfo, chunkSize int) [][]FileInfo {
	var chunks [][]FileInfo
	for i := 0; i < len(files); i += chunkSize {
		end := i + chunkSize
		if end > len(files) {
			end = len(files)
		}
		chunks = append(chunks, files[i:end])
	}
	return chunks
}

type OrganizedFile struct {
	Name        string `json:"name"`
	CurrentPath string `json:"current_path"`
	NewPath     string `json:"new_path"`
}

func callGPT(jsonFile string, rootDir string) []OrganizedFile {
	apikey := os.Getenv("OPENAI_API_KEY")
	if apikey == "" {
		fmt.Println("API key not found in environment variables.")
		return nil
	}

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Printf("There's an error2!%v\n", err)
		return nil
	}
	prompt := fmt.Sprintf(`
You are a file organizer bot. Your task is to analyze and categorize a list of files based on the following factors:

1. File Type:
   - Images: .jpg, .jpeg, .png, .gif, .bmp, .webp, .heic
   - Videos: .mp4, .avi, .mov, .mkv
   - Documents: .pdf, .doc, .docx, .txt, .ppt, .pptx, .xls, .xlsx
   - Audio: .mp3, .wav, .aac
   - Code: .py, .go, .java, .js, .html, .css, .c, .cpp
   - System Files: .sys, .dll, .ini, .log
   - Archives: .zip, .rar, .tar, .gz

2. File Size (if available):
   - Small: <1MB
   - Medium: 1MB - 100MB
   - Large: >100MB

3. Date Created(if available):
	- Today
	- Last Week/This week
	- Past Month/This Month
	- Past Year/This Year
	- By Year (eg, 2023,2024, etc.)

5. By Content(if available):
	- Government Documents
	- Resumes
	- Educational Materials
	- Presentations
	- Legal/Confidential Documents

6. By Useage(if available)
	- Recently Accessed
	- Frequently Accessed
	- Rarely Used
	- Duplicates

7. By Sensitivity(if available)
	- Public
	- Private
	- Encrypted
	- Contains Passwords
	- Confidential

8. By Labels(if available)
	- Project
	- Personal
	- Family
	- Official

9. By Source(if available)
	- WhatsApp
	- Instagram
	- Telegram
	- Chrome or any Browser Download
	- Screenshot
	- Camera Upload

10. Logical Grouping (if applicable):
   - Use any common prefixes in filenames.
   - Use existing folder structures or hints from the current path to suggest subcategorization.
   - Group files that belong to the same context, project, or module if patterns suggest so.
Rules:
- Use the given "path" argument as the root folder. All new paths must remain inside this folder.
- Reuse folders if they already exist (e.g., "/root/Images" if already present).
- Always preserve the original filename.
- Be logical and concise; avoid unnecessary nesting.
- Analyse which attribute shows highes variability or contrast and use that attribute for the current data set to organize and group.  

Output Format: Respond ONLY with a JSON array in the following format:

[
  {
    "name": "example.jpg",
    "current_path": "/home/user/downloads/example.jpg",
    "new_path": "/home/user/downloads/Images/example.jpg"
  },
  ...
]

Given:
- Root path: %s
- File list: %s
`, rootDir, string(data))

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

	if resp.StatusCode != 200 {
		if errMsg, ok := result["error"].(map[string]interface{}); ok {
			fmt.Printf("API Error [%v]: %v\n", errMsg["type"], errMsg["message"])
		} else {
			fmt.Printf("Unexpected response:\n%s\n", string(bodyBytes))
		}
		return nil
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		fmt.Printf("No choices returned by API\n")
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

	for i, f := range organizedFiles {
		newRelPath := filepath.Base(filepath.Dir(f.NewPath))
		newFullPath := filepath.Join(rootDir, newRelPath, f.Name)
		organizedFiles[i].NewPath = newFullPath
	}

	return organizedFiles
}

func organizeFiles(files []OrganizedFile, rootDir string) {
	for _, file := range files {
		newDir := filepath.Dir(file.NewPath)
		if !strings.HasPrefix(file.NewPath, rootDir) {
			fmt.Printf("Skipping invalid new path (outside root): %s\n", file.NewPath)
			continue
		}

		os.MkdirAll(newDir, os.ModePerm)

		if _, err := os.Stat(file.NewPath); err == nil {
			fmt.Printf("File already exists at destination: %s\n", file.NewPath)
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

func deleteEmptyDirs(path string) {
	filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return nil
	})

	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			entries, _ := os.ReadDir(p)
			if len(entries) == 0 {
				err := os.Remove(p)
				if err == nil {
					fmt.Printf("Deleted empty folder: %s\n", p)
				}
			}
		}
		return nil
	})
}
