package golang

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"

	common "github.com/jha-naman/tree-tags/common"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type Processor struct {
	Tags        []common.TagEntry
	FileBytes   [][]byte
	FileName    string
	packageName string
	cursor      *sitter.TreeCursor
}

func GetFileTags(fileName string) []common.TagEntry {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("error while trying to read file:", file, err.Error())
	}

	var fileBytes [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileBytes = append(fileBytes, slices.Clone(scanner.Bytes()))
	}

	p := Processor{FileName: fileName, FileBytes: fileBytes}
	return p.GetTags()
}

func (p *Processor) GetTags() []common.TagEntry {
	parser := getGolangParser()
	tree, err := parser.ParseCtx(context.TODO(), nil, bytes.Join(p.FileBytes, []byte("\n")))
	if err != nil {
		log.Fatal("error while parsing file:", p.FileName, err.Error())
	}

	p.cursor = sitter.NewTreeCursor(tree.RootNode())
	p.extractTags()

	return p.Tags
}

func (p *Processor) extractTags() {
	cursor := p.cursor
	node := cursor.CurrentNode()

	switch node.Type() {
	case "package_clause":
		p.processPackageDeclaration()
	case "import_declaration":
		p.processImportDeclaration()
	case "function_declaration":
		p.processFunctionDeclaration()
	case "var_declaration":
		p.processVarDeclaration()
	case "const_declaration":
		p.processConstDeclaration()
	case "type_declaration":
		p.processTypeDeclaration()
	case "method_declaration":
		p.processMethodDeclaration()
	}

	if cursor.GoToNextSibling() {
		p.extractTags()
	} else if cursor.GoToFirstChild() {
		p.extractTags()
	}
}

func getGolangParser() *sitter.Parser {
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	return parser
}

var charsEscapeRegex = regexp.MustCompile("([$/])")
var replaceRegex = []byte("\\${1}")

func (p *Processor) addressStringFromBytes(rawBytes []byte) string {
	return fmt.Sprintf("/^%s$/%s", string(charsEscapeRegex.ReplaceAll(rawBytes, replaceRegex)), ";\"")
}

func (p *Processor) stringFromByteRange(fileBytes [][]byte, nodeRange sitter.Range) string {
	rowStart, rowEnd := nodeRange.StartPoint.Row, nodeRange.EndPoint.Row

	if rowStart == rowEnd {
		return string(fileBytes[rowStart][nodeRange.StartPoint.Column:nodeRange.EndPoint.Column])
	}

	var byteStr []byte
	for i := rowStart; i <= rowEnd; i++ {
		if i == rowStart {
			byteStr = append(byteStr, fileBytes[i][nodeRange.StartPoint.Column:]...)
		} else if i == rowEnd {
			byteStr = append(byteStr, fileBytes[i][:nodeRange.EndPoint.Column]...)
		} else {
			byteStr = append(byteStr, fileBytes[i]...)
		}
	}

	return string(byteStr)
}
