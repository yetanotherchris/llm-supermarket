# llm-supermarket

<img src="./mw28k.jpg" width="25%" height="25%">

This is a set of basic coding tests used to compare LLM coding models and find out if:

1. A model can complete the task.
2. How the model costs to complete the task.

The task are not overly complicated and already have examples on GitHub that models can borrow from (this is removed in the 'hard' mode tests). 
The task is intended to be something models can complete, but at the same time isn't a basic kata. The intention is to see how far a model can autonomously 
complete the task, and how much it costs.

If models and the tooling gradually improve the price should drop, number of tokens used should drop, and completeness should improve.

## Coding tests

There's currently 4 tests:

1. Create a CLI tool that encrypts/decrypts files using the rclone crypt format (which hardcodes its salt). Create a PR, test decrypting two files.
2. Markdown - turn bullets points into a table.
3. Markdown - create links in a markdown file from repositories in a GitHub organisation.
4. Markdown - convert links in a markdown table into inline format.

### Results

All models used medium effort unless stated.  "% used" shows the context it used, typically a percentage of 256k.

#### CLI test

Tests results are from June 2026 using OpenCode and Claude Code. 

| Model                            | Language  | Tokens              | Pass/Fail | Cost  | Time taken           | Notes                                                                                                      |
|----------------------------------|-----------|---------------------|-----------|-------|--------------------- |------------------------------------------------------------------------------------------------------------|
| [**Claude 4.5 Haiku high**][1]   | go        | 71.9k               | ❌        | $6.19  | ~18 mins            | No auto mode;Didn't create the PR: "You're right - I apologize! I committed directly to main instead of..." (cost 10c to correct, linked wrong PR url after that); "5 commits with complete implementation, security hardening, and production-ready code.": 2 tests failed in the PR;merge conflict then it pushed to main randomnly; |
| [**Claude 4.6 Sonnet**][2]       | go        | 146.4k              | ✅        | $3.49  | 34 mins 16 secs     | Failed to decrypt files; asked about rclone custom salt. Passed after nudge to look at example repo. *(Need to re-run using Github.com)* |
| **Claude 4.8 Opus**              | go        | 87.8k               | ✅        | $6.84  | 15 mins 42 secs     | Successfully decrypted both files. Needs public GitHub for Scoop. Offered to merge PR. Wrote clear TODO list. *(Need to re-run Github.com)* |
| [**Claude 4.5 Haiku high**][3]   | csharp    | ~82.3k (95%)        | ❌        | $6.19  | 27 mins 40 secs     | Didn't create a PR, had to be prompted: "You're absolutely right. Let me create a proper PR workflow. I'll reset main, create a feature branch, redo the work, and create a proper PR.";Merge issues with the PR it created;PR fixed, merged, then it discovered it had used the wrong algorithm.;It failed with the test files, but confidently. declared "The CLI works perfectly with its own format, as proven by the successful test." |
| [**Claude 4.6 Sonnet high**][4]  | csharp    | ~189k (33k Haiku)   | ✅        | $6.63  | 38m 40s             |  |
| [**Claude 4.8 Opus high**][5]    | csharp    | 393k (20k Haiku )   | ✅        | $20.05 | 1 hour 56 mins      |  |
| [**DeepSeek V4 Flash**][6]       | go        | 99,914 (10% used)   | ✅        | $0.17  | 55 mins 21 secs     | Confused PR merging with Scoop installation and became stuck. *(Need to re-run)* |
| [**DeepSeek V4 Flash**][7]       | csharp    | 235,746 (24% used)  | ✅        | $0.86  | ~45 mins            |  |
| [**DeepSeek V4 Flash**][8]       | python    | 102,594 (10% used)  | ✅        | $0.05  | ~19 mins            | Had to ask it to verify it had installed via pip before completing. |
| [**Gemini 3.5 Flash**][9]        | go        | 324,349 (31% used)  | ✅        | $6.82  | 40 mins 25 secs     |  |
| [**GLM-5.2**][10]                | go        | 178,847 (18% used)  | ✅        | $4.71  | ~40 mins 31 secs    | Worked out how to merge automatically and that the app name was incorrect. |
| [**GPT-5.1 Codex Mini**][11]     | go        | 76,671 (19% used)   | ✅        | $0.52  | Not provided        | Didn't merge changes; I merged and had to re-prompt. Prompt was missing this. |
| [**GPT-5.1 Codex Mini**][12]     | csharp    | 154,561 (39% used)  | ❌        | $4.10  | ~1 hour             | Didn't finish: stopped before completion. Didn't create a .gitignore. |
| [**GPT-5.3 Codex**][13]          | csharp    | 153,000 (39%)       | ❌        | $3.57  | ~50 mins            | Didn't finish: had to prompt 3 times, then succeeded with Scoop. |
| [**Grok Build 0.1**][14]         | go        | 74,238 (29% used)   | ✅        | $1.16  | ~15 mins            |  |
| [**Grok Build 0.1**][15]         | csharp    | 196,769 (77% used)  | ✅        | $5.14  | 47 mins 40 secs     |  |
| [**Grok Build 0.1**][16]         | python    | 142,577 (56% used)  | ✅        | $1.94  | ~20 mins            |  |
| [**Kimi 2.7 code**][17]          | go        | 132,799 (51% usage) | ✅        | $2.18  | ~42 mins            | Prompted for PR to be merged. Ran two code reviews without prompting. |
| [**Kimi 2.7 code**][18]          | csharp    | 27,871 (11% usage)  | ✅        | $3.90  | ~1 hour             | Prompted for PR to be merged. Darwin build failed on GitHub; fixed upon prompting. Asked to create a tag. |
| [**Kimi 2.7 code**][19]          | python    | 73,671 (28% usage)  | ✅        | $0.91  | ~40 mins            | Didn't need a PR but created one per instructions. Went slowly for a while. |
| [**MiniMax M2.7**][20]           | go        | 118,479 (58% usage) | ❌        | $0.99  | 28 mins             | Failed - didn't infer the default rclone salt (I changed the prompt after this). It did prompt with 3 options for me to merge the PR |
| [**Qwen 3.6 plus**][21]          | go        | 123,123 (50% usage) | ✅        | $0.51  | 30 mins             | It auto merged the PR, didn't prompt for it to be merged. |
| [**Qwen 3.6 plus**][22]          | csharp    | 225,871 (23% usage) | ✅        | $3.78  | 55 mins             | It auto merged the PR, didn't prompt for it to be merged. |


