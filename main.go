package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)

	repoURL, branch, err := parseArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("%v", err)
	}
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
		if err := runCmd(repoPath, "git", "fetch", "--all", "--prune"); err != nil {
			log.Fatalf("git fetch failed: %v", err)
		}
		log.Printf("checking out branch %s", branch)
		if hasLocalRef(repoPath, "refs/heads/"+branch) {
			if err := runCmd(repoPath, "git", "checkout", branch); err != nil {
				log.Fatalf("git checkout failed: %v", err)
			}
		} else if hasLocalRef(repoPath, "refs/remotes/origin/"+branch) {
			if err := runCmd(repoPath, "git", "checkout", "-b", branch, "--track", "origin/"+branch); err != nil {
				log.Fatalf("git checkout failed: %v", err)
			}
		} else {
			log.Fatalf("branch not found locally or on origin: %s", branch)
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

func parseArgs(args []string) (string, string, error) {
	var repoURL string
	var branch string

	for i := 0; i < len(args); i++ {
		arg := strings.TrimSpace(args[i])
		if arg == "" {
			continue
		}
		if arg == "-b" {
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return "", "", errors.New("usage: goi <repo_url> [-b branch]")
			}
			branch = strings.TrimSpace(args[i+1])
			i++
			continue
		}
		if strings.HasPrefix(arg, "-b=") {
			branch = strings.TrimSpace(strings.TrimPrefix(arg, "-b="))
			if branch == "" {
				return "", "", errors.New("usage: goi <repo_url> [-b branch]")
			}
			continue
		}
		if strings.HasPrefix(arg, "-") {
			return "", "", errors.New("usage: goi <repo_url> [-b branch]")
		}
		if repoURL == "" {
			repoURL = arg
		} else {
			return "", "", errors.New("usage: goi <repo_url> [-b branch]")
		}
	}

	if repoURL == "" {
		return "", "", errors.New("usage: goi <repo_url> [-b branch]")
	}
	return repoURL, branch, nil
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

func hasLocalRef(dir, ref string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", ref)
	cmd.Dir = dir
	return cmd.Run() == nil
}
