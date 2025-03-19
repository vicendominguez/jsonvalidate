# JSONValidate

This is a simple command-line tool for validating JSON files or input. Nothing else. 


Just checking if it parses.
It can validate individual files, read JSON from standard input, or recursively validate all JSON files in a directory.

This is not a linter.

## Why?

I didn't want to install any runtime like:

* [zaach/jsonlint](https://github.com/zaach/jsonlint) requires node.
* [Seldaek/jsonlint](https://github.com/Seldaek/jsonlint) requires php.

Some times I needed to test a lot of JSON files in a tree directory. This is pure speed for me. Just trying to parse them.

Just a binary. And it is fast!. Example: 

jsonlint:
```
❯ find . -name "*.json" |wc -l
     156
❯ time find . -name "*.json" -exec jsonlint -c -q {} \;
find . -name "*.json" -exec jsonlint -c -q {} \;  9,02s user 2,14s system 91% cpu 12,218 total
```

jsonvalidate:

```
 ❯ time jsonvalidate -r
jsonvalidate -r  0,06s user 0,12s system 108% cpu 0,163 total
```

Success!

## Usage

Run the tool with the following options:

- **Validate a single file**:  
  `jsonvalidate <file.json>    # /*/*/file.json is valid`

- **Validate JSON from standard input**:  
  `cat <file.json> | jsonvalidate`

- **Recursively validate JSON files in a directory**:  
  `jsonvalidate -r <directory>    # not directory run as recursive`

- **Print the version**:  
  `jsonvalidate -v`

- **Enable debug logging**:  
  `jsonvalidate -d <file.json>`

- **Print help information**:  
  `jsonvalidate -h`

## Key Features

- **Single file validation**: Validate a specific JSON file.
- **Recursive validation**: Validate all JSON files in a directory and its subdirectories.
- **Standard input support**: Validate JSON piped from standard input.
- **Debug mode**: Enable verbose logging for troubleshooting.

## Examples

- Validate a single file:  
  `jsonvalidate /*/*/example.json`

- Validate all JSON files in a directory:  
  `jsonvalidate -r ./data`

- Validate JSON from standard input:  
  `cat data.json | jsonvalidate`

This tool is designed to be simple and efficient for parsing JSON files in various scenarios.
