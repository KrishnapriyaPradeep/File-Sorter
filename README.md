GPT-Integrated File Organiser

This is a Go-based intelligent file organisation tool that scans, classifies, and rearranges files within a given folder using OpenAI's GPT-3.5 Turbo model.
It recursively lists all files — including those inside subfolders — gathers their metadata, splits them into manageable chunks, leverages GPT for smart categorisation, and physically sorts them into logical folders. Categorisation is based on file type, size, origin, and usage context.

HOW IT WORKS

      1. Input Folder Detection
      
            Run the program with a given folder path.
          
            It recursively fetches all files inside the folder and its subfolders.
          
      2. Metadata Collection
      
            For each file, it collects:
          
              File name
          
              Full path
          
              File size
          
              Permissions
          
              File extension
      
      3. Chunking & JSON Export
      
            Metadata is split into chunks of 30 files.
          
            Each chunk is saved as a .json file for processing.
          
      4. AI Categorisation
      
            Each .json chunk is sent to OpenAI’s GPT API with a custom prompt.
          
            GPT analyses and returns:
          
            What each file likely is.
          
            Where it logically belongs.
          
            A list of current_path → new_path mappings.
      
      5. File Organisation
      
            Files are moved to their suggested folders inside the root directory.
          
            Files that aren't recognised or categorised are moved to an Unsorted/ folder.
          
      6. Cleanup
      
            Any empty directories are removed for a clean final structure.
    

TECH INFO

      Language	Go (Golang)
      
      AI Backend	OpenAI GPT-3.5 Turbo (API)
      
      Libraries	os, filepath, io, encoding/json, net/http

USEAGE

      go run main.go <folder_path>

REQUIREMENTS

      Export your OpenAI API key before running by:  export OPENAI_API_KEY=your-api-key-here


HOW TO IMPLEMENT

      PREREQUISITES
            OS - TERMINAL  
            GO (Install)
            OpenAI Account - API Keys - Create new secret key - Copy and safely save key (Won't be able to see again)
      RUN THE PROGRAM
           1. A project directory - Initialise Go module - Save code in the directory
           2. Open the current directory's Terminal
           3. Set API Key in Terminal (Command Prompt or PowerShell)
                For Windows 
                      $env:OPENAI_API_KEY="your-secret-key"
                For macOS/ Linux
                      export OPENAI_API_KEY="your-secret-key"
                Run this command each time you open a new terminal
           4. Run Program
                go run Program.go <Folderpath>
                Replace "Program" with the name you have saved your code.
                Replace Folderpath with the path of directory you want to organize
![Screenshot 2025-06-02 154114](https://github.com/user-attachments/assets/77ffb31d-6d72-4197-a883-883b2ea90f04)

![Screenshot 2025-06-02 154720](https://github.com/user-attachments/assets/8fbc7684-8070-4254-b6e8-f2888d9f9e41)
     
![Screenshot 2025-06-02 154738](https://github.com/user-attachments/assets/2387b44a-dd4b-461d-9213-b55c8d7a4e87)
![Screenshot 2025-06-02 154756](https://github.com/user-attachments/assets/90ba2deb-d22f-47dc-b78d-6ed57a066237)
