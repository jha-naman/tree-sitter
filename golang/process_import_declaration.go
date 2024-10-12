package golang

import (
	common "github.com/jha-naman/tree-tags/common"
)

func (p *Processor) processImportDeclaration() {
	cursor := p.cursor
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	if !cursor.GoToNextSibling() {
		return
	}

	switch cursor.CurrentNode().Type() {
	case "import_spec":
		p.processImportSpec()
	case "import_spec_list":
		if !cursor.GoToFirstChild() {
			break
		}
		defer cursor.GoToParent()

		p.processImportSpec()
		for cursor.GoToNextSibling() {
			p.processImportSpec()
		}
	}
}

func (p *Processor) processImportSpec() {
	var tag common.TagEntry
	cursor := p.cursor
	node := cursor.CurrentNode()

	if node.Type() != "import_spec" || !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	node = cursor.CurrentNode()
	line := p.FileBytes[node.StartPoint().Row]
	if node.Type() != "package_identifier" {
		return
	}

	tag = common.TagEntry{
		Name:     string(line[node.StartPoint().Column:node.EndPoint().Column]),
		Kind:     "P",
		Address:  p.addressStringFromBytes(line),
		FileName: p.FileName,
	}

	cursor.GoToNextSibling()
	node = cursor.CurrentNode()

	tag.ExtensionFields = map[string]string{
		"package": string(line[node.StartPoint().Column+1:node.EndPoint().Column-1]),
	}

	p.Tags = append(p.Tags, tag)
}

