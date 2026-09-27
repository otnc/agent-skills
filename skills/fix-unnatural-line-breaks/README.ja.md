[English](README.md) | 日本語

# fix-unnatural-line-breaks

文章の途中で、一定の文字数(カラム幅)に合わせて改行を挿入する「不自然な改行」を検出し、直すための Claude Skill です。

## これは何を直すのか

こういう改行を捕まえます。

```
吾輩は猫である。名前はまだ無い。どこで生れたかとんと
見当がつかぬ。何でも薄暗いじめじめした所でニャーニャー
泣いていた事だけは記憶している。
```

改行は段落の区切り・箇条書きの項目・見出しの前後といった、意味のある境界にだけ入れる、というのが基本方針です。

```
吾輩は猫である。名前はまだ無い。どこで生れたかとんと見当がつかぬ。何でも薄暗いじめじめした所でニャーニャー泣いていた事だけは記憶している。
```

対象はMarkdownなどのドキュメント本文と、コード中のコメント・docstring・JSDocの両方です。日本語・英語のどちらにも対応します。

## なぜ直すのか

多くのMarkdownレンダラーは段落内の単一改行を空白1つとして扱うため、見た目には影響しないことが多いですが、次のような場面で実害があります。

- ソース自体が読みにくい(diffやレビューで、文の途中で行が切れているのは不自然)
- 1文字の修正のつもりで段落全体を書き直すと、diffが無意味に大きくなる
- コードコメントやターミナル、コミットメッセージなど、改行がそのまま表示される環境では、見た目にもそのまま崩れる

## 構成

```
fix-unnatural-line-breaks/
├── SKILL.md              # コアルール + 進め方 + クイックチェック
├── references/
│   ├── detection.md      # 不自然な改行の見分け方(日本語/英語)
│   ├── rules.md          # 直し方のルール
│   └── examples.md       # before/after の対比例
├── scripts/
│   ├── lint.mjs          # 疑わしい改行を機械的に検出するスクリプト(Node.js、依存パッケージなし)
│   └── lint.test.mjs     # lint.mjsのnode:testテストスイート(`node --test scripts/*.test.mjs`)
├── README.md
├── README.ja.md
└── LICENSE
```

SKILL.md・references・scripts の中身はすべて英語で書かれています。Skill自体が日本語話者にも英語話者にも同じように使えるようにするための判断です(トリガーとなる`description`も英語・日本語両方のフレーズを含めています)。日本語のドキュメント・コメントに対する検出ルールや実例は、英語の説明文の中で日本語のサンプル文を示す形で扱っています。

## インストール

**Claude Code(個人用)**

```bash
git clone https://github.com/otnc/skills_fix-unnatural-line-breaks ~/.claude/skills/fix-unnatural-line-breaks
```

**Claude Code(プロジェクト単位)**

```bash
git clone https://github.com/otnc/skills_fix-unnatural-line-breaks <project>/.claude/skills/fix-unnatural-line-breaks
```

## 使い方

検出はNode.jsだけで機械的に行えます。インストールも依存パッケージも不要です(Claude Code自体がNode.js製のCLIとして配布されているため、このSkillが実際に動く環境には最初からNode.jsが入っていると考えて問題ありません)。

```bash
node scripts/lint.mjs path/to/README.md
node scripts/lint.mjs --json path/to/file.ts
```

Node.jsが無い環境では、Claude が `references/detection.md` の観点で目視チェックします。

検出結果はあくまで疑いの提示であり、網羅的なリストではありません。スクリプトは少数の正規表現ヒューリスティックにすぎないため、実際は問題ない行(コードブロック・テーブル・意図的な改行など)を誤検出することもあれば、どのルールにも一致しない本物の不自然な改行を見逃すこともあります。実際に直すかどうかは文脈で判断したうえで、Claudeはフラグの立った行だけでなくファイル全体を目視でも読み、スクリプトが見逃した箇所まで拾います。

## ファイル

**SKILL.md**
大原則、やること/やらないこと、進め方、クイックセルフチェック。

**references/detection.md**
不自然な改行の見分け方。日本語特有の観点(助詞の直後・直前など)と英語特有の観点(前置詞句の途中など)、誤検出しやすい箇所を整理しています。

**references/rules.md**
検出したあとの直し方。段落単位でまとめる、残すべき改行を先に確定させる、コードコメントでの扱いなど。

**references/examples.md**
README・箇条書き・コードコメント(TypeScript/Python)のbefore/after対比。

**scripts/lint.mjs**
依存パッケージなしで動くNode.js製の検出スクリプト。行末の文末記号・助詞・接続語・次の行の書き出しに加え、コメントブロック内で段落区切りとして使われている裸の`//`・`#`・`*`行からも、機械的な折り返しの疑いを洗い出します。Markdown以外のファイルでは、コメント行(`//`・`#`・`*`・`///`・`;;`)同士だけを比較するので、`#`でコメントを書く言語(Python・Shell・Ruby・YAMLなど)でも、C系の`//`・`/** */`コメントでも動作します。出力メッセージは(チェック対象のファイルの言語に関わらず)英語です。検出結果は網羅的ではありません(`references/detection.md`参照)。

**scripts/lint.test.mjs**
`lint.mjs`の中核的な検出動作・誤検出回避動作を検証する`node:test`テストスイート。`node --test scripts/*.test.mjs`で実行できます。

## クレジット

このSkillが対象とする規則は、[otoneko.](https://github.com/otnc) が個人のCLAUDE.mdで運用していた「ドキュメントやコードコメントで不自然な改行を入れない」というルールを、汎用のSkillとして切り出したものです。

## ライセンス

MIT。詳細は `LICENSE` を参照。
