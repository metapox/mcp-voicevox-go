# MCP VOICEVOX Go - Model Context Protocol 仕様書

## 概要

MCP VOICEVOX Goは、VOICEVOXエンジンを使用してテキストを音声に変換するModel Context Protocol (MCP) サーバーです。

## サポートするプロトコル

- **Stdio**: 標準入出力経由でのJSON-RPC 2.0通信
- **HTTP/WebSocket**: HTTP APIとWebSocket経由でのMCP通信

## MCP メソッド

### 1. initialize

サーバーの初期化を行います。

**リクエスト:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {}
}
```

**レスポンス:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "tools": {}
    },
    "serverInfo": {
      "name": "mcp-voicevox-go",
      "version": "1.0.0"
    }
  }
}
```

### 2. tools/list

利用可能なツール一覧を取得します。

**リクエスト:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

**レスポンス:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "tools": [
      {
        "name": "text_to_speech",
        "description": "テキストを音声に変換します",
        "inputSchema": {
          "type": "object",
          "properties": {
            "text": {
              "type": "string",
              "description": "音声に変換するテキスト"
            },
            "speaker_id": {
              "type": "integer",
              "description": "話者ID（省略時はデフォルト話者を使用）"
            }
          },
          "required": ["text"]
        }
      },
      {
        "name": "get_speakers",
        "description": "利用可能な話者一覧を取得します",
        "inputSchema": {
          "type": "object",
          "properties": {}
        }
      }
    ]
  }
}
```

### 3. tools/call

ツールを実行します。

#### text_to_speech ツール

**リクエスト:**
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "text_to_speech",
    "arguments": {
      "text": "こんにちは、世界！",
      "speaker_id": 3
    }
  }
}
```

**レスポンス:**
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "音声合成が完了しました。\nテキスト: こんにちは、世界！\n話者ID: 3\nファイル: /tmp/mcp-voicevox/speech_3_1699123456.wav\n状態: ファイルに保存され、音声を再生しました"
      }
    ]
  }
}
```

#### get_speakers ツール

**リクエスト:**
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "get_speakers",
    "arguments": {}
  }
}
```

**レスポンス:**
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "利用可能な話者一覧:\n[\n  {\n    \"name\": \"四国めたん\",\n    \"speaker_uuid\": \"7ffcb7ce-00ec-4bdc-82cd-45a8889e43ff\",\n    \"styles\": [\n      {\n        \"name\": \"ノーマル\",\n        \"id\": 2\n      }\n    ]\n  }\n]"
      }
    ]
  }
}
```

## エラーレスポンス

エラーが発生した場合、以下の形式でレスポンスが返されます：

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32603,
    "message": "Internal error"
  }
}
```

### エラーコード

| コード | 説明 |
|--------|------|
| -32700 | Parse error - JSON解析エラー |
| -32600 | Invalid Request - 無効なリクエスト |
| -32601 | Method not found - メソッドが見つからない |
| -32602 | Invalid params - 無効なパラメータ |
| -32603 | Internal error - 内部エラー |
| -40001 | VOICEVOX connection error - VOICEVOX接続エラー |
| -40002 | Audio synthesis error - 音声合成エラー |
| -40003 | Audio playback error - 音声再生エラー |
| -40004 | File operation error - ファイル操作エラー |
| -40005 | Configuration error - 設定エラー |

## 設定

### 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| `MCP_VOICEVOX_URL` | VOICEVOXのAPIエンドポイント | `http://localhost:50021` |
| `MCP_VOICEVOX_PORT` | サーバーのポート番号（serverモードのみ） | `8080` |
| `MCP_VOICEVOX_TEMP_DIR` | 一時ファイルディレクトリ | システムの一時ディレクトリ |
| `MCP_VOICEVOX_DEFAULT_SPEAKER` | デフォルトの話者ID | `3` |
| `MCP_VOICEVOX_ENABLE_PLAYBACK` | 音声の自動再生を有効にする | `false` |

### コマンドラインオプション

#### server サブコマンド

```bash
mcp-voicevox server [flags]
```

| フラグ | 短縮形 | 説明 | デフォルト値 |
|--------|--------|------|-------------|
| `--port` | `-p` | サーバーのポート番号 | `8080` |
| `--voicevox-url` | `-u` | VOICEVOXのAPIエンドポイント | `http://localhost:50021` |
| `--temp-dir` | `-t` | 一時ファイルを保存するディレクトリ | システムの一時ディレクトリ |
| `--default-speaker` | `-s` | デフォルトの話者ID | `3` |

#### stdio サブコマンド

```bash
mcp-voicevox stdio [flags]
```

| フラグ | 短縮形 | 説明 | デフォルト値 |
|--------|--------|------|-------------|
| `--voicevox-url` | `-u` | VOICEVOXのAPIエンドポイント | `http://localhost:50021` |
| `--temp-dir` | `-t` | 一時ファイルを保存するディレクトリ | システムの一時ディレクトリ |
| `--default-speaker` | `-s` | デフォルトの話者ID | `3` |
| `--enable-playback` | | 音声の自動再生を有効にする | `false` |

## 使用例

### Amazon Q CLI との連携

```bash
# MCPサーバーとして登録
q mcp add --name voicevox --command "mcp-voicevox stdio"

# 環境変数を使用する場合
q mcp add --name voicevox --command "mcp-voicevox stdio" \
  --env MCP_VOICEVOX_DEFAULT_SPEAKER=1 \
  --env MCP_VOICEVOX_ENABLE_PLAYBACK=true
```

### Docker での使用

```bash
# サーバーモードで起動
docker run -p 8080:8080 mcp-voicevox:latest

# 環境変数を指定
docker run -p 8080:8080 \
  -e MCP_VOICEVOX_URL=http://voicevox:50021 \
  -e MCP_VOICEVOX_DEFAULT_SPEAKER=1 \
  mcp-voicevox:latest
```

## 制限事項

- テキストの最大長: 1000文字
- 同時接続数: 制限なし（リソースに依存）
- 音声ファイル形式: WAV形式のみ
- サポートプラットフォーム: Linux, macOS, Windows

## トラブルシューティング

### よくある問題

1. **VOICEVOX接続エラー**
   - VOICEVOXエンジンが起動していることを確認
   - `MCP_VOICEVOX_URL` の設定を確認

2. **音声再生エラー**
   - 音声再生コマンドがインストールされていることを確認
   - macOS: `afplay`, Linux: `paplay`/`aplay`/`mpv`, Windows: PowerShell

3. **ファイル権限エラー**
   - 一時ディレクトリの書き込み権限を確認
   - `MCP_VOICEVOX_TEMP_DIR` の設定を確認
