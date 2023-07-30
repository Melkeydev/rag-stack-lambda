deploy:
	cdk deploy

# this has to cd in the lambda folder and then execute command
build:
	go build -o main