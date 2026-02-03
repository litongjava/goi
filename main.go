package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)

	var branch string
	flag.StringVar(&branch, "b", "", "branch name")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("usage: mvni <repo_url> [-b branch]")
	}

	repoURL := flag.Arg(0)
	repoName, err := repoDirName(repoURL)
	if err != nil {
		log.Fatalf("invalid repo url: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cannot resolve home dir: %v", err)
	}
	codeDir := filepath.Join(home, "code")
	if err := os.MkdirAll(codeDir, 0o755); err != nil {
		log.Fatalf("cannot create code dir: %v", err)
	}

	repoPath := filepath.Join(codeDir, repoName)

	if _, err := os.Stat(repoPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("cloning %s into %s", repoURL, repoPath)
		if err := runCmd(codeDir, "git", "clone", repoURL); err != nil {
			log.Fatalf("git clone failed: %v", err)
		}
	} else if err != nil {
		log.Fatalf("cannot stat repo path: %v", err)
	} else {
		log.Printf("repo exists at %s", repoPath)
		if err := runCmd(repoPath, "git", "reset", "--hard"); err != nil {
			log.Fatalf("git reset failed: %v", err)
		}
	}

	if branch != "" {
		log.Printf("checking out branch %s", branch)
		if err := runCmd(repoPath, "git", "checkout", branch); err != nil {
			log.Fatalf("git checkout failed: %v", err)
		}
	}

	log.Printf("pulling latest changes")
	if err := runCmd(repoPath, "git", "pull"); err != nil {
		log.Fatalf("git pull failed: %v", err)
	}

	log.Printf("running build")
	if err := runCmd(repoPath, "build"); err != nil {
		log.Fatalf("build failed: %v", err)
	}
}

func repoDirName(repoURL string) (string, error) {
	repoURL = strings.TrimSpace(repoURL)
	if repoURL == "" {
		return "", errors.New("empty url")
	}
	repoURL = strings.TrimSuffix(repoURL, "/")
	base := filepath.Base(repoURL)
	base = strings.TrimSuffix(base, ".git")
	if base == "." || base == "/" || base == "" {
		return "", errors.New("cannot parse repo name")
	}
	return base, nil
}

func runCmd(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
