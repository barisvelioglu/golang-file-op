package main

import (
	"fmt"
	"os"
	"path/filepath"
)

//https://medium.com/@Flashleo/mounting-windows-file-share-on-docker-container-ac930092c0a5

func main() {

	var files []string

	root := `//10.214.6.148/Maxymos/`
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
