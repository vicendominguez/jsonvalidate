package main

import (
	"flag"
	"jsonvalidate/logger"
	"jsonvalidate/validate"
	"jsonvalidate/version"
	"os"

	"github.com/pterm/pterm"
)

const (
	appVersion = "0.1.1"
)

func main() {
	// Define flags
	versionFlag := flag.Bool("v", false, "Print the version of the application")
	helpFlag := flag.Bool("h", false, "Print help information")
	recursiveFlag := flag.Bool("r", false, "Recursively validate JSON files in the directory")
	debugFlag := flag.Bool("d", false, "Enable debug logging (verbose output)")

	flag.Parse()

	// Configure logger
	logLevel := pterm.LogLevelInfo // Default log level
	if *debugFlag {
		logLevel = pterm.LogLevelDebug // Enable debug logging if -d or --debug is used
	}

	logger.SetLogLevel(logLevel)

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
	if len(paths) == 0 && !*recursiveFlag {
		// No paths provided, use standard input
		logger.Log.Debug("Using sding...")
		if err := validate.ValidateJSONFromStdin(); err != nil {
			logger.Log.Error("Error validating JSON from stdin", logger.Log.Args("error", err))
			os.Exit(1)
		}
		return
	} else if len(paths) == 0 && *recursiveFlag {
		path := "."
		err := validate.ValidateJSONFilesRecursively(path)
		if err != nil {
			logger.Log.Error("Error validating JSON", logger.Log.Args("path", path, "error", err))
			os.Exit(1)
		}
	}

	var finalErr error
	for _, p := range paths {
		fi, err := os.Stat(p)
		if err != nil {
			logger.Log.Error("Error reading path", logger.Log.Args("path", p, "error", err))
			finalErr = err
			continue
		}

		if *recursiveFlag && fi.IsDir() {
			err = validate.ValidateJSONFilesRecursively(p)
		} else if *recursiveFlag && !fi.IsDir() {
			// If recursive flag is used on a file, treat the file as a directory
			// by validating just that file.
			err = validate.ValidateJSONFilesWithPattern(p)
		} else if !*recursiveFlag {
			// If not recursive, assume file input
			err = validate.ValidateJSONFile(p)
		}

		if err != nil {
			logger.Log.Error("Error validating JSON", logger.Log.Args("path", p, "error", err))
			finalErr = err
		}
	}

	if finalErr != nil {
		os.Exit(1)
	}
}

func printHelp() {
	pterm.DefaultBasicText.Println("Usage: jsonvalidate [OPTIONS] [FILE]")
	pterm.DefaultBasicText.Println("Validate JSON files or input.")
	pterm.Println()

	pterm.DefaultBasicText.Println("Options:")
	pterm.DefaultBasicText.Println("  -h, --help       Print this help message")
	pterm.DefaultBasicText.Println("  -v, --version    Print the version of the application")
	pterm.DefaultBasicText.Println("  -r, --recursive  Recursively validate JSON files in the directory")
	pterm.DefaultBasicText.Println("  -d, --debug      Enable debug logging (verbose output)")
	pterm.Println()

	pterm.DefaultBasicText.Println("Examples:")
	pterm.DefaultBasicText.Println("  jsonvalidate test.json")
	pterm.DefaultBasicText.Println("  cat test.json | jsonvalidate")
	pterm.DefaultBasicText.Println("  jsonvalidate -r *.json")
	pterm.DefaultBasicText.Println("  jsonvalidate -r folder/*.json")
}
