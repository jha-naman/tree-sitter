package common

import (
	"fmt"
	"log"
	"strings"
)

type TagEntry struct {
	Name, FileName, Address, Kind string
	ExtensionFields               map[string]string
}

func (t TagEntry) Bytes() []byte {
	tagFields := []string{
		t.Name,
		t.FileName,
		t.Address,
		t.Kind,
	}

	for k, v := range t.ExtensionFields {
		tagFields = append(tagFields, fmt.Sprintf("%s:%s", k, v))
	}

	return []byte(strings.Join(tagFields, "\t"))
}

var allowedFieldNames = []string{"Name", "FileName", "Address", "Kind", "ExtensionFields"}

func (t *TagEntry) SetFieldByName(fieldName string, value interface{}) {
	switch fieldName {
	case "Name":
		t.Name = value.(string)
	case "FileName":
		t.FileName = value.(string)
	case "Address":
		t.Address = value.(string)
	case "Kind":
		t.Kind = value.(string)
	case "ExtensionFields":
		t.ExtensionFields = value.(map[string]string)
	default:
		log.Fatalf("invalid field name %s. should be one of %v", fieldName, allowedFieldNames)
	}
}
