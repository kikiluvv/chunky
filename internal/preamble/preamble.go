package preamble

import (
	"fmt"
	"log"

	"github.com/kikiluvv/chunky/internal/config"
)

// Preamble returns the preamble text, tries to load from .chunkyconfig first.
func Preamble(repoPath string, totalChunks int) string {
	cfg, err := config.LoadConfig(repoPath)
	if err != nil {
		log.Printf("[preamble] error loading .chunkyconfig: %v", err)
	} else if cfg != nil && cfg.Preamble.Text != "" {
		return cfg.Preamble.Text
	}

	// fallback default
	return fmt.Sprintf(
		"### Repo: %s\n### Total Chunks: %d\n\n"+
			"Below are the extracted chunks from the repo for your analysis.\n\n",
		repoPath, totalChunks)
}

// Postamble returns the postamble text, tries to load from .chunkyconfig first.
func Postamble(repoPath string) string {
	cfg, err := config.LoadConfig(repoPath)
	if err != nil {
		log.Printf("[preamble] error loading .chunkyconfig: %v", err)
	} else if cfg != nil && cfg.Postamble.Text != "" {
		return cfg.Postamble.Text
	}

	// fallback default
	return "\n---\nEnd of chunks. Please analyze accordingly.\n"
}
