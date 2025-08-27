package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles contains all the styles for the UI
type Styles struct {
	// Colors
	PrimaryColor      lipgloss.Color
	SecondaryColor    lipgloss.Color
	AccentColor       lipgloss.Color
	TextColor         lipgloss.Color
	DimTextColor      lipgloss.Color
	BorderColor       lipgloss.Color
	ActiveBorderColor lipgloss.Color

	// Base styles
	Base             lipgloss.Style
	Header           lipgloss.Style
	Panel            lipgloss.Style
	ActivePanel      lipgloss.Style
	PanelTitle       lipgloss.Style
	ActivePanelTitle lipgloss.Style

	// Input styles
	InputField       lipgloss.Style
	ActiveInputField lipgloss.Style

	// List styles
	List         lipgloss.Style
	ActiveList   lipgloss.Style
	SelectedItem lipgloss.Style

	// Button styles
	Button       lipgloss.Style
	ActiveButton lipgloss.Style

	// Help styles
	Help     lipgloss.Style
	HelpKey  lipgloss.Style
	HelpDesc lipgloss.Style

	// Text styles
	DimText lipgloss.Style
}

// DefaultStyles returns the default style configuration
func DefaultStyles() *Styles {
	s := &Styles{}

	// Colors
	s.PrimaryColor = lipgloss.Color("#7D56F4")
	s.SecondaryColor = lipgloss.Color("#9B59B6")
	s.AccentColor = lipgloss.Color("#F39C12")
	s.TextColor = lipgloss.Color("#FAFAFA")
	s.DimTextColor = lipgloss.Color("#626262")
	s.BorderColor = lipgloss.Color("#383838")
	s.ActiveBorderColor = s.PrimaryColor

	// Base styles
	s.Base = lipgloss.NewStyle().
		Background(lipgloss.Color("#1A1A1A")).
		Foreground(s.TextColor)

	s.Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(s.TextColor).
		Background(s.PrimaryColor).
		Padding(0, 1).
		Width(80).
		Align(lipgloss.Center)

	// Panel styles
	s.Panel = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(s.BorderColor).
		Padding(0, 1).
		MarginRight(1)

	s.ActivePanel = s.Panel.Copy().
		BorderForeground(s.ActiveBorderColor).
		BorderStyle(lipgloss.ThickBorder())

	s.PanelTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(s.DimTextColor).
		Padding(0, 1)

	s.ActivePanelTitle = s.PanelTitle.Copy().
		Foreground(s.AccentColor)

	// Input styles
	s.InputField = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.BorderColor).
		Padding(0, 1).
		Foreground(s.TextColor)

	s.ActiveInputField = s.InputField.Copy().
		BorderForeground(s.ActiveBorderColor)

	// List styles
	s.List = lipgloss.NewStyle().
		Foreground(s.TextColor)

	s.ActiveList = s.List.Copy().
		BorderForeground(s.ActiveBorderColor)

	s.SelectedItem = lipgloss.NewStyle().
		Background(s.PrimaryColor).
		Foreground(s.TextColor).
		Bold(true)

	// Button styles
	s.Button = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.BorderColor).
		Padding(0, 1).
		Foreground(s.TextColor)

	s.ActiveButton = s.Button.Copy().
		Background(s.PrimaryColor).
		BorderForeground(s.ActiveBorderColor)

	// Help styles
	s.Help = lipgloss.NewStyle().
		Foreground(s.DimTextColor).
		MarginTop(1).
		Padding(0, 1)

	s.HelpKey = lipgloss.NewStyle().
		Foreground(s.AccentColor).
		Bold(true)

	s.HelpDesc = lipgloss.NewStyle().
		Foreground(s.DimTextColor)

	s.DimText = lipgloss.NewStyle().
		Foreground(s.DimTextColor)

	return s
}

// GetPanelStyle returns the appropriate panel style based on active state
func (s *Styles) GetPanelStyle(isActive bool) lipgloss.Style {
	if isActive {
		return s.ActivePanel
	}
	return s.Panel
}

// GetPanelTitleStyle returns the appropriate panel title style based on active state
func (s *Styles) GetPanelTitleStyle(isActive bool) lipgloss.Style {
	if isActive {
		return s.ActivePanelTitle
	}
	return s.PanelTitle
}

// GetInputStyle returns the appropriate input style based on active state
func (s *Styles) GetInputStyle(isActive bool) lipgloss.Style {
	if isActive {
		return s.ActiveInputField
	}
	return s.InputField
}

// Panel dimensions
const (
	PanelMinWidth  = 20
	PanelMinHeight = 5
	PanelMaxHeight = 15
)
