package chunker

import (
	"bytes"
	"regexp"
)

// Chunk represents a chunk of code from a file.
type Chunk struct {
	FilePath string
	Content  string
}

// Simple regex patterns for stripping comments (supports JS, Go, Python style)
// For more languages, this can be expanded or replaced with a proper parser.
var (
	singleLineComment = regexp.MustCompile(`(?m)^\s*//.*$|^\s*#.*$`)
	multiLineComment  = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// ChunkFile splits a file content into chunks of maxChunkSize lines,
// optionally stripping comments.
func ChunkFile(filePath string, content []byte, maxChunkSize int, stripComments bool) ([]Chunk, error) {
	if stripComments {
		content = stripCommentsFromContent(content)
	}

	lines := bytes.Split(content, []byte("\n"))
	var chunks []Chunk

	for start := 0; start < len(lines); start += maxChunkSize {
		end := start + maxChunkSize
		if end > len(lines) {
			end = len(lines)
		}

		chunkLines := lines[start:end]
		chunkContent := bytes.Join(chunkLines, []byte("\n"))

		chunks = append(chunks, Chunk{
			FilePath: filePath,
			Content:  string(chunkContent),
		})
	}

	return chunks, nil
}

// stripCommentsFromContent removes single-line and multi-line comments from content.
func stripCommentsFromContent(content []byte) []byte {
	noMulti := multiLineComment.ReplaceAll(content, []byte(""))
	noSingle := singleLineComment.ReplaceAll(noMulti, []byte(""))
	return noSingle
}
