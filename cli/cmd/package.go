package cmd

import (
	"embed"
	"io/fs"
	"os"

	"github.com/spf13/rag-cli/template"
)

type Options struct {
	Deploy  string
	Git     bool
	Redis   bool
	CORS    string
	AppName string
}

type Project struct {
	AppName     string
	AbsolutPath string
	Options     *Options
}

func (p *Project) Create() error {
	if _, err := os.Stat(p.AbsolutPath); err == nil {
		if err := os.Mkdir(p.AbsolutPath, 0755); err != nil {
			return err
		}
	}
	files, err := getAllFilenames(&template.Template)
	// todo
	// for now just correctly copy template as is w copy package
	return nil
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
