package input

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Value() string
	Blur() tea.Msg
	View() string
	Update(tea.Msg) (Input, tea.Cmd)
}

func (sa *ShortAnswerField) Value() string {
	return sa.field.Value()
}

func (la *LongAnswerField) Value() string {
	return la.field.Value()
}

func (sa *ShortAnswerField) Blur() tea.Msg {
	return sa.field.Blur
}

func (la *LongAnswerField) Blur() tea.Msg {
	return la.field.Blur
}

func (sa *ShortAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	sa.field, cmd = sa.field.Update(msg)
	return sa, cmd
}

func (la *LongAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	la.field, cmd = la.field.Update(msg)
	return la, cmd
}

func (sa *ShortAnswerField) View() string {
	return sa.field.View()
}

func (la *LongAnswerField) View() string {
	return la.field.View()
}

type ShortAnswerField struct {
	field textinput.Model
}

type LongAnswerField struct {
	field textarea.Model
}

// textinput
func NewShortAnswerField() ShortAnswerField {
	ti := textinput.New()
	ti.Placeholder = "Your answer here"
	ti.Focus()
	return ShortAnswerField{ti}
}

// textarea
func NewLongAnswerField() LongAnswerField {
	ta := textarea.New()
	ta.Placeholder = "Your answer here"
	ta.Focus()
	return LongAnswerField{ta}
}
