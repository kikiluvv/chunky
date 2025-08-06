package tokenizer

import (
	"strings"
)

// EstimateTokensApprox estimates token count using word count heuristic.
// avg 1.3 tokens per word is a rough GPT-3/4 estimate.
func EstimateTokensApprox(text string) int {
	words := strings.Fields(text)
	return int(float64(len(words))*1.3 + 0.5) // round to nearest int
}

// EstimateTokensByChars estimates tokens by character count / 4 heuristic.
func EstimateTokensByChars(text string) int {
	return (len(text) + 3) / 4
}

// SplitChunkByTokenLimit splits content into smaller parts by approx token limit,
// trying to split by line boundaries for neatness.
func SplitChunkByTokenLimit(filePath, content string, maxTokens int) []string {
	lines := strings.Split(content, "\n")
	var parts []string
	var currentLines []string
	currentCount := 0

	for _, line := range lines {
		lineTokens := EstimateTokensApprox(line)
		if currentCount+lineTokens > maxTokens && len(currentLines) > 0 {
			parts = append(parts, strings.Join(currentLines, "\n"))
			currentLines = []string{}
			currentCount = 0
		}
		currentLines = append(currentLines, line)
		currentCount += lineTokens
	}
	if len(currentLines) > 0 {
		parts = append(parts, strings.Join(currentLines, "\n"))
	}

	return parts
}
