package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type githubClient struct {
	owner  string
	token  string
	client *http.Client
}

func newGitHubClient(owner, token string) *githubClient {
	return &githubClient{
		owner:  owner,
		token:  token,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

type repository struct {
	Name          string `json:"name"`
	Private       bool   `json:"private"`
	Fork          bool   `json:"fork"`
	Archived      bool   `json:"archived"`
	DefaultBranch string `json:"default_branch"`
}

type release struct {
	TagName     string    `json:"tag_name"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
}

type treeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// getJSON performs a GET request and returns the status code and body. Transport errors and rate-limit exhaustion are returned as errors; HTTP status codes (including 404) are not, so callers can branch on them.
func (g *githubClient) getJSON(rawURL string) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	if resp.StatusCode == http.StatusForbidden && resp.Header.Get("X-RateLimit-Remaining") == "0" {
		return resp.StatusCode, body, fmt.Errorf("GitHub API rate limit exceeded; set GITHUB_TOKEN or GH_TOKEN")
	}
	return resp.StatusCode, body, nil
}

func (g *githubClient) apiURL(format string, args ...any) string {
	return "https://api.github.com" + fmt.Sprintf(format, args...)
}

// listUserRepos returns the owner's public repositories (the /users endpoint never returns private ones; the Private field is still checked downstream as defense in depth).
func (g *githubClient) listUserRepos() ([]repository, error) {
	var all []repository
	for page := 1; page <= 20; page++ {
		u := g.apiURL("/users/%s/repos?per_page=100&page=%d", url.PathEscape(g.owner), page)
		status, body, err := g.getJSON(u)
		if err != nil {
			return nil, err
		}
		if status != http.StatusOK {
			return nil, fmt.Errorf("list repositories of %s: status %d", g.owner, status)
		}
		var repos []repository
		if err := json.Unmarshal(body, &repos); err != nil {
			return nil, err
		}
		all = append(all, repos...)
		if len(repos) < 100 {
			return all, nil
		}
	}
	return all, nil
}

// findSkillPaths returns the SKILL.md paths in the repository's default branch: a root-level "SKILL.md" or "skills/<name>/SKILL.md". It returns nil when the repository has none (or has no branch / is empty).
func (g *githubClient) findSkillPaths(repoName, branch string) ([]string, error) {
	if branch == "" {
		branch = "main"
	}
	u := g.apiURL("/repos/%s/%s/git/trees/%s?recursive=1",
		url.PathEscape(g.owner), url.PathEscape(repoName), url.PathEscape(branch))
	status, body, err := g.getJSON(u)
	if err != nil {
		return nil, err
	}
	switch status {
	case http.StatusOK:
	case http.StatusNotFound, http.StatusConflict: // no such branch or empty repository
		return nil, nil
	default:
		return nil, fmt.Errorf("read tree of %s/%s: status %d", g.owner, repoName, status)
	}
	var tr struct {
		Tree []treeEntry `json:"tree"`
	}
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}
	var paths []string
	for _, e := range tr.Tree {
		if e.Type != "blob" {
			continue
		}
		if e.Path == "SKILL.md" || isNestedSkillPath(e.Path) {
			paths = append(paths, e.Path)
		}
	}
	return paths, nil
}

func isNestedSkillPath(p string) bool {
	parts := strings.Split(p, "/")
	return len(parts) == 3 && parts[0] == "skills" && parts[2] == "SKILL.md"
}

// listPublishedReleases returns the repository's published (non-draft, non-prerelease) releases, oldest first.
func (g *githubClient) listPublishedReleases(repoName string) ([]release, error) {
	var all []release
	for page := 1; page <= 10; page++ {
		u := g.apiURL("/repos/%s/%s/releases?per_page=100&page=%d",
			url.PathEscape(g.owner), url.PathEscape(repoName), page)
		status, body, err := g.getJSON(u)
		if err != nil {
			return nil, err
		}
		if status != http.StatusOK {
			return nil, fmt.Errorf("list releases of %s/%s: status %d", g.owner, repoName, status)
		}
		var rels []release
		if err := json.Unmarshal(body, &rels); err != nil {
			return nil, err
		}
		for _, r := range rels {
			if !r.Draft && !r.Prerelease {
				all = append(all, r)
			}
		}
		if len(rels) < 100 {
			break
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].PublishedAt.Before(all[j].PublishedAt) })
	return all, nil
}

// fetchLicense returns the repository's license: its SPDX identifier and the link to the license file. It returns nil when the repository has no license file, and also when the lookup fails (after logging why), so a transient API error never fails the sync.
func (g *githubClient) fetchLicense(repoName string) *LicenseInfo {
	u := g.apiURL("/repos/%s/%s/license", url.PathEscape(g.owner), url.PathEscape(repoName))
	status, body, err := g.getJSON(u)
	if err != nil {
		logf("license of %s/%s: %v (leaving it unlinked)", g.owner, repoName, err)
		return nil
	}
	switch status {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil // no license file
	default:
		logf("license of %s/%s: status %d (leaving it unlinked)", g.owner, repoName, status)
		return nil
	}
	var lr struct {
		HTMLURL string `json:"html_url"`
		License *struct {
			SpdxID string `json:"spdx_id"`
		} `json:"license"`
	}
	if err := json.Unmarshal(body, &lr); err != nil {
		logf("license of %s/%s: %v (leaving it unlinked)", g.owner, repoName, err)
		return nil
	}
	if lr.License == nil || lr.License.SpdxID == "" || lr.HTMLURL == "" {
		return nil
	}
	return &LicenseInfo{SpdxID: lr.License.SpdxID, URL: lr.HTMLURL}
}
