package golang

import (
	common "github.com/jha-naman/tree-tags/common"
)

func (p *Processor) processFunctionDeclaration() {
	cursor := p.cursor
	parentNode := cursor.CurrentNode()

	cursor.GoToFirstChild()
	defer cursor.GoToParent()

	childCount := 1
	var fnName, result string
	var line []byte

	for cursor.GoToNextSibling() {
		currentNode := cursor.CurrentNode()

		switch parentNode.FieldNameForChild(childCount) {
		case "name":
			line = p.FileBytes[currentNode.StartPoint().Row]
			fnName = string(line[currentNode.StartPoint().Column:currentNode.EndPoint().Column])
		case "result":
			result = p.stringFromByteRange(p.FileBytes, currentNode.Range())
		}

		childCount++
	}

	tag := common.TagEntry{
		Name:            fnName,
		FileName:        p.FileName,
		Address:         p.addressStringFromBytes(line),
		Kind:            "f",
		ExtensionFields: map[string]string{"package": p.packageName},
	}

	if result != "" {
		tag.ExtensionFields["typeref:typename"] = result
	}

	p.Tags = append(p.Tags, tag)
}

