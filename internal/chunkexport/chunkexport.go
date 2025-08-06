package chunkexport

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kikiluvv/chunky/internal/chunker"
)

type ManifestEntry struct {
	FileName string `json:"file_name"`
	Source   string `json:"source_file"`
}

func ExportChunks(chunks []chunker.Chunk, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	manifest := []ManifestEntry{}

	for i, chunk := range chunks {
		safeName := sanitizeFileName(chunk.FilePath)
		fileName := fmt.Sprintf("%s_chunk_%d.txt", safeName, i+1)
		fullPath := filepath.Join(outDir, fileName)

		if err := os.WriteFile(fullPath, []byte(chunk.Content), 0644); err != nil {
			return fmt.Errorf("failed writing chunk file %s: %w", fileName, err)
		}

		manifest = append(manifest, ManifestEntry{
			FileName: fileName,
			Source:   chunk.FilePath,
		})
	}

	manifestPath := filepath.Join(outDir, "manifest.json")
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed writing manifest: %w", err)
	}

	return nil
}

// sanitizeFileName replaces slashes and unsafe chars so it’s filename-friendly
func sanitizeFileName(path string) string {
	name := strings.ReplaceAll(path, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	return name
}
