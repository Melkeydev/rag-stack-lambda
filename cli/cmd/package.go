package cmd

import (
	"bytes"
	"fmt"
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

	if err := p.executeCmd("git",
		[]string{"clone", "--depth", "1", "-b", "main", "https://github.com/Melkeydev/ragStack.git", "."},
		appDir); err != nil {
		return err
	}

	if err := os.RemoveAll(fmt.Sprintf("%s/.git", appDir)); err != nil {
		return err
	}

	if p.Options.Git {
		if err := p.executeCmd("git", []string{"init"}, appDir); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) executeCmd(name string, args []string, dir string) error {
	command := exec.Command(name, args...)
	command.Dir = dir
	var out bytes.Buffer
	command.Stdout = &out
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}
