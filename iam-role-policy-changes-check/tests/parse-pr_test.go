package tests

import (
	"iam-role-policy-changes-check/identifyiam"
	"testing"
)

func TestGoodPr(t *testing.T) {
	goodFile := "good"

	goodTest, _ := identifyiam.ParsePR(goodFile)
	if goodTest == false {
		t.Errorf("Parsing the UI failed; want pass, got %t", goodTest)
	}

}

func TestBadPr(t *testing.T) {
	badFile := "bad"

	badTest, _ := identifyiam.ParsePR(badFile)
	if badTest == true {
		t.Errorf("Parsing the UI failed; want fail, got %t", badTest)
	}
}
