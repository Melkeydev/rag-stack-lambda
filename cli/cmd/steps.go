package cmd

type StepSchema struct {
	StepName string
	Options  []string
	Headers  string
	Field    *string
}

type Steps struct {
	Steps []StepSchema
}

func initSteps(options *Options) *Steps {
	steps := &Steps{
		Steps: []StepSchema{
			{
				StepName: "Git",
				Options:  []string{"Yes", "No thanks"},
				Headers:  "Do you want to initialize a new Git project?",
				Field:    &options.Git,
			},
			{
				StepName: "Server",
				Options:  []string{"AWS Lambda", "AWS EC2"},
				Headers:  "How do you want to deploy",
				Field:    &options.Deploy,
			},
			{
				StepName: "Redis",
				Options:  []string{"Yes", "No thanks"},
				Headers:  "Do yo want to use Redis?",
				Field:    &options.Redis,
			},
			{
				StepName: "CORS",
				Options:  []string{"Domain", "Protocol", "Port"},
				Headers:  "Which CORS will you need?",
				Field:    &options.CORS,
			},
		},
	}

	return steps
}
