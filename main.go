package main

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path"
	"sort"

	common "github.com/jha-naman/tree-tags/common"
	golang "github.com/jha-naman/tree-tags/golang"
)

func main() {
	fileNames, err := getFileNames()
	if err != nil {
		log.Fatalf("error getting filenames: %s", err.Error())
	}

	var tags []common.TagEntry

	for _, fileName := range fileNames {
		tags = append(tags, golang.GetFileTags(fileName)...)
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	tagFile, err := os.Create("tags")
	if err != nil {
		log.Fatal("error while trying to create tag file:", err.Error())
	}

	writer := bufio.NewWriter(tagFile)

	for _, tag := range tags {
		if _, err = writer.Write(append(tag.Bytes(), []byte("\n")...)); err != nil {
			log.Fatal("error while trying to write tag file:", err.Error())
		}
	}

	if err = writer.Flush(); err != nil {
		log.Fatal("error while trying to write tag file:", err.Error())
	}
}

func getFileNames() ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var matchingFiles []string

	fs.WalkDir(os.DirFS(wd), ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		ext := path.Ext(filePath)
		if !d.IsDir() && ext == ".go" {
			matchingFiles = append(matchingFiles, filePath)
		}

		return nil
	})

	return matchingFiles, nil
}
