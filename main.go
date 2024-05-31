package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gammazero/deque"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer f.Close()

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	input   textinput.Model
	result  string
	active  bool
	history History
}

func initialModel() model {
	var input = textinput.New()
	input.Focus()
	input.CharLimit = 5

	return model{
		input:  input,
		result: "",
		active: false,
		history: History{
			queue: deque.New[string](),
			iter:  -1,
		},
	}

}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func updateResult(in string) string {
	if notes, fret, err := parse(in); err == nil {
		return half_step_plus(notes, fret)
	}

	return ""
}

// Set textinput based on the history iter.
func (m *model) setInput() {
	m.active = false
	m.input.SetValue(m.history.Entry())
	m.input.CursorEnd()
}

func (m *model) UpdateInput(msg tea.Msg) {
	old := m.input.Value()
	old_result := updateResult(old)

	m.input, _ = m.input.Update(msg)

	in := m.input.Value()
	m.result = updateResult(in)

	if in == "" {
		m.active = false
	}

	if m.result == "" {
		m.history.iter = -1
	} else if in != old {
		if m.active && old_result != "" && strings.HasPrefix(in, old) {
			m.history.Pop()
		}
		m.active = m.history.Push(in)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyUp:
			m.history.Up()
			m.setInput()
		case tea.KeyDown:
			m.history.Down()
			m.setInput()
		case tea.KeyLeft:
			m.input.Reset()
		}
	}

	m.UpdateInput(msg)

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"E A D G [0-21]\n%v\n%v",
		m.input.View(),
		m.result)
}
