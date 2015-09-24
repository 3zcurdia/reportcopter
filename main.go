package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/russross/blackfriday"
)

const layout = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Changelog</title>
  </head>
  <body>
	  {{.Content}}
	</body>
</html>
`

// ChangeLog represents a collection of changes between versions
type ChangeLog struct {
	NameTags string
	Commits  []CommitMessage
}

// CommitMessage contains the commit information
type CommitMessage struct {
	ShortCommit string `json:"shortcommit"`
	Commit      string `json:"commit"`
	Author      string `json:"author"`
	Email       string `json:"email"`
	Date        string `json:"date"`
	Message     string `json:"message"`
}

// Template for html format
type Template struct {
	Content string
}

// Report full report structure
type Report struct {
	pattern   string
	originURL string
	commitON  bool
	authorON  bool
	ChangeLog []ChangeLog
	JSON      string
	Markdown  string
	HTML      string
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

func fetchChanges(releasePattern string, limit int) []ChangeLog {
	var changeLogs []ChangeLog
	tags := fetchTags(releasePattern)

	var diffs []string
	for i := 0; i < len(tags)-1; i++ {
		diffs = append(diffs, fmt.Sprintf("%v..%v", tags[i+1], tags[i]))
		if len(diffs) >= limit {
			break
		}
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
	}

	return changeLogs
}

func (c *CommitMessage) toMarkdown(originURL string, commitON, authorON bool) string {
	space := regexp.MustCompile(`-`)
	var commitLink, authorLink string
	if commitON {
		commitLink = fmt.Sprintf("[%v](%vcommit/%v)", c.ShortCommit, originURL, c.Commit)
	}
	if authorON {
		authorLink = fmt.Sprintf("[%v](mailto:%v)", c.Author, c.Email)
	}
	return fmt.Sprintf("* %v %v %v\n", commitLink, space.ReplaceAllString(c.Message, " "), authorLink)
}

func fetchProjectURL() string {
	originURL, _ := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	re := regexp.MustCompile(`:`)
	originURL = []byte(re.ReplaceAllString(string(originURL), `/`))
	re = regexp.MustCompile(`\.git\n$`)
	originURL = []byte(re.ReplaceAllString(string(originURL), `/`))
	re = regexp.MustCompile(`^git@`)
	originURL = []byte(re.ReplaceAllString(string(originURL), `https://`))

	return string(originURL)
}

func buildReport(releasePattern string, commitLinksON, authorLinksON bool, limit int) Report {
	report := Report{
		pattern:   releasePattern,
		originURL: fetchProjectURL(),
		commitON:  commitLinksON,
		authorON:  authorLinksON,
	}
	report.ChangeLog = fetchChanges(releasePattern, limit)

	byteJSON, _ := json.Marshal(report.ChangeLog)
	report.JSON = string(byteJSON)

	return report
}

func (r *Report) getMarkdown() string {
	if len(r.Markdown) > 0 {
		return r.Markdown
	}
	var mdBuffer bytes.Buffer
	mdBuffer.WriteString("# Changelog\n")

	for _, change := range r.ChangeLog {
		mdBuffer.WriteString(fmt.Sprintf("\n## %v \n\n", change.NameTags))
		for _, commit := range change.Commits {
			mdBuffer.WriteString(commit.toMarkdown(string(r.originURL), r.commitON, r.authorON))
		}
	}
	r.Markdown = mdBuffer.String()
	return r.Markdown
}

func (r *Report) getHTML() string {
	if len(r.HTML) > 0 {
		return r.HTML
	}
	out := bytes.NewBuffer(nil)

	html := blackfriday.MarkdownCommon([]byte(r.getMarkdown()))

	t := template.Must(template.New("layout").Parse(layout))
	err := t.Execute(out, Template{Content: string(html)})
	if err != nil {
		panic(err)
	}
	r.HTML = out.String()
	return r.HTML
}

func main() {
	app := cli.NewApp()
	app.Name = "Changes reporter"
	app.Usage = "Generate changelog report from git commits and tag releases"
	app.Version = "0.0.2"
	app.Author = "Luis Ezcurdia"
	app.Email = "ing.ezcurdia@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pattern, p",
			Value: `v[\d{1,4}\.]{1,}`,
			Usage: "Regular expresion for release tags",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "markdown",
			Usage: "Output format for report",
		},
		cli.StringFlag{
			Name:  "limit, l",
			Value: "500",
			Usage: "Limit the report to a determinated number of versions",
		},
		cli.StringFlag{
			Name:  "no-commit",
			Value: "false",
			Usage: "Omit commit links on markdown and HTML reports",
		},
		cli.StringFlag{
			Name:  "no-author",
			Value: "false",
			Usage: "Omit author links on markdown and HTML reports",
		},
		cli.StringFlag{
			Name:  "only-message",
			Value: "false",
			Usage: "Only show commit messages on markdown and HTML reports",
		},
	}
	app.Action = func(c *cli.Context) {
		limit, _ := strconv.Atoi(c.String("limit"))
		commitLinksON := strings.ToLower(c.String("no-commit")) == "false"
		authorLinksON := strings.ToLower(c.String("no-author")) == "false"
		if strings.ToLower(c.String("only-message")) != "false" {
			commitLinksON = false
			authorLinksON = false
		}
		report := buildReport(c.String("pattern"), commitLinksON, authorLinksON, limit)
		switch strings.ToLower(c.String("format")) {
		case "json":
			fmt.Println(report.JSON)
		case "html":
			fmt.Print(report.getHTML())
		default:
			fmt.Println(report.getMarkdown())
		}
	}
	app.Run(os.Args)
}
