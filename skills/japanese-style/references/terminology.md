# 技術用語の例

開発者向けの文章で、左の表現が右の概念を指す場合の例。確認した公開資料には表記の揺れもある。表の言い換え例には、資料にある表記と、直訳に気づくために作った例の両方を含む。

`lint.js --terminology` はこの表を読み、左の語句を見直し候補として報告する。「文脈」「分岐」「画像」のように日常でも使う語は、同じ行の近くに対応する技術用語がある場合に絞る。たとえば「文脈に合う」はそのまま、「LLM に渡す文脈」なら「コンテキスト」を候補に出す。書き手が意図して選んだ語や、対象の意味に合う語はそのまま使う。

## ソフトウェアの配布と実行

| 見直し候補 | 用語の例 | 場面 |
| --- | --- | --- |
| 版、版番号 | バージョン、バージョン番号 | ソフトウェアや API のバージョンを特定する。 |
| 新版を公開する | 新しいバージョンをリリースする | 利用者向けに新しいバージョンを公開する。 |
| 更新履歴 | リリースノート | リリースに添える告知を指す。変更履歴全体を指す場合は「変更履歴」が自然。 |
| 導入する | インストールする | コマンドでソフトウェアを入れる手順を指す。組織で使い始める話なら「導入」が自然。 |
| 実行環境 | ランタイム | Node.js など実行を担うソフトウェア自体を指す。OS や設定まで含む場合は「実行環境」が自然。 |
| 構築する | ビルドする | ソースから配布ファイルやコンテナイメージを作る。 |

表記の確認に使った資料。[Node.js のリリース](https://nodejs.org/ja/about/previous-releases)、[Node.js のダウンロード](https://nodejs.org/ja/download)、[GitHub のリリースノート](https://docs.github.com/ja/contributing/style-guide-and-content-model/release-note-content-type)、[Azure Container Instances の FAQ](https://learn.microsoft.com/ja-jp/azure/container-instances/container-instances-faq)

## Git と CI

| 見直し候補 | 用語の例 | 場面 |
| --- | --- | --- |
| プル要求、変更取り込み依頼 | プルリクエスト | GitHub 上で変更を提案する機能を指す。プロジェクトで定着していれば「PR」とも書ける。 |
| 変更を記録する | コミットする | Git のコミットを作る操作を指す。 |
| 分岐 | ブランチ | Git で作成・切り替え・削除する対象を指す。処理の分岐には「分岐」を使う。 |
| 変更を統合する | マージする | Git のマージを指す。変更を取り込んだ方法が分からない場合は「取り込む」で伝えられる。 |
| 作業手順 | ワークフロー | GitHub Actions の workflow を指す。 |
| 作業、手順 | ジョブ、ステップ | GitHub Actions の設定上の単位を指す。ステップはジョブの中で実行する。 |
| 実行機、起動条件 | ランナー、トリガー | CI の実行を担うランナーや、実行のきっかけとなるトリガーを指す。 |
| 成果物 | アーティファクト | `upload-artifact` で保存したファイルを指す。「ビルド成果物」も、プロジェクトで定着していれば自然。 |

表記の確認に使った資料。[GitHub のプルリクエスト](https://docs.github.com/ja/pull-requests/get-started/about-pull-requests)、[GitHub の Hello World](https://docs.github.com/ja/get-started/using-github/hello-world)、[GitHub Actions のワークフロー例](https://docs.github.com/ja/actions/tutorials/create-an-example-workflow)、[GitHub Actions のアーティファクト](https://docs.github.com/ja/actions/concepts/workflows-and-actions/workflow-artifacts)

## HTTP と Web API

| 見直し候補 | 用語の例 | 場面 |
| --- | --- | --- |
| 要求、応答 | リクエスト、レスポンス | HTTP で送受信するメッセージを指す。人が出す要求やサーバーが応答する動作には、元の表現も使える。 |
| 要求の見出し | リクエストヘッダー | HTTP メッセージのヘッダーフィールドを指す。 |
| 接続先、入口 | エンドポイント | API が公開する特定の呼び出し先を指す。 |
| 積載データ | ペイロード | API が受け取る対象のデータを指す。HTTP ボディの話なら「リクエストボディ」と書ける。 |
| 再試行する | リトライする | 開発者向けの文章で通信失敗後の retry を指す。製品が「再試行」を使っているならその表記に合わせる。利用者への案内では「もう一度試す」も自然。 |
| 時間切れ | タイムアウト | クライアントやサーバーの待ち時間の上限、またはその上限に達した事象を指す。 |
| 一時保存 | キャッシュ | 再利用するために保持したデータ、またはキャッシュする処理を指す。 |
| 接続期間 | セッション | ログイン状態などを複数のリクエストにわたって維持する仕組みを指す。 |

表記の確認に使った資料。[MDN の HTTP ヘッダー](https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Headers)、[MDN のセッション管理](https://developer.mozilla.org/ja/docs/Web/Security/Authentication/Session_management)、[Microsoft Learn の Web API ガイド](https://learn.microsoft.com/ja-jp/azure/architecture/best-practices/api-implementation)、[AWS SDK for Java のリトライ設定](https://docs.aws.amazon.com/ja_jp/sdk-for-java/latest/developer-guide/retry-strategy.html)

## インフラと運用

| 見直し候補 | 用語の例 | 場面 |
| --- | --- | --- |
| 配備する、配置する | デプロイする | アプリを本番環境や検証環境に反映する。 |
| 展開する、元に戻す | ロールアウトする、ロールバックする | 更新を順次反映する操作と、以前の状態に戻す操作を指す。 |
| 複製 | レプリカ | Kubernetes などが維持するインスタンスを指す。 |
| 容器、画像、登録簿 | コンテナ、イメージ、レジストリ | Docker の実行単位、実行に使うイメージ、イメージの保管先を指す。 |
| 注釈 | アノテーション | Kubernetes のリソースに付与するメタデータを指す。本文への注釈なら「注釈」を使う。 |

表記の確認に使った資料。[Kubernetes の Deployment](https://kubernetes.io/ja/docs/concepts/workloads/controllers/deployment/)、[Kubernetes のワークロード管理](https://kubernetes.io/ja/docs/concepts/workloads/management/)、[Azure Container Registry](https://learn.microsoft.com/ja-jp/azure/container-registry/)

## AI エージェント

| 見直し候補 | 用語の例 | 場面 |
| --- | --- | --- |
| 指示文 | プロンプト | モデルに渡す入力を一つの単位として指す。個々の指示だけを指すなら「指示」でよい。 |
| 文脈、前後の事情 | コンテキスト | 会話履歴やツールの結果など、モデルが参照する情報を指す。文章の前後関係なら「文脈」が自然。 |
| 文字数、単語数 | トークン数 | モデルがテキストを処理する単位を数える。 |
| 代理処理、処理担当 | エージェント | AI エージェントやサブエージェントという機能を指す。 |

表記の確認に使った資料。[GitHub Copilot のコンテキスト管理](https://docs.github.com/ja/copilot/concepts/prompting-and-context/context-management)、[GitHub Copilot のプロンプト](https://docs.github.com/ja/copilot/concepts/prompting-and-context/prompt-engineering)、[GitHub Copilot SDK のカスタムエージェント](https://docs.github.com/ja/copilot/how-tos/copilot-sdk/features/custom-agents)
