package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tom555-555/better-mv/internal/model"
)

// handleKeyPress processes key presses and returns updated model and commands
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global key bindings (work in any panel)
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		// Clear search mode or selections
		if m.state.IsSearching {
			m.state.IsSearching = false
			m.searchInput.field.SetValue("")
			m.state.SearchQuery = ""
			m.state.SearchResults = []model.DirectoryInfo{}
		}
		return m, nil

	// Vim-style panel navigation
	case "h":
		return m.moveToPanel(m.getPanelToLeft()), nil
	case "j":
		return m.moveToPanel(m.getPanelBelow()), nil
	case "k":
		return m.moveToPanel(m.getPanelAbove()), nil
	case "l":
		return m.moveToPanel(m.getPanelToRight()), nil

	// Tab navigation
	case "tab":
		return m.moveToPanel(m.getNextPanel()), nil
	case "shift+tab":
		return m.moveToPanel(m.getPrevPanel()), nil

	// Search activation
	case "/":
		if !m.state.IsSearching {
			m.state.IsSearching = true
			m.SetActivePanel(model.SearchResultsList)
		}
		return m, nil
	}

	// Panel-specific key bindings
	return m.handlePanelSpecificKeys(msg)
}

// Panel navigation logic

// getPanelToLeft returns the panel to the left of current active panel
func (m *Model) getPanelToLeft() model.PanelType {
	switch m.state.ActivePanel {
	case model.TargetDirInput:
		return m.targetDirInput.surroundings.Left
	case model.TargetFilesList:
		return m.targetFilesList.surroundings.Left
	case model.SelectedFilesList:
		return m.selectedFilesList.surroundings.Left
	case model.SearchResultsList:
		return m.searchResultsList.surroundings.Left
	default:
		return m.state.ActivePanel // No movement
	}
}

// getPanelToRight returns the panel to the right of current active panel
func (m *Model) getPanelToRight() model.PanelType {
	switch m.state.ActivePanel {
	case model.CurrentDirInput:
		return m.currentDirInput.surroundings.Right
	case model.CurrentFilesList:
		return m.currentFilesList.surroundings.Right
	case model.TargetDirInput:
		return m.targetDirInput.surroundings.Right
	case model.SelectedFilesList:
		return m.searchResultsList.surroundings.Right
	default:
		return m.state.ActivePanel // No movement
	}
}

// getPanelAbove returns the panel above current active panel
func (m *Model) getPanelAbove() model.PanelType {
	switch m.state.ActivePanel {
	case model.CurrentFilesList:
		return m.currentDirInput.surroundings.Top
	case model.TargetFilesList:
		return m.targetDirInput.surroundings.Top
	default:
		return m.state.ActivePanel // No movement
	}
}

// getPanelBelow returns the panel below current active panel
func (m *Model) getPanelBelow() model.PanelType {
	switch m.state.ActivePanel {
	case model.CurrentDirInput:
		return m.currentFilesList.surroundings.Bottom
	case model.TargetDirInput:
		return m.targetFilesList.surroundings.Bottom
	default:
		return m.state.ActivePanel // No movement
	}
}

// getNextPanel returns the next panel in tab order
func (m *Model) getNextPanel() model.PanelType {
	panels := []model.PanelType{
		model.CurrentDirInput,
		model.CurrentFilesList,
		model.TargetDirInput,
		model.TargetFilesList,
		model.SelectedFilesList,
		model.SearchResultsList,
	}

	for i, panel := range panels {
		if panel == m.state.ActivePanel {
			return panels[(i+1)%len(panels)]
		}
	}
	return model.CurrentDirInput
}

// getPrevPanel returns the previous panel in tab order
func (m *Model) getPrevPanel() model.PanelType {
	panels := []model.PanelType{
		model.CurrentDirInput,
		model.CurrentFilesList,
		model.TargetDirInput,
		model.TargetFilesList,
		model.SelectedFilesList,
		model.SearchResultsList,
	}

	for i, panel := range panels {
		if panel == m.state.ActivePanel {
			prevIndex := (i - 1 + len(panels)) % len(panels)
			return panels[prevIndex]
		}
	}
	return model.SearchResultsList
}

