package main

import "fmt"

// Idea:
// Create a gh action that does..
// - checks out the pr code
// - gets the changes and save them to json
// - copy them to your gh action file path
// - go package to read the number of changes
// - if the changes are <= 1 and the pr owner is ministryofjustice/webops
// - approve the pr

func main() {
	fmt.Println("stated")
}
