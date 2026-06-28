# llm-supermarket

<img src="./mw28k.jpg" width="25%" height="25%">

This is a set of basic coding tests used to compare LLM coding models and find out if:

1. A model can complete the task.
2. How the model costs to complete the task.

If models and the tooling gradually improve the price should drop, number of tokens used should drop, and completeness should improve.

## Coding tests

There's currently 4 tests:

1. Create a CLI tool that encrypts/decrypts files using the rclone crypt format (which hardcodes its salt). Create a PR, test decrypting two files.
2. Markdown - turn bullets points into a Markdown table.
3. Markdown - create links in a Markdown table, using repositories in a GitHub organisation.
4. Markdown - convert links in a Markdown table into inline format.

The idea is not to have SWE style tests, or "create a Rust compiler", but a few boring everyday coding tests to compare costs/tokens. The CLI test 
was designed so that there's a relatively hard to find answer for how to decrypt the example files, which the model has to discover itself (e.g 
the rclone default salt value).

### Results

All models used medium effort unless stated.  "% used" shows the context it used, typically a percentage of 256k.
Tests results are from June 2026 using OpenCode and Claude Code. 

#### CLI test

| Model                            | Language  | Tokens              | Pass/Fail | Cost    | Time taken           | Notes                                                                                                      |
|----------------------------------|-----------|---------------------|-----------|---------|--------------------- |------------------------------------------------------------------------------------------------------------|
| [**Claude 4.5 Haiku high**][1]   | go        | 71.9k               | ❌        | $6.19  | ~18 mins             | No auto mode;Didn't create the PR: "You're right - I apologize! I committed directly to main instead of..." (cost 10c to correct, linked wrong PR url after that); "5 commits with complete implementation, security hardening, and production-ready code.": 2 tests failed in the PR;merge conflict then it pushed to main randomnly; |
| [**Claude 4.6 Sonnet**][2]       | go        | 146.4k              | ✅        | $3.49  | 34 mins 16 secs      | Failed to decrypt files; asked about rclone custom salt. Passed after nudge to look at example repo.  |
| [**Claude 4.8 Opus**][3]         | go        | 107.1k              | ✅        | $6.20  | 17 mins 43 secs      | Followed all instructions exactly. Produced a neat table of decrypted filename and its contents when finishing. |
| [**Claude 4.5 Haiku high**][4]   | csharp    | 82.3k (95%)         | ❌        | $6.19  | 27 mins 40 secs      | Didn't create a PR, had to be prompted: "You're absolutely right. Let me create a proper PR workflow. I'll reset main, create a feature branch, redo the work, and create a proper PR.";Merge issues with the PR it created;PR fixed, merged, then it discovered it had used the wrong algorithm.;It failed with the test files, but confidently. declared "The CLI works perfectly with its own format, as proven by the successful test." |
| [**Claude 4.6 Sonnet high**][5]  | csharp    | 189k (33k Haiku)    | ✅        | $6.63  | 38m 40s              |  |
| [**Claude 4.6 Sonnet**][6]       | csharp    | 126.5k (12k Haiku)  | ✅        | $5.97  | 32m 38s              |  |
| [**Claude 4.8 Opus**][7]         | csharp    | 141.1k (15.8k Haiku)| ✅        | $10.00 | 26 min 42 secs       |  |
| [**Claude 4.8 Opus high**][8]    | csharp    | 393k (20k Haiku )   | ✅        | $20.05 | 1 hour 56 mins       |  |
| [**DeepSeek V4 Flash**][9]       | go        | 118,806 (12% used)  | ✅        | $0.08  | 17 mins 3 secs       | Created the PR, it auto-merged the PR, installed via Scoop and unencrypted successfully. |
| [**DeepSeek V4 Flash**][10]      | csharp    | 235,746 (24% used)  | ✅        | $0.86  | ~45 mins             |  |
| [**DeepSeek V4 Flash**][11]      | python    | 102,594 (10% used)  | ✅        | $0.05  | ~19 mins             | Had to ask it to verify it had installed via pip before completing. |
| [**Gemini 3.5 Flash**][12]       | go        | 324,349 (31% used)  | ✅        | $6.82  | 40 mins 25 secs      |  |
| [**GLM-5.2**][13]                | go        | 178,847 (18% used)  | ✅        | $4.71  | ~40 mins 31 secs     | Worked out how to merge automatically and that the app name was incorrect. |
| [**GPT-5.1 Codex Mini**][14]     | go        | 76,671 (19% used)   | ✅        | $0.52  | Not provided         | Didn't merge changes; I merged and had to re-prompt. Prompt was missing this. |
| [**GPT-5.1 Codex Mini**][15]     | csharp    | 154,561 (39% used)  | ❌        | $4.10  | ~1 hour              | Didn't finish: stopped before completion. Didn't create a .gitignore. Scanned non-repo directories.|
| [**GPT-5.3 Codex**][16]          | csharp    | 153,000 (39%)       | ❌        | $3.57  | ~50 mins             | Didn't finish: had to prompt 3 times, then succeeded with Scoop. Scanned non-repo directories.|
| [**Grok Build 0.1**][17]         | go        | 74,238 (29% used)   | ✅        | $1.16  | ~15 mins             |  |
| [**Grok Build 0.1**][18]         | csharp    | 196,769 (77% used)  | ✅        | $5.14  | 47 mins 40 secs      |  |
| [**Grok Build 0.1**][19]         | python    | 142,577 (56% used)  | ✅        | $1.94  | ~20 mins             |  |
| [**Kimi 2.7 code**][20]          | go        | 132,799 (51% usage) | ✅        | $2.18  | ~42 mins             | Prompted for PR to be merged. Ran two code reviews without prompting. |
| [**Kimi 2.7 code**][21]          | csharp    | 27,871 (11% usage)  | ✅        | $3.90  | ~1 hour              | Prompted for PR to be merged. Darwin build failed on GitHub; fixed upon prompting. Asked to create a tag. |
| [**Kimi 2.7 code**][22]          | python    | 73,671 (28% usage)  | ✅        | $0.91  | ~40 mins             | Didn't need a PR but created one per instructions. Went slowly for a while. |
| [**MiniMax M2.7**][23]           | go        | 118,479 (58% usage) | ❌        | $0.99  | 28 mins              | Failed - didn't infer the default rclone salt (I changed the prompt after this). It did prompt with 3 options for me to merge the PR |
| [**Qwen 3.6 plus**][24]          | go        | 123,123 (50% usage) | ✅        | $0.51  | 30 mins              | It auto merged the PR, didn't prompt for it to be merged. |
| [**Qwen 3.6 plus**][25]          | csharp    | 225,871 (23% usage) | ✅        | $3.78  | 55 mins              | It auto merged the PR, didn't prompt for it to be merged. |


#### 2. Markdown - bullet points to table

| Model                | Time taken         | Tokens used    | Cost   | Pass/Fail | Notes               |
|----------------------|--------------------|----------------|--------|-----------|---------------------|
| Deepseek v4 Flash    | 39 seconds         | 32,217         | $0.18  | ✅        | |
| GLM-52               | 2 mins             | 20,764         | $0.09  | ✅        | |
| Kimi                 | 18 mins            | 127,085 48%    | $1.13  | ✅        | |
| Claude Sonnet 4.6    | 2 mins 50 seconds  | 27,838 3%      | $0.23  | ✅        | |
| GPT 5.4 mini         | 15 mins 30 seconds | 132, 167 33%   | $0.31  | ❌        |Asked for ask permissions for the temp folder, then the root repository folder; Didn't know how to use "gh" command line. |
| Gemini 3.5 Flash     | 2 mins 38 seconds  | 46,002 4%      | $0.70  | ✅        | |
| Nemotron 3 Ultra     | 2 min2 43 seconds  | 27,219 3%      | Free   | ❌        |Messed the links up |

[1]: https://github.com/llm-supermarket/cli-claude45-haiku-go
[2]: https://github.com/llm-supermarket/cli-claude46-sonnet-go
[3]: https://github.com/llm-supermarket/cli-claude48-opus-go
[4]: https://github.com/llm-supermarket/cli-claude45-haiku-high-csharp
[5]: https://github.com/llm-supermarket/cli-claude45-sonnet-high-csharp
[6]: https://github.com/llm-supermarket/cli-claude46-sonnet-medium-csharp/
[7]: https://github.com/llm-supermarket/cli-claude48-opusmedium-csharp
[8]: https://github.com/llm-supermarket/cli-claude45-opushigh-csharp
[9]: https://github.com/llm-supermarket/cli-deepseekv4-flash-go
[10]: https://github.com/llm-supermarket/cli-deepseekv4-flash-csharp
[11]: https://github.com/llm-supermarket/cli-deepseekv4-flash-python
[12]: https://github.com/llm-supermarket/cli-gemini-go
[13]: https://github.com/llm-supermarket/cli-glm-go
[14]: https://github.com/llm-supermarket/cli-chatgpt-go
[15]: https://github.com/llm-supermarket/cli-chatgpt-csharp
[16]: https://github.com/llm-supermarket/cli-chatgpt-csharp
[17]: https://github.com/llm-supermarket/cli-grok-go
[18]: https://github.com/llm-supermarket/cli-grok-csharp
[19]: https://github.com/llm-supermarket/cli-grok-python
[20]: https://github.com/llm-supermarket/cli-kimi-go
[21]: https://github.com/llm-supermarket/cli-kimi-csharp
[22]: https://github.com/llm-supermarket/cli-kimi-python
[23]: https://github.com/llm-supermarket/cli-minimaxm27-go
[24]: https://github.com/llm-supermarket/cli-qwen36plus-go
[25]: https://github.com/llm-supermarket/cli-qwen36plus-csharp