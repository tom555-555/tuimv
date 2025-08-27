# tuimv - ファイル移動ツール

> [English version is here](README.md) / [英語版はこちら](README.md)

tuimvは、ターミナルベースのファイル移動ツールで、直感的なUIと効率的なファイル操作を提供します。

## 概要

tuimvは、Bubble Tea（Go言語のTUIフレームワーク）を使用して構築された、6つのパネルレイアウトを持つファイル管理ツールです。ユーザーは現在のディレクトリからファイルを選択し、ターゲットディレクトリに移動することができます。

## アーキテクチャ

### パッケージ構造

```
better-mv/
├── cmd/bmv/           # メインエントリーポイント
├── internal/
│   ├── model/         # データモデルとアプリケーション状態
│   └── ui/            # ユーザーインターフェースとTUIロジック
└── specifications/     # プロジェクト仕様書
```

### コアコンポーネント

#### 1. アプリケーション状態管理 (`internal/model/`)

- **`AppState`**: アプリケーションの全体状態を管理
- **`PanelType`**: 6つのパネルの種類を定義
- **`AppMode`**: アプリケーションの動作モード（通常、検索、移動）
- **`FileInfo`**: ファイル情報の構造体
- **`DirectoryInfo`**: ディレクトリ情報の構造体

#### 2. ユーザーインターフェース (`internal/ui/`)

- **`Model`**: メインUIモデルと状態管理
- **`View`**: 画面レンダリングロジック
- **`Keybindings`**: キーボード入力処理
- **`Messages`**: カスタムメッセージタイプ
- **`Styles`**: UIスタイリング

## パネルレイアウト

アプリケーションは6つのパネルで構成されています：

