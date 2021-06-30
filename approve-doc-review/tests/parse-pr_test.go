package tests

import (
	"approve-doc-review/check"
	"testing"
)

func TestGoodPr(t *testing.T) {
	goodFile := "good"

	goodTest, _ := check.ParsePR(goodFile)
	if goodTest == false {
		t.Errorf("Parsing the UI failed; want pass, got %t", goodTest)
	}

}

func TestBadPr(t *testing.T) {
	badFile := "bad"

	badTest, _ := check.ParsePR(badFile)
	if badTest == true {
		t.Errorf("Parsing the UI failed; want fail, got %t", badTest)
	}
}
