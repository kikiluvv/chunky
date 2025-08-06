package gitpuller

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CloneOrUpdate(repoURL string) (string, error) {
	// create a cache dir inside user's home (e.g. ~/.chunky/repos)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(homeDir, ".chunky", "repos")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	// create a safe dir name from repo url (simple slug)
	dirName := slugifyRepoURL(repoURL)
	repoPath := filepath.Join(cacheDir, dirName)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		// clone if not exists
		fmt.Println("Cloning repo:", repoURL)
		cmd := exec.Command("git", "clone", repoURL, repoPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("git clone failed: %w", err)
		}
	} else {
		// else pull latest
		fmt.Println("Pulling latest changes in:", repoPath)
		cmd := exec.Command("git", "-C", repoPath, "pull")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("git pull failed: %w", err)
		}
	}

	return repoPath, nil
}

func slugifyRepoURL(repoURL string) string {
	// remove protocol
	s := strings.TrimPrefix(repoURL, "https://")
	s = strings.TrimPrefix(s, "git@")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.TrimSuffix(s, ".git")
	return s
}
