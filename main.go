package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func processYAMLFiles(directory, searchPattern, replaceString string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	ymlCount := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yml" {
			ymlCount++

			// Open the file
			filePath := filepath.Join(directory, file.Name())
			f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			// Read, modify, and write back each line
			scanner := bufio.NewScanner(f)
			var newLines []string
			for scanner.Scan() {
				line := scanner.Text()
				newLines = append(newLines, replaceAfterColon(line, searchPattern, replaceString))
			}

			// Truncate and write the modified content
			f.Truncate(0)
			f.Seek(0, 0)
			for _, line := range newLines {
				f.WriteString(line + "\n")
			}

			fmt.Printf("Processed %s\n", file.Name())
		}
	}

	fmt.Printf("Found %d YAML files in the directory.\n", ymlCount)

	return nil
}

func replaceAfterColon(line, searchPattern, replaceString string) string {
	re := regexp.MustCompile(`(?P<key>\w+): *(?P<value>.*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 2 { // Ensure we have key and value groups
		value := strings.ReplaceAll(matches[2], searchPattern, replaceString)
		return fmt.Sprintf("%s: %s", matches[1], value)
	}
	return line
}

func main() {
	dirPtr := flag.String("dir", ".", "Directory containing the YAML files")
	flag.Parse()

	var searchPattern, replaceString string
	fmt.Print("Enter the string to search for (From): ")
	fmt.Scanln(&searchPattern)
	fmt.Print("Enter the string to replace it with (To): ")
	fmt.Scanln(&replaceString)

	err := processYAMLFiles(*dirPtr, searchPattern, replaceString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Processing complete.")
}
