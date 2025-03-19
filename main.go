package main

import (
	"flag"
	"fmt"
	"jsonvalidate/validate"
	"jsonvalidate/version"
	"os"
)

const (
	appVersion = "0.1.1"
)

func main() {
	// Define flags
	versionFlag := flag.Bool("v", false, "Print the version of the application")
	helpFlag := flag.Bool("h", false, "Print help information")
	recursiveFlag := flag.Bool("r", false, "Recursively validate JSON files in the directory")

	flag.Parse()

	// Handle version flag
	if *versionFlag {
		version.PrintVersion(appVersion)
		return
	}

	// Handle help flag
	if *helpFlag {
		printHelp()
		return
	}

	paths := flag.Args()
	if len(paths) == 0  && !*recursiveFlag {
		// No paths provided, use standard input
		if err := validate.ValidateJSONFromStdin(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		return
	} else if len(paths) == 0  && *recursiveFlag {
		path := "."
		err := validate.ValidateJSONFilesRecursively(path)
		if err != nil {
			os.Exit(1)
		}

	}


	var finalErr error
	for _, p := range paths {
		fi, err := os.Stat(p)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", p, err)
			finalErr = err
			continue
		}

		if *recursiveFlag && fi.IsDir() {
			err = validate.ValidateJSONFilesRecursively(p)
		} else if *recursiveFlag && !fi.IsDir() {
			// If recursive flag is used on a file, treat the file as a directory
			// by validating just that file.
			err = validate.ValidateJSONFile(p)
		} else if !*recursiveFlag {
			// If not recursive, assume file input
			err = validate.ValidateJSONFile(p)
		}

		if err != nil {
			fmt.Printf("Error validating %s: %v\n", p, err)
			finalErr = err
		}
	}

	if finalErr != nil {
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage: jsonvalidate [OPTIONS] [FILE]")
	fmt.Println("Validate JSON files or input.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help       Print this help message")
	fmt.Println("  -v, --version    Print the version of the application")
	fmt.Println("  -r, --recursive  Recursively validate JSON files in the directory")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  jsonvalidate test.json")
	fmt.Println("  cat test.json | jsonvalidate")
	fmt.Println("  jsonvalidate -r *.json")
	fmt.Println("  jsonvalidate -r folder/*.json")
}
