package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// ChangeLog represents a collection of changes between versions
type ChangeLog struct {
	NameTags string
	Commits  []CommitMessage
}

// CommitMessage contains the commit information
type CommitMessage struct {
	ShortCommit string
	Commit      string
	Author      string
	Email       string
	Date        string
	Message     string
}

func fetchTags(releasePattern string) []string {
	out, _ := exec.Command("git", "log", "--tags", "--simplify-by-decoration", "--pretty=\"%ai @%d\"").Output()
	uncuratedTags := strings.Split(string(out), "\n")

	var tags []string
	validTag := regexp.MustCompile(fmt.Sprintf(`(.*\s)@\s.*tag:\s(%v)`, releasePattern))

	var match [][]string
	for _, tag := range uncuratedTags {
		match = validTag.FindAllStringSubmatch(tag, -1)
		if len(match) > 0 && len(match[0]) > 2 {
			// fmt.Printf("%v => %v \n", match[0][1], match[0][2])
			tags = append(tags, strings.Replace(match[0][2], ",", "", -1))
		}
	}
	return tags
}

func fetchChanges(releasePattern string) []ChangeLog {
	var changeLogs []ChangeLog
	tags := fetchTags(releasePattern)

	var diffs []string
	for i := 0; i < len(tags)-1; i++ {
		diffs = append(diffs, fmt.Sprintf("%v..%v", tags[i+1], tags[i]))
	}

	format := "--pretty=format:{\"shortcommit\":\"%h\", \"commit\":\"%H\", \"author\":\"%an\", \"email\":\"%ae\", \"date\":\"%ad\", \"message\":\"%f\"},"

	var out []byte
	var outs string
	var data []CommitMessage
	re := regexp.MustCompile(`,$`)
	for _, diff := range diffs {
		out, _ = exec.Command("git", "log", format, diff).Output()
		outs = fmt.Sprintf("[%v]", re.ReplaceAllString(string(out), ""))
		// fmt.Printf("%v\n", outs)
		if err := json.Unmarshal([]byte(outs), &data); err != nil {
			panic(err)
		}
		changeLogs = append(changeLogs, ChangeLog{NameTags: diff, Commits: data})
		// fmt.Println(data)
	}

	return changeLogs
}

func main() {
	releasePattern := `v[\d{1,4}\.]{1,}`
	// releasePattern := `release-v?[\d{1,4}\.]{1,}`
	changeLog := fetchChanges(releasePattern)
	fmt.Println(changeLog)
}
