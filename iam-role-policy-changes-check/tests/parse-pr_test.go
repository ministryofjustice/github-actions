package tests

import (
	"iam-role-policy-changes-check/identify-iam"
	"testing"
)

func TestGoodPr(t *testing.T) {
	goodFile := "good"

	goodTest, _ := identify-iam.ParsePR(goodFile)
	if goodTest == false {
		t.Errorf("Parsing the UI failed; want pass, got %t", goodTest)
	}

}

func TestBadPr(t *testing.T) {
	badFile := "bad"

	badTest, _ := identify-iam.ParsePR(badFile)
	if badTest == true {
		t.Errorf("Parsing the UI failed; want fail, got %t", badTest)
	}
}
