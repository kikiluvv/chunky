package walker

import (
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

type Walker struct {
	root         string
	ignorer      *ignore.GitIgnore
	includeGlobs []string
}

func New(root string, ignoreFilePath string, includeGlobs []string) (*Walker, error) {
	ignorer, err := ignore.CompileIgnoreFile(ignoreFilePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return &Walker{
		root:         root,
		ignorer:      ignorer,
		includeGlobs: includeGlobs,
	}, nil
}

func (w *Walker) Walk() ([]string, error) {
	var files []string
	err := filepath.Walk(w.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// get relative path from root
		rel, err := filepath.Rel(w.root, path)
		if err != nil {
			return err
		}

		// normalize to forward slashes for gitignore matching
		relUnix := filepath.ToSlash(rel)

		// ignore files that match the .chunkyignore patterns
		if w.ignorer != nil && w.ignorer.MatchesPath(relUnix) {
			return nil
		}

		if strings.Contains(relUnix, "/.git/") || strings.HasSuffix(relUnix, "/.git") {
			return nil
		}

		// if includeGlobs set, only include matched files
		if len(w.includeGlobs) > 0 {
			match := false
			for _, pattern := range w.includeGlobs {
				// use ToSlash here too
				matched, _ := filepath.Match(pattern, relUnix)
				if matched {
					match = true
					break
				}
			}
			if !match {
				return nil
			}
		}

		files = append(files, path)
		return nil
	})
	return files, err
}
