package ui

import "github.com/tom555-555/better-mv/internal/model"

// Custom messages for the application

// PanelSwitchMsg represents a request to switch to a different panel
type PanelSwitchMsg struct {
	Panel model.PanelType
}

// FileSelectedMsg represents a file selection event
type FileSelectedMsg struct {
	File  *model.FileInfo
	Index int
}

// FileDeselectedMsg represents a file deselection event
type FileDeselectedMsg struct {
	File  *model.FileInfo
	Index int
}

// DirectoryChangedMsg represents a directory change event
type DirectoryChangedMsg struct {
	Path     string
	IsTarget bool // true for target directory, false for current directory
}

// SearchQueryChangedMsg represents a search query change
type SearchQueryChangedMsg struct {
	Query string
}

// SearchResultsMsg represents search results
type SearchResultsMsg struct {
	Results []model.DirectoryInfo
}

// FileListUpdatedMsg represents file list update event
type FileListUpdatedMsg struct {
	Files    []model.FileInfo
	IsTarget bool // true for target directory, false for current directory
}

// FileMoveRequestMsg represents a request to move files
type FileMoveRequestMsg struct {
	Files      []model.FileInfo
	TargetPath string
}

// FileMoveCompletedMsg represents completion of file move operation
type FileMoveCompletedMsg struct {
	Success bool
	Error   error
}

// ErrorMsg represents an error event
type ErrorMsg struct {
	Error error
}

// StatusMsg represents a status message
type StatusMsg struct {
	Message string
}

// InitialLoadMsg represents initial data loading completion
type InitialLoadMsg struct {
	CurrentFiles []model.FileInfo
	TargetFiles  []model.FileInfo
}