```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│ Panel 1:       │ Panel 3:       │ Panel 5:       │ Panel 6:       │
│ Current Dir    │ Target Dir      │ Selected Files │ Fuzzy Search   │
│ Input          │ Input           │ List           │ Results        │
├─────────────────┼─────────────────┼─────────────────┼─────────────────┤
│ Panel 2:       │ Panel 4:       │                 │                 │
│ Current Files  │ Target Files    │                 │                 │
│ List           │ List            │                 │                 │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

### パネルの詳細

1. **Current Directory Input** - 現在のディレクトリパス入力
2. **Current Files List** - 現在のディレクトリのファイル一覧
3. **Target Directory Input** - ターゲットディレクトリパス入力
4. **Target Files List** - ターゲットディレクトリのファイル一覧
5. **Selected Files List** - 移動対象として選択されたファイル一覧
6. **Fuzzy Search Results** - ディレクトリ検索結果

## 動作フロー

### 1. アプリケーション起動

```go
func main() {
    m := ui.NewModel()           // UIモデルの初期化
    p := tea.NewProgram(m)       // Bubble Teaプログラムの作成
    p.Run()                      // プログラムの実行
}
```

### 2. モデル初期化

```go
func NewModel() *Model {
    state := model.NewAppState()  // アプリケーション状態の初期化
    
    // 各UIコンポーネントの初期化
    currentDirInput := NewShortTextInput(...)
    targetDirInput := NewShortTextInput(...)
    // ... その他のコンポーネント
    
    return &Model{...}
}
```

### 3. イベント処理

アプリケーションは以下のメッセージタイプを処理します：

- **`tea.WindowSizeMsg`**: ウィンドウサイズの変更
- **`tea.KeyMsg`**: キーボード入力
- **`PanelSwitchMsg`**: パネル切り替え
- **`FileSelectedMsg`**: ファイル選択
- **`DirectoryChangedMsg`**: ディレクトリ変更
- **`SearchQueryChangedMsg`**: 検索クエリ変更

### 4. パネルナビゲーション

ユーザーは以下の方法でパネル間を移動できます：

- **Vimスタイル**: `h`, `j`, `k`, `l`キー
- **Tabナビゲーション**: `Tab`, `Shift+Tab`
- **直接指定**: カスタムメッセージ

```go
func (m *Model) getPanelToLeft() model.PanelType {
    switch m.state.ActivePanel {
    case model.TargetDirInput:
        return model.CurrentDirInput
    case model.TargetFilesList:
        return model.CurrentFilesList
    // ... その他のケース
    }
}
```

### 5. ファイル操作

#### ファイル選択
```go
func (m Model) handleFileSelectedMsg(msg FileSelectedMsg) (tea.Model, tea.Cmd) {
    // ファイルが既に選択されているかチェック
    for _, selected := range m.state.SelectedFiles {
        if selected.AbsPath == msg.File.AbsPath {
            return m, nil // 既に選択済み
        }
    }
    
    // ファイルを選択済みリストに追加
    msg.File.IsSelected = true
    m.state.SelectedFiles = append(m.state.SelectedFiles, *msg.File)
    
    return m, nil
}
```

#### ディレクトリ変更
```go
func (m Model) handleDirectoryChangedMsg(msg DirectoryChangedMsg) (tea.Model, tea.Cmd) {
    if msg.IsTarget {
        m.state.TargetDir = msg.Path
        m.targetDirInput.SetValue(msg.Path)
    } else {
        m.state.CurrentDir = msg.Path
        m.currentDirInput.SetValue(msg.Path)
    }
    
    // TODO: ディレクトリスキャンのトリガー
    return m, nil
}
```

## キーバインド

### グローバルキーバインド

- **`Ctrl+C` / `q`**: アプリケーション終了
- **`Esc`**: 検索モードのクリア、選択の解除
- **`/`**: 検索モードのアクティベート

### ナビゲーション

- **`h`, `j`, `k`, `l`**: Vimスタイルのパネル移動
- **`Tab`**: 次のパネルに移動
- **`Shift+Tab`**: 前のパネルに移動

### パネル固有のキーバインド

#### ファイルリストパネル
- **`↑`, `↓`**: カーソル移動
- **`Space`**: ファイル選択/選択解除
- **`Enter`**: ディレクトリナビゲーション

#### 選択済みファイルパネル
- **`Delete` / `Backspace`**: 選択からファイルを削除
- **`Cmd+Enter` / `Ctrl+Enter`**: ファイル移動の実行

## 状態管理

### AppState構造体

```go
type AppState struct {
    CurrentDir    string            // 現在のディレクトリパス
    TargetDir     string            // ターゲットディレクトリパス
    CurrentFiles  []FileInfo        // 現在のディレクトリのファイル
    TargetFiles   []FileInfo        // ターゲットディレクトリのファイル
    SelectedFiles []FileInfo        // 移動対象として選択されたファイル
    SearchResults []DirectoryInfo   // 検索結果
    SearchQuery   string            // 現在の検索クエリ
    ActivePanel   PanelType         // 現在アクティブなパネル
    CurrentIndex  map[PanelType]int // 各パネルのカーソル位置
    Mode          AppMode           // 現在のアプリケーションモード
    IsSearching   bool              // 検索モードの状態
}
```

### パネル状態の追跡

各パネルのカーソル位置は`CurrentIndex`マップで管理され、パネル切り替え時にも保持されます。

## メッセージシステム

Bubble Teaのメッセージシステムを使用して、UIコンポーネント間の通信を実現しています：

```go
// パネル切り替えメッセージ
type PanelSwitchMsg struct {
    Panel model.PanelType
}

// ファイル選択メッセージ
type FileSelectedMsg struct {
    File  *model.FileInfo
    Index int
}

// ディレクトリ変更メッセージ
type DirectoryChangedMsg struct {
    Path     string
    IsTarget bool
}
```

## 拡張性

### 新しいパネルの追加

1. `model.PanelType`に新しい定数を追加
2. `Model`構造体にUIコンポーネントを追加
3. `View()`関数でレンダリングロジックを実装
4. キーバインド処理を追加
5. ナビゲーションロジックを更新

### 新しいメッセージタイプの追加

1. `messages.go`に新しいメッセージ構造体を定義
2. `Update()`関数でメッセージハンドラーを実装
3. 必要に応じてUIコンポーネントを更新

## 今後の実装予定

- [ ] ファイル移動操作の実装
- [ ] ディレクトリスキャンの実装
- [ ] エラー処理とユーザー通知
- [ ] 設定ファイルのサポート
- [ ] プラグインシステム
- [ ] 履歴機能

## 技術スタック

- **言語**: Go
- **TUIフレームワーク**: Bubble Tea
- **スタイリング**: Lip Gloss
- **アーキテクチャ**: MVCパターン
- **状態管理**: カスタム状態マシン

## 開発環境のセットアップ

```bash
# リポジトリのクローン
git clone <repository-url>
cd better-mv

# 依存関係のインストール
go mod download

# アプリケーションの実行
go run cmd/bmv/main.go

# ビルド
go build -o bmv cmd/bmv/main.go
```

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。 
