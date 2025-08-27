package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tom555-555/better-mv/internal/model"
)

// View renders the main application UI with 6 panels
func (m Model) View() string {
	// Handle window size updates
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	header := m.headerView()

	// Top row: Panels 1, 3, 5, 6
	topRow := lipgloss.JoinHorizontal(lipgloss.Top,
		m.currentDirInputView(), // Panel 1: Current Directory Input
		m.targetDirInputView(),  // Panel 3: Target Directory Input
		m.selectedFilesView(),   // Panel 5: Selected Files List
		m.fuzzySearchView(),     // Panel 6: Fuzzy Search Results
	)

	// Bottom row: Panels 2, 4, and empty spaces
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top,
		m.currentFilesView(), // Panel 2: Current Files List
		m.targetFilesView(),  // Panel 4: Target Files List
		createEmptyPanel(),   // Empty space to align with Panel 5
		createEmptyPanel(),   // Empty space to align with Panel 6
	)

	// Main content area
	content := lipgloss.JoinVertical(lipgloss.Left,
		topRow,
		bottomRow,
	)

	// Footer with help
	help := m.helpView()

	// Combine all elements
	return m.styles.Base.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			"",
			content,
			"",
			help,
		),
	)
}

// Update handles all message types and updates the model accordingly
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case PanelSwitchMsg:
		return m.handlePanelSwitchMsg(msg)

	case FileSelectedMsg:
		return m.handleFileSelectedMsg(msg)

	case FileDeselectedMsg:
		return m.handleFileDeselectedMsg(msg)

	case DirectoryChangedMsg:
		return m.handleDirectoryChangedMsg(msg)

	case SearchQueryChangedMsg:
		return m.handleSearchQueryChangedMsg(msg)

	case SearchResultsMsg:
		return m.handleSearchResultsMsg(msg)

	case FileListUpdatedMsg:
		return m.handleFileListUpdatedMsg(msg)

	case FileMoveRequestMsg:
		return m.handleFileMoveRequestMsg(msg)

	case FileMoveCompletedMsg:
		return m.handleFileMoveCompletedMsg(msg)

	case ErrorMsg:
		return m.handleErrorMsg(msg)

	case StatusMsg:
		return m.handleStatusMsg(msg)

	case InitialLoadMsg:
		return m.handleInitialLoadMsg(msg)
	}

	// Update active input field with remaining messages
	return m.updateInputFields(msg)
}

// Message handlers

func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	// Update header width
	m.styles.Header = m.styles.Header.Width(msg.Width - 4)

	return m, nil
}

func (m Model) handlePanelSwitchMsg(msg PanelSwitchMsg) (tea.Model, tea.Cmd) {
	m.SetActivePanel(msg.Panel)
	return m, nil
}

func (m Model) handleFileSelectedMsg(msg FileSelectedMsg) (tea.Model, tea.Cmd) {
	// Add file to selected files if not already selected
	for _, selected := range m.state.SelectedFiles {
		if selected.AbsPath == msg.File.AbsPath {
			return m, nil // Already selected
		}
	}

	// Mark file as selected
	msg.File.IsSelected = true
	m.state.SelectedFiles = append(m.state.SelectedFiles, *msg.File)

	return m, nil
}

func (m Model) handleFileDeselectedMsg(msg FileDeselectedMsg) (tea.Model, tea.Cmd) {
	// Remove file from selected files
	var newSelectedFiles []model.FileInfo
	for _, selected := range m.state.SelectedFiles {
		if selected.AbsPath != msg.File.AbsPath {
			newSelectedFiles = append(newSelectedFiles, selected)
		}
	}

	// Mark file as not selected
	msg.File.IsSelected = false
	m.state.SelectedFiles = newSelectedFiles

	return m, nil
}

func (m Model) handleDirectoryChangedMsg(msg DirectoryChangedMsg) (tea.Model, tea.Cmd) {
	if msg.IsTarget {
		m.state.TargetDir = msg.Path
		m.targetDirInput.field.SetValue(msg.Path)
	} else {
		m.state.CurrentDir = msg.Path
		m.currentDirInput.field.SetValue(msg.Path)
	}

	// TODO: Trigger directory scan
	return m, nil
}

func (m Model) handleSearchQueryChangedMsg(msg SearchQueryChangedMsg) (tea.Model, tea.Cmd) {
	m.state.SearchQuery = msg.Query
	m.searchInput.field.SetValue(msg.Query)

	// TODO: Trigger fuzzy search
	return m, nil
}

func (m Model) handleSearchResultsMsg(msg SearchResultsMsg) (tea.Model, tea.Cmd) {
	m.state.SearchResults = msg.Results
	return m, nil
}

func (m Model) handleFileListUpdatedMsg(msg FileListUpdatedMsg) (tea.Model, tea.Cmd) {
	if msg.IsTarget {
		m.state.TargetFiles = msg.Files
	} else {
		m.state.CurrentFiles = msg.Files
	}
	return m, nil
}

func (m Model) handleFileMoveRequestMsg(msg FileMoveRequestMsg) (tea.Model, tea.Cmd) {
	// TODO: Implement file move operation
	return m, nil
}

func (m Model) handleFileMoveCompletedMsg(msg FileMoveCompletedMsg) (tea.Model, tea.Cmd) {
	if msg.Success {
		// Clear selected files on successful move
		m.state.SelectedFiles = []model.FileInfo{}
		// TODO: Refresh file lists
	}
	return m, nil
}

func (m Model) handleErrorMsg(msg ErrorMsg) (tea.Model, tea.Cmd) {
	// TODO: Show error in status or notification area
	return m, nil
}

func (m Model) handleStatusMsg(msg StatusMsg) (tea.Model, tea.Cmd) {
	// TODO: Show status message
	return m, nil
}

func (m Model) handleInitialLoadMsg(msg InitialLoadMsg) (tea.Model, tea.Cmd) {
	m.state.CurrentFiles = msg.CurrentFiles
	m.state.TargetFiles = msg.TargetFiles
	return m, nil
}

// updateInputFields updates the active input field with the message
func (m Model) updateInputFields(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state.ActivePanel {
	case model.CurrentDirInput:
		m.currentDirInput.field, cmd = m.currentDirInput.field.Update(msg)
	case model.TargetDirInput:
		m.targetDirInput.field, cmd = m.targetDirInput.field.Update(msg)
	case model.SearchResultsList:
		if m.state.IsSearching {
			m.searchInput.field, cmd = m.searchInput.field.Update(msg)
		}
	}

	return m, cmd
}

// createEmptyPanel creates an empty panel for layout alignment
func createEmptyPanel() string {
	return lipgloss.NewStyle().
		Width(panelWidth).
		Height(panelHeight).
		Render("")
}
