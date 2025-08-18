package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tom555-555/better-mv/internal/ui/input"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(80)
	return s
}

type Question struct {
	question string
	answer   string
	input    input.Input
}

func NewQuestion(question string) Question {
	return Question{question: question}
}

func newShortQuestion(question string) Question {
	q := NewQuestion(question)
	field := input.NewShortAnswerField()
	q.input = &field
	return q
}

func newLongQuestion(question string) Question {
	q := NewQuestion(question)
	field := input.NewLongAnswerField()
	q.input = &field
	return q
}

type Model struct {
	questions []Question
	done      bool
	height    int
	width     int
	index     int
	styles    *Styles
}

func New(questions []Question) *Model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Width = 80
	answerField.Placeholder = "Type your answer here..."
	answerField.Focus()
	return &Model{
		questions: questions,
		styles:    styles,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		case "enter":
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			current.answer = current.input.Value()
			log.Printf("question %s, answer: %s", current.question, current.answer)
			m.Next()
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	current := &m.questions[m.index]
	if m.done {
		var output string
		for _, q := range m.questions {
			output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
		}
		return output
	}
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			current.question,
			m.styles.InputField.Render(current.input.View()),
		),
	)
}

func (m *Model) Next() {
	log.Printf("next: %d", m.index)
	log.Printf("len: %d", len(m.questions)-1)
	log.Printf("if smaller than len: %t", m.index < len(m.questions)-1)
	if m.index < len(m.questions)-1 {
		log.Printf("incrementing index")
		m.index++
	} else {
		log.Print("resetting index")
		m.index = 0
	}
}

func main() {
	questions := []Question{
		newShortQuestion("what is your name?"),
		newShortQuestion("what is your age?"),
		newLongQuestion("what is your favorite quote?"),
	}
	m := New(questions)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("Err: %w", err)
	}
	defer f.Close()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
