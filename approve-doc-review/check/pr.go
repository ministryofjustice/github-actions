// Check is used to perform bool checks on whether a PR and its user
// are valid to be auto approved.
package check

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

// ParsePR takes a fileName argument defined by a GitHub action pipeline.
// It opens fileName and adds each line to a collection for parsing. The file in fileName
// is the output of a `git diff` command in a pull request. The function then
// checks the collection and returns true if numOfAdds >= 1 and the additions in the output
// of the `git diff` contain `+last_reviewed_on`.
func ParsePR(fileName string) (bool, error) {
	// text will contain parsened scanner results as strings.
	var text []string
	// numOfAdds is a counter for the number of additions in a PR.
	var numOfAdds int

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open %s: %s", fileName, err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "-") && !strings.HasPrefix(scanner.Text(), "---") {
			text = append(text, scanner.Text())
		}
		if strings.HasPrefix(scanner.Text(), "+") && !strings.HasPrefix(scanner.Text(), "+++") {
			numOfAdds++
			text = append(text, scanner.Text())
		}
	}

	file.Close()

	// If the text collection contains anything other than `+/-last_reviewed_on`, it'll fail.
	for _, line := range text {
		if strings.HasPrefix(line, "+last_reviewed_on") || strings.HasPrefix(line, "-last_reviewed_on") {
			continue
		} else {
			return false, errors.New("Change found that isn't a review date: " + line)
		}
	}

	// If the PR contains only removals, not additions, it'll fail.
	if numOfAdds >= 1 {
		return true, nil
	}

	return false, nil
}
