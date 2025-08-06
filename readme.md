chunky/
├── cmd/
│   └── chunky/              # main CLI entrypoint (main.go lives here)
├── internal/
│   ├── gitpuller/           # clone/pull git repos
│   ├── walker/              # repo traversal & file filtering
│   ├── chunker/             # split files into LLM-sized chunks
│   ├── formatter/           # formats chunks into txt, md, json
│   ├── preamble/            # handles pre/postamble text templates
│   ├── tokenizer/           # estimates tokens for chunks
│   └── flags/               # CLI flag parsing and validation
├── configs/                 # default config files (.chunkyignore example)
├── scripts/                 # helper scripts (maybe CI or build)
├── go.mod
├── go.sum
└── README.md
