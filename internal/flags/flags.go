package flags

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	RepoURL      string
	LocalPath    string
	ChunkSize    int
	Format       string
	NoComments   bool
	IncludeGlobs []string
}

func Parse() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.RepoURL, "repo-url", "", "Git repository URL to clone (optional if --local-path provided)")
	flag.StringVar(&cfg.LocalPath, "local-path", "", "Local path to a git repo (optional if --repo-url provided)")
	flag.IntVar(&cfg.ChunkSize, "chunk-size", 3000, "Approximate max tokens per chunk")
	flag.StringVar(&cfg.Format, "format", "txt", "Output format: txt, md, json")
	flag.BoolVar(&cfg.NoComments, "no-comments", false, "Strip comments from code before chunking")
	var includes string
	flag.StringVar(&includes, "include-globs", "", "Comma-separated glob patterns to include (default: all)")

	flag.Parse()

	if cfg.RepoURL == "" && cfg.LocalPath == "" {
		return nil, errors.New("must provide either --repo-url or --local-path")
	}

	if cfg.ChunkSize <= 0 {
		return nil, errors.New("--chunk-size must be positive")
	}

	if includes != "" {
		cfg.IncludeGlobs = strings.Split(includes, ",")
		for i := range cfg.IncludeGlobs {
			cfg.IncludeGlobs[i] = strings.TrimSpace(cfg.IncludeGlobs[i])
		}
	}

	if cfg.Format != "txt" && cfg.Format != "md" && cfg.Format != "json" {
		return nil, fmt.Errorf("invalid --format %q; must be one of: txt, md, json", cfg.Format)
	}

	return cfg, nil
}
