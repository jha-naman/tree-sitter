package golang

import (
	"strings"
	"testing"

	"github.com/jha-naman/tree-tags/common"
	"github.com/stretchr/testify/assert"
)

func TestPackageDeclaration(t *testing.T) {
	input := "package treetags\n"
	expectedTags := []common.TagEntry{
		{
			Name:            "treetags",
			FileName:        "",
			Address:         "/^package treetags$/;\"",
			Kind:            "p",
			ExtensionFields: nil,
		},
	}

	assert.Equal(t, expectedTags, extractTagsFromString(input))
}

func TestImportDeclaration(t *testing.T) {
	tests := []struct {
		input        string
		expectedTags []common.TagEntry
	}{
		{
			input: `import assert "github.com/stretchr/testify/assert"`,
			expectedTags: []common.TagEntry{
				{Name: "assert", FileName: "", Address: `/^import assert "github.com\/stretchr\/testify\/assert"$/;"`, Kind: "P", ExtensionFields: map[string]string{"package": "github.com/stretchr/testify/assert"}},
			},
		},
		{
			input: `
			import (
				"sync"

				assert "github.com/stretchr/testify/assert"
			)
			`,
			expectedTags: []common.TagEntry{
				{
					Name:            "assert",
					FileName:        "",
					Address:         "/^\t\t\t\tassert \"github.com\\/stretchr\\/testify\\/assert\"$/;\"",
					Kind:            "P",
					ExtensionFields: map[string]string{"package": "github.com/stretchr/testify/assert"},
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedTags, extractTagsFromString(test.input))
	}
}

func TestFunctionDeclaration(t *testing.T) {
	tests := []struct {
		input        string
		expectedTags []common.TagEntry
	}{
		{
			input: `package main; func main() {}`,
			expectedTags: []common.TagEntry{
				{
					Name:            "main",
					FileName:        "",
					Address:         `/^package main; func main() {}$/;"`,
					Kind:            "p",
					ExtensionFields: nil,
				},
				{
					Name:            "main",
					FileName:        "",
					Address:         "/^package main; func main() {}$/;\"",
					Kind:            "f",
					ExtensionFields: map[string]string{"package": "main"},
				},
			},
		},
		{
			input: `package main; func foo(bar, baz string, arr []string) (error, map[string]string) {}`,
			expectedTags: []common.TagEntry{
				{
					Name:            "main",
					FileName:        "",
					Address:         `/^package main; func foo(bar, baz string, arr []string) (error, map[string]string) {}$/;"`,
					Kind:            "p",
					ExtensionFields: nil,
				},
				{
					Name:            "foo",
					FileName:        "",
					Address:         `/^package main; func foo(bar, baz string, arr []string) (error, map[string]string) {}$/;"`,
					Kind:            "f",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "(error, map[string]string)"},
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedTags, extractTagsFromString(test.input))
	}

}

func TestVarDeclaration(t *testing.T) {
	tests := []struct {
		input        string
		expectedTags []common.TagEntry
	}{
		{
			input: "package main; var x, y int",
			expectedTags: []common.TagEntry{
				{
					Name:            "main",
					FileName:        "",
					Address:         `/^package main; var x, y int$/;"`,
					Kind:            "p",
					ExtensionFields: nil,
				},

				{
					Name:            "x",
					FileName:        "",
					Address:         `/^package main; var x, y int$/;"`,
					Kind:            "v",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "int"},
				},
				{
					Name:            "y",
					FileName:        "",
					Address:         `/^package main; var x, y int$/;"`,
					Kind:            "v",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "int"},
				},
			},
		},
		{
			input: `
package main
var (
	a, b int
	x map[string]string
	i interface{}
	z = "zed"
)
			`,
			expectedTags: []common.TagEntry{
				{
					Name:            "main",
					FileName:        "",
					Address:         "/^package main$/;\"",
					Kind:            "p",
					ExtensionFields: nil,
				},
				{
					Name:            "a",
					FileName:        "",
					Address:         "/^\ta, b int$/;\"",
					Kind:            "v",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "int"},
				},
				{Name: "b",
					FileName:        "",
					Address:         "/^\ta, b int$/;\"",
					Kind:            "v",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "int"},
				},
				{
					Name:            "x",
					FileName:        "",
					Address:         "/^\tx map[string]string$/;\"",
					Kind:            "v",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "map[string]string"},
				},
				{
					Name:            "i",
					FileName:        "",
					Address:         "/^\ti interface{}$/;\"",
					Kind:            "v",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "interface{}"},
				},
				{
					Name:            "z",
					FileName:        "",
					Address:         "/^\tz = \"zed\"$/;\"",
					Kind:            "v",
					ExtensionFields: map[string]string{"package": "main"},
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedTags, extractTagsFromString(test.input))
	}
}

