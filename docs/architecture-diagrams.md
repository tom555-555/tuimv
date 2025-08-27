# Better MV - アーキテクチャ図

このドキュメントでは、Better MVの動作フローとstateの変更を図で説明します。

## 1. アプリケーション起動フロー

```mermaid
flowchart TD
    A[main関数開始] --> B[tea.LogToFile設定]
    B --> C[ui.NewModel呼び出し]
    C --> D[model.NewAppState作成]
    D --> E[UIコンポーネント初期化]
    E --> F[tea.NewProgram作成]
    F --> G[プログラム実行開始]
    G --> H[Model.Init呼び出し]
    H --> I[メインループ開始]
    I --> J[イベント待機]
    J --> K{イベント発生}
    K --> L[Update関数呼び出し]
    L --> M[View関数呼び出し]
    M --> J
```

## 2. 状態管理の構造

```mermaid
classDiagram
    class AppState {
        +string CurrentDir
        +string TargetDir
        +[]FileInfo CurrentFiles
        +[]FileInfo TargetFiles
        +[]FileInfo SelectedFiles
        +[]DirectoryInfo SearchResults
        +string SearchQuery
        +PanelType ActivePanel
        +map[PanelType]int CurrentIndex
        +AppMode Mode
        +bool IsSearching
        +NewAppState() *AppState
        +GetCurrentIndex(PanelType) int
        +SetCurrentIndex(PanelType, int)
        +GetFileCount(PanelType) int
    }

    class Model {
        +*AppState state
        +model.PanelType currentPanel
        +textinput.Model currentDirInput
        +textinput.Model targetDirInput
        +textinput.Model searchInput
        +list.Model currentFilesList
        +list.Model targetFilesList
        +list.Model selectedFilesList
        +list.Model searchResultsList
        +int width
        +int height
        +*Styles styles
        +NewModel() *Model
        +Init() tea.Cmd
        +Update(tea.Msg) (tea.Model, tea.Cmd)
        +View() string
        +GetActivePanel() model.PanelType
        +SetActivePanel(model.PanelType)
    }

    class PanelType {
        <<enumeration>>
        CurrentDirInput
        CurrentFilesList
        TargetDirInput
        TargetFilesList
        SelectedFilesList
        SearchResultsList
    }

    class AppMode {
        <<enumeration>>
        NormalMode
        SearchMode
        MoveMode
    }

    Model --> AppState : contains
    AppState --> PanelType : uses
    AppState --> AppMode : uses
```

## 3. イベント処理フロー

```mermaid
flowchart TD
    A[イベント受信] --> B{イベントタイプ}
    
    B -->|WindowSizeMsg| C[ウィンドウサイズ更新]
    B -->|KeyMsg| D[キー入力処理]
    B -->|PanelSwitchMsg| E[パネル切り替え]
    B -->|FileSelectedMsg| F[ファイル選択]
    B -->|DirectoryChangedMsg| G[ディレクトリ変更]
    B -->|SearchQueryChangedMsg| H[検索クエリ変更]
    
    C --> I[幅・高さ更新]
    I --> J[ヘッダー幅調整]
    
    D --> K{キーの種類}
    K -->|グローバルキー| L[グローバル処理]
    K -->|パネル固有キー| M[パネル固有処理]
    
    E --> N[アクティブパネル更新]
    N --> O[フォーカス設定]
    
    F --> P[選択済みファイル追加]
    G --> Q[ディレクトリパス更新]
    H --> R[検索クエリ更新]
    
    L --> S[状態更新完了]
    M --> S
    O --> S
    P --> S
    Q --> S
    R --> S
    J --> S
    
    S --> T[View関数呼び出し]
    T --> U[画面更新]
```

## 4. パネルナビゲーション状態遷移

```mermaid
stateDiagram-v2
    [*] --> CurrentDirInput : 初期状態
    
    CurrentDirInput --> CurrentFilesList : j (下)
    CurrentDirInput --> TargetDirInput : l (右)
    
    CurrentFilesList --> CurrentDirInput : k (上)
    CurrentFilesList --> TargetFilesList : l (右)
    
    TargetDirInput --> CurrentDirInput : h (左)
    TargetDirInput --> TargetFilesList : j (下)
    TargetDirInput --> SelectedFilesList : l (右)
    
    TargetFilesList --> CurrentFilesList : h (左)
    TargetFilesList --> TargetDirInput : k (上)
    
    SelectedFilesList --> TargetDirInput : h (左)
    SelectedFilesList --> SearchResultsList : l (右)
    
    SearchResultsList --> SelectedFilesList : h (左)
    
    CurrentDirInput --> TargetDirInput : Tab
    TargetDirInput --> CurrentFilesList : Tab
    CurrentFilesList --> TargetFilesList : Tab
    TargetFilesList --> SelectedFilesList : Tab
    SelectedFilesList --> SearchResultsList : Tab
    SearchResultsList --> CurrentDirInput : Tab
```

