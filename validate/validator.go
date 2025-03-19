package validate

import (
	"encoding/json"
	"fmt"
	"io"
	"jsonvalidate/logger"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ValidateJSONFile checks if a single JSON file is valid.
func ValidateJSONFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()
	logger.Log.Debug("Validating file", logger.Log.Args("Path:", filePath))
	decoder := json.NewDecoder(file)
	return validateJSON(decoder)
}

// ValidateJSONFromStdin reads and validates JSON from standard input.
func ValidateJSONFromStdin() error {
	decoder := json.NewDecoder(os.Stdin)
	return validateJSON(decoder)
}

// processFilesWithWorkers runs a worker pool to validate files concurrently.
func processFilesWithWorkers(files []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(files))
	fileChan := make(chan string, len(files))

	numWorkers := 10 // Number of parallel workers

	// Start worker goroutines
	for i := range numWorkers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for filePath := range fileChan {
				logger.Log.Debug("Found file!", logger.Log.Args("Path:", filePath), logger.Log.Args("WorkerID", workerID))
				if err := ValidateJSONFile(filePath); err != nil {
					errChan <- fmt.Errorf("error in %s: %w", filePath, err)
				}
			}
		}(i)
	}

	// Send file paths to workers
	for _, filePath := range files {
		fileChan <- filePath
	}

	// Close channels when done
	close(fileChan)
	wg.Wait()
	close(errChan)

	// Collect errors
	var finalErr error
	for e := range errChan {
		finalErr = e
	}

	return finalErr
}

// ValidateJSONFilesRecursively finds and validates all JSON files in a directory.
func ValidateJSONFilesRecursively(rootPath string) error {
	var files []string

	// Walk through the directory and collect JSON files
	err := filepath.Walk(rootPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue with other files
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".json") {
			files = append(files, filePath)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error during file traversal: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	return processFilesWithWorkers(files)
}

// ValidateJSONFilesWithPattern validates JSON files matching a shell pattern.
func ValidateJSONFilesWithPattern(pattern string) error {
	// Find files using shell-like pattern matching
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern %s: %w", pattern, err)
	}

	if len(matches) == 0 {
		return nil
	}

	return processFilesWithWorkers(matches)
}

// validateJSON ensures JSON validity.
func validateJSON(decoder *json.Decoder) error {
	var obj json.RawMessage
	for {
		err := decoder.Decode(&obj)
		if err == io.EOF {
			return nil // JSON is valid
		}
		if err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
	}
}
