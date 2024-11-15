package common

import (
	"errors"
	"strings"
)

var ErrStringIsAComment = errors.New("cannot create tag for a comment")

func TagFromString(text string) (TagEntry, error) {
	if strings.HasPrefix(text, "!_TAG_") {
		return TagEntry{}, ErrStringIsAComment
	}

	fields := []string{"Name", "FileName", "Address", "Kind", "ExtensionFields"}
	fieldIndex := 0
	tag := TagEntry{}
	var theOneBeforeChar, previousChar rune
	var fieldAggregator string
	var extensionFields = make(map[string]string)

	for _, runeValue := range text {
		switch fields[fieldIndex] {
		case "Name", "FileName", "Kind":
			if runeValue != '\t' {
				fieldAggregator += string(runeValue)
				continue
			}

			tag.SetFieldByName(fields[fieldIndex], fieldAggregator)
			fieldIndex++
			fieldAggregator = ""
		case "Address":
			if theOneBeforeChar == ';' && previousChar == '"' && runeValue == '\t' {
				tag.SetFieldByName(fields[fieldIndex], fieldAggregator)
				fieldAggregator = ""
				fieldIndex++
				continue
			}

			fieldAggregator += string(runeValue)
			theOneBeforeChar = previousChar
			previousChar = runeValue
		case "ExtensionFields":
			if runeValue != '\t' {
				fieldAggregator += string(runeValue)
				continue
			}

			extensionFieldKey, extextensionFieldVal := extensionFieldFromAggregator(fieldAggregator)
			extensionFields[extensionFieldKey] = extextensionFieldVal

			fieldAggregator = ""
		}
	}

	if fieldAggregator != "" {
		switch fields[fieldIndex] {
		case "ExtensionFields":
			extensionFieldKey, extextensionFieldVal := extensionFieldFromAggregator(fieldAggregator)
			extensionFields[extensionFieldKey] = extextensionFieldVal

			tag.SetFieldByName("ExtensionFields", extensionFields)
		default:
			tag.SetFieldByName(fields[fieldIndex], fieldAggregator)
		}
	}

	return tag, nil
}

func extensionFieldFromAggregator(fieldAggregator string) (key, value string) {
	splits := strings.Split(fieldAggregator, ":")
	splitsCount := len(splits)
	if splitsCount == 2 {
		return splits[0], splits[1]
	} else {
		return strings.Join(splits[:2], ":"), splits[splitsCount-1]
	}
}
