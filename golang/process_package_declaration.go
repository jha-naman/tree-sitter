package golang

import (
	common "github.com/jha-naman/tree-tags/common"
)

func (p *Processor) processPackageDeclaration() {
	cursor := p.cursor
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	if !cursor.GoToNextSibling() {
		return
	}

	node := cursor.CurrentNode()
	lineBytes := p.FileBytes[node.StartPoint().Row]
	p.packageName = string(lineBytes[node.StartPoint().Column:node.EndPoint().Column])

	tag := common.TagEntry{
		Name:     p.packageName,
		FileName: p.FileName,
		Address:  p.addressStringFromBytes(lineBytes),
		Kind:     "p",
	}

	p.Tags = append(p.Tags, tag)
}

