package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tom555-555/better-mv/internal/model"
)

// Panel dimensions
const (
	panelWidth  = 25
	panelHeight = 12
)

// Panel 1: Current Directory Input
func (m *Model) currentDirInputView() string {
	isActive := m.state.ActivePanel == model.CurrentDirInput

	title := m.styles.GetPanelTitleStyle(isActive).Render("Current Dir")
	input := m.styles.GetInputStyle(isActive).
		Width(panelWidth - 4).
		Render(m.currentDirInput.field.Value())

	content := lipgloss.JoinVertical(lipgloss.Left, title, input)

	// Add empty space to maintain height
	for i := lipgloss.Height(content); i < panelHeight-2; i++ {
		content = lipgloss.JoinVertical(lipgloss.Left, content, "")
	}

	return m.styles.GetPanelStyle(isActive).
		Width(panelWidth).
		Height(panelHeight).
		Render(content)
}

// Panel 2: Current Files List
func (m *Model) currentFilesView() string {
	isActive := m.state.ActivePanel == model.CurrentFilesList

	title := m.styles.GetPanelTitleStyle(isActive).Render("Current Files")

	// Build file list content
	var fileContent strings.Builder
	for i, file := range m.state.CurrentFiles {
		cursor := " "
		if isActive && i == m.state.GetCurrentIndex(model.CurrentFilesList) {
			cursor = ">"
		}

		selected := " "
		if file.IsSelected {
			selected = "✓"
		}

		line := fmt.Sprintf("%s%s %s", cursor, selected, file.GetDisplayName())
		if i < len(m.state.CurrentFiles)-1 {
			line += "\n"
		}
		fileContent.WriteString(line)
	}

	listContent := fileContent.String()
	if len(m.state.CurrentFiles) == 0 {
		listContent = m.styles.DimText.Render("(empty)")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, listContent)

	// Add space bar hint
	if isActive {
		hint := m.styles.HelpDesc.Render("[␣ Select]")
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", hint)
	}

	return m.styles.GetPanelStyle(isActive).
		Width(panelWidth).
		Height(panelHeight).
		Render(content)
}

// Panel 3: Target Directory Input
func (m *Model) targetDirInputView() string {
	isActive := m.state.ActivePanel == model.TargetDirInput

	title := m.styles.GetPanelTitleStyle(isActive).Render("Target Dir")
	input := m.styles.GetInputStyle(isActive).
		Width(panelWidth - 4).
		Render(m.targetDirInput.field.Value())

	content := lipgloss.JoinVertical(lipgloss.Left, title, input)

	// Add empty space to maintain height
	for i := lipgloss.Height(content); i < panelHeight-2; i++ {
		content = lipgloss.JoinVertical(lipgloss.Left, content, "")
	}

	return m.styles.GetPanelStyle(isActive).
		Width(panelWidth).
		Height(panelHeight).
		Render(content)
}

// Panel 4: Target Files List
func (m *Model) targetFilesView() string {
	isActive := m.state.ActivePanel == model.TargetFilesList

	title := m.styles.GetPanelTitleStyle(isActive).Render("Target Files")

	// Build file list content
	var fileContent strings.Builder
	for i, file := range m.state.TargetFiles {
		cursor := " "
		if isActive && i == m.state.GetCurrentIndex(model.TargetFilesList) {
			cursor = ">"
		}

		line := fmt.Sprintf("%s %s", cursor, file.GetDisplayName())
		if i < len(m.state.TargetFiles)-1 {
			line += "\n"
		}
		fileContent.WriteString(line)
	}

	listContent := fileContent.String()
	if len(m.state.TargetFiles) == 0 {
		listContent = m.styles.DimText.Render("(empty)")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, listContent)

	return m.styles.GetPanelStyle(isActive).
		Width(panelWidth).
		Height(panelHeight).
		Render(content)
}

// Panel 5: Selected Files List
func (m *Model) selectedFilesView() string {
	isActive := m.state.ActivePanel == model.SelectedFilesList

	title := m.styles.GetPanelTitleStyle(isActive).Render("Selected Files")

	// Build selected files content
	var fileContent strings.Builder
	for i, file := range m.state.SelectedFiles {
		cursor := " "
		if isActive && i == m.state.GetCurrentIndex(model.SelectedFilesList) {
			cursor = ">"
		}

		line := fmt.Sprintf("%s✓ %s", cursor, file.GetDisplayName())
		if i < len(m.state.SelectedFiles)-1 {
			line += "\n"
		}
		fileContent.WriteString(line)
	}

	listContent := fileContent.String()
	if len(m.state.SelectedFiles) == 0 {
		listContent = m.styles.DimText.Render("(no files selected)")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, listContent)

	// Add action buttons
	if len(m.state.SelectedFiles) > 0 {
		moveButton := m.styles.Button.Render("[⌘↵ Move]")
		clearButton := m.styles.Button.Render("[⎋ Clear]")
		buttons := lipgloss.JoinVertical(lipgloss.Left, moveButton, clearButton)
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", buttons)
	}

	return m.styles.GetPanelStyle(isActive).
		Width(panelWidth).
		Height(panelHeight).
		Render(content)
}

