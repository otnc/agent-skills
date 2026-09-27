package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 5 * time.Minute}

type skillDir struct {
	name string // directory name under skills/ in this repository
	path string // extracted source directory in the temporary dir
}

// applyRelease vendors the skill directories of the given release into skills/ and returns the directory names now owned by the repository. With dryRun set, nothing is written and only the plan is reported.
func applyRelease(owner, repoName, tag string, prev *SourceState, dryRun bool) ([]string, error) {
	tmp, err := os.MkdirTemp("", "agent-skills-sync-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)

	root, err := downloadAndExtractTarball(owner, repoName, tag, tmp)
	if err != nil {
		return nil, err
	}

	newDirs, err := findSkillDirs(root, repoName)
	if err != nil {
		return nil, err
	}
	if len(newDirs) == 0 {
		return nil, fmt.Errorf("release %s of %s contains no skill directories", tag, repoName)
	}

	seen := make(map[string]bool, len(newDirs))
	for _, d := range newDirs {
		if seen[d.name] {
			return nil, fmt.Errorf("duplicate skill name %q in release %s", d.name, tag)
		}
		seen[d.name] = true
	}

	names := make([]string, len(newDirs))
	for i, d := range newDirs {
		names[i] = d.name
	}

	if dryRun {
		for _, d := range newDirs {
			logf("  would vendor %s -> %s/%s", d.name, skillsDir, d.name)
		}
		return names, nil
	}

	// Remove the directories owned by the previous sync of this repository, so renamed and deleted skills disappear too.
	if prev != nil {
		for _, d := range prev.Dirs {
			if err := os.RemoveAll(filepath.Join(skillsDir, d)); err != nil {
				return nil, err
			}
		}
	}
	for _, d := range newDirs {
		dst := filepath.Join(skillsDir, d.name)
		if _, err := os.Lstat(dst); err == nil {
			return nil, fmt.Errorf("%s already exists but is not owned by %s", filepath.ToSlash(dst), repoName)
		}
		if err := copyTree(d.path, dst); err != nil {
			return nil, err
		}
	}
	return names, nil
}

// findSkillDirs locates the skill directories in an extracted release: a root-level SKILL.md (the repository root is the skill, named after the frontmatter "name" field) or every skills/<name>/ directory that contains a SKILL.md.
func findSkillDirs(root, repoName string) ([]skillDir, error) {
	rootSkill := filepath.Join(root, "SKILL.md")
	if fileExists(rootSkill) {
		name := frontmatterName(rootSkill)
		if name == "" {
			name = repoName
		}
		return []skillDir{{name: name, path: root}}, nil
	}
	entries, err := os.ReadDir(filepath.Join(root, "skills"))
	if err != nil {
		return nil, nil // no root SKILL.md and no skills/ directory
	}
	var dirs []skillDir
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		p := filepath.Join(root, "skills", e.Name())
		if fileExists(filepath.Join(p, "SKILL.md")) {
			dirs = append(dirs, skillDir{name: e.Name(), path: p})
		}
	}
	return dirs, nil
}

// downloadAndExtractTarball fetches the release tarball from GitHub and extracts it into destDir. It returns the tarball's single top-level directory (the "<repo>-<tag>" prefix GitHub adds).
func downloadAndExtractTarball(owner, repo, tag, destDir string) (string, error) {
	u := fmt.Sprintf("https://github.com/%s/%s/archive/refs/tags/%s.tar.gz",
		url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(tag))
	resp, err := httpClient.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download %s: status %d", u, resp.StatusCode)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", err
	}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		clean := filepath.Clean(hdr.Name)
		if clean == "." || strings.HasPrefix(clean, "..") {
			continue
		}
		target := filepath.Join(destDir, clean)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return "", err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return "", err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode)&0o777)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return "", err
			}
			f.Close()
		default:
			// Symlinks and other entry types are skipped; skill payloads are plain text files.
		}
	}

	entries, err := os.ReadDir(destDir)
	if err != nil {
		return "", err
	}
	if len(entries) != 1 || !entries[0].IsDir() {
		return "", fmt.Errorf("unexpected tarball layout in %s", u)
	}
	return filepath.Join(destDir, entries[0].Name()), nil
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}

// frontmatterName reads the "name" field from a SKILL.md's YAML frontmatter, or returns "" when there is no frontmatter or no name.
func frontmatterName(skillMdPath string) string {
	data, err := os.ReadFile(skillMdPath)
	if err != nil {
		return ""
	}
	lines := strings.SplitN(string(data), "\n", 200)
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return ""
	}
	for _, l := range lines[1:] {
		if strings.TrimSpace(l) == "---" {
			break
		}
		if strings.HasPrefix(l, "name:") {
			v := strings.TrimSpace(strings.TrimPrefix(l, "name:"))
			v = strings.Trim(v, "\"'")
			return v
		}
	}
	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
