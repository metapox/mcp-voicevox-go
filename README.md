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

### HTTPサーバーとしての起動

```bash
go run cmd/*.go server
```

### Stdio経由での起動（MCP用）

```bash
go run cmd/*.go stdio
```

### オプション

#### serverサブコマンド

- `--port`, `-p`: サーバーのポート番号（デフォルト: 8080）
- `--voicevox-url`, `-u`: VOICEVOXのAPIエンドポイント（デフォルト: http://localhost:50021）
- `--temp-dir`, `-t`: 一時ファイルを保存するディレクトリ（デフォルト: システムの一時ディレクトリ）
- `--default-speaker`, `-s`: デフォルトの話者ID（デフォルト: 3）

#### stdioサブコマンド

- `--voicevox-url`, `-u`: VOICEVOXのAPIエンドポイント（デフォルト: http://localhost:50021）
- `--temp-dir`, `-t`: 一時ファイルを保存するディレクトリ（デフォルト: システムの一時ディレクトリ）
- `--default-speaker`, `-s`: デフォルトの話者ID（デフォルト: 3）
- `--enable-playback`: 音声の自動再生を有効にする

### 環境変数

以下の環境変数で設定を変更できます：

- `MCP_VOICEVOX_PORT`: サーバーのポート番号（serverサブコマンドのみ）
- `MCP_VOICEVOX_URL`: VOICEVOXのAPIエンドポイント
- `MCP_VOICEVOX_TEMP_DIR`: 一時ファイルディレクトリ
- `MCP_VOICEVOX_DEFAULT_SPEAKER`: デフォルトの話者ID
- `MCP_VOICEVOX_ENABLE_PLAYBACK`: 音声の自動再生を有効にする（true/false）

例:

```bash
export MCP_VOICEVOX_DEFAULT_SPEAKER=1
export MCP_VOICEVOX_PORT=8888
export MCP_VOICEVOX_ENABLE_PLAYBACK=true
go run cmd/*.go server
```

### Amazon Q CLIとの連携

Amazon Q CLIで使用するには、MCPサーバーとして登録します:

```bash
q mcp add --name voicevox --command "mcp-voicevox stdio"
```

環境変数を使用する場合:

```bash
q mcp add --name voicevox --command "mcp-voicevox stdio" --env MCP_VOICEVOX_DEFAULT_SPEAKER=1 --env MCP_VOICEVOX_ENABLE_PLAYBACK=true
```

これにより、Amazon Q内で `voicevox___text_to_speech` と `voicevox___get_speakers` ツールが利用可能になります。

## 機能

- `text_to_speech`: テキストを音声に変換します
  - パラメータ:
    - `text`: 音声に変換するテキスト（必須）
    - `speaker_id`: 話者ID（省略時はデフォルト話者を使用）
  - 音声再生が有効な場合、合成後に自動で音声を再生します

- `get_speakers`: 利用可能な話者一覧を取得します

## 音声再生対応プラットフォーム

- **macOS**: `afplay` コマンドを使用
- **Linux**: `paplay`、`aplay`、`mpv` のいずれかを使用（利用可能なものを自動検出）
- **Windows**: PowerShell の `Media.SoundPlayer` を使用

## ビルド

```bash
go build -o mcp-voicevox cmd/*.go
```

## ライセンス

MIT

