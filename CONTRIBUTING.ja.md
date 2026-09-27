# Contributing

[English](CONTRIBUTING.md) | **日本語**

このリポジトリは、otnc の `skills_*` スキルリポジトリの生成されたミラーです。`skills/` と `sync-state.json` の中身はすべて同期ワークフローが上書きするため、スキルの中身に関するIssueやPull Requestは元のリポジトリ(下記)で開いてください。**ここ**にIssueやPull Requestを開くのは、同期の仕組みそのもの(`cmd/sync`、ワークフロー、生成されるドキュメント)が壊れている・バグがある場合だけにしてください:

<!-- sources-list:start -->

- [otnc/skills\_ai-agent-saving-tech](https://github.com/otnc/skills_ai-agent-saving-tech) ([MIT](https://github.com/otnc/skills_ai-agent-saving-tech/blob/main/LICENSE))
- [otnc/skills\_fix-unnatural-line-breaks](https://github.com/otnc/skills_fix-unnatural-line-breaks) ([MIT](https://github.com/otnc/skills_fix-unnatural-line-breaks/blob/main/LICENSE))

<!-- sources-list:end -->

## 新しいスキルリポジトリを追加する

1. `otnc` の下に `skills_<何か>` という名前の公開リポジトリを作り、スキルを `skills/<名前>/SKILL.md` に置く(ルート直下の `SKILL.md` でも可)。
2. リリースを公開する(例: `v0.1.0`)。draft と prerelease のリリースは対象外。
3. 次の同期実行(1日2回)を待つか、即時に実行する:

   ```bash
   gh api repos/otnc/agent-skills/dispatches -f event_type=sync
   ```

同期は候補ごとに検証する: 公開されていること、forkでないこと、アーカイブされていないこと、`SKILL.md` があること(ルートまたは `skills/` 直下)、公開済みのリリースが1つ以上あること。検証を通らなかったリポジトリはスキップされ、ワークフローのログに報告されます。同期されると、READMEのスキル一覧とこのファイルの元リポジトリ一覧は自動的に再生成されます。

特定のリポジトリを恒久的に除外したい場合は、[`sources.json`](sources.json) の `exclude` に追加してください。

## エージェント向けスキル

AIコーディングエージェントは、[`skills` CLI](https://www.npmjs.com/package/skills)(`npx skills`)経由で [`.agents/skills/`](.agents/skills/) から追加の指示(「スキル」)を読み込みます。コミットされているのは `.agents/skills/` と [`skills-lock.json`](skills-lock.json) だけで、これが正本です。各エージェントが実際に参照する他の場所(`.claude/skills/`、`agent/skills/` など)は生成されたシンボリックリンクやコピーであり、gitignoreされています。クローン後は以下で復元してください。

```sh
npx skills experimental_install
```

インストールされているスキル([kiritan](.agents/skills/kiritan/SKILL.md))は、このリポジトリが使っているKiritanのi18n構成での作業のしかたをエージェントに教えます。`base/` 配下を編集する前に読み込んでください。

## 同期をローカルで実行する

```bash
go run ./cmd/sync --dry-run   # 計画のみ、変更なし
go run ./cmd/sync             # 適用、コミットとタグ付けをローカルで
go run ./cmd/sync --replay    # 初回セットアップ: 過去の全リリースから履歴を再構築
```

`GITHUB_TOKEN`/`GH_TOKEN` は公開リポジトリでは任意ですが、設定するとレート制限を回避できます。コミットの作成者は、`SYNC_GIT_NAME`/`SYNC_GIT_EMAIL` を設定した場合はその値(ワークフローでは `github-actions[bot]`)、未設定の場合はローカルのgit設定になります。

## READMEとCONTRIBUTING

`README.md`、`README.ja.md`、`CONTRIBUTING.md`、`CONTRIBUTING.ja.md` は、[`base/`](base/) 配下のファイルから、devDependenciesに入っている [kiritan](https://github.com/otnc/kiritan) で生成されます。ベースファイルを編集したら、以下を実行します。

```sh
npm install
npm run build      # kiritan build
npm run docs:check # kiritan check: 生成済みドキュメントが古くなっていれば失敗する
```

`base/README.base.md` のスキル一覧の節と `base/CONTRIBUTING.base.md` の元リポジトリ一覧は同期が再生成するので、手で編集しないでください。
