English | [日本語](README.ja.md)

# fix-unnatural-line-breaks

A Claude Skill that detects and fixes "unnatural line breaks": newlines inserted mid-sentence just to keep lines under some fixed column width.

## What this fixes

It catches things like this:

```
Call me Ishmael. Some years ago—never mind how long
precisely—having little or no money in my purse, and
nothing particular to interest me on shore, I thought
I would sail about a little and see the watery part
of the world.
```

The rule is simple: only break lines at meaningful boundaries, such as paragraph breaks, list item boundaries, or around headings.

```
Call me Ishmael. Some years ago—never mind how long precisely—having little or no money in my purse, and nothing particular to interest me on shore, I thought I would sail about a little and see the watery part of the world.
```

It targets both prose documentation (Markdown, README files) and code comments (docstrings, JSDoc, line/block comments). It handles both English and Japanese.

## Why this matters

Most Markdown renderers treat a single newline inside a paragraph as a space, so the rendered output often looks fine either way. Still, the raw source suffers in a few concrete ways.

- The source itself is harder to read and review; a sentence cut off mid-way looks broken in a diff.
- A one-word edit turns into a full-paragraph reflow, making diffs needlessly large.
- In places where line breaks are shown as-is (code comments, terminals, commit messages), the awkward wrapping is visible directly.

## Layout

```
fix-unnatural-line-breaks/
├── SKILL.md              # Core rules + workflow + quick self-check
├── references/
│   ├── detection.md      # How to spot unnatural breaks (Japanese/English)
│   ├── rules.md          # How to fix them
│   └── examples.md       # Before/after pairs
├── scripts/
│   ├── lint.mjs          # Heuristic detector (plain Node.js, no dependencies)
│   └── lint.test.mjs     # node:test suite for lint.mjs (`node --test scripts/*.test.mjs`)
├── README.md
├── README.ja.md
└── LICENSE
```

`SKILL.md`, `references/`, and `scripts/` are all written in English so the Skill itself is equally usable and reviewable by Japanese and English speakers alike (the `description` field that triggers it also includes both English and Japanese phrasing). Detection rules and examples specific to Japanese text are covered within that English prose, using actual Japanese sample sentences where relevant.

## Install

**Claude Code (personal)**

```bash
git clone https://github.com/otnc/skills_fix-unnatural-line-breaks ~/.claude/skills/fix-unnatural-line-breaks
```

**Claude Code (per-project)**

```bash
git clone https://github.com/otnc/skills_fix-unnatural-line-breaks <project>/.claude/skills/fix-unnatural-line-breaks
```

## Usage

Detection runs mechanically with plain Node.js — no install, no third-party dependency. (Node.js is a safe assumption here: Claude Code itself is distributed as a Node.js CLI, so wherever this Skill would actually run, Node.js is already present.)

```bash
node scripts/lint.mjs path/to/README.md
node scripts/lint.mjs --json path/to/file.ts
```

Without Node.js, Claude reviews the file manually using the criteria in `references/detection.md`.

Detections are only suggestions, and the list is not exhaustive: the script is a small set of regex heuristics, so it flags some lines that are actually fine (code blocks, tables, intentional line breaks) and misses real unnatural breaks that don't match any of its rules. Whether to fix a flagged line is a judgment call based on context, and Claude also reads the whole file by eye — not just the flagged lines — to catch what the script missed.

## Files

**SKILL.md**
Core principles, what to do / what not to do, workflow, and a quick self-check.

**references/detection.md**
How to spot unnatural line breaks: Japanese-specific cues (right after/before a particle, etc.), English-specific cues (mid-prepositional-phrase, etc.), and places where false positives are common.

**references/rules.md**
How to fix what's detected: merging at the paragraph level, deciding upfront which breaks to keep, and how to handle code comments.

**references/examples.md**
Before/after pairs for README prose, bullet lists, and code comments (TypeScript/Python).

**scripts/lint.mjs**
A dependency-free Node.js detector. It flags likely mechanical wraps based on sentence-ending punctuation, trailing particles/conjunctions, and how the next line starts, plus a bare `//`/`#`/`*` line used as a paragraph separator inside a comment block. For non-Markdown files, it only ever compares comment lines (`//`, `#`, `*`, `///`, `;;`) against each other, so it works across `#`-comment languages (Python, Shell, Ruby, YAML, ...) and C-style `//`/`/** */` comments alike. Its own output is always in English, regardless of the language of the file being checked. Its findings are not exhaustive — see `references/detection.md`.

**scripts/lint.test.mjs**
A `node:test` suite covering `lint.mjs`'s core detection and false-positive-avoidance behavior. Run it with `node --test scripts/*.test.mjs`.

## Credits

The rule this Skill encodes started as a personal `CLAUDE.md` convention used by [otoneko.](https://github.com/otnc): never insert unnatural mid-sentence line breaks in documentation or code comments. This Skill packages that rule for general use.

## License

MIT. See `LICENSE`.
