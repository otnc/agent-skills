import { Buffer, isUtf8 } from "node:buffer";
import { readFile } from "node:fs/promises";
import { stdin } from "node:process";
import { parseArgs } from "node:util";

const usage = `Usage: node lint.js [--terminology] [--kangxi-radicals] [--encoding] <file>...
Use - to read from standard input. With no check options, all checks run.`;
const terminologyFile = new URL("../references/terminology.md", import.meta.url);

const terminologyContexts = [
  { terms: ["版", "更新履歴", "導入する", "実行環境", "構築する"], nearby: /Node\.js|Deno|Bun|npm|CLI|SDK|API|ソフトウェア|パッケージ|ビルド|リリース|バイナリ|コンテナ/iu },
  { terms: ["分岐", "変更を記録する", "変更を統合する"], nearby: /\bGit\b|git (?:branch|switch|merge|commit)|リポジトリ|ブランチ|コミット/iu },
  { terms: ["作業手順", "作業", "手順", "起動条件", "成果物"], nearby: /GitHub Actions|\bCI\b|ワークフロー|ジョブ|ランナー|upload-artifact/iu },
  { terms: ["要求", "応答", "接続先", "入口", "再試行する", "時間切れ", "一時保存", "接続期間"], nearby: /HTTP|\bAPI\b|\bSDK\b|fetch\(|サーバー|クライアント|リクエスト|レスポンス|通信|ログイン|ブラウザ/iu },
  { terms: ["配置する", "展開する", "元に戻す", "複製", "注釈"], nearby: /Kubernetes|\bkubectl\b|デプロイ|ロールアウト|本番環境|コンテナ|\bPod\b/iu },
  { terms: ["容器", "画像", "登録簿"], nearby: /Docker|Dockerfile|コンテナイメージ|\bOCI\b/iu },
  { terms: ["文脈", "前後の事情", "指示文", "文字数", "単語数", "処理担当"], nearby: /\bLLM\b|言語モデル|AI エージェント|ChatGPT|Claude|プロンプト|トークン/iu },
];

function terminologySuggestions(markdown) {
  const suggestions = new Map();

  // The reference table pairs alternatives in order or shares one across the row.
  for (const row of markdown.matchAll(/^\| (.+?) \| (.+?) \| .+? \|\r?$/gm)) {
    const [, candidates, alternatives] = row;
    if (candidates === "見直し候補" || candidates === "---") continue;

    const terms = candidates.split("、");
    const replacements = alternatives.split("、");
    if (replacements.length !== 1 && terms.length !== replacements.length) {
      throw new Error(`Cannot pair terminology candidates with alternatives: ${candidates}`);
    }

    for (const [index, term] of terms.entries()) {
      if (suggestions.has(term)) throw new Error(`Duplicate terminology candidate: ${term}`);
      suggestions.set(term, replacements.length === 1 ? replacements[0] : replacements[index]);
    }
  }

  if (suggestions.size === 0) throw new Error("No terminology candidates found.");
  return suggestions;
}

function diagnosticPattern(suggestions, { checkKangxiRadicals, checkEncoding }) {
  // Match longer phrases first and keep 版 inside compounds and edition numbers.
  const terms = [...suggestions.keys()]
    .sort((left, right) => right.length - left.length)
    .map((term) => term.replace(/[.*+?^${}()|[\]\\]/g, "\\$&"))
    .map((term) => term === "版" ? "(?<![\\p{Script=Han}0-9０-９])版(?!\\p{Script=Han})" : term);

  if (checkKangxiRadicals) terms.unshift("[\\u2F00-\\u2FD5]");
  if (checkEncoding) terms.unshift("\\uFFFD");
  return new RegExp(terms.join("|"), "gu");
}

function fitsTerminologyContext(line, match) {
  const rule = terminologyContexts.find(({ terms }) => terms.includes(match[0]));
  if (!rule) return true;

  // Keep a technical cue in a neighboring sentence from triggering a warning.
  const start = match.index;
  const end = start + match[0].length;
  const before = line.slice(Math.max(0, start - 32), start).split(/[。！？!?]/u).at(-1);
  const after = line.slice(end, end + 32).split(/[。！？!?]/u)[0];
  return rule.nearby.test(before) || rule.nearby.test(after);
}

function codePointLabel(character) {
  return `U+${character.codePointAt(0).toString(16).toUpperCase().padStart(4, "0")}`;
}

function lint(text, label, pattern, suggestions) {
  let hasErrors = false;

  for (const [index, line] of text.split(/\r\n?|\n/).entries()) {
    let column = 1;
    let previousIndex = 0;

    for (const match of line.matchAll(pattern)) {
      const [found] = match;
      column += [...line.slice(previousIndex, match.index)].length;
      previousIndex = match.index;
      const location = `${label}:${index + 1}:${column}`;
      const suggestion = suggestions.get(found);

      if (found === "\uFFFD") {
        console.log(`${location}: error [encoding] Unicode replacement character (U+FFFD) is already present in the input.`);
        hasErrors = true;
      } else if (suggestion === undefined) {
        const ideograph = found.normalize("NFKC");
        console.log(`${location}: error [kangxi-radical] 「${found}」 (${codePointLabel(found)}) is a radical symbol. In ordinary prose, consider 「${ideograph}」 (${codePointLabel(ideograph)}).`);
        hasErrors = true;
      } else if (fitsTerminologyContext(line, match)) {
        console.log(`${location}: warn [terminology] 「${found}」 → 「${suggestion}」 if it fits the context. Keep intentional wording.`);
      }
    }
  }

  return hasErrors;
}

async function readStdin() {
  const chunks = [];
  for await (const chunk of stdin) chunks.push(chunk);
  return Buffer.concat(chunks);
}

async function main() {
  const { values, positionals } = parseArgs({
    options: {
      help: { type: "boolean", short: "h" },
      terminology: { type: "boolean" },
      "kangxi-radicals": { type: "boolean" },
      encoding: { type: "boolean" },
    },
    allowPositionals: true,
  });

  if (values.help) {
    console.log(usage);
    return;
  }

  if (!positionals.length) throw new Error("Pass text files or use - to read stdin.");
  if (positionals.filter((file) => file === "-").length > 1) {
    throw new Error("Standard input can only be read once.");
  }

  const { terminology, "kangxi-radicals": kangxiRadicals, encoding } = values;
  const useDefaultChecks = !terminology && !kangxiRadicals && !encoding;
  const suggestions = terminology || useDefaultChecks
    ? terminologySuggestions(await readFile(terminologyFile, "utf8"))
    : new Map();
  const checkEncoding = encoding || useDefaultChecks;
  const pattern = diagnosticPattern(suggestions, {
    checkKangxiRadicals: kangxiRadicals || useDefaultChecks,
    checkEncoding,
  });

  for (const file of positionals) {
    try {
      const input = file === "-" ? await readStdin() : await readFile(file);
      const label = file === "-" ? "<stdin>" : file;
      // Check the bytes before decoding would introduce replacement characters.
      if (checkEncoding && !isUtf8(input)) {
        console.log(`${label}: error [encoding] Input is not valid UTF-8. Check the original encoding before editing.`);
        process.exitCode ??= 1;
        continue;
      }

      if (lint(input.toString("utf8"), label, pattern, suggestions)) {
        process.exitCode ??= 1;
      }
    } catch (error) {
      console.error(error.message);
      process.exitCode = 2;
    }
  }
}

try {
  await main();
} catch (error) {
  console.error(error.message);
  console.error(usage);
  process.exitCode = 2;
}
