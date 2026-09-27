// Command sync aggregates otnc's skills_* repositories into this repository's skills/ directory, so every skill installs from a single source with the skills CLI (npx skills add otnc/agent-skills).
//
// Behavior:
//   - Discovers the owner's public, non-fork, non-archived repositories whose name starts with the configured prefix (sources.json).
//   - A repository is synced only if it passes validation: it is public, it contains a SKILL.md (at the root or directly under skills/), and it has at least one published (non-draft, non-prerelease) release.
//   - Each synced release is committed as "sync: <repo> <tag>" and tagged "<repo>@<tag>", so past versions remain installable from this repository.
//   - --replay rebuilds that history by applying every past release, oldest first (initial setup).
//   - Repositories that were synced before but no longer pass validation (deleted, turned private, ...) have their directories removed.
//   - The license of each source repository is recorded in sync-state.json and linked in the generated skill list and source list.
//   - After syncing, the skill list sections of the kiritan base documents (base/README.base.md, base/CONTRIBUTING.base.md) are regenerated from the sync state and the docs are rebuilt with kiritan (npm run build).
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var (
	configPath = "sources.json"
	statePath  = "sync-state.json"
	skillsDir  = "skills"
)

func main() {
	log.SetFlags(0)
	dryRun := flag.Bool("dry-run", false, "plan only: no git operations, no changes to skills/ or sync-state.json")
	replay := flag.Bool("replay", false, "apply every past release oldest-first, one commit and tag each (initial setup)")
	push := flag.Bool("push", false, "push commits and tags to origin after syncing")
	flag.Parse()

	cfg := mustLoadConfig(configPath)
	token := firstNonEmpty(os.Getenv("GITHUB_TOKEN"), os.Getenv("GH_TOKEN"))
	gh := newGitHubClient(cfg.Owner, token)
	state := mustLoadState(statePath)

	valid := discoverAndValidate(gh, cfg)

	// The license of every valid source repository, linked in the generated docs. nil means the repository has no license file (or it could not be read), which renders as no license.
	licenses := make(map[string]*LicenseInfo, len(valid))
	for _, v := range valid {
		licenses[cfg.Owner+"/"+v.repo.Name] = gh.fetchLicense(v.repo.Name)
	}

	// Jobs: (repository, release) pairs to apply.
	type job struct {
		full string // owner/repo
		name string // repo name
		rel  release
	}
	var jobs []job
	if *replay {
		for _, v := range valid {
			for _, rel := range v.releases {
				jobs = append(jobs, job{cfg.Owner + "/" + v.repo.Name, v.repo.Name, rel})
			}
		}
		sort.Slice(jobs, func(i, j int) bool { return jobs[i].rel.PublishedAt.Before(jobs[j].rel.PublishedAt) })
	} else {
		for _, v := range valid {
			full := cfg.Owner + "/" + v.repo.Name
			latest := v.releases[len(v.releases)-1]
			if st := state.Sources[full]; st != nil && st.Tag == latest.TagName {
				continue // already up to date
			}
			jobs = append(jobs, job{full, v.repo.Name, latest})
		}
	}

	// Removals: previously synced repositories that no longer pass validation (deleted, turned private, renamed away, ...).
	validSet := make(map[string]bool, len(valid))
	for _, v := range valid {
		validSet[cfg.Owner+"/"+v.repo.Name] = true
	}
	var removals []string
	for full := range state.Sources {
		if !validSet[full] {
			removals = append(removals, full)
		}
	}
	sort.Strings(removals)

	if len(jobs) == 0 && len(removals) == 0 && len(state.Sources) == 0 {
		logf("everything up to date; nothing to do")
		return
	}

	for _, j := range jobs {
		logf("sync %s %s", j.full, j.rel.TagName)
		prev := state.Sources[j.full]
		dirs, err := applyRelease(cfg.Owner, j.name, j.rel.TagName, prev, *dryRun)
		if err != nil {
			log.Fatalf("apply %s %s: %v", j.full, j.rel.TagName, err)
		}
		if *dryRun {
			continue
		}
		state.Sources[j.full] = &SourceState{
			Tag:      j.rel.TagName,
			Dirs:     dirs,
			License:  licenses[j.full],
			SyncedAt: time.Now().UTC().Format(time.RFC3339),
		}
		saveState(statePath, state)
		if gitCommit(fmt.Sprintf("sync: %s %s", j.name, j.rel.TagName)) {
			tag := j.name + "@" + j.rel.TagName
			if !tagExists(tag) {
				gitOrDie("tag", tag)
			}
		} else {
			logf("  no changes (identical content)")
		}
	}

	for _, full := range removals {
		st := state.Sources[full]
		if st == nil {
			continue
		}
		logf("remove %s (no longer public or valid; removing %d dir(s))", full, len(st.Dirs))
		if *dryRun {
			continue
		}
		for _, d := range st.Dirs {
			if err := os.RemoveAll(filepath.Join(skillsDir, d)); err != nil {
				log.Fatalf("remove %s: %v", d, err)
			}
		}
		delete(state.Sources, full)
		saveState(statePath, state)
		gitCommit(fmt.Sprintf("sync: remove %s (source no longer public or valid)", full))
	}

	if *dryRun {
		logf("dry run: no changes were written")
		return
	}

	// Refresh the recorded license of repositories that were already up to date (freshly synced ones got theirs while their release was applied).
	licensesChanged := false
	for _, v := range valid {
		full := cfg.Owner + "/" + v.repo.Name
		st := state.Sources[full]
		if st == nil || licenseEqual(st.License, licenses[full]) {
			continue
		}
		st.License = licenses[full]
		licensesChanged = true
	}
	if licensesChanged {
		saveState(statePath, state)
	}

	// Reflect the current set of skills in the kiritan base documents, always rebuild the generated docs (README, CONTRIBUTING) with kiritan, and commit whatever changed. This also commits sync-state.json when only the recorded licenses changed.
	if changed, err := updateDocs(state); err != nil {
		log.Fatalf("update docs: %v", err)
	} else if changed {
		logf("docs updated (base documents regenerated, docs rebuilt with kiritan)")
	}

	if *push {
		pushAll()
	}
}
