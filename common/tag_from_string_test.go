package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagFromString(t *testing.T) {
	text := `FileName	golang/extract_tags.go	/^	FileName    string$/;"	m	struct:golang.Processor	typeref:typename:string`
	tag, _ := TagFromString(text)
	expectedTag := TagEntry{
		"FileName",
		"golang/extract_tags.go",
		`/^	FileName    string$/;"`,
		"m",
		map[string]string{
			"struct": "golang.Processor",
			"typeref:typename": "string",
		},
	}

	assert.Equal(t, expectedTag, tag)
}
