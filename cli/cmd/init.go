/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/rag-cli/cmd/creator"
	multi "github.com/spf13/rag-cli/cmd/ui/multiSelect"
	"github.com/spf13/rag-cli/cmd/ui/spinner"
	textinput "github.com/spf13/rag-cli/cmd/ui/textInput"
)

type ProjectSchema struct {
	Name   string
	Deploy string
	Redis  string
	CORS   string
	Git    string
}

var (
	logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true).Padding(1)
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logo := `
██████╗  █████╗  ██████╗ 
██╔══██╗██╔══██╗██╔════╝ 
██████╔╝███████║██║  ███╗
██╔══██╗██╔══██║██║   ██║
██║  ██║██║  ██║╚██████╔╝
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ 
		`
		UserWantsToExit := creator.ExitProgram{Value: false}
		fmt.Printf("%s\n", logoStyle.Render(logo))
		var program *tea.Program
		myProject := ProjectSchema{}
		projectName := &textinput.Output{}

		program = tea.NewProgram(textinput.InitialModelTextInput("myAwesomeApp", "What is your project name?", projectName, &UserWantsToExit))
		if _, err := program.Run(); err != nil {
			cobra.CheckErr(err)
		}

		UserWantsToExit.CheckExitStatus(program)
		myProject.Name = projectName.Output

		spec := Options{}
		Steps := initSteps(&spec)

		for _, step := range Steps.Steps {
			s := &multi.Selection{}
			program = tea.NewProgram(multi.InitialModelMulti(step.Options, s, step.Headers, &UserWantsToExit))
			if _, err := program.Run(); err != nil {
				cobra.CheckErr(err)
			}

			UserWantsToExit.CheckExitStatus(program)

			*step.Field = s.Choice
		}

		project := Project{
			AppName: myProject.Name,
			Options: &spec,
		}

		currentWorkingDir, err := os.Getwd()
		if err != nil {
			cobra.CheckErr(err)
		}

		project.AbsolutPath = currentWorkingDir

		var initIOWg sync.WaitGroup
		initIOWg.Add(1)
		projectLoading := spinner.LoadingState{Loading: true}
		go project.Create(&initIOWg, &projectLoading)
		program = tea.NewProgram(spinner.InitModelSpinner("Generating project", &projectLoading, &UserWantsToExit))
		if _, err := program.Run(); err != nil {
			cobra.CheckErr(err)
		}

		UserWantsToExit.CheckExitStatus(program)

		initIOWg.Wait()
		program.ReleaseTerminal()

		if err != nil {
			cobra.CheckErr(err)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
