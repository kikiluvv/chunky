# 🧩 chunky

> *cut your codebase into bite-sized pieces for large language models to swallow whole... or choke on.*

**chunky** is a CLI tool written in Go for parsing, filtering, chunking, and formatting source code repositories into LLM-digestible chunks.  

It’s built for developers working on AI tooling, embeddings, summarization, search, or just staring into the abyss of their own spaghetti code.

---

## ✨ Features

- 🌀 **Recursive repo traversal** with `.chunkyignore` and glob filtering
- ✂️ **Line-based chunking** with optional comment stripping
- 📏 **Token estimation** (via heuristics)
- 📦 **Formatter output**: plain text, Markdown, or JSON
- 📋 **Preamble/postamble support** to wrap context
- 💅 Easily extensible for more formats and chunking logic

---

## 🛠️ Usage

```sh
chunky --repo /path/to/repo --out output_dir
```

### 🏴 CLI Flags

| Flag    | Type | Description |
| -------- | -------- | ------- |
| **--repo-url**  | string | Remote Git URL to clone (optional if --local-path is set)    |
| **--local-path** | string | Local path to an existing repo (optional if --repo-url is set)     |
| **--chunk-size** | int | Max tokens per chunk (approximate) — default: 3000     |
| **--format** | string | Output format: txt, md, or json — default: txt     |
| **--no-comments** | bool | Strip single- and multi-line comments before chunking     |
| **--include-globs** | string | Comma-separated glob patterns to include specific files (e.g. *.go,*.md)     |


ℹ️ *You must provide either --repo-url or --local-path, or the void will protest.*

---

## 🗂 Project Structure
```txt
chunky/
├── cmd/
│   └── chunky/              # main CLI entrypoint
├── internal/
│   ├── gitpuller/           # optional repo cloning/pulling
│   ├── walker/              # file discovery & ignore filtering
│   ├── chunker/             # line-based chunk splitting
│   ├── formatter/           # txt / md / json output
│   ├── preamble/            # pre/postamble injection
│   ├── tokenizer/           # rough token estimators
│   └── flags/               # CLI flag parsing
├── configs/                 # default .chunkyignore file
├── scripts/                 # (optional) helper scripts
├── go.mod
└── README.md
```

---

## 🧠 Token Estimation
chunky uses simple heuristics:

- `~1.3 tokens` per word

or 

- `len(text) / 4` based on average token length

This isn’t perfect — but it’s good enough for quick splits.

Use proper tools like OpenAI’s `tiktoken` or `transformers` if you need precision.

---

## 🌫️ Example Output
```txt
📄 repo-name/
├── chunk_001.md
├── chunk_002.md
├── ...
```
Each chunk contains:
- File path metadata
- Source code (formatted based on flags)
- Markdown fences or JSON structure (if configured)

---

## 🔥 Example Use Cases

- Prepping code for GPT-4 or Claude-3 ingestion
- Repo embeddings for vector databases (e.g., Chroma, Pinecone)
- Context windows for agentic AI systems
- Training dataset prep
- Aesthetic pain

---

## 📄 .chunkyignore
just like `.gitignore`, this file defines what gets skipped.
supports globs, directories, files, etc. Ignore things that LLMs 
have no use for, such as `node_modules`.

```gitignore
node_modules/
build/
*.lock
.DS_Store
```

---

## 🧪 Sample Run
```sh
chunky \
  --repo ./my-project \
  --out ./chunks \
  --strip-comments \
  --format md \
  --max-lines 80 \
  --preamble ./preamble.txt \
  --postamble ./postamble.txt
```

---

## 💀 Limitations
- Comment stripping is regex-based (and naive)
- Not token-aware — lines ≠ tokens
- Splits mid-function unless handled manually
- Only supports code-as-text; no AST parsing yet

---

## 💌 Future Ideas
- Tree-sitter based syntax-aware splitting
- Function/class boundary chunking
- Chunk hashing + manifesting
- LLM preview/diffing tools
- Plugin system for alternate strategies

---

## 📜 License
MIT — free to fork, remix, and glitchify.

---

## 🌑 Final Note
chunky doesn’t pretend to understand your code.
it just carves it up, piece by piece, like memory loss.
the LLMs will try to make sense of it.
they may succeed. or hallucinate.

either way, the silence in the chunks is yours to fill.

---

appeal2heaven 2025