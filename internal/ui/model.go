package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tom555-555/better-mv/internal/model"
)

// Model represents the main UI model for the application
type Model struct {
	// Application state
	state *model.AppState

	// UI components for inputs
	currentDirInput struct {
		field        textinput.Model
		surroundings Surroundings
	}
	targetDirInput struct {
		field        textinput.Model
		surroundings Surroundings
	}
	searchInput struct {
		field        textinput.Model
		surroundings Surroundings
	}

	// UI components for lists
	currentFilesList struct {
		field        list.Model
		surroundings Surroundings
	}
	targetFilesList struct {
		field        list.Model
		surroundings Surroundings
	}
	selectedFilesList struct {
		field        list.Model
		surroundings Surroundings
	}
	searchResultsList struct {
		field        list.Model
		surroundings Surroundings
	}

	// UI state
	width  int
	height int

	// Styles
	styles *Styles
}

type Surroundings struct {
	Top    model.PanelType
	Left   model.PanelType
	Right  model.PanelType
	Bottom model.PanelType
}

type ShortTextInput struct {
	field        textinput.Model
	surroundings Surroundings
}

type ListElement struct {
	field        list.Model
	surroundings Surroundings
}

func NewShortTextInput(placeholder string, initVal string, width int, surroundings Surroundings) ShortTextInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Width = width
	ti.SetValue(initVal)
	return ShortTextInput{ti, surroundings}
}

func NewList(
	items []list.Item,
	delegate list.DefaultDelegate,
	title string,
	width int,
	height int,
	showStatusBar bool,
	isFiltered bool,
	surroundings Surroundings,
) ListElement {
	listEl := list.New(items, delegate, width, height)
	listEl.Title = title
	listEl.SetShowStatusBar(showStatusBar)
	listEl.SetFilteringEnabled(isFiltered)
	return ListElement{listEl, surroundings}
}

// NewModel creates a new UI model
func NewModel() *Model {
	state := model.NewAppState()

	// Initialize text inputs
	currentDirInput := NewShortTextInput("Enter current directory path...", ".", 30, Surroundings{
		Top:    model.CurrentDirInput,
		Left:   model.CurrentDirInput,
		Right:  model.TargetDirInput,
		Bottom: model.CurrentFilesList,
	})
	targetDirInput := NewShortTextInput("Enter target directory path...", ".", 30, Surroundings{
		Top:    model.TargetDirInput,
		Left:   model.CurrentDirInput,
		Right:  model.SelectedFilesList,
		Bottom: model.TargetFilesList,
	})
	searchInput := NewShortTextInput("Search directories...", "", 30, Surroundings{
		Top:    model.SearchResultsList,
		Left:   model.SelectedFilesList,
		Right:  model.SearchResultsList,
		Bottom: model.SearchResultsList,
	})

	// Initialize lists
	currentFilesList := NewList([]list.Item{}, newFileItemDelegate(), "Current Files", 30, 10, false, false, Surroundings{
		Top:    model.CurrentDirInput,
		Left:   model.CurrentFilesList,
		Right:  model.TargetFilesList,
		Bottom: model.CurrentFilesList,
	})
	targetFilesList := NewList([]list.Item{}, newFileItemDelegate(), "Target Files", 30, 10, false, false, Surroundings{
		Top:    model.TargetDirInput,
		Left:   model.CurrentFilesList,
		Right:  model.TargetFilesList,
		Bottom: model.TargetFilesList,
	})
	selectedFilesList := NewList([]list.Item{}, newFileItemDelegate(), "Selected Files", 30, 10, false, false, Surroundings{
		Top:    model.SelectedFilesList,
		Left:   model.TargetDirInput,
		Right:  model.SearchResultsList,
		Bottom: model.SelectedFilesList,
	})
	searchResultsList := NewList([]list.Item{}, newDirectoryItemDelegate(), "Search Results", 30, 10, false, false, Surroundings{
		Top:    model.SearchResultsList,
		Left:   model.SelectedFilesList,
		Right:  model.SearchResultsList,
		Bottom: model.SearchResultsList,
	})

	// Focus the first panel
	currentDirInput.field.Focus()

	return &Model{
		state:             state,
		currentDirInput:   currentDirInput,
		targetDirInput:    targetDirInput,
		searchInput:       searchInput,
		currentFilesList:  currentFilesList,
		targetFilesList:   targetFilesList,
		selectedFilesList: selectedFilesList,
		searchResultsList: searchResultsList,
		styles:            DefaultStyles(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// GetActivePanel returns the currently active panel
func (m *Model) GetActivePanel() model.PanelType {
	return m.state.ActivePanel
}

func (m *Model) BlurPanel(panel model.PanelType) {
	switch panel {
	case model.CurrentDirInput:
		m.currentDirInput.field.Blur()
	case model.TargetDirInput:
		m.targetDirInput.field.Blur()
	case model.SearchResultsList:
		m.searchInput.field.Blur()
	}
}

// SetActivePanel sets the active panel and updates focus
func (m *Model) SetActivePanel(panel model.PanelType) {
	// Blur all inputs first
	m.BlurPanel(m.state.ActivePanel)

	// Update state
	m.state.ActivePanel = panel

	// Focus the appropriate input
	switch panel {
	case model.CurrentDirInput:
		m.currentDirInput.field.Focus()
	case model.TargetDirInput:
		m.targetDirInput.field.Focus()
	case model.SearchResultsList:
		if m.state.IsSearching {
			m.searchInput.field.Focus()
		}
	}
}

// Item delegates for lists

// fileItem represents a file item in a list
type fileItem struct {
	info *model.FileInfo
}

func (f fileItem) FilterValue() string { return f.info.Name }
func (f fileItem) Title() string       { return f.info.GetDisplayName() }
func (f fileItem) Description() string { return f.info.GetSizeString() }

// directoryItem represents a directory item in search results
type directoryItem struct {
	info *model.DirectoryInfo
}

func (d directoryItem) FilterValue() string { return d.info.Name }
func (d directoryItem) Title() string       { return d.info.GetDisplayName() }
func (d directoryItem) Description() string { return d.info.GetDisplayPath() }

// newFileItemDelegate creates a new delegate for file items
func newFileItemDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	return delegate
}

// newDirectoryItemDelegate creates a new delegate for directory items
func newDirectoryItemDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	return delegate
}
