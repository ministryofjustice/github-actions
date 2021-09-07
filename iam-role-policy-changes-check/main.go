package main

import (
	"flag"
	ghaction "github.com/sethvargo/go-githubactions"
	"iam-role-policy-changes-check/check"
	"log"
)

var (
	fileName = flag.String("filename", "", "Personal access token from GitHub.")
)

func main() {
	flag.Parse()

	// prRelevant will return true or false depending on the contents of fileName. We don't want
	// the GH action to error here so we just log the error and take no action.
	prRelevant, err := check.ParsePR(*fileName)
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
		ghaction.SetOutput("review_pr", "false")
	}
}
