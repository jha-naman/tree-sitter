package golang

import (
	common "github.com/jha-naman/tree-tags/common"
	sitter "github.com/smacker/go-tree-sitter"
)

func (p *Processor) processVarDeclaration() {
	cursor := p.cursor
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	p.processVarSpec()

	for cursor.GoToNextSibling() {
		switch cursor.CurrentNode().Type() {
		case "var_spec":
			p.processVarSpec()
		case "var_spec_list":
			if !cursor.GoToFirstChild() {
				continue
			}
			defer cursor.GoToParent()

			p.processVarSpec()
			for cursor.GoToNextSibling() {
				p.processVarSpec()
			}
		}
	}
}

func (p *Processor) processVarSpec() {
	identifierTags := []common.TagEntry{}
	var typeIdentifier string

	processIdentifier := func(node *sitter.Node) {
		line := p.FileBytes[node.StartPoint().Row]
		identifierTags = append(identifierTags, common.TagEntry{
			Name:            string(line[node.StartPoint().Column:node.EndPoint().Column]),
			FileName:        p.FileName,
			Address:         p.addressStringFromBytes(line),
			Kind:            "v",
			ExtensionFields: map[string]string{"package": p.packageName},
		})
	}

	cursor := p.cursor
	childCount := 0
	parentNode := cursor.CurrentNode()
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	node := cursor.CurrentNode()
	if parentNode.FieldNameForChild(childCount) == "name" {
		processIdentifier(node)
	}

	for cursor.GoToNextSibling() {
		childCount++

		node = cursor.CurrentNode()
		switch parentNode.FieldNameForChild(childCount) {
		case "name":
			processIdentifier(node)
		case "type":
			typeIdentifier = p.stringFromByteRange(p.FileBytes, node.Range())
		}
	}

	if typeIdentifier != "" {
		for _, tag := range identifierTags {
			tag.ExtensionFields["typeref:typename"] = typeIdentifier
		}
	}

	p.Tags = append(p.Tags, identifierTags...)
}
