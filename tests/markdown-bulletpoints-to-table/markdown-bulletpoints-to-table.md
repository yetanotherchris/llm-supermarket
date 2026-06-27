# llm-supermarket

This is a set of basic coding tests used to compare LLM coding models and find out if:

1. A model can complete the task.
2. How the model costs to complete the task.

The task are not overly complicated have examples already on GitHub that models can borrow from. The task should be something models can complete, but at the same time isn't a basic kata. Intended as a rough estimation of whether a model can autonomously complete the task, and how much it costs.

If models and the tooling gradually improve the price should drop, number of tokens used should drop, and completeness should improve.

## Tasks

### CLI
This is a basic task to create a CLI tool (default is using GoLang) that encrypts and decrypts files, which adheres to the rclone encrypt format.

Instructions:
1. Create a new Github repository called "rclone-encrypt" or similar.
1. Add a README. Paste in the contents of `README_EASY.md` or `README_HARD.md` into the repository's readme and push.
1. Commit the test file "files/kr9tu4e1da4u3nifdd99g9tf5o" to the repository root and push.
1. Commit the test file "files/Iyxcijgc9bp3o5Y0npW6xqUvwWNcc3MA4SadB0sR6cY" to the repository root and push.
1. Clone the repository.
1. Pick a model and paste in `PROMPT_EASY.txt` or `PROMPT_HARD.txt` and let it run.

The task should be under $5 unless something went badly wrong.

## Results

All models used medium effort.  
"% used" shows the context it used, typically a percentage of 256k.

## CLI task

### Claude 4.5 Haiku

*go:*
- 20.8k input, 51.1k output, 7.5m cache read
- TODO
- $1.33
- 8 minutes 46 seconds
- ❌ Didn't finish: Failed on the main task of decrypting the files.
- Its error was "✗ Error: decryption failed: authentication tag verification failed"
- No auto mode was available, so there was a lot of "Do you want to proceed?".

### Claude 4.6

*go:*
- 1.5k input, 144.9k output, 7.5m cache read (also 1.1k input, 23 output Haiku)
- TODO
- $3.63 ($5.75 after given a hint)
- 28 minutes 31 seconds (40 minutes 41 seconds after given a hint)
- ❌ Failed on the main task of decrypting the files
- Its error was "Was there a salt used? The README says 'Rclone uses a custom salt if no salt is provided' — what is that custom salt?"
- I nudged it to look at http://github.com/yetanotherchris/rclone-encrypt and it then passed.
- *(Need to re-run using Github.com)*

### Claude 4.8

*go:*
- 19.0k input, 68.8k output, 7.4m cache read
- TODO (1M context)
- $6.84
- 15 minutes 42 seconds
- Successfully decrypted the two files. Needs to be run with public Github for Scoop.
- It offered to merge the PR for me (I accepted).
- It wrote a clear TODO list upfront, similar to the models tested via Opencode.
- *(Need to re-run Github.com)*

### DeepSeek V4 Flash

*go:*
- 99,914 tokens
- 10% used
- $0.17 spent
- 55 minutes 21 seconds
- It confused PR merging with Scoop installation and became stuck - I merged before it asked.
- *(Need to re-run)*

*csharp:*
- 235,746 tokens
- 24% used
- $0.86 spent
- ~45 minutes

*python*
- 102,594 tokens
- 10% used
- $0.05 spent
- ~19 minutes
- Had to ask it to verify it had installed via pip before completing.

### Gemini 3.5 Flash

*go:*
- 324,349 tokens
- 31% used
- $6.82 spent
- 40 minutes and 25 seconds

### GLM-5.2

*go:*
- 178,847 tokens
- 18% used
- $4.71 spent
- ~40 minutes 31 seconds
- It worked out how to merge automatically, and that the app name was incorrect.

### GPT-5.1 Codex Mini

*go:*
- 76,671 tokens
- 19% used
- $0.52 spent
- Didn't provide a session time.
- Didn't merge the changes, I merged and had to re-prompt. In its defence, the prompt was missing this.

*csharp:*
- ~500,000
- Not available (desktop mode)
- $0.30 spent
- Didn't provide a session time (around 1 hour).
- ❌ Didn't finish: stopped before completion "If you’d like me to move toward that goal now, I can keep building out the CLI/key derivation"
- Didn't create a .gitignore

### GPT-5.3 Codex

*csharp:*
- ~153,000
- 39%
- $3.57 spent
- Didn't provide a session time (around 50 minutes).
- ❌ Didn't finish: I had to prompt it 3 times. It then succeeded with Scoop

### Grok Build 0.1

*go:*
- 74,238 tokens
- 29% used
- $1.16 spent
- ~15 minutes

*csharp:*
- 196,769 tokens
- 77% used
- $5.14 spent
- 47 minutes 40 seconds

*python*
- ~142,577 tokens
- 56% used
- $1.94
- ~20 minutes

### Kimi 2.7 code

*go:*
- 132,799 tokens
- 51% usage
- $2.18
- ~42 minutes
- Prompted for the PR to be merged, with a link
- Ran two code reviews without any prompting.

*csharp:*
- 27,871 tokens
- 11% usage
- $3.90
- ~1 hour
- Prompted for the PR to be merged, with a link
- One of the builds failed on Github (Darwin build)
- It fixed it upon prompting. Politely asked to create a tag on the repo.

*python:*
- 73,671 tokens
- 28% usage
- $0.91
- ~40 minutes
- Didn't need a PR but created one per the instructions.
- Went slowly for a while.