import * as cdk from "aws-cdk-lib";
import { RagStack } from "./rag";
import { RagStackFargate } from "./ragFargate";

const app = new cdk.App();
new RagStack(app, "RagStack");
new RagStackFargate(app, "RagStackFargate");

app.synth();
