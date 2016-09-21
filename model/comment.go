package model

import (
	"regexp"
	"strings"
)

type Comment struct {
	Title string
	Body  string
	Args  map[string]string
}

var (
	sectionRegex = regexp.MustCompile(`\n{2,}`)
	newlineRegex = regexp.MustCompile(`\n`)
	argsRegex    = regexp.MustCompile(`Args:`)
	argRegex     = regexp.MustCompile(`^\s+(.*?)\s+(.*)`)
)

func parseArgs(section string, args map[string]string) {
	lines := newlineRegex.Split(section, -1)
	// Toss first line ("Args:")
	for _, line := range lines[1:] {
		matches := argRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			args[matches[1]] = matches[2]
		}
	}
}

func ParseComment(text string) *Comment {
	sections := sectionRegex.Split(text, -1)
	var bodySections []string
	args := make(map[string]string)

	for _, section := range sections[1:] {
		if argsRegex.FindString(section) != "" {
			parseArgs(section, args)
		} else {
			bodySections = append(bodySections, section)
		}
	}

	return &Comment{
		Title: sections[0],
		Body:  strings.Join(bodySections, "\n\n"),
		Args:  args,
	}
}