func TestConstDeclaration(t *testing.T) {
	tests := []struct {
		input        string
		expectedTags []common.TagEntry
	}{
		{
			input: `package main; const foo = "foo"`,
			expectedTags: []common.TagEntry{
				{
					Name:            "main",
					FileName:        "",
					Address:         `/^package main; const foo = "foo"$/;"`,
					Kind:            "p",
					ExtensionFields: nil,
				},

				{
					Name:            "foo",
					FileName:        "",
					Address:         `/^package main; const foo = "foo"$/;"`,
					Kind:            "c",
					ExtensionFields: map[string]string{"package": "main"},
				},
			},
		},
		{
			input: `
package main
const (
	foo = "foo"
	bar = 1
)
	`,
			expectedTags: []common.TagEntry{
				{
					Name:            "main",
					FileName:        "",
					Address:         `/^package main$/;"`,
					Kind:            "p",
					ExtensionFields: nil,
				},
				{
					Name:            "foo",
					FileName:        "",
					Address:         "/^\tfoo = \"foo\"$/;\"",
					Kind:            "c",
					ExtensionFields: map[string]string{"package": "main"},
				},
				{
					Name:            "bar",
					FileName:        "",
					Address:         "/^\tbar = 1$/;\"",
					Kind:            "c",
					ExtensionFields: map[string]string{"package": "main"},
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedTags, extractTagsFromString(test.input))
	}
}

func TestTypeDeclaration(t *testing.T) {
	tests := []struct {
		input        string
		expectedTags []common.TagEntry
	}{

		{
			input: "package main; type Alias int; type AnotherOne Alias",
			expectedTags: []common.TagEntry{
				{
					Name:            "main",
					FileName:        "",
					Address:         "/^package main; type Alias int; type AnotherOne Alias$/;\"",
					Kind:            "p",
					ExtensionFields: nil,
				},
				{
					Name:            "Alias",
					FileName:        "",
					Address:         "/^package main; type Alias int; type AnotherOne Alias$/;\"",
					Kind:            "t",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "int"},
				},
				{
					Name:            "AnotherOne",
					FileName:        "",
					Address:         "/^package main; type Alias int; type AnotherOne Alias$/;\"",
					Kind:            "t",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "Alias"},
				},
			},
		},
		{
			input: "package main; type Alias = map[string]string",
			expectedTags: []common.TagEntry{
				{
					Name:            "main",
					FileName:        "",
					Address:         "/^package main; type Alias = map[string]string$/;\"",
					Kind:            "p",
					ExtensionFields: nil,
				},
				{
					Name:            "Alias",
					FileName:        "",
					Address:         "/^package main; type Alias = map[string]string$/;\"",
					Kind:            "a",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "map[string]string"},
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedTags, extractTagsFromString(test.input))
	}
}

func TestMethodDeclaration(t *testing.T) {
	tests := []struct {
		input        string
		expectedTags []common.TagEntry
	}{

		{
			input: `
package main
type foo int
func (f foo) String() {}
func (f *foo) Bar(baz string) map[string]string { return nil }
			`,
			expectedTags: []common.TagEntry{
				{
					Name:            "main",
					FileName:        "",
					Address:         "/^package main$/;\"",
					Kind:            "p",
					ExtensionFields: nil,
				},
				{
					Name:            "foo",
					FileName:        "",
					Address:         "/^type foo int$/;\"",
					Kind:            "t",
					ExtensionFields: map[string]string{"package": "main", "typeref:typename": "int"},
				},
				{
					Name:            "String",
					FileName:        "",
					Address:         "/^func (f foo) String() {}$/;\"",
					Kind:            "f",
					ExtensionFields: map[string]string{"unkown": "main.foo"},
				},
				{
					Name:            "Bar",
					FileName:        "",
					Address:         "/^func (f *foo) Bar(baz string) map[string]string { return nil }$/;\"",
					Kind:            "f",
					ExtensionFields: map[string]string{"unkown": "main.*foo", "typeref:typename": "map[string]string"},
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedTags, extractTagsFromString(test.input))
	}
}

func extractTagsFromString(codeStr string) []common.TagEntry {
	var codeBytes [][]byte
	for _, line := range strings.Split(codeStr, "\n") {
		codeBytes = append(codeBytes, []byte(line))
	}

	p := Processor{FileBytes: codeBytes}
	return p.GetTags()
}