#### 2. Markdown - bullet points to table


#### Deepseek - 39 seconds
32,217
$0.18

#### GLM-52 - 2 mins
20,764
$0.09

#### kimi - 18 mins
127,085 48%
$1.13

#### Claude Sonnet 4.6 - 2 mins 50 seconds
27,838 3%
$0.23

#### GPT 5.4 mini -  15 mins 30 seconds
132, 167 33%
$0.31
Start asking for ask to the temp folder, root repository folder
Didn't know how to use "gh" command line.

#### Gemini 3.5 flash - 2 mins 38 seconds
46,002 4%
$0.70

#### Nemotron 3 Ultra - 2 min2 43 seconds
27,219 3%
Free
Messed the links up

[1]: https://github.com/llm-supermarket/cli-claude45-haiku-go
[2]: https://github.com/llm-supermarket/cli-claude46-sonnet-go
[3]: https://github.com/llm-supermarket/cli-claude45-haiku-csharp
[4]: https://github.com/llm-supermarket/cli-claude45-sonnet-csharp
[5]: https://github.com/llm-supermarket/cli-claude45-opushigh-csharp
[6]: https://github.com/llm-supermarket/cli-deepseek-go
[7]: https://github.com/llm-supermarket/cli-csharp
[8]: https://github.com/llm-supermarket/cli-deepseek-python
[9]: https://github.com/llm-supermarket/cli-gemini-go
[10]: https://github.com/llm-supermarket/cli-glm-go
[11]: https://github.com/llm-supermarket/cli-chatgpt-go
[12]: https://github.com/llm-supermarket/cli-chatgpt-csharp
[13]: https://github.com/llm-supermarket/cli-chatgpt-csharp
[14]: https://github.com/llm-supermarket/cli-grok-go
[15]: https://github.com/llm-supermarket/cli-grok-csharp
[16]: https://github.com/llm-supermarket/cli-grok-python
[17]: https://github.com/llm-supermarket/cli-kimi-go
[18]: https://github.com/llm-supermarket/cli-kimi-csharp
[19]: https://github.com/llm-supermarket/cli-kimi-python
[20]: https://github.com/llm-supermarket/cli-minimaxm27-go
[21]: https://github.com/llm-supermarket/cli-qwen36plus-go
[22]: https://github.com/llm-supermarket/cli-qwen36plus-csharp