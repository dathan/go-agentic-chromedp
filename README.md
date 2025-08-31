# Agentic Go Example with Local Code LLMs

This repository contains a simple **Go** program that demonstrates how an agentic workflow can control a web browser.  The example uses the [chromedp](https://github.com/chromedp/chromedp) library to launch Chrome, navigate to [random.org/integers](https://www.random.org/integers/), click the **“Get Numbers”** button and wait for the results.  It is intended as a starting point for building autonomous agents that use local large language models (LLMs) to plan or generate actions.

## Recommended local code models (mid‑2025)

Several open‑source code LLMs can run on Apple silicon hardware without sending your code to the cloud.  Based on recent reviews【326645243015398†L84-L98】, models to consider on a 14‑inch MacBook Pro (M3 Pro, 36 GB RAM) include:

| Model | VRAM requirement (quantised) | Strengths | Notes |
|------|-------------------------------|----------|------|
| **StarCoder2** (7‑15 B) | ~8‑24 GB【326645243015398†L84-L92】 | Strong general‑purpose coding; trained on 600+ programming languages【231297671492799†L149-L176】 including Go; good for scripting and research. | Use `ollama pull starcoder2:7b` for a smaller footprint. |
| **Qwen 2.5 Coder** (7 B or 14 B) | 12‑16 GB【326645243015398†L93-L95】 | Multilingual coding assistant; performs well across 40+ languages【607765842191462†L57-L61】 and excels at fill‑in‑the‑middle editing. | Use `ollama pull qwen2.5-coder:7b` for Go coding tasks. |
| **DeepSeek‑Coder** (6.7 B/33 B) | 12‑16 GB【326645243015398†L88-L90】 | Fast model with advanced parallel token prediction; handles long contexts well. | Larger variants require >24 GB of GPU memory. |
| **Phi‑3 Mini** (3.8 B) | 4‑8 GB【326645243015398†L96-L97】 | Compact model with solid logic‑reasoning abilities; suitable for lightweight coding tasks or running entirely on CPU. | Ideal when memory is constrained or for quick prototyping. |

The MacBook Pro specification shown in the screenshot supports quantised versions of **StarCoder2** (7 B or 15 B) and **Qwen 2.5 Coder** (14 B) comfortably.  Larger models like **Code Llama 70 B** may be possible with aggressive quantisation but are less practical on a laptop【326645243015398†L109-L116】.  Use quantised *GGUF* or *GPTQ* formats to reduce memory footprint【326645243015398†L109-L117】.

## Setting up your environment

Follow these steps to install the necessary tools and run the example program:

1. **Install Go.**  If Go is not already installed, use Homebrew:
   ```bash
   brew install go
   ```

2. **Install Ollama (for local LLMs).**  Ollama provides a CLI and background service for running local models.  On macOS, you can install it via Homebrew【99469359986873†L56-L69】:
   ```bash
   brew install ollama
   ollama --version # verify installation
   ```
   Ollama automatically starts a background service.  List available models with `ollama list`【99469359986873†L70-L71】 and pull one of the recommended coding models, for example:
   ```bash
   ollama pull starcoder2:7b
   ollama pull qwen2.5-coder:7b
   ```

3. **(Optional) Use LM Studio.**  For a graphical user interface, download [LM Studio](https://lmstudio.ai) and drop the downloaded `.dmg` into your Applications folder.  LM Studio lets you chat with local models and manage downloads; however, the CLI via Ollama is sufficient for programmatic access【326645243015398†L119-L124】.

4. **Set up the Go project.**
   ```bash
   # create a new project directory (if cloning this repo, skip mkdir)
   mkdir go-agentic-example && cd go-agentic-example

   # initialise a Go module
   go mod init example.com/randomagent

   # add chromedp dependency
   go get github.com/chromedp/chromedp
   ```
   Copy the `main.go` file from this repository into the project directory.  It contains the automation logic.

5. **Run the agent.**
   ```bash
   go run main.go
   ```

   A Chrome window will open, navigate to **random.org**, click the **Get Numbers** button and then exit.  You should see `Random.org automation completed successfully.` printed to your terminal.  The browser will remain open, handing control back to you after the automation.

6. **Experiment with the local model.**  You can query your local LLM through `curl` or via the Ollama CLI.  For example, to ask **Qwen 2.5 Coder** to generate a Go function:
   ```bash
   curl http://localhost:11434/api/generate \
     -d '{"model": "qwen2.5-coder:7b", "prompt": "Write a Go function to compute factorial.", "stream": false}'
   ```
   This returns JSON with the model’s response.  You can incorporate such responses into your agent to decide what actions to perform or to generate code on the fly.

## Project structure

```
go-agentic-example/
├── main.go    # Go program that drives the browser via chromedp
└── README.md  # this file with setup instructions and model recommendations
```

## Notes

* The program uses chromedp, which requires a Chrome/Chromium installation.  Chrome comes pre‑installed on macOS; if not, download it from [google.com/chrome](https://www.google.com/chrome/).
* When running LLMs locally, monitor your system’s memory and CPU usage.  Quantised model formats reduce VRAM requirements【326645243015398†L109-L117】 but may trade off some accuracy.
* Qwen 2.5 Coder’s rich model sizes (0.5B up to 32B) provide options for different resource budgets【607765842191462†L87-L100】.  Its 7B and 14B variants are ideal for laptop use, and the model performs well across a wide range of programming languages【607765842191462†L57-L61】.

Happy coding!
