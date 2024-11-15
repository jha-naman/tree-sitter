package main

import (
	"bufio"
	"errors"
	"flag"
	"io/fs"
	"log"
	"os"
	"path"
	"slices"
	"sort"

	common "github.com/jha-naman/tree-tags/common"
	golang "github.com/jha-naman/tree-tags/golang"
)

var options = common.Options{}

func main() {

	initOptions()

	fileNames, err := getFileNames()
	if err != nil {
		log.Fatalf("error getting filenames: %s", err.Error())
	}

	tags, err := initTags(fileNames)
	if err != nil {
		log.Fatal("error while initialising tags:", err.Error())
	}

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

func initOptions() {
	flag.BoolVar(&options.AppendMode, "a", false, "shorthand form for 'append' option")
	flag.BoolVar(&options.AppendMode, "append", false, "add this flag to re-generate tags for given list of files instead of re-generating the tags file from scratch for the whole project, will remove stale tags belonging to the given list of files")

	flag.Parse()
}

func getFileNames() ([]string, error) {
	if options.AppendMode {
		fileNames := flag.Args()
		if len(fileNames) == 0 {
			log.Fatal("need to supply file names when used in append mode (using -a as command line option)")
		}

		return fileNames, nil
	}

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

func initTags(fileNamesToSkip []string) ([]common.TagEntry, error) {
	tags := []common.TagEntry{}

	if !options.AppendMode {
		return tags, nil
	}

	file, err := os.Open("tags")
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		tag, err := common.TagFromString(text)
		if errors.Is(err, common.ErrStringIsAComment) {
			continue
		}

		if !slices.Contains(fileNamesToSkip, tag.FileName) {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}