// Panel 6: Fuzzy Search Results
func (m *Model) fuzzySearchView() string {
	isActive := m.state.ActivePanel == model.SearchResultsList
	isSearching := m.state.IsSearching

	title := m.styles.GetPanelTitleStyle(isActive).Render("Fuzzy Search")

	var content string

	if isSearching {
		// Show search input
		searchInput := m.styles.GetInputStyle(isActive).
			Width(panelWidth - 6).
			Render("search: " + m.searchInput.field.Value())

		// Show search results
		var resultsContent strings.Builder
		for i, result := range m.state.SearchResults {
			cursor := " "
			if isActive && i == m.state.GetCurrentIndex(model.SearchResultsList) {
				cursor = ">"
			}

			line := fmt.Sprintf("%s %s", cursor, result.GetDisplayPath())
			if i < len(m.state.SearchResults)-1 {
				line += "\n"
			}
			resultsContent.WriteString(line)
		}

		resultsList := resultsContent.String()
		if len(m.state.SearchResults) == 0 && m.state.SearchQuery != "" {
			resultsList = m.styles.DimText.Render("(no results)")
		} else if len(m.state.SearchResults) == 0 {
			resultsList = m.styles.DimText.Render("(type to search)")
		}

		content = lipgloss.JoinVertical(lipgloss.Left, title, searchInput, "", resultsList)
	} else {
		// Show instructions when not searching
		instructions := strings.Join([]string{
			"Press '/' to start",
			"fuzzy search for",
			"directories",
		}, "\n")
		content = lipgloss.JoinVertical(lipgloss.Left, title, "", m.styles.DimText.Render(instructions))
	}

	return m.styles.GetPanelStyle(isActive).
		Width(panelWidth).
		Height(panelHeight).
		Render(content)
}

// Header view
func (m *Model) headerView() string {
	return m.styles.Header.Render("better-mv (bmv)")
}

// Help view
func (m *Model) helpView() string {
	var helpItems []string

	switch m.state.ActivePanel {
	case model.CurrentDirInput, model.TargetDirInput:
		helpItems = []string{
			m.styles.HelpKey.Render("Enter") + " " + m.styles.HelpDesc.Render("Change directory"),
			m.styles.HelpKey.Render("hjkl/Tab") + " " + m.styles.HelpDesc.Render("Navigate panels"),
		}
	case model.CurrentFilesList:
		helpItems = []string{
			m.styles.HelpKey.Render("↑↓") + " " + m.styles.HelpDesc.Render("Navigate"),
			m.styles.HelpKey.Render("Space") + " " + m.styles.HelpDesc.Render("Select/Unselect"),
			m.styles.HelpKey.Render("Enter") + " " + m.styles.HelpDesc.Render("Open directory"),
		}
	case model.TargetFilesList:
		helpItems = []string{
			m.styles.HelpKey.Render("↑↓") + " " + m.styles.HelpDesc.Render("Navigate"),
		}
	case model.SelectedFilesList:
		helpItems = []string{
			m.styles.HelpKey.Render("↑↓") + " " + m.styles.HelpDesc.Render("Navigate"),
			m.styles.HelpKey.Render("⌘↵") + " " + m.styles.HelpDesc.Render("Move files"),
			m.styles.HelpKey.Render("Del") + " " + m.styles.HelpDesc.Render("Remove from selection"),
		}
	case model.SearchResultsList:
		if m.state.IsSearching {
			helpItems = []string{
				m.styles.HelpKey.Render("↑↓") + " " + m.styles.HelpDesc.Render("Navigate results"),
				m.styles.HelpKey.Render("Enter") + " " + m.styles.HelpDesc.Render("Select directory"),
				m.styles.HelpKey.Render("Esc") + " " + m.styles.HelpDesc.Render("Exit search"),
			}
		} else {
			helpItems = []string{
				m.styles.HelpKey.Render("/") + " " + m.styles.HelpDesc.Render("Start search"),
			}
		}
	}

	// Add global help
	helpItems = append(helpItems,
		m.styles.HelpKey.Render("Ctrl+C")+" "+m.styles.HelpDesc.Render("Quit"),
		m.styles.HelpKey.Render("q")+" "+m.styles.HelpDesc.Render("Quit"),
	)

	return m.styles.Help.Render(strings.Join(helpItems, " • "))
}
