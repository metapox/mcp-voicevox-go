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
- `--default-speed-scale`: デフォルトの話速（0.5-2.0、デフォルト: 1.0）
- `--default-pitch-scale`: デフォルトの音高（-0.15-0.15、デフォルト: 0.0）
- `--default-intonation-scale`: デフォルトの抑揚（0.0-2.0、デフォルト: 1.0）
- `--default-volume-scale`: デフォルトの音量（0.0-2.0、デフォルト: 1.0）

#### stdioサブコマンド

- `--voicevox-url`, `-u`: VOICEVOXのAPIエンドポイント（デフォルト: http://localhost:50021）
- `--temp-dir`, `-t`: 一時ファイルを保存するディレクトリ（デフォルト: システムの一時ディレクトリ）
- `--default-speaker`, `-s`: デフォルトの話者ID（デフォルト: 3）
- `--enable-playback`: 音声の自動再生を有効にする
- `--default-speed-scale`: デフォルトの話速（0.5-2.0、デフォルト: 1.0）
- `--default-pitch-scale`: デフォルトの音高（-0.15-0.15、デフォルト: 0.0）
- `--default-intonation-scale`: デフォルトの抑揚（0.0-2.0、デフォルト: 1.0）
- `--default-volume-scale`: デフォルトの音量（0.0-2.0、デフォルト: 1.0）

### 環境変数

以下の環境変数で設定を変更できます：

- `MCP_VOICEVOX_PORT`: サーバーのポート番号（serverサブコマンドのみ）
- `MCP_VOICEVOX_URL`: VOICEVOXのAPIエンドポイント
- `MCP_VOICEVOX_TEMP_DIR`: 一時ファイルディレクトリ
- `MCP_VOICEVOX_DEFAULT_SPEAKER`: デフォルトの話者ID
- `MCP_VOICEVOX_ENABLE_PLAYBACK`: 音声の自動再生を有効にする（true/false）
- `MCP_VOICEVOX_DEFAULT_SPEED_SCALE`: デフォルトの話速（0.5-2.0）
- `MCP_VOICEVOX_DEFAULT_PITCH_SCALE`: デフォルトの音高（-0.15-0.15）
- `MCP_VOICEVOX_DEFAULT_INTONATION_SCALE`: デフォルトの抑揚（0.0-2.0）
- `MCP_VOICEVOX_DEFAULT_VOLUME_SCALE`: デフォルトの音量（0.0-2.0）

例:

```bash
export MCP_VOICEVOX_DEFAULT_SPEAKER=1
export MCP_VOICEVOX_PORT=8888
export MCP_VOICEVOX_ENABLE_PLAYBACK=true
export MCP_VOICEVOX_DEFAULT_SPEED_SCALE=1.5
go run cmd/*.go server
```

### Amazon Q CLIとの連携

Amazon Q CLIで使用するには、MCPサーバーとして登録します:

```bash
q mcp add --name voicevox --command "mcp-voicevox stdio"
```

環境変数を使用する場合:

```bash
q mcp add --name voicevox --command "mcp-voicevox stdio" --env MCP_VOICEVOX_DEFAULT_SPEAKER=1 --env MCP_VOICEVOX_ENABLE_PLAYBACK=true --env MCP_VOICEVOX_DEFAULT_SPEED_SCALE=1.5
```

これにより、Amazon Q内で `voicevox___text_to_speech` と `voicevox___get_speakers` ツールが利用可能になります。

## 機能

- `text_to_speech`: テキストを音声に変換します
  - パラメータ:
    - `text`: 音声に変換するテキスト（必須）
    - `speaker_id`: 話者ID（省略時はデフォルト話者を使用）
    - `speed_scale`: 話速（0.5-2.0、省略時はデフォルト値を使用）
    - `pitch_scale`: 音高（-0.15-0.15、省略時はデフォルト値を使用）
    - `intonation_scale`: 抑揚（0.0-2.0、省略時はデフォルト値を使用）
    - `volume_scale`: 音量（0.0-2.0、省略時はデフォルト値を使用）
  - 音声再生が有効な場合、合成後に自動で音声を再生します

- `get_speakers`: 利用可能な話者一覧を取得します

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

## 使用例

### 1.5倍速で起動
```bash
# コマンドライン引数
mcp-voicevox stdio --default-speed-scale 1.5 --enable-playback

# 環境変数
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

### Amazon Q CLI での設定例
```bash
# 1.5倍速設定
q mcp add --name voicevox-fast \
  --command "mcp-voicevox stdio" \
  --env MCP_VOICEVOX_DEFAULT_SPEED_SCALE=1.5 \
  --env MCP_VOICEVOX_ENABLE_PLAYBACK=true

# 2倍速設定
q mcp add --name voicevox-2x \
  --command "mcp-voicevox stdio --default-speed-scale 2.0 --enable-playback"
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

