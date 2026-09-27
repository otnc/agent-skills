// Black-box tests for lint.mjs: write a temp file, run the CLI against it, and check what it reports.
// Testing through the CLI (rather than importing internals) keeps this in sync with what `node scripts/lint.mjs <file>` actually does for a user.
import assert from "node:assert/strict";
import { execFileSync } from "node:child_process";
import { mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "node:test";

const SCRIPT = join(dirname(fileURLToPath(import.meta.url)), "lint.mjs");

function runLint(content, filename) {
  const dir = mkdtempSync(join(tmpdir(), "lint-test-"));
  const file = join(dir, filename);
  writeFileSync(file, content, "utf-8");
  try {
    return execFileSync("node", [SCRIPT, "--json", file], {
      encoding: "utf-8",
    });
  } finally {
    rmSync(dir, { recursive: true, force: true });
  }
}

function findingsFor(content, filename) {
  return JSON.parse(runLint(content, filename));
}

test("flags a Japanese sentence wrapped mid-way at a particle", () => {
  const content =
    "石炭をば早や積み果てつ。中等室の卓のほとりはいと静かにて、\n" +
    "熾熱灯の光の晴れがましきも徒なり。\n";
  const findings = findingsFor(content, "sample.md");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].line, 1);
});

test("flags an English sentence wrapped mid-way with no terminal punctuation", () => {
  const content =
    "It is a truth universally acknowledged, that a single\n" +
    "man in possession of a good fortune, must be in want of a wife.\n";
  const findings = findingsFor(content, "sample.md");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].line, 1);
});

test("does not flag a properly joined paragraph", () => {
  const content =
    "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.\n";
  assert.deepEqual(findingsFor(content, "sample.md"), []);
});

test("does not flag list items, headings, or blank-line paragraph boundaries", () => {
  const content =
    "# A heading\n" +
    "\n" +
    "- first item\n" +
    "- second item\n" +
    "\n" +
    "A short paragraph that ends cleanly.\n";
  assert.deepEqual(findingsFor(content, "sample.md"), []);
});

test("does not flag content inside a fenced code block", () => {
  const content = "```\nfoo bar of\nbaz\n```\n";
  assert.deepEqual(findingsFor(content, "sample.md"), []);
});

test("does not flag a Markdown table row", () => {
  const content = "| a | b of |\n| - | - |\n";
  assert.deepEqual(findingsFor(content, "sample.md"), []);
});

test("does not flag YAML key: value lines in a .yml file", () => {
  const content = "name: example\ndescription: a wrapped value continuing\n";
  assert.deepEqual(findingsFor(content, "sample.yml"), []);
});

test("does not flag YAML frontmatter in a Markdown file", () => {
  const content =
    "---\n" + "title: a value that looks wrapped\n" + "---\n" + "Body text.\n";
  assert.deepEqual(findingsFor(content, "sample.md"), []);
});

test("only compares comment lines in non-Markdown source files, never code", () => {
  const content =
    "// This comment sentence is split across two lines and\n" +
    "// continues here.\n" +
    "const of = 1;\n" +
    "const bar = of +\n" +
    "  2;\n";
  const findings = findingsFor(content, "sample.ts");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].line, 1);
});

test("flags a wrapped sentence in a # comment even when the next line is also a # comment", () => {
  // Regression test: HEADING_RE ("#{1,6} ") used to be applied to every file type, so in a #-comment language a "# ..." continuation line matched it and was mistaken for a Markdown heading, silently suppressing the finding.
  const content =
    "# This script installs the dependencies needed to build the\n" +
    "# project, then runs the test suite before exiting.\n";
  const findings = findingsFor(content, "sample.py");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].line, 1);
});

test("flags a wrapped sentence in a JSDoc/block comment even when the next line starts with *", () => {
  // Regression test: LIST_ITEM_RE ("[-*+] ") used to be applied to every file type, so a " * ..." block-comment continuation line matched it and was mistaken for a Markdown bullet, silently suppressing the finding.
  const content =
    "/**\n" +
    " * This function shuffles the given array in place using the\n" +
    " * Fisher-Yates algorithm, and does not mutate the original.\n" +
    " */\n";
  const findings = findingsFor(content, "sample.ts");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].line, 2);
});

test("flags a Japanese comment continuation that only got missed because the next line's comment prefix was never stripped", () => {
  // Regression test (GitHub issue #1, root cause): the continuation checks compared `body` (prefix stripped) against the *raw* trimmed next line, which for a non-Markdown file is almost always still "* ...", "// ...", or "# ...", not the real text. Both continuation rules were effectively dead code for code comments.
  const content =
    "/**\n" +
    " * OS標準のユーザーデータディレクトリ配下に結果保存用フォルダを決定する\n" +
    " * 実装になっている\n" +
    " */\n";
  const findings = findingsFor(content, "sample.ts");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].line, 2);
});

