package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
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
	appDir := fmt.Sprintf("%s/%s", p.AbsolutPath, p.AppName)
	if _, err := os.Stat(p.AbsolutPath); err == nil {
		if err := os.Mkdir(appDir, 0755); err != nil {
			return err
		}
	}

	gitPull := exec.Command("git", "clone", "https://github.com/Melkeydev/ragStack.git", ".")
	gitPull.Dir = appDir
	var out bytes.Buffer
	gitPull.Stdout = &out
	if err := gitPull.Run(); err != nil {
		return err
	}
	fmt.Println(out.String())

	if err := os.RemoveAll(fmt.Sprintf("%s/.git", appDir)); err != nil {
		return err
	}
	return nil
}

