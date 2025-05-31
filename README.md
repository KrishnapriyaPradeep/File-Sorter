GPT-Integrated File Organizer

This is a Go-based intelligent file organization tool that scans, classifies, and rearranges files within a given folder using OpenAI's GPT-3.5 Turbo model.
It recursively lists all files — including those inside subfolders — gathers their metadata, splits them into manageable chunks, leverages GPT for smart categorization, and physically sorts them into logical folders. Categorization is based on file type, size, origin, and usage context.

How It Works

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

4. AI Categorization

  Each .json chunk is sent to OpenAI’s GPT API with a custom prompt.

  GPT analyzes and returns:

  What each file likely is.

  Where it logically belongs.

  A list of current_path → new_path mappings.

5. File Organization

  Files are moved to their suggested folders inside the root directory.

  Files that aren't recognized or categorized are moved to an Unsorted/ folder.

6. Cleanup

  Any empty directories are removed for a clean final structure.

Key Advantages

  Avoids deep nesting — folders are clean, intuitive, and flat.

  Considers file size, origin (like WhatsApp), sensitivity, file type, and usage.

  Practical design — perfect for organizing messy Downloads folders or project directories.

Tech Info

  Language	Go (Golang)
  
  AI Backend	OpenAI GPT-3.5 Turbo (API)
  
  Libraries	os, filepath, io, encoding/json, net/http

Usage

  go run main.go <folder_path>

Requirements

  Export your OpenAI API key before running by:  export OPENAI_API_KEY=your-api-key-here
