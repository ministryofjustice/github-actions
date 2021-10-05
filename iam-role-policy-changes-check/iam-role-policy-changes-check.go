package main

import (
	"flag"
	ghaction "github.com/sethvargo/go-githubactions"
	"iam-role-policy-changes-check/identifyiam"
	"log"
)

func main() {
	flag.Parse()
        // fileName is the file created by a GitHub action, it contains the output of a git diff.
        fileName := "changes"
	// prRelevant will return true or false depending on the contents of fileName. We don't want
	// the GH action to error here so we just log the error and take no action.
	prRelevant, err := identifyiam.ParsePR(fileName)
	if err != nil {
		log.Println("Unable to parse the PR - ", err)
	}

	// Conditional check to see if we should pass or fail the step. We don't want a hard fail so we set
	// the output to false and log.
	if prRelevant {
		log.Println("Success: The changes in this PR are not related IAM roles/Policies.")
		 ghaction.SetOutput("review_pr", "true")
	} else {
		log.Println("Fail: Attention - Either the PR contains changes that potentially relate to IAM roles or IAM Policies .")
		ghaction.SetOutput("review_pr", "true")
	}
}
