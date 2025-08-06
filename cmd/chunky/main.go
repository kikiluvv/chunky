package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kikiluvv/chunky/internal/chunker"
	"github.com/kikiluvv/chunky/internal/flags"
	"github.com/kikiluvv/chunky/internal/formatter"
	"github.com/kikiluvv/chunky/internal/gitpuller"
	"github.com/kikiluvv/chunky/internal/preamble"
	"github.com/kikiluvv/chunky/internal/tokenizer"
	"github.com/kikiluvv/chunky/internal/walker"
)

func main() {
	cfg, err := flags.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	var repoPath string
	if cfg.RepoURL != "" {
		repoPath, err = gitpuller.CloneOrUpdate(cfg.RepoURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error cloning/pulling repo: %v\n", err)
			os.Exit(1)
		}
	} else if cfg.LocalPath != "" {
		repoPath = cfg.LocalPath
	} else {
		fmt.Fprintf(os.Stderr, "must provide either --repo-url or --local-path\n")
		os.Exit(1)
	}

	// FIX: Make repoPath absolute so .chunkyconfig loads properly
	repoPath, err = filepath.Abs(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving absolute repo path: %v\n", err)
		os.Exit(1)
	}

	repoName := filepath.Base(repoPath)
	outputDir := filepath.Join("chunk_output", repoName)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output directory: %v\n", err)
		os.Exit(1)
	}

	ignoreFile := filepath.Join(repoPath, ".chunkyignore")
	w, err := walker.New(repoPath, ignoreFile, cfg.IncludeGlobs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating walker: %v\n", err)
		os.Exit(1)
	}

	files, err := w.Walk()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error walking repo: %v\n", err)
		os.Exit(1)
	}

	var allChunks []chunker.Chunk

	for _, f := range files {
		fmt.Println("chunking file:", f)

		content, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", f, err)
			continue
		}

		if isBinary(content) {
			fmt.Fprintf(os.Stderr, "skipping binary file %s\n", f)
			continue
		}

		lineChunks, err := chunker.ChunkFile(f, content, cfg.ChunkSize, cfg.NoComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error chunking file %s: %v\n", f, err)
			continue
		}

		maxTokens := cfg.ChunkSize
		for _, chunk := range lineChunks {
			parts := tokenizer.SplitChunkByTokenLimit(chunk.FilePath, chunk.Content, maxTokens)
			for _, part := range parts {
				allChunks = append(allChunks, chunker.Chunk{
					FilePath: chunk.FilePath,
					Content:  part,
				})
			}
		}
	}

	// Load preamble/postamble from preamble package (which tries .chunkyconfig)
	pre := preamble.Preamble(repoPath, len(allChunks))
	post := preamble.Postamble(repoPath)

	// 🍰 Insert preamble/postamble into chunk 0 and last chunk
	if len(allChunks) > 0 {
		allChunks[0].Content = fmt.Sprintf(
			"-- PREAMBLE START --\n%s\n-- PREAMBLE END --\n\n%s",
			strings.TrimSpace(pre),
			allChunks[0].Content,
		)

		lastIdx := len(allChunks) - 1
		allChunks[lastIdx].Content = fmt.Sprintf(
			"%s\n\n-- POSTAMBLE START --\n%s\n-- POSTAMBLE END --",
			strings.TrimSpace(allChunks[lastIdx].Content),
			strings.TrimSpace(post),
		)
	}

	// 💾 write each chunk to its own file
	for i, c := range allChunks {
		formatted, err := formatter.FormatChunks([]chunker.Chunk{c}, formatter.FormatType(cfg.Format))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error formatting chunk: %v\n", err)
			continue
		}

		fileName := fmt.Sprintf("%03d_chunk.%s", i, strings.ToLower(cfg.Format))
		fullPath := filepath.Join(outputDir, fileName)
		if err := ioutil.WriteFile(fullPath, []byte(formatted), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing chunk file: %v\n", err)
		}
	}

	// 🖨️ Print full concatenated output to stdout
	formatted, err := formatter.FormatChunks(allChunks, formatter.FormatType(cfg.Format))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error formatting all chunks: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(formatted)
}

// isBinary does a quick check for null bytes to skip binary files
func isBinary(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}
