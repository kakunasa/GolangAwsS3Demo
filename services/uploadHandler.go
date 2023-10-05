package services

import (
	"io"
	"net/http"
)

// UploadHandler This function handles the file upload http request.
func UploadHandler(w http.ResponseWriter, r *http.Request) map[string][]byte {

	// Parse the form data.
	// Limit the total file size to 10MB.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return nil
	}

	// Get the list of uploaded files.
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return nil
	}

	// Limit the number of uploaded files to 10.
	if len(files) > 10 {
		http.Error(w, "Maximum 10 files allowed for upload", http.StatusBadRequest)
		return nil
	}

	// Create a map to store the file name and the file's byte slice.
	fileMap := make(map[string][]byte)

	// Iterate over each uploaded file.
	for _, fileHeader := range files {

		// Check the file size.
		if fileHeader.Size > 1024*1024 {
			http.Error(w, "File size exceeds the maximum limit of 1.0MB", http.StatusBadRequest)
			return nil
		}

		// Open the uploaded file.
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Unable to open file", http.StatusInternalServerError)
			return nil
		}

		// Read the file data to a byte slice.
		fileData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file data", http.StatusInternalServerError)
			return nil
		}

		// Store the file data in the map.
		fileMap[fileHeader.Filename] = fileData
	}

	// Return the map containing the file name and the file byte slice.
	return fileMap
}