// moveToPanel switches to the specified panel
func (m Model) moveToPanel(panel model.PanelType) tea.Model {
	if panel == m.state.ActivePanel {
		return m // No change if same panel
	}

	m.SetActivePanel(panel)
	return m
}

// handlePanelSpecificKeys processes keys specific to each panel
func (m Model) handlePanelSpecificKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state.ActivePanel {
	case model.CurrentDirInput, model.TargetDirInput:
		return m.handleDirectoryInputKeys(msg)
	case model.CurrentFilesList:
		return m.handleCurrentFilesKeys(msg)
	case model.TargetFilesList:
		return m.handleTargetFilesKeys(msg)
	case model.SelectedFilesList:
		return m.handleSelectedFilesKeys(msg)
	case model.SearchResultsList:
		return m.handleSearchResultsKeys(msg)
	}

	return m, nil
}

// handleDirectoryInputKeys processes keys for directory input panels
func (m Model) handleDirectoryInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// TODO: Implement directory change logic
		// For now, just update the display
		var path string
		if m.state.ActivePanel == model.CurrentDirInput {
			path = m.currentDirInput.field.Value()
			m.state.CurrentDir = path
		} else {
			path = m.targetDirInput.field.Value()
			m.state.TargetDir = path
		}
		return m, nil
	}

	return m, nil
}

// handleCurrentFilesKeys processes keys for current files list panel
func (m Model) handleCurrentFilesKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		index := m.state.GetCurrentIndex(model.CurrentFilesList)
		if index > 0 {
			m.state.SetCurrentIndex(model.CurrentFilesList, index-1)
		}
	case "down":
		index := m.state.GetCurrentIndex(model.CurrentFilesList)
		maxIndex := len(m.state.CurrentFiles) - 1
		if index < maxIndex {
			m.state.SetCurrentIndex(model.CurrentFilesList, index+1)
		}
	case " ", "space":
		// TODO: Implement file selection toggle
		return m, nil
	case "enter":
		// TODO: Implement directory navigation
		return m, nil
	}

	return m, nil
}

// handleTargetFilesKeys processes keys for target files list panel
func (m Model) handleTargetFilesKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		index := m.state.GetCurrentIndex(model.TargetFilesList)
		if index > 0 {
			m.state.SetCurrentIndex(model.TargetFilesList, index-1)
		}
	case "down":
		index := m.state.GetCurrentIndex(model.TargetFilesList)
		maxIndex := len(m.state.TargetFiles) - 1
		if index < maxIndex {
			m.state.SetCurrentIndex(model.TargetFilesList, index+1)
		}
	}

	return m, nil
}

// handleSelectedFilesKeys processes keys for selected files list panel
func (m Model) handleSelectedFilesKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		index := m.state.GetCurrentIndex(model.SelectedFilesList)
		if index > 0 {
			m.state.SetCurrentIndex(model.SelectedFilesList, index-1)
		}
	case "down":
		index := m.state.GetCurrentIndex(model.SelectedFilesList)
		maxIndex := len(m.state.SelectedFiles) - 1
		if index < maxIndex {
			m.state.SetCurrentIndex(model.SelectedFilesList, index+1)
		}
	case "delete", "backspace":
		// TODO: Implement file removal from selection
		return m, nil
	case "cmd+enter", "ctrl+enter":
		// TODO: Implement file move operation
		return m, nil
	}

	return m, nil
}

// handleSearchResultsKeys processes keys for search results panel
func (m Model) handleSearchResultsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.state.IsSearching {
		return m, nil
	}

	switch msg.String() {
	case "up":
		index := m.state.GetCurrentIndex(model.SearchResultsList)
		if index > 0 {
			m.state.SetCurrentIndex(model.SearchResultsList, index-1)
		}
	case "down":
		index := m.state.GetCurrentIndex(model.SearchResultsList)
		maxIndex := len(m.state.SearchResults) - 1
		if index < maxIndex {
			m.state.SetCurrentIndex(model.SearchResultsList, index+1)
		}
	case "enter":
		// TODO: Implement directory selection from search results
		return m, nil
	}

	return m, nil
}
