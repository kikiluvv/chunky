package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kikiluvv/chunky/internal/chunker"
)

// FormatType is the output format enum
type FormatType string

const (
	FormatTxt  FormatType = "txt"
	FormatMd   FormatType = "md"
	FormatJson FormatType = "json"
)

// FormatChunks formats chunks according to the format type
func FormatChunks(chunks []chunker.Chunk, format FormatType) (string, error) {
	switch format {
	case FormatTxt:
		return formatTxt(chunks), nil
	case FormatMd:
		return formatMd(chunks), nil
	case FormatJson:
		return formatJson(chunks)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func formatTxt(chunks []chunker.Chunk) string {
	var sb strings.Builder
	for _, c := range chunks {
		sb.WriteString(fmt.Sprintf("== FILE: %s ==\n", c.FilePath))
		sb.WriteString(c.Content)
		sb.WriteString("\n\n")
	}
	return sb.String()
}

func FormatSingleChunkWithNumber(chunk chunker.Chunk, index int, format FormatType) (string, error) {
	switch format {
	case FormatTxt:
		return fmt.Sprintf("== Chunk %d — FILE: %s ==\n%s\n\n", index, chunk.FilePath, chunk.Content), nil
	case FormatMd:
		ext := fileExtension(chunk.FilePath)
		lang := langFromExt(ext)
		return fmt.Sprintf("## Chunk %d — %s\n\n```%s\n%s\n```\n\n", index, chunk.FilePath, lang, chunk.Content), nil
	case FormatJson:
		// json formatting of single chunk with number embedded? simplest is to marshal normally
		// but you could wrap in a struct if needed
		return formatJson([]chunker.Chunk{chunk})
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func formatMd(chunks []chunker.Chunk) string {
	var sb strings.Builder
	for i, c := range chunks {
		ext := fileExtension(c.FilePath)
		lang := langFromExt(ext)
		sb.WriteString(fmt.Sprintf("## Chunk %d — %s\n\n```%s\n%s\n```\n\n", i, c.FilePath, lang, c.Content))
	}
	return sb.String()
}

func formatJson(chunks []chunker.Chunk) (string, error) {
	data, err := json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// fileExtension extracts the file extension, lowercased without dot
func fileExtension(path string) string {
	dot := strings.LastIndex(path, ".")
	if dot == -1 || dot == len(path)-1 {
		return ""
	}
	return strings.ToLower(path[dot+1:])
}

// langFromExt maps common extensions to markdown code fence languages
func langFromExt(ext string) string {
	switch ext {
	case "js", "jsx":
		return "javascript"
	case "ts", "tsx":
		return "typescript"
	case "go":
		return "go"
	case "py":
		return "python"
	case "java":
		return "java"
	case "c":
		return "c"
	case "cpp", "cc", "cxx", "hpp", "h":
		return "cpp"
	case "json":
		return "json"
	case "html", "htm":
		return "html"
	case "css":
		return "css"
	case "sh", "bash":
		return "bash"
	default:
		return ""
	}
}
