package golang

import (
	"fmt"

	common "github.com/jha-naman/tree-tags/common"
)

// Example tree:
//
//	(type_declaration
//	    (type_spec
//	        name: (type_identifier)
//	        type: (struct_type
//	            (field_declaration_list
//	                (field_declaration
//	                    name: (field_identifier)
//	                    name: (field_identifier)
//	                    type: (type_identifier))
//	                (field_declaration
//	                    name: (field_identifier)
//	                    type: (type_identifier))))))
func (p *Processor) processTypeDeclaration() {
	cursor := p.cursor
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	for cursor.GoToNextSibling() {
		switch cursor.CurrentNode().Type() {
		case "type_spec":
			p.processTypeSpec()
		case "type_alias":
			p.processTypeAlias()
		}
	}
}

func (p *Processor) processTypeSpec() {
	cursor := p.cursor
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	node := cursor.CurrentNode()
	typeName := string(p.FileBytes[node.StartPoint().Row][node.StartPoint().Column:node.EndPoint().Column])

	for cursor.GoToNextSibling() {
		node = cursor.CurrentNode()
		switch node.Type() {
		case "type_identifier":
			parentNode := node.Parent()
			p.Tags = append(p.Tags, common.TagEntry{
				Name:            typeName,
				FileName:        p.FileName,
				Address:         p.addressStringFromBytes(p.FileBytes[parentNode.StartPoint().Row]),
				Kind:            "t",
				ExtensionFields: map[string]string{"package": p.packageName, "typeref:typename": p.stringFromByteRange(p.FileBytes, node.Range())},
			})
		case "struct_type":
			parentNode := node.Parent()
			p.Tags = append(p.Tags, common.TagEntry{
				Name:            typeName,
				FileName:        p.FileName,
				Address:         p.addressStringFromBytes(p.FileBytes[parentNode.StartPoint().Row]),
				Kind:            "s",
				ExtensionFields: map[string]string{"package": p.packageName},
			})
			p.processStructType(typeName)
		case "interface_type":
			parentNode := node.Parent()
			p.Tags = append(p.Tags, common.TagEntry{
				Name:            typeName,
				FileName:        p.FileName,
				Address:         p.addressStringFromBytes(p.FileBytes[parentNode.StartPoint().Row]),
				Kind:            "i",
				ExtensionFields: map[string]string{"package": p.packageName},
			})
			p.processInterfaceMethods(typeName)
		default:
			parentNode := node.Parent()
			p.Tags = append(p.Tags, common.TagEntry{
				Name:            typeName,
				FileName:        p.FileName,
				Address:         p.addressStringFromBytes(p.FileBytes[parentNode.StartPoint().Row]),
				Kind:            "a",
				ExtensionFields: map[string]string{"package": p.packageName, "typeref:typename": p.stringFromByteRange(p.FileBytes, node.Range())},
			})
		}
	}
}

func (p *Processor) processTypeAlias() {
	cursor := p.cursor
	parentNode := cursor.CurrentNode()
	childCount := 0
	var typeName, aliasedTypeName, address string

	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	node := cursor.CurrentNode()
	typeName = string(p.FileBytes[node.StartPoint().Row][node.StartPoint().Column:node.EndPoint().Column])
	address = p.addressStringFromBytes(p.FileBytes[node.StartPoint().Row])

	for cursor.GoToNextSibling() {
		childCount++
		node = cursor.CurrentNode()
		if parentNode.FieldNameForChild(childCount) == "type" {
			aliasedTypeName = p.stringFromByteRange(p.FileBytes, node.Range())
			break
		}
	}

	p.Tags = append(p.Tags, common.TagEntry{
		Name:            typeName,
		FileName:        p.FileName,
		Address:         address,
		Kind:            "a",
		ExtensionFields: map[string]string{"package": p.packageName, "typeref:typename": aliasedTypeName},
	})
}

// Example tree:
//
//	(struct_type
//	    (field_declaration_list
//	        (field_declaration
//	            name: (field_identifier)
//	            name: (field_identifier)
//	            type: (type_identifier))
//	        (field_declaration
//	            name: (field_identifier)
//	            type: (type_identifier))))
func (p *Processor) processStructType(typeName string) {
	cursor := p.cursor
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	// move cursor  to the 'field_declaration_list' node
	if !cursor.GoToNextSibling() {
		return
	}

	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	p.processFieldDeclaration(typeName)

	for cursor.GoToNextSibling() {
		p.processFieldDeclaration(typeName)
	}
}

func (p *Processor) processFieldDeclaration(typeName string) {
	cursor := p.cursor
	parentNode := cursor.CurrentNode()

	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	childCount := 0
	structFieldTags := []common.TagEntry{}
	var typeString string

	if parentNode.FieldNameForChild(childCount) == "name" {
		structFieldTags = append(structFieldTags, p.processFieldIdentifier(typeName))
	}

	for cursor.GoToNextSibling() {
		childCount++
		node := cursor.CurrentNode()
		switch parentNode.FieldNameForChild(childCount) {
		case "name":
			structFieldTags = append(structFieldTags, p.processFieldIdentifier(typeName))
		case "type":
			typeString = p.stringFromByteRange(p.FileBytes, node.Range())
		}
	}

	for _, tag := range structFieldTags {
		tag.ExtensionFields["typeref:typename"] = typeString
	}

	p.Tags = append(p.Tags, structFieldTags...)
}

func (p *Processor) processFieldIdentifier(typeName string) common.TagEntry {
	node := p.cursor.CurrentNode()
	line := p.FileBytes[node.StartPoint().Row]
	return common.TagEntry{
		Name:            string(line[node.StartPoint().Column:node.EndPoint().Column]),
		FileName:        p.FileName,
		Address:         p.addressStringFromBytes(line),
		Kind:            "m",
		ExtensionFields: map[string]string{"struct": fmt.Sprintf("%s.%s", p.packageName, typeName)},
	}
}

func (p *Processor) processInterfaceMethods(typeName string) {
	cursor := p.cursor
	if !cursor.GoToFirstChild() {
		return
	}
	defer cursor.GoToParent()

	for cursor.GoToNextSibling() {
		parentNode := cursor.CurrentNode()
		if parentNode.Type() != "method_elem" {
			continue
		}

		if !cursor.GoToFirstChild() {
			continue
		}
		defer cursor.GoToParent()

		node := cursor.CurrentNode()
		childCount := 0
		fnName, resultStr := p.stringFromByteRange(p.FileBytes, node.Range()), ""
		line := p.FileBytes[node.StartPoint().Row]

		for cursor.GoToNextSibling() {
			childCount++
			switch parentNode.FieldNameForChild(childCount) {
			case "parameters":
			case "result":
				resultStr = p.stringFromByteRange(p.FileBytes, cursor.CurrentNode().Range())
			}
		}

		tag := common.TagEntry{
			Name:            fnName,
			FileName:        p.FileName,
			Address:         p.addressStringFromBytes(line),
			Kind:            "n",
			ExtensionFields: map[string]string{"interface": fmt.Sprintf("%s.%s", p.packageName, typeName)},
		}

		if resultStr != "" {
			tag.ExtensionFields["typeref:typename"] = resultStr
		}

		p.Tags = append(p.Tags, tag)
	}
}