test("flags a line ending in the reading-pause comma even when no tracked particle matches exactly", () => {
  // Regression test (GitHub issue #1): "、" itself is the strongest signal of a mid-sentence break, but there was no dedicated rule for it — a particle followed by "、" (e.g. "実装で、") didn't match the exact-suffix particle check either.
  const content =
    "/**\n" +
    " * npmパッケージの更新では消えない場所に置くための実装で、\n" +
    " * 外部パッケージには依存しない。\n" +
    " */\n";
  const findings = findingsFor(content, "sample.ts");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].reason, 'line ends with the reading-pause comma "、"');
});

test("flags a line ending in one of the broadened particles (だけ)", () => {
  // Regression test (GitHub issue #1): JP_TRAILING_PARTICLES was missing several common continuing particles found in real audited code, including だけ/とも/ため/のみ/ほど/くらい/ぐらい/など/って/たり/ながら/つつ/やら/し.
  const content =
    "/**\n" +
    " * 未指定のフィールドだけ\n" +
    " * ここでの設定にフォールバックする。\n" +
    " */\n";
  const findings = findingsFor(content, "sample.ts");
  assert.equal(findings.length, 1);
  assert.equal(
    findings[0].reason,
    'line ends with the particle/conjunction "だけ"',
  );
});

test("flags a line whose continuation is a JSDoc inline tag like {@link ...}", () => {
  // Regression test (GitHub issue #1): a next line opening with {@link ...} continues the previous line by construction, even when the previous line ends on a dictionary-form verb with no particle and no trailing punctuation.
  const content =
    "/**\n" +
    " * 環境変数から設定を読み取り、CLI全体で使い回す\n" +
    " * {@link ConnpassClient} を組み立てる。\n" +
    " */\n";
  const findings = findingsFor(content, "sample.ts");
  assert.equal(findings.length, 1);
  assert.equal(
    findings[0].reason,
    "next line continues with an inline tag (e.g. {@link ...})",
  );
});

test("does not pair a shebang line with the real comment line after it", () => {
  // Regression test (GitHub issue #1, edge case): "#!/usr/bin/env node" matches COMMENT_PREFIX_RE (bare "#"), so without an exclusion it gets treated as a prose comment line and can be falsely paired with the next real comment.
  const content =
    "#!/usr/bin/env node\n" +
    "// A short, complete comment.\n";
  assert.deepEqual(findingsFor(content, "sample.ts"), []);
});

test("still flags a real wrap on the line right after a shebang", () => {
  const content =
    "#!/usr/bin/env node\n" +
    "// This comment sentence is split across two lines and\n" +
    "// continues here.\n";
  const findings = findingsFor(content, "sample.ts");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].line, 2);
});

test("does not scan an example code block fenced inside a JSDoc comment", () => {
  // Regression test (GitHub issue #1, edge case): CODE_FENCE_RE used to require the fence at the very start of the line, so a fence prefixed with " * " (nested inside /** */) was never recognized, and the example code inside it got scanned as if it were comment prose.
  const content =
    "/**\n" +
    " * Usage example:\n" +
    " * ```ts\n" +
    " * const of = shuffle([1, 2, 3]);\n" +
    " * ```\n" +
    " */\n";
  assert.deepEqual(findingsFor(content, "sample.ts"), []);
});

test("flags a bare // marker used as a paragraph separator between two comment lines", () => {
  // Regression test (GitHub issue #2): a lone "//" (or "* ", "#") line between two real comment lines isn't how comments actually get written by hand — real code uses an unmarked blank line or no separator at all.
  const content =
    "// First paragraph of the comment.\n" +
    "//\n" +
    "// Second paragraph of the comment.\n";
  const findings = findingsFor(content, "sample.js");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].line, 2);
  assert.equal(
    findings[0].reason,
    "bare comment marker used as a paragraph separator — use an unmarked blank line or no separator at all",
  );
});

test("flags a bare * marker used as a paragraph separator inside a JSDoc block", () => {
  const content =
    "/**\n" +
    " * First paragraph.\n" +
    " *\n" +
    " * Second paragraph.\n" +
    " */\n";
  const findings = findingsFor(content, "sample.ts");
  assert.equal(findings.length, 1);
  assert.equal(findings[0].line, 3);
});

test("does not flag an unmarked blank line used as a comment paragraph separator", () => {
  const content =
    "// First paragraph.\n" + "\n" + "// Second paragraph.\n";
  assert.deepEqual(findingsFor(content, "sample.js"), []);
});

test("does not flag two adjacent comment lines with no separator at all", () => {
  const content = "// First paragraph.\n" + "// Second paragraph.\n";
  assert.deepEqual(findingsFor(content, "sample.js"), []);
});

test("does not flag a bare comment marker that isn't sandwiched between two real comment lines", () => {
  const content = "//\n" + "const x = 1;\n";
  assert.deepEqual(findingsFor(content, "sample.js"), []);
});

test("prints a plain-text summary when no issues are found", () => {
  const dir = mkdtempSync(join(tmpdir(), "lint-test-"));
  const file = join(dir, "clean.md");
  writeFileSync(file, "A single clean paragraph.\n", "utf-8");
  try {
    const output = execFileSync("node", [SCRIPT, file], {
      encoding: "utf-8",
    });
    assert.match(output, /No suspicious line breaks found\./);
  } finally {
    rmSync(dir, { recursive: true, force: true });
  }
});
