# Coding Winner - 競プロ支援Discord Bot

AtCoderの競技プログラミングを支援するDiscord Botです。モチベーション向上と継続的な学習を促進するための様々な機能を提供します。

## 主な機能

### 1. ユーザー登録
- `/register <atcoder_username>` - AtCoderのユーザー名を登録
- 提出履歴を自動同期

### 2. コンテスト通知
- `/contest-notify <channel> [enable_reminder]` - コンテスト通知を設定
- 24時間以内に開始予定のコンテストを自動通知
- リアクションを付けたユーザーには開始30分前にDMでリマインド

### 3. 週次精進レポート
- `/weekly-report <channel>` - 週次レポートを設定
- 毎週月曜日の朝9時に先週のAC数をランキング形式で表示
- 難易度別のAC数も表示

### 4. 今日の一問
- `/daily-problem <channel> [difficulty_min] [difficulty_max]` - 今日の一問を設定
- 毎日朝9時に指定難易度範囲からランダムに問題を配信
- デフォルト難易度: 400〜800

### 5. バーチャルコンテスト
- `/virtual-create <title> <duration> <problems>` - バーチャルコンテストを作成
- `/virtual-start <contest_id>` - コンテストを開始
- `/virtual-standings <contest_id>` - 順位表を表示

### 6. 統計情報
- `/mystats` - 自分の今週の統計情報を表示

## 技術スタック

- **言語**: Go 1.21+
- **Discord Bot**: discordgo
- **データベース**: PostgreSQL
- **スケジューラ**: robfig/cron
- **インフラ**: Fly.io (完全無料)

## セットアップ

### 前提条件

- Go 1.21以上
- PostgreSQL
- Discord Bot Token

### ローカル開発

1. リポジトリをクローン

```bash
git clone https://github.com/yourusername/coding-winner.git
cd coding-winner
```

2. 依存関係をインストール

```bash
go mod download
```

3. 環境変数を設定

```bash
cp .env.example .env
# .envファイルを編集して必要な情報を入力
```

4. データベースをセットアップ

```bash
# PostgreSQLを起動してデータベースを作成
createdb coding_winner
```

5. アプリケーションを実行

```bash
go run cmd/bot/main.go
```

### Discord Botの作成

1. [Discord Developer Portal](https://discord.com/developers/applications)にアクセス
2. "New Application"をクリック
3. Botセクションに移動して"Add Bot"をクリック
4. Bot Tokenをコピーして`.env`の`DISCORD_BOT_TOKEN`に設定
5. OAuth2 > URL Generatorで以下を選択:
   - Scopes: `bot`, `applications.commands`
   - Bot Permissions: `Send Messages`, `Add Reactions`, `Read Message History`, `Use Slash Commands`
6. 生成されたURLでBotをサーバーに招待

## Fly.ioへのデプロイ

### 1. Fly.io CLIのインストール

```bash
curl -L https://fly.io/install.sh | sh
```

### 2. Fly.ioにログイン

```bash
fly auth login
```

### 3. PostgreSQLデータベースの作成

```bash
fly postgres create --name coding-winner-db --region nrt
```

データベースの接続情報を保存してください。

### 4. アプリケーションの作成

```bash
fly launch --name coding-winner-bot --region nrt --no-deploy
```

### 5. シークレットの設定

```bash
fly secrets set DISCORD_BOT_TOKEN=your_bot_token_here
fly secrets set DATABASE_URL=postgres://user:pass@host:5432/dbname
```

### 6. デプロイ

```bash
fly deploy
```

### 7. ログの確認

```bash
fly logs
```

## プロジェクト構造

```
coding-winner/
├── cmd/
│   └── bot/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── bot/                     # Discord Bot関連
│   │   ├── bot.go
│   │   ├── commands.go
│   │   └── handlers/            # コマンドハンドラー
│   ├── scheduler/               # スケジューラー
│   ├── atcoder/                 # AtCoder API クライアント
│   ├── database/                # データベース操作
│   └── models/                  # データモデル
├── migrations/                  # SQLマイグレーション
├── Dockerfile
├── fly.toml
├── go.mod
└── README.md
```

## データベーススキーマ

### テーブル

- `users` - ユーザー情報
- `contest_notifications` - コンテスト通知設定
- `submissions` - 提出履歴
- `problems` - 問題情報
- `daily_problem_config` - 今日の一問設定
- `virtual_contests` - バーチャルコンテスト
- `virtual_contest_submissions` - バーチャルコンテスト提出
- `weekly_report_config` - 週次レポート設定

詳細は `migrations/` フォルダを参照してください。

## 自動実行タスク

- **15分ごと**:
  - ユーザーの提出データを同期
  - コンテスト情報をチェックして通知
- **毎日朝3時**: 問題データを同期
- **毎日朝9時**: 今日の一問を配信
- **毎週月曜日朝9時**: 週次精進レポートを送信

## トラブルシューティング

### Goがインストールされていない

```bash
# macOS
brew install go

# Linux
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### データベース接続エラー

`.env`ファイルの`DATABASE_URL`が正しく設定されているか確認してください。

### Botがオフラインになる

Fly.ioの無料枠では、アプリがスリープする可能性があります。定期的なタスクがあるため、基本的にはスリープしません。

## コスト

完全無料で運用可能：

- **Fly.io**: 3 Shared CPU VMs（無料枠）
- **PostgreSQL**: Fly.io PostgreSQL 3GB（無料枠）
- **AtCoder Problems API**: 無料
- **Discord Bot**: 無料

## ライセンス

MIT License

## 参考

- [AtCoder Problems API](https://github.com/kenkoooo/AtCoderProblems)
- [discordgo](https://github.com/bwmarrin/discordgo)
- [Fly.io Documentation](https://fly.io/docs/)

## 貢献

プルリクエストを歓迎します！バグ報告や機能リクエストはIssueで報告してください。

## サポート

問題が発生した場合は、GitHubのIssueで報告してください。
