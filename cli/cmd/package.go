package cmd

import (
	"fmt"
	cp "github.com/otiai10/copy"
	"os"
)

type Options struct {
	EC2   bool
	Lamda bool
	Redis bool
	CORS  bool
}

type Project struct {
	AppName     string
	AbsolutPath string
	Options
}

func (p  *Project) Create() error {
	if _, err := os.Stat(p.AbsolutPath); err == nil {
		if err := os.Mkdir(p.AbsolutPath, 0755); err != nil {
			return err
		}
	}

	// todo
	// for now just correctly copy template as is w copy package
}
