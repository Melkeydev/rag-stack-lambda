import * as cdk from "aws-cdk-lib";
import { RagStack } from "./rag";

const app = new cdk.App();
new RagStack(app, "RagStack");

app.synth();
