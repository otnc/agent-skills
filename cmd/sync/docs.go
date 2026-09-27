package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

const (
	readmeBasePath  = "base/README.base.md"
	contribBasePath = "base/CONTRIBUTING.base.md"

	docsCommitMessage = "docs: update the generated skill list"
)

type skillRow struct {
	skill   string       // skill directory name under skills/
	repo    string       // source repository, "owner/repo"
	license *LicenseInfo // license of the source repository
}

// skillRows returns one row per vendored skill directory, sorted by source repository then skill name, so the generated tables are stable.
func skillRows(state *State) []skillRow {
	var rows []skillRow
	for full, st := range state.Sources {
		for _, d := range st.Dirs {
			rows = append(rows, skillRow{skill: d, repo: full, license: st.License})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].repo != rows[j].repo {
			return rows[i].repo < rows[j].repo
		}
		return rows[i].skill < rows[j].skill
	})
	return rows
}

// sourceRow is one row of the source repository list.
type sourceRow struct {
	repo    string       // source repository, "owner/repo"
	license *LicenseInfo // license of the source repository
}

// sourceRows returns the synced source repositories with their licenses, sorted by name.
func sourceRows(state *State) []sourceRow {
	var rows []sourceRow
	for full, st := range state.Sources {
		rows = append(rows, sourceRow{repo: full, license: st.License})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].repo < rows[j].repo })
	return rows
}

func repoLink(full string) string {
	return fmt.Sprintf("[%s](https://github.com/%s)", full, full)
}

// licenseLabel renders a license as a "[SPDX](url)" link. A repository without a license file renders as "-"; one whose file GitHub cannot classify ("NOASSERTION") keeps the link, labeled "Other" like GitHub itself does.
func licenseLabel(l *LicenseInfo) string {
	if l == nil {
		return "-"
	}
	if l.SpdxID == "NOASSERTION" {
		return fmt.Sprintf("[Other](%s)", l.URL)
	}
	return fmt.Sprintf("[%s](%s)", l.SpdxID, l.URL)
}

// licenseEqual reports whether two license records are equal; nil (no license) only equals nil.
func licenseEqual(a, b *LicenseInfo) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// renderTable renders the skill table for a locale ("en" or "ja").
func renderTable(locale string, rows []skillRow) string {
	var b strings.Builder
	if locale == "ja" {
		b.WriteString("| スキル | 元のリポジトリ | ライセンス |\n")
	} else {
		b.WriteString("| Skill | Source repository | License |\n")
	}
	b.WriteString("| --- | --- | --- |\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "| `%s` | %s | %s |\n", r.skill, repoLink(r.repo), licenseLabel(r.license))
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// renderList renders the source repository list with each repository's license (locale-neutral).
func renderList(rows []sourceRow) string {
	var b strings.Builder
	for _, r := range rows {
		b.WriteString("- ")
		b.WriteString(repoLink(r.repo))
		if l := licenseLabel(r.license); l != "-" {
			fmt.Fprintf(&b, " (%s)", l)
		}
		b.WriteString("\n")
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// replaceSection replaces the content between every pair of "<!-- marker:start -->" and "<!-- marker:end -->" with body. It returns the new content and whether anything changed.
func replaceSection(content, marker, body string) (string, bool, error) {
	start := "<!-- " + marker + ":start -->"
	end := "<!-- " + marker + ":end -->"
	var b strings.Builder
	rest := content
	found := false
	for {
		i := strings.Index(rest, start)
		if i < 0 {
			b.WriteString(rest)
			break
		}
		j := strings.Index(rest[i+len(start):], end)
		if j < 0 {
			return content, false, fmt.Errorf("end marker %q missing after %q", end, start)
		}
		j += i + len(start)
		found = true
		b.WriteString(rest[:i+len(start)])
		b.WriteString("\n")
		b.WriteString(body)
		b.WriteString("\n")
		rest = rest[j:]
	}
	if !found {
		return content, false, fmt.Errorf("start marker %q not found", start)
	}
	out := b.String()
	return out, out != content, nil
}

// updateDocs regenerates the marker-delimited sections of the kiritan base documents from the sync state, always rebuilds the generated docs with kiritan (npm run build), and commits whatever changed. Running the build unconditionally also repairs generated docs that drifted (hand edits, kiritan version changes) even when the sections themselves did not change. It returns true when a commit was created.
func updateDocs(state *State) (bool, error) {
	rows := skillRows(state)
	repos := sourceRows(state)

	sections := []struct {
		path   string
		marker string
		body   string
	}{
		{readmeBasePath, "skills-table:en", renderTable("en", rows)},
		{readmeBasePath, "skills-table:ja", renderTable("ja", rows)},
		{contribBasePath, "sources-list", renderList(repos)},
	}

	for _, s := range sections {
		data, err := os.ReadFile(s.path)
		if err != nil {
			return false, err
		}
		newContent, _, err := replaceSection(string(data), s.marker, s.body)
		if err != nil {
			return false, fmt.Errorf("%s: %v", s.path, err)
		}
		if err := os.WriteFile(s.path, []byte(newContent), 0o644); err != nil {
			return false, err
		}
	}

	// Always rebuild the generated docs with kiritan, then commit whatever changed (base documents, generated docs, sync-state.json).
	npm := "npm"
	if runtime.GOOS == "windows" {
		npm = "npm.cmd"
	}
	out, err := exec.Command(npm, "run", "build").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("kiritan build: %s", strings.TrimSpace(string(out)))
	}
	return gitCommit(docsCommitMessage), nil
}