## 5. ファイル選択状態の変化

```mermaid
flowchart LR
    A[ファイルリスト表示] --> B[Spaceキー押下]
    B --> C{ファイル選択状態チェック}
    C -->|未選択| D[ファイル選択]
    C -->|選択済み| E[ファイル選択解除]
    
    D --> F[SelectedFiles配列に追加]
    F --> G[ファイルのIsSelected = true]
    G --> H[選択済みファイルリスト更新]
    
    E --> I[SelectedFiles配列から削除]
    I --> J[ファイルのIsSelected = false]
    J --> H
    
    H --> K[画面再描画]
```

## 6. 検索モードの状態遷移

```mermaid
stateDiagram-v2
    [*] --> NormalMode : アプリ起動
    
    NormalMode --> SearchMode : /キー押下
    SearchMode --> NormalMode : Escキー押下
    
    SearchMode --> Searching : 検索クエリ入力
    Searching --> SearchMode : 検索結果表示
    
    SearchMode --> NormalMode : 検索結果選択
    Searching --> NormalMode : 検索キャンセル
```

## 7. メッセージの流れ

```mermaid
sequenceDiagram
    participant User
    participant Model
    participant AppState
    participant UI
    
    User->>Model: キー入力
    Model->>Model: handleKeyPress
    Model->>Model: パネル判定
    
    alt パネル切り替え
        Model->>Model: moveToPanel
        Model->>Model: SetActivePanel
        Model->>AppState: ActivePanel更新
        Model->>UI: フォーカス設定
    else ファイル選択
        Model->>Model: handleFileSelectedMsg
        Model->>AppState: SelectedFiles更新
    else ディレクトリ変更
        Model->>Model: handleDirectoryChangedMsg
        Model->>AppState: CurrentDir/TargetDir更新
    end
    
    Model->>UI: View関数呼び出し
    UI->>User: 画面更新
```

## 8. データフローの概要

```mermaid
flowchart TB
    subgraph "入力層"
        A[キーボード入力]
        B[ウィンドウサイズ変更]
    end
    
    subgraph "処理層"
        C[イベントハンドラー]
        D[状態更新ロジック]
        E[パネル管理]
    end
    
    subgraph "状態層"
        F[AppState]
        G[Model.currentPanel]
        H[UIコンポーネント状態]
    end
    
    subgraph "表示層"
        I[View関数]
        J[レンダリング]
        K[画面出力]
    end
    
    A --> C
    B --> C
    C --> D
    C --> E
    D --> F
    E --> G
    F --> I
    G --> I
    H --> I
    I --> J
    J --> K
```

## 9. エラーハンドリングフロー

```mermaid
flowchart TD
    A[操作実行] --> B{エラー発生?}
    B -->|No| C[正常完了]
    B -->|Yes| D[エラーメッセージ作成]
    
    D --> E[ErrorMsg作成]
    E --> F[エラーハンドラー呼び出し]
    F --> G[エラー状態保存]
    G --> H[ユーザー通知]
    
    H --> I{リトライ可能?}
    I -->|Yes| J[リトライオプション表示]
    I -->|No| K[エラー状態継続]
    
    J --> L[ユーザー選択]
    L --> M{リトライ実行?}
    M -->|Yes| A
    M -->|No| K
    
    K --> N[エラー状態クリア]
    N --> O[通常状態復帰]
```

## 10. パフォーマンス最適化のポイント

```mermaid
flowchart LR
    A[イベント受信] --> B{更新が必要?}
    B -->|No| C[スキップ]
    B -->|Yes| D[状態更新]
    
    D --> E{View更新が必要?}
    E -->|No| F[完了]
    E -->|Yes| G[View関数実行]
    
    G --> H[差分レンダリング]
    H --> I[画面更新]
    
    C --> F
    F --> J[次のイベント待機]
    I --> J
```

これらの図により、Better MVの複雑な状態管理とイベント処理の流れが視覚的に理解できます。各コンポーネントがどのように連携し、状態がどのように変化するかが明確になります。
