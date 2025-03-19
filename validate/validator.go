package validate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ValidateJSONFile validates a single JSON file.
func ValidateJSONFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return validateJSON(decoder)
}

// ValidateJSONFromStdin validates JSON from standard input.
func ValidateJSONFromStdin() error {
	decoder := json.NewDecoder(os.Stdin)
	return validateJSON(decoder)
}

// ValidateJSONFilesRecursively validates all JSON files in a directory concurrently.
func ValidateJSONFilesRecursively(rootPath string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 100)   // Buffered channel for errors.
	fileChan := make(chan string, 100) // Buffered channel for file paths.

	// Worker pool to validate files concurrently.
	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for filePath := range fileChan {
				fmt.Printf("Worker %d validating: %s\n", workerID, filePath)
				if err := ValidateJSONFile(filePath); err != nil {
					errChan <- fmt.Errorf("error in %s: %w", filePath, err)
				}
			}
		}(i)
	}

	// Walk the directory tree and send JSON files to the fileChan.
	walkErr := filepath.Walk(rootPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			// Log error and continue with other files.
			fmt.Printf("Error accessing %s: %v\n", filePath, err)
			return nil
		}
		// Use case-insensitive check if necessary:
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".json") {
			fmt.Printf("Found JSON file: %s\n", filePath)
			fileChan <- filePath
		}
		return nil
	})

	close(fileChan) // Signal that no more files will be sent.
	wg.Wait()       // Wait for all workers to finish.
	close(errChan)  // Close error channel after processing is complete.

	// Report any errors encountered during file validation.
	var finalErr error
	for e := range errChan {
		fmt.Println(e)
		finalErr = e // Last error stored (if needed).
	}

	if walkErr != nil {
		return fmt.Errorf("error during file traversal: %w", walkErr)
	}
	return finalErr
}

// validateJSON uses a streaming decoder to check JSON validity.
func validateJSON(decoder *json.Decoder) error {
	var obj json.RawMessage
	for {
		err := decoder.Decode(&obj)
		if err == io.EOF {
			return nil // All JSON decoded successfully.
		}
		if err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
	}
}
