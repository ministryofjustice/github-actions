package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Idea:
// Create a gh action that does..
// - checks out the pr code
// - gets the changes and save them to json
// - copy them to your gh action file path
// - go package to read the number of changes
// - if the changes are <= 1 and the pr owner is ministryofjustice/webops
// - approve the pr

// fileName is created by a GitHub action that outputs a git diff to a file
// and copies it to the container running this package.
const fileName = "changes"

func main() {
	// read the file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open %s: %s", fileName, err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	for _, line := range text {
		fmt.Println(line)
	}
	// for every line that starts with a + and not +++
	// check to see if that line starts with last_reviewed_on:
}
