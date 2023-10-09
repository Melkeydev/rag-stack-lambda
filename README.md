# Welcome to Rag Stack Serverless Service

This is an example of an end-to-end stack that uses React, AWS and Go for a fully scalable and hosted applicaton
This stacks comes with middleware, protected routes, logging in and registering users to a Dynamo Database.

This stack consists of:

- Vite, React, Tailwind on the Frontend
- Go, AWS SDK on the Backend
- DynamoDB, Lambda, API Gateway and CloudFront on the infrastructure

Lambda is the deployment of choice. If you want to deploy your code to a Fargate service, check out [this repo](https://github.com/Melkeydev/rag-stack-fargate)

## Prerequisites

- AWS CDK and Typescript should be installed on your system.

- AWS credentials should be configured on your system.

- A domain name registered with Route53 with a hosted zone.

## Deployment

Make sure to build the main binary file with: `go build -o main`

in the `lambda/cmd` route

## Useful commands

- `npm run build` compile typescript to js
- `npm run watch` watch for changes and compile
- `npm run test` perform the jest unit tests
- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template

## Clean up

Run the following command to delete the stack.

- `cdk destroy`
