package golang

import (
	common "github.com/jha-naman/tree-tags/common"
	sitter "github.com/smacker/go-tree-sitter"
)

func (p *Processor) processConstDeclaration() {
	cursor := p.cursor
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	p.processConstSpec()

	for cursor.GoToNextSibling() {
		p.processConstSpec()
	}
}

func (p *Processor) processConstSpec() {
	cursor := p.cursor
	identifierTags := []common.TagEntry{}

	processIdentifier := func(node *sitter.Node) {
		line := p.FileBytes[node.StartPoint().Row]
		identifierTags = append(identifierTags, common.TagEntry{
			Name:            string(line[node.StartPoint().Column:node.EndPoint().Column]),
			FileName:        p.FileName,
			Address:         p.addressStringFromBytes(line),
			Kind:            "c",
			ExtensionFields: map[string]string{"package": p.packageName},
		})
	}

	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	node := cursor.CurrentNode()
	if node.Type() == "identifier" {
		processIdentifier(node)
	}

	for cursor.GoToNextSibling() {
		node := cursor.CurrentNode()
		switch node.Type() {
		case "identifier":
			processIdentifier(node)
		}
	}

	p.Tags = append(p.Tags, identifierTags...)
}

