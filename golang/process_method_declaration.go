package golang

import (
	"fmt"

	common "github.com/jha-naman/tree-tags/common"
	sitter "github.com/smacker/go-tree-sitter"
)

func (p *Processor) processMethodDeclaration() {
	cursor := p.cursor
	parentNode := cursor.CurrentNode()
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	var name, address, typerefName, receiverType string

	processReceiver := func(node *sitter.Node) string {
		if !cursor.GoToFirstChild() {
			return ""
		}
		defer cursor.GoToParent()

		cursor.GoToNextSibling()

		node = cursor.CurrentNode()
		if !cursor.GoToFirstChild() {
			return ""
		}
		defer cursor.GoToParent()

		count := 0
		for cursor.GoToNextSibling() {
			count++
			childNode := cursor.CurrentNode()
			if node.FieldNameForChild(count) == "type" {
				return p.stringFromByteRange(p.FileBytes, childNode.Range())
			}
		}
		return ""
	}

	childCount := 1
	for cursor.GoToNextSibling() {
		node := cursor.CurrentNode()
		switch parentNode.FieldNameForChild(childCount) {
		case "name":
			name = p.stringFromByteRange(p.FileBytes, node.Range())
			address = p.addressStringFromBytes(p.FileBytes[node.StartPoint().Row])
		case "receiver":
			receiverType = processReceiver(node)
		case "result":
			typerefName = p.stringFromByteRange(p.FileBytes, node.Range())
		case "body":
			break
		}
		childCount++
	}

	tagEntry := common.TagEntry{
		Name:            name,
		FileName:        p.FileName,
		Address:         address,
		Kind:            "f",
		ExtensionFields: map[string]string{},
	}

	if typerefName != "" {
		tagEntry.ExtensionFields["typeref:typename"] = typerefName
	}

	if receiverType != "" {
		tagEntry.ExtensionFields["unkown"] = fmt.Sprintf("%s.%s", p.packageName, receiverType)
	}

	p.Tags = append(p.Tags, tagEntry)
}
