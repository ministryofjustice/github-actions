package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

// Idea:
// Create a gh action that does..
// - checks out the pr code
// - gets the changes and save them to json
// - copy them to your gh action file path
// - go package to read the number of changes
// - if the changes are <= 1 and the pr owner is ministryofjustice/webops
// - approve the pr

const fileName = "changes"

func main() {
	_, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("File wasn't found:", err)
	} else {
		fmt.Println("File was found")
	}
}
