package extract

import (
	"bytes"
	"fmt"
)

func Metadata(input []byte) ([]byte, []byte, error) {
	openingLoc := bytes.Index(input, frontmatterSep)
	closingLoc := bytes.LastIndex(input, frontmatterSep)

	if openingLoc == closingLoc {
		return nil, nil, fmt.Errorf("could not determine metadata location: could not find opening and closing front matter separators\n input: \n%s", string(input))
	}

	metadataBytes := input[openingLoc:closingLoc]
	remainder := input[closingLoc+len(frontmatterSep):]

	return metadataBytes, remainder, nil
}

var frontmatterSep = []byte("---")
