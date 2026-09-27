package main

import (
	"sort"
	"strings"
)

type validatedSource struct {
	repo     repository
	releases []release // published releases, oldest first
}

// discoverAndValidate lists the owner's repositories, keeps the ones whose name starts with cfg.Prefix, and validates every candidate. A repository is accepted only when it is public, not a fork, not archived, not excluded, contains a SKILL.md (at the root or directly under skills/), and has at least one published release. Every skip is reported with its reason.
func discoverAndValidate(gh *githubClient, cfg *Config) []validatedSource {
	repos, err := gh.listUserRepos()
	if err != nil {
		logf("list repositories: %v", err)
		return nil
	}
	excluded := makeExclusionSet(cfg.Exclude)

	var valid []validatedSource
	for _, r := range repos {
		full := cfg.Owner + "/" + r.Name
		if !strings.HasPrefix(r.Name, cfg.Prefix) {
			continue
		}
		switch {
		case r.Fork:
			logf("skip %s (fork)", full)
			continue
		case r.Archived:
			logf("skip %s (archived)", full)
			continue
		case r.Private:
			// The /users/{owner}/repos endpoint only lists public repositories, but check anyway: a private repository must never be vendored in.
			logf("skip %s (private)", full)
			continue
		case excluded[full] || excluded[r.Name]:
			logf("skip %s (excluded in sources.json)", full)
			continue
		}
		paths, err := gh.findSkillPaths(r.Name, r.DefaultBranch)
		if err != nil {
			logf("skip %s (could not read tree: %v)", full, err)
			continue
		}
		if len(paths) == 0 {
			logf("skip %s (no SKILL.md at the root or under skills/*/)", full)
			continue
		}
		releases, err := gh.listPublishedReleases(r.Name)
		if err != nil {
			logf("skip %s (could not list releases: %v)", full, err)
			continue
		}
		if len(releases) == 0 {
			logf("skip %s (no published releases)", full)
			continue
		}
		latest := releases[len(releases)-1].TagName
		logf("ok   %s (%d skill(s), %d release(s), latest %s)", full, len(paths), len(releases), latest)
		valid = append(valid, validatedSource{repo: r, releases: releases})
	}
	sort.Slice(valid, func(i, j int) bool { return valid[i].repo.Name < valid[j].repo.Name })
	return valid
}

func makeExclusionSet(exclude []string) map[string]bool {
	set := make(map[string]bool, len(exclude))
	for _, e := range exclude {
		set[strings.TrimSpace(e)] = true
	}
	return set
}
