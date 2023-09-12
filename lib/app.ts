import * as cdk from "aws-cdk-lib";

import { RagStack } from "./rag";
// TODO: remove this in the future
import { RagStackCdkStack } from "./rag_stack_cdk-stack";

const app = new cdk.App();

new RagStackCdkStack(app, "RagStackOld");
new RagStack(app, "RagStack");

app.synth();
