/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
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

		fmt.Printf("%s\n", logoStyle.Render(logo))
		myProject := ProjectSchema{}
		projectName := &textinput.Output{}
		textinput.TextInputRun("myAwesomeApp", "What is your project name?", projectName)
		myProject.Name = projectName.Output

		selectedGit := &multi.Selection{}
		gitOptions := []string{"Yes", "No thanks"}
		gitHeader := fmt.Sprintf("Do you want to initialize git with %s?", myProject.Name)
		multi.MultiBoxSelectRun(gitOptions, selectedGit, gitHeader)
		myProject.Git = selectedGit.Choice

		step := 0
		for step <= 2 {

			s := &multi.Selection{}
			var options []string
			var header string
			switch step {
			case 0:
				options = []string{"AWS Lambda", "AWS EC2"}
				header = fmt.Sprintf("How do you want to deploy %s?", myProject.Name)
				multi.MultiBoxSelectRun(options, s, header)
				myProject.Deploy = s.Choice
			case 1:
				options = []string{"Yes", "No thanks"}
				header = fmt.Sprintf("Do you want to use redis with %s?", myProject.Name)
				multi.MultiBoxSelectRun(options, s, header)
				myProject.Redis = s.Choice
			case 2:
				options = []string{"Domain", "Protocol", "Port"}
				header = fmt.Sprintf("Which CORS policy will %s use?", myProject.Name)
				multi.MultiBoxSelectRun(options, s, header)
				myProject.CORS = s.Choice
			}
			step++
		}
		spec := Options{
			Deploy: myProject.Deploy,
			Redis:  true,
			CORS:   myProject.CORS,
			Git:    myProject.Git == "Yes",
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

		var wg sync.WaitGroup
		wg.Add(1)
		projectLoading := spinner.LoadingState{Loading: true}
		loadingProgram := spinner.SpineMe("Generating project", &projectLoading, &wg)

		err = project.Create()

		projectLoading.Loading = false
		wg.Wait()
		loadingProgram.ReleaseTerminal()
		loadingProgram.RestoreTerminal()

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
