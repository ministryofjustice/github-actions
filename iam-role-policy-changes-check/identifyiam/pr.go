// Check is used to perform bool checks on whether a PR
// contain potential IAM Role and Policy changes/additions 
// to alert a reviewer to the above content
package identifyiam

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
// of the `git diff` contain any of the following combination - a prefix `+` or `-` in addition to:
// `aws_iam` or `"Effect":` or `"Action":` or `"s3:` or `ec2:` or `"iam:` or `"sqs:` somewhere in the line
// If it finds one - it will fail at the first occurence
// If it does not contain any - it will succeed
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

	// If the text collection contains anything other than `+` `-` in addtion to:
	// `aws_iam` or `"Effect":` or `"Action":` or `"s3:` or `ec2:` or `"iam:` or `"SQS:` or `"sqs:``, it'll pass.
	        for _, line := range text {
                if (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-")) &&

                                                (strings.Contains(line, "aws_iam") || strings.Contains(line, "\"Effect\":") || strings.Contains(line, "\"Action\":") || strings.Contains(line, "\"s3:") || strings.Contains(line, "\"ec2:") || strings.Contains(line, "\"IAM") || strings.Contains(line, "\"iam:") || strings.Contains(line, "\"SQS:") || strings.Contains(line, "\"sqs:")) {
                        return false, errors.New("Reviewer to check - Change(s) found that are potentially Iam policy/role related: " + line)
                } else {
                        continue
                }
	}

        return true, nil
}
