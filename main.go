package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//https://medium.com/@Flashleo/mounting-windows-file-share-on-docker-container-ac930092c0a5

type CsvConnection struct {
	ID                     string
	Name                   string
	ConnectionType         string
	FileDirectory          string
	ProcessedFileDirectory string
	ErrorDirectory         string
	Inputs                 []CsvConnectionInput
}

type CsvConnectionInput struct {
	ID        string
	Name      string
	FileName  string
	Delimeter string
	IndexFile bool
	Parts     []CsvConnectionInputPart
}

type CsvConnectionInputPart struct {
	SkipLines          int
	MaxRows            int
	MaxRowsPerRead     int
	ReplacementHeaders string
}

func main() {

	csvConnection1 := CsvConnection{
		ID:                     "1234",
		Name:                   "Maxymos",
		ConnectionType:         "CSV",
		FileDirectory:          "//10.214.6.148/Maxymos/",
		ProcessedFileDirectory: "//10.214.6.148/Maxymos/Succeed/",
		ErrorDirectory:         "//10.214.6.148//Maxymos/Failed/",
		Inputs:                 []CsvConnectionInput{},
	}

	csvConnection1Input1 := CsvConnectionInput{
		ID:        "123456",
		Name:      "Bearing",
		FileName:  `(Bearing\\)(2021-10-12\\)[a-zA-Z0-9\-\_]+.csv`,
		Delimeter: ";",
		IndexFile: false,
		Parts:     []CsvConnectionInputPart{},
	}

	csvConnection1Input1.Parts = append(csvConnection1Input1.Parts, CsvConnectionInputPart{
		SkipLines:          10,
		MaxRows:            10,
		MaxRowsPerRead:     10,
		ReplacementHeaders: "",
	})

	//p := filepath.FromSlash(csvConnection1.FileDirectory)
	//fileSearchPattern := "(" + strings.ReplaceAll(p, "\\", "\\\\") + ")" + csvConnection1Input1.FileName
	fileSearchPattern := `(\\\\10.214.6.148\\Maxymos\\)(Bearing\\)(2021-10-12\\)[a-zA-Z0-9-_]+.csv`
	fmt.Println(fileSearchPattern)

	re := regexp.MustCompile(fileSearchPattern)
	fileNames, err := filteredSearchOfDirectoryTree(re, csvConnection1.FileDirectory, 1)

	fmt.Println(fileNames)
	fmt.Println(err)
}

func filteredSearchOfDirectoryTree(re *regexp.Regexp, dir string, limit int) ([]string, error) {
	files := []string{}

	// Function variable that can be used to filter
	// files based on the pattern.
	// Note that it uses re internally to filter.
	// Also note that it populates the files variable with
	// the files that matches the pattern.
	walk := func(fn string, fi os.FileInfo, err error) error {
		if len(files) == limit {
			return nil
		}

		fmt.Println("****************")
		fmt.Println(re)
		fmt.Println(fn)
		if re.MatchString(strings.ReplaceAll(fn, "/", "\\")) == false {
			return nil
		}
		if fi.IsDir() {
			fmt.Println("2222222222222")
			fmt.Println(fn + string(os.PathSeparator))
		} else {
			fmt.Println("-------------")
			fmt.Println(fn)
			files = append(files, fn)
		}
		return nil
	}
	filepath.Walk(dir, walk)
	fmt.Printf("Found %[1]d files.\n", len(files))
	return files, nil
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	destPathParts := strings.Split(destPath, string(os.PathSeparator))
	destPathPartsLen := len(destPathParts)
	destPathWithoutFileName := ""

	for i, part := range destPathParts {
		if i > 0 && i < destPathPartsLen-1 {
			destPathWithoutFileName += string(os.PathSeparator) + part
		}
	}

	if _, err := os.Stat(destPathWithoutFileName); os.IsNotExist(err) {
		if err := os.MkdirAll(destPathWithoutFileName, 0700); err != nil {
			return fmt.Errorf("Couldn't create file: %s", err)
		}
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
