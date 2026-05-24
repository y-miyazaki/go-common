## tutorial

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# <Tutorial Title: e.g., Getting Started with X / Building Your First Y>

<A brief, welcoming introduction. State exactly what the reader will build or achieve by the end of this tutorial, who it is for, and how long it should take (e.g., 10-minute guide).>

Focus on:
- learning-by-doing (practical experience)
- a single, reliable, failure-proof path
- immediate, visible results at each step
- concrete instructions (exact commands and expected outputs)

Avoid:
- deep conceptual or architectural explanations (leave this to Explanation)
- multiple choices or alternatives ("You could also use X or Y..." configures confusion)
- exhaustive reference material or edge cases (leave this to Reference)
- abstract or theoretical hand-waving

## Prerequisites

Before starting, ensure you have the following installed and configured:
- <Dependency/Tool 1 (e.g., Node.js v18+)>
- <Dependency/Tool 2 (e.g., A running local Docker instance)>
- <Initial state/Knowledge assumption (e.g., Basic familiarity with terminal)>

## Goal

By completing this tutorial, you will:
1. <Outcome 1: e.g., Set up a local development environment>
2. <Outcome 2: e.g., Create and deploy a single endpoint>
3. <Outcome 3: e.g., Verify the result via a curl request>

---

## Step 1: <Clear, Action-Oriented Title>

<Describe the minimal context needed for this step.>

Run the following command to <action>:

```bash
<exact command>
```

**Expected Output:**
```text
<paste the exact successful log/output here so the user knows they are on track>
```

## Step 2: <Clear, Action-Oriented Title>

<Provide the code snippet or configuration change. Keep it minimal and copy-pasteable.>

Create a file named `<filename>` and add the following content:

```<language>
<minimal, complete, and working code snippet>
```

## Step 3: <Clear, Action-Oriented Title>

<The final friction point—running or compiling the project.>

```bash
<exact command to run/verify>
```

---

## Verification (How to check your success)

To make sure everything is working perfectly, execute:

```bash
<verification command, e.g., curl http://localhost:8080/health>
```

You should see the following response:
```json
<expected successful response>
```

🎉 **Congratulations!** You have successfully completed the tutorial.

## Next Steps

Now that you have the basics down, explore these resources to go deeper:
- How-To: `<Title>` (`<relative-path-to-how-to.md>`) - For specific task recipes.
- Reference: `<Title>` (`<relative-path-to-reference.md>`) - To look up detailed configurations/API specs.
- Explanation: `<Title>` (`<relative-path-to-explanation.md>`) - To understand the architecture and design philosophy.

## Decision Prompts

Consider:
- Is this guide 100% reliable? If a beginner follows it exactly, will it work without errors?
- Did you eliminate all alternative choices? (e.g., choosing between npm vs yarn, Docker vs bare-metal—pick *one* for the tutorial).
- Are explanations kept to the absolute minimum required to complete the task?
- Is there an explicit verification step where the user gets immediate feedback?
```
