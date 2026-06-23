# llm-supermarket

This is a basic coding test for a real world CLI task, used to compare LLM coding models and find out if:

1. A model can complete the task.
2. How the model costs to complete the task.

The task is not overly complicated and has examples on GitHub already that models can borrow from. The task should be something models can complete, but at the same time isn't a basic kata. Obviously it's only testing how good models are creating CLIs for a particular domain, but it's intended as a rough estimation of whether a model can autonomously create the command line tool, and how much it costs.

If models and the tooling gradually improve the price should drop, number of tokens used should drop, and completeness should improve.

## Instructions

1. Create a new Github repository called "rclone-encrypt" or similar.
1. Add a README. Paste in the contents of `TEST_README.md` or `HARDMODE_TEST_README.md` into the repository's readme and push.
1. Commit the test file "test-files/kr9tu4e1da4u3nifdd99g9tf5o" to the repository and push.
1. Commit the test file "test-files/Iyxcijgc9bp3o5Y0npW6xqUvwWNcc3MA4SadB0sR6cY" to the repository and push.
1. Clone the repository.
1. Pick a model and paste in `TEST_PROMPT.txt` or `HARDMODE_TEST_PROMPT.txt` and let it run.

See how much this cost in tokens - it won't be expensive (unless something went wrong), a few dollars at most.

## Results

Some comparisons will live here.

**Grok Build 0.1**

- 74,238 tokens
- 29% used
- $1.16 spent
- ~15 minutes
- Need to be re-run as it may have used 3xGrok 4.3 agents

**GLM-5.2**

- 178,847 tokens
- 18% used
- $4.71 spent
- ~40 minutes 31 seconds
- It worked out how to merge automatically, and that the app name was incorrect.

**GPT-5.1 Codex Mini**

76,671 tokens
19% used
$0.52 spent

- Didn't provide a session time.
- Didn't merge the changes, I merged and had to re-prompt. In its defence, the prompt was missing this.