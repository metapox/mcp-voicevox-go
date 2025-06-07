# MCP VOICEVOX Go

VOICEVOXを使用したModel Context Protocol (MCP) サーバーの実装です。このサーバーを使用することで、LLMからVOICEVOXの音声合成機能を利用できるようになります。

## 前提条件

- Go 1.21以上
- [VOICEVOX](https://voicevox.hiroshiba.jp/) がローカルで実行されていること（デフォルトでは `http://localhost:50021` でアクセス可能）

## インストール

```bash
git clone https://github.com/metapox/mcp-voicevox-go.git
cd mcp-voicevox-go
go mod tidy
```

## 使い方

### サーバーの起動

```bash
go run cmd/server/main.go server
```

### オプション

- `--port`, `-p`: サーバーのポート番号（デフォルト: 8080）
- `--voicevox-url`, `-u`: VOICEVOXのAPIエンドポイント（デフォルト: http://localhost:50021）
- `--temp-dir`, `-t`: 一時ファイルを保存するディレクトリ（デフォルト: システムの一時ディレクトリ）
- `--default-speaker`, `-s`: デフォルトの話者ID（デフォルト: 3）

### 環境変数

以下の環境変数で設定を変更できます：

- `MCP_VOICEVOX_PORT`: サーバーのポート番号
- `MCP_VOICEVOX_URL`: VOICEVOXのAPIエンドポイント
- `MCP_VOICEVOX_TEMP_DIR`: 一時ファイルディレクトリ
- `MCP_VOICEVOX_DEFAULT_SPEAKER`: デフォルトの話者ID

例:

```bash
export MCP_VOICEVOX_DEFAULT_SPEAKER=1
export MCP_VOICEVOX_PORT=8888
go run cmd/server/main.go server
```

### Amazon Q CLIとの連携

Amazon Q CLIで使用するには、MCPサーバーとして登録します:

```bash
q mcp add --name voicevox --command "./mcp-voicevox server"
```

環境変数を使用する場合:

```bash
q mcp add --name voicevox --command "./mcp-voicevox server" --env MCP_VOICEVOX_DEFAULT_SPEAKER=1
```

これにより、Amazon Q内で `voicevox___text_to_speech` と `voicevox___get_speakers` ツールが利用可能になります。

## 機能

- `text_to_speech`: テキストを音声に変換します
  - パラメータ:
    - `text`: 音声に変換するテキスト（必須）
    - `speaker_id`: 話者ID（省略時はデフォルト話者を使用）

- `get_speakers`: 利用可能な話者一覧を取得します

## ビルド

```bash
go build -o mcp-voicevox cmd/server/main.go
```

## ライセンス

MIT

