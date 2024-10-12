package common

import (
	"fmt"
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
