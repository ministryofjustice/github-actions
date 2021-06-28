package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	ghaction "github.com/sethvargo/go-githubactions"
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
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open %s: %s", fileName, err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var text []string
	var numOfAdds int
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "+") && !strings.HasPrefix(scanner.Text(), "+++") {
			numOfAdds++
			text = append(text, scanner.Text())
		}
	}

	file.Close()

	for _, line := range text {
		if !strings.HasPrefix(line, "+last_reviewed_on") {
			ghaction.SetOutput("review_pr", "false")
			log.Println("This PR contains more than review changes. A human must intervene. Changes in this PR:", text)
			// Exit softly as to not fail the GitHub action.
			os.Exit(0)
		}
	}

	if numOfAdds >= 1 {
		log.Println("This PR only contains reviews, so the check will pass.")
		ghaction.SetOutput("review_pr", "true")
	} else {
		log.Println("This PR contains no additions, so the check will fail.")
		ghaction.SetOutput("review_pr", "false")
	}
}
