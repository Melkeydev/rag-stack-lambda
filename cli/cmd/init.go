/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	multi "github.com/spf13/rag-cli/cmd/ui/multiSelect"
	textinput "github.com/spf13/rag-cli/cmd/ui/textInput"
)

type ProjectSchema struct {
	Name   string
	Deploy string
	Redis  string
	CORS   string
	Git    string
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//workDir, err := os.Getwd()
		//if err != nil {
		//	return
		//}
		//fmt.Println(workDir)

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
		fmt.Println(myProject)
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
