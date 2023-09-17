/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"embed"
	"github.com/spf13/rag-cli/cmd"
	"io/fs"
	"fmt"
)

//go:embed template
var Template embed.FS

func main() {
	fmt.Println(getAllFilenames(&Template))
	cmd.Execute()
}

func getAllFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}
