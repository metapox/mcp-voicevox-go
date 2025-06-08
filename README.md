# MCP VOICEVOX Go

VOICEVOXを使用したModel Context Protocol (MCP) サーバーの実装です。このサーバーを使用することで、LLMからVOICEVOXの音声合成機能を利用できるようになります。

## 前提条件

- Go 1.21以上
- [VOICEVOX](https://voicevox.hiroshiba.jp/) がローカルで実行されていること（デフォルトでは `http://localhost:50021` でアクセス可能）

## インストール

### Go installを使用（推奨）

```bash
go install github.com/metapox/mcp-voicevox-go@latest
```

### ソースからビルド

```bash
git clone https://github.com/metapox/mcp-voicevox-go.git
cd mcp-voicevox-go
go mod tidy
go build -o mcp-voicevox .
```

## クイックスタート

### 1. VOICEVOXエンジンの起動

VOICEVOXアプリケーションを起動するか、VOICEVOXエンジンを直接起動してください。

### 2. Amazon Q CLIでの設定

`~/.aws/amazonq/mcp.json` に以下の設定を追加：

```json
{
  "mcpServers": {
    "voicevox": {
      "command": "mcp-voicevox",
      "args": ["stdio"],
      "env": {
        "MCP_VOICEVOX_ENABLE_PLAYBACK": "true"
      }
    }
  }
}
```

### 3. Amazon Q CLIの再起動

```bash
# Amazon Q CLIを再起動して設定を反映
q chat
```

これで `voicevox___text_to_speech` と `voicevox___get_speakers` ツールが利用可能になります。

## 設定テンプレート

### 基本設定

```json
{
  "mcpServers": {
    "voicevox": {
      "command": "mcp-voicevox",
      "args": ["stdio"],
      "env": {
        "MCP_VOICEVOX_ENABLE_PLAYBACK": "true"
      },
      "autoApprove": [
        "get_speakers",
        "text_to_speech"
      ]
    }
  }
}
```

### 1.5倍速設定

```json
{
  "mcpServers": {
    "voicevox-fast": {
      "command": "mcp-voicevox",
      "args": ["stdio"],
      "env": {
        "MCP_VOICEVOX_ENABLE_PLAYBACK": "true",
        "MCP_VOICEVOX_DEFAULT_SPEED_SCALE": "1.5"
      },
      "autoApprove": [
        "get_speakers",
        "text_to_speech"
      ]
    }
  }
}
```

### 2倍速設定

```json
{
  "mcpServers": {
    "voicevox-2x": {
      "command": "mcp-voicevox",
      "args": ["stdio"],
      "env": {
        "MCP_VOICEVOX_ENABLE_PLAYBACK": "true",
        "MCP_VOICEVOX_DEFAULT_SPEED_SCALE": "2.0"
      },
      "autoApprove": [
        "get_speakers",
        "text_to_speech"
      ]
    }
  }
}
```

### カスタム音声設定

```json
{
  "mcpServers": {
    "voicevox-custom": {
      "command": "mcp-voicevox",
      "args": ["stdio"],
      "env": {
        "MCP_VOICEVOX_ENABLE_PLAYBACK": "true",
        "MCP_VOICEVOX_DEFAULT_SPEAKER": "1",
        "MCP_VOICEVOX_DEFAULT_SPEED_SCALE": "1.3",
        "MCP_VOICEVOX_DEFAULT_PITCH_SCALE": "0.05",
        "MCP_VOICEVOX_DEFAULT_INTONATION_SCALE": "1.2"
      },
      "autoApprove": [
        "get_speakers",
        "text_to_speech"
      ]
    }
  }
}
```

### 複数設定の併用

```json
{
  "mcpServers": {
    "voicevox-normal": {
      "command": "mcp-voicevox",
      "args": ["stdio"],
      "env": {
        "MCP_VOICEVOX_ENABLE_PLAYBACK": "true"
      },
      "autoApprove": ["get_speakers", "text_to_speech"]
    },
    "voicevox-fast": {
      "command": "mcp-voicevox",
      "args": ["stdio"],
      "env": {
        "MCP_VOICEVOX_ENABLE_PLAYBACK": "true",
        "MCP_VOICEVOX_DEFAULT_SPEED_SCALE": "1.5"
      },
      "autoApprove": ["get_speakers", "text_to_speech"]
    },
    "voicevox-2x": {
      "command": "mcp-voicevox",
      "args": ["stdio"],
      "env": {
        "MCP_VOICEVOX_ENABLE_PLAYBACK": "true",
        "MCP_VOICEVOX_DEFAULT_SPEED_SCALE": "2.0"
      },
      "autoApprove": ["get_speakers", "text_to_speech"]
    }
  }
}
```

## 使い方

### HTTPサーバーとしての起動

```bash
mcp-voicevox server
```

### Stdio経由での起動（MCP用）

```bash
mcp-voicevox stdio
```

## オプション

### 共通オプション

| オプション | 短縮形 | 説明 | デフォルト |
|------------|--------|------|------------|
| `--voicevox-url` | `-u` | VOICEVOXのAPIエンドポイント | `http://localhost:50021` |
| `--temp-dir` | `-t` | 一時ファイルディレクトリ | システムの一時ディレクトリ |
| `--default-speaker` | `-s` | デフォルトの話者ID | `3` |
| `--default-speed-scale` | | デフォルトの話速（0.5-2.0） | `1.0` |
| `--default-pitch-scale` | | デフォルトの音高（-0.15-0.15） | `0.0` |
| `--default-intonation-scale` | | デフォルトの抑揚（0.0-2.0） | `1.0` |
| `--default-volume-scale` | | デフォルトの音量（0.0-2.0） | `1.0` |

### serverサブコマンド専用

| オプション | 短縮形 | 説明 | デフォルト |
|------------|--------|------|------------|
| `--port` | `-p` | サーバーのポート番号 | `8080` |

### stdioサブコマンド専用

| オプション | 説明 | デフォルト |
|------------|------|------------|
| `--enable-playback` | 音声の自動再生を有効にする | `false` |

## 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|------------|
| `MCP_VOICEVOX_PORT` | サーバーのポート番号（serverのみ） | `8080` |
| `MCP_VOICEVOX_URL` | VOICEVOXのAPIエンドポイント | `http://localhost:50021` |
| `MCP_VOICEVOX_TEMP_DIR` | 一時ファイルディレクトリ | システムの一時ディレクトリ |
| `MCP_VOICEVOX_DEFAULT_SPEAKER` | デフォルトの話者ID | `3` |
| `MCP_VOICEVOX_ENABLE_PLAYBACK` | 音声の自動再生（true/false） | `false` |
| `MCP_VOICEVOX_DEFAULT_SPEED_SCALE` | デフォルトの話速（0.5-2.0） | `1.0` |
| `MCP_VOICEVOX_DEFAULT_PITCH_SCALE` | デフォルトの音高（-0.15-0.15） | `0.0` |
| `MCP_VOICEVOX_DEFAULT_INTONATION_SCALE` | デフォルトの抑揚（0.0-2.0） | `1.0` |
| `MCP_VOICEVOX_DEFAULT_VOLUME_SCALE` | デフォルトの音量（0.0-2.0） | `1.0` |

例:

```bash
export MCP_VOICEVOX_DEFAULT_SPEAKER=1
export MCP_VOICEVOX_PORT=8888
export MCP_VOICEVOX_ENABLE_PLAYBACK=true
export MCP_VOICEVOX_DEFAULT_SPEED_SCALE=1.5
mcp-voicevox server
```

### Amazon Q CLIとの連携

Amazon Q CLIで使用するには、上記のmcp.json設定を使用してください。設定後、以下のツールが利用可能になります：

- `voicevox___text_to_speech`: テキストを音声に変換
- `voicevox___get_speakers`: 利用可能な話者一覧を取得

## 機能

### text_to_speech
テキストを音声に変換します。

**パラメータ:**
- `text`: 音声に変換するテキスト（必須）
- `speaker_id`: 話者ID（省略時はデフォルト話者を使用）
- `speed_scale`: 話速（0.5-2.0、省略時はデフォルト値を使用）
- `pitch_scale`: 音高（-0.15-0.15、省略時はデフォルト値を使用）
- `intonation_scale`: 抑揚（0.0-2.0、省略時はデフォルト値を使用）
- `volume_scale`: 音量（0.0-2.0、省略時はデフォルト値を使用）

音声再生が有効な場合、合成後に自動で音声を再生します。

### get_speakers
利用可能な話者一覧を取得します。

## 音声パラメータの詳細

### 話速（Speed Scale）
- **範囲**: 0.5 - 2.0
- **デフォルト**: 1.0
- **効果**: 
  - 0.5: 半分の速度（ゆっくり）
  - 1.0: 通常速度
  - 1.5: 1.5倍速（少し早口）
  - 2.0: 2倍速（かなり早口）

### 音高（Pitch Scale）
- **範囲**: -0.15 - 0.15
- **デフォルト**: 0.0
- **効果**:
  - 負の値: 低い声
  - 0.0: 通常の音高
  - 正の値: 高い声

### 抑揚（Intonation Scale）
- **範囲**: 0.0 - 2.0
- **デフォルト**: 1.0
- **効果**:
  - 0.0: 平坦な話し方
  - 1.0: 通常の抑揚
  - 2.0: 抑揚を強調

### 音量（Volume Scale）
- **範囲**: 0.0 - 2.0
- **デフォルト**: 1.0
- **効果**:
  - 0.0: 無音
  - 1.0: 通常の音量
  - 2.0: 大音量

## コマンドライン使用例

### 基本的な使用方法

```bash
# 1.5倍速で起動
mcp-voicevox stdio --default-speed-scale 1.5 --enable-playback

# 環境変数で設定
export MCP_VOICEVOX_DEFAULT_SPEED_SCALE=1.5
mcp-voicevox stdio
```

### 複数パラメータの組み合わせ

```bash
# 少し早口で高めの声、抑揚強め
mcp-voicevox stdio \
  --default-speed-scale 1.3 \
  --default-pitch-scale 0.08 \
  --default-intonation-scale 1.4 \
  --enable-playback
```

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

