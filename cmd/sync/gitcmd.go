package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func gitOrDie(args ...string) {
	out, err := git(args...)
	if err != nil {
		log.Fatalf("git %s: %s", strings.Join(args, " "), out)
	}
}

// gitIdentity returns the author identity for sync commits: SYNC_GIT_NAME/SYNC_GIT_EMAIL when set (the sync workflow sets them), otherwise the local git config. It fails when neither is available.
func gitIdentity() (name, email string) {
	name = firstNonEmpty(os.Getenv("SYNC_GIT_NAME"), gitConfig("user.name"))
	email = firstNonEmpty(os.Getenv("SYNC_GIT_EMAIL"), gitConfig("user.email"))
	if name == "" || email == "" {
		log.Fatalf("commit identity is not configured: set SYNC_GIT_NAME/SYNC_GIT_EMAIL (the sync workflow does) or git user.name/user.email")
	}
	return
}

// gitConfig returns the given git config value, or "" when it is unset.
func gitConfig(key string) string {
	out, err := git("config", "--get", key)
	if err != nil {
		return ""
	}
	return out
}

// gitCommit stages everything and creates a commit with the given message. It returns true when a commit was created, false when there was nothing to commit.
func gitCommit(message string) bool {
	gitOrDie("add", "-A")
	if _, err := git("diff", "--cached", "--quiet"); err == nil {
		return false // nothing staged
	}
	name, email := gitIdentity()
	gitOrDie("-c", "user.name="+name, "-c", "user.email="+email, "commit", "-m", message)
	return true
}

func tagExists(tag string) bool {
	_, err := git("rev-parse", "-q", "--verify", "refs/tags/"+tag)
	return err == nil
}

func pushAll() {
	if _, err := git("remote", "get-url", "origin"); err != nil {
		logf("no origin remote; skipping push")
		return
	}
	gitOrDie("push", "origin", "HEAD")
	gitOrDie("push", "origin", "--tags")
}
