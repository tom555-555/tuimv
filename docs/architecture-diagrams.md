# tuimv - Architecture Diagrams

> [日本語版はこちら](architecture-diagrams_ja.md) / [Japanese version is here](architecture-diagrams_ja.md)

This document explains tuimv's operation flow and state changes using diagrams.

## 1. Application Startup Flow

```mermaid
flowchart TD
    A[main function start] --> B[tea.LogToFile setup]
    B --> C[ui.NewModel call]
    C --> D[model.NewAppState creation]
    D --> E[UI component initialization]
    E --> F[tea.NewProgram creation]
    F --> G[Program execution start]
    G --> H[Model.Init call]
    H --> I[Main loop start]
    I --> J[Event waiting]
    J --> K{Event occurred}
    K --> L[Update function call]
    L --> M[View function call]
    M --> J
```

## 2. State Management Structure

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

## 3. Event Processing Flow

```mermaid
flowchart TD
    A[Event received] --> B{Event type}
    
    B -->|WindowSizeMsg| C[Window size update]
    B -->|KeyMsg| D[Key input processing]
    B -->|PanelSwitchMsg| E[Panel switching]
    B -->|FileSelectedMsg| F[File selection]
    B -->|DirectoryChangedMsg| G[Directory change]
    B -->|SearchQueryChangedMsg| H[Search query change]
    
    C --> I[Width/height update]
    I --> J[Header width adjustment]
    
    D --> K{Key type}
    K -->|Global keys| L[Global processing]
    K -->|Panel-specific keys| M[Panel-specific processing]
    
    E --> N[Active panel update]
    N --> O[Focus setting]
    
    F --> P[Selected files addition]
    G --> Q[Directory path update]
    H --> R[Search query update]
    
    L --> S[State update complete]
    M --> S
    O --> S
    P --> S
    Q --> S
    R --> S
    J --> S
    
    S --> T[View function call]
    T --> U[Screen update]
```

## 4. Panel Navigation State Transitions

```mermaid
stateDiagram-v2
    [*] --> CurrentDirInput : Initial state
    
    CurrentDirInput --> CurrentFilesList : j (down)
    CurrentDirInput --> TargetDirInput : l (right)
    
    CurrentFilesList --> CurrentDirInput : k (up)
    CurrentFilesList --> TargetFilesList : l (right)
    
    TargetDirInput --> CurrentDirInput : h (left)
    TargetDirInput --> TargetFilesList : j (down)
    TargetDirInput --> SelectedFilesList : l (right)
    
    TargetFilesList --> CurrentFilesList : h (left)
    TargetFilesList --> TargetDirInput : k (up)
    
    SelectedFilesList --> TargetDirInput : h (left)
    SelectedFilesList --> SearchResultsList : l (right)
    
    SearchResultsList --> SelectedFilesList : h (left)
    
    CurrentDirInput --> TargetDirInput : Tab
    TargetDirInput --> CurrentFilesList : Tab
    CurrentFilesList --> TargetFilesList : Tab
    TargetFilesList --> SelectedFilesList : Tab
    SelectedFilesList --> SearchResultsList : Tab
    SearchResultsList --> CurrentDirInput : Tab
```

## 5. File Selection State Changes

```mermaid
flowchart LR
    A[File list display] --> B[Space key press]
    B --> C{File selection state check}
    C -->|Unselected| D[File selection]
    C -->|Selected| E[File deselection]
    
    D --> F[Add to SelectedFiles array]
    F --> G[File IsSelected = true]
    G --> H[Selected files list update]
    
    E --> I[Remove from SelectedFiles array]
    I --> J[File IsSelected = false]
    J --> H
    
    H --> K[Screen redraw]
```

## 6. Search Mode State Transitions

```mermaid
stateDiagram-v2
    [*] --> NormalMode : Application startup
    
    NormalMode --> SearchMode : / key press
    SearchMode --> NormalMode : Esc key press
    
    SearchMode --> Searching : Search query input
    Searching --> SearchMode : Search results display
    
    SearchMode --> NormalMode : Search result selection
    Searching --> NormalMode : Search cancel
```

## 7. Message Flow

```mermaid
sequenceDiagram
    participant User
    participant Model
    participant AppState
    participant UI
    
    User->>Model: Key input
    Model->>Model: handleKeyPress
    Model->>Model: Panel determination
    
    alt Panel switching
        Model->>Model: moveToPanel
        Model->>Model: SetActivePanel
        Model->>AppState: ActivePanel update
        Model->>UI: Focus setting
    else File selection
        Model->>Model: handleFileSelectedMsg
        Model->>AppState: SelectedFiles update
    else Directory change
        Model->>Model: handleDirectoryChangedMsg
        Model->>AppState: CurrentDir/TargetDir update
    end
    
    Model->>UI: View function call
    UI->>User: Screen update
```

## 8. Data Flow Overview

```mermaid
flowchart TB
    subgraph "Input Layer"
        A[Keyboard input]
        B[Window size change]
    end
    
    subgraph "Processing Layer"
        C[Event handlers]
        D[State update logic]
        E[Panel management]
    end
    
    subgraph "State Layer"
        F[AppState]
        G[Model.currentPanel]
        H[UI component states]
    end
    
    subgraph "Display Layer"
        I[View function]
        J[Rendering]
        K[Screen output]
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

## 9. Error Handling Flow

```mermaid
flowchart TD
    A[Operation execution] --> B{Error occurred?}
    B -->|No| C[Success completion]
    B -->|Yes| D[Error message creation]
    
    D --> E[ErrorMsg creation]
    E --> F[Error handler call]
    F --> G[Error state saving]
    G --> H[User notification]
    
    H --> I{Retry possible?}
    I -->|Yes| J[Retry option display]
    I -->|No| K[Error state continuation]
    
    J --> L[User selection]
    L --> M{Execute retry?}
    M -->|Yes| A
    M -->|No| K
    
    K --> N[Error state clear]
    N --> O[Normal state recovery]
```

## 10. Performance Optimization Points

```mermaid
flowchart LR
    A[Event received] --> B{Update needed?}
    B -->|No| C[Skip]
    B -->|Yes| D[State update]
    
    D --> E{View update needed?}
    E -->|No| F[Complete]
    E -->|Yes| G[View function execution]
    
    G --> H[Differential rendering]
    H --> I[Screen update]
    
    C --> F
    F --> J[Next event waiting]
    I --> J
```

These diagrams provide a visual understanding of tuimv's complex state management and event processing flow. They clearly show how each component collaborates and how states change.
