package chunker

import (
	"bytes"
	"regexp"
)

type Chunk struct {
	FilePath string
	Content  string
}

var (
	singleLineComment = regexp.MustCompile(`(?m)^\s*(//|#).*?$`)
	inlineComment     = regexp.MustCompile(`(?m)([^:"'])\s+(//|#).*?$`)
	multiLineComment  = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

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

func stripCommentsFromContent(content []byte) []byte {
	// 🩸 Remove multiline /* ... */
	content = multiLineComment.ReplaceAll(content, []byte(""))

	// 🌫 Remove single-line // and #
	content = singleLineComment.ReplaceAll(content, []byte(""))

	// 🌒 Remove inline comments (but keep code before them)
	content = inlineComment.ReplaceAll(content, []byte("$1"))

	return content
}
