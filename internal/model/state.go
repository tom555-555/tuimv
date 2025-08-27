package model

// PanelType represents the different panels in the UI
type PanelType int

const (
	CurrentDirInput   PanelType = iota // Panel 1: Current directory input
	CurrentFilesList                   // Panel 2: Current directory files list
	TargetDirInput                     // Panel 3: Target directory input
	TargetFilesList                    // Panel 4: Target directory files list
	SelectedFilesList                  // Panel 5: Selected files list
	SearchResultsList                  // Panel 6: Fuzzy search results list
)

// AppMode represents the current application mode
type AppMode int

const (
	NormalMode AppMode = iota // Normal file browsing mode
	SearchMode                // Directory fuzzy search mode
	MoveMode                  // File moving mode
)

// AppState represents the overall application state
type AppState struct {
	CurrentDir    string            // Current directory path
	TargetDir     string            // Target directory path
	CurrentFiles  []FileInfo        // Files in current directory
	TargetFiles   []FileInfo        // Files in target directory
	SelectedFiles []FileInfo        // Files selected for moving
	SearchResults []DirectoryInfo   // Fuzzy search results
	SearchQuery   string            // Current search query
	ActivePanel   PanelType         // Currently active panel
	CurrentIndex  map[PanelType]int // Current cursor position for each panel
	Mode          AppMode           // Current application mode
	IsSearching   bool              // Whether fuzzy search is active
}

// NewAppState creates a new application state with default values
func NewAppState() *AppState {
	return &AppState{
		CurrentDir:    ".",
		TargetDir:     ".",
		CurrentFiles:  make([]FileInfo, 0),
		TargetFiles:   make([]FileInfo, 0),
		SelectedFiles: make([]FileInfo, 0),
		SearchResults: make([]DirectoryInfo, 0),
		SearchQuery:   "",
		ActivePanel:   CurrentDirInput,
		CurrentIndex:  make(map[PanelType]int),
		Mode:          NormalMode,
		IsSearching:   false,
	}
}

// GetCurrentIndex returns the current cursor position for a panel
func (s *AppState) GetCurrentIndex(panel PanelType) int {
	if index, exists := s.CurrentIndex[panel]; exists {
		return index
	}
	return 0
}

// SetCurrentIndex sets the cursor position for a panel
func (s *AppState) SetCurrentIndex(panel PanelType, index int) {
	s.CurrentIndex[panel] = index
}

// GetFileCount returns the number of items in a panel
func (s *AppState) GetFileCount(panel PanelType) int {
	switch panel {
	case CurrentFilesList:
		return len(s.CurrentFiles)
	case TargetFilesList:
		return len(s.TargetFiles)
	case SelectedFilesList:
		return len(s.SelectedFiles)
	case SearchResultsList:
		return len(s.SearchResults)
	default:
		return 0
	}
}
