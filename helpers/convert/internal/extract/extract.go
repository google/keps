package extract

import (
	"bytes"
	"regexp"
	"strings"
)

// breadcrumb: https://play.golang.org/p/_trR0hWp69M
const (
	markdownTitle           = `# [\w\s]+\n`
	levelTwoMarkdownHeading = `\s## [\w\s]+\n|\s## [\w\s]+ \[optional\]\n`
)

// TODO add Title(input []byte) ([]byte, error)

func Sections(input []byte) (map[string][]byte, error) {

	// if you look at the KEP template
	//	https://raw.githubusercontent.com/kubernetes/community/master/keps/0000-kep-template.md
	// you'll notice that all the headings we're likely to want to extract are at the second
	// heading level
	topLevelSection := regexp.MustCompile(levelTwoMarkdownHeading)
	headers := topLevelSection.FindAllIndex(input, -1)

	sections := map[string][]byte{}

	for i, headingLoc := range headers {
		headingStart := headingLoc[0]
		headingEnd := headingLoc[1]

		headingName := string(input[headingStart:headingEnd])
		headingName = strings.Replace(headingName, "[optional]", "", -1)
		headingName = strings.TrimSpace(headingName)

		var content []byte
		switch i {
		case len(headers) - 1:
			content = input[headingEnd:]
		default:
			nextHeadingLoc := headers[i+1]
			nextHeadingStart := nextHeadingLoc[0]
			content = input[headingEnd:nextHeadingStart]
		}

		sections[headingName] = bytes.TrimSpace(content)
	}

	return sections, nil
}

const (
	TableOfContentsHeading       = "## Table of Contents"
	SummaryHeading               = "## Summary"
	MotivationHeading            = "## Motivation"
	ProposalHeading              = "## Proposal"
	GraduationCriteriaHeading    = "## Graduation Criteria"
	ImplementationHistoryHeading = "## Implementation History"
	DrawbacksHeading             = "## Drawbacks"
	AlternativesHeading          = "## Alternatives"
	InfrastructureNeededHeading  = "## Infrastructure Needed"
)
