package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	input       textinput.Model
	result      string
	fresh       bool
	history     []string
	historyIter int
}

func initialModel() model {
	var input = textinput.New()
	input.Focus()
	input.CharLimit = 5

	return model{
		input:  input,
		result: "",
		fresh:  false,
		// TODO: History should be own struct with double-ended queue.
		// AND TESTS.
		history:     []string{},
		historyIter: -1,
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

// Return if we inserted.
func (m *model) historyPush() bool {
	in := m.input.Value()
	if len(m.history) > 0 {
		last := m.history[len(m.history)-1]
		if in == last {
			return false
		}
	}
	m.history = append(m.history, in)
	m.historyIter = 0

	if len(m.history) > 50 {
		_, m.history = m.history[0], m.history[1:]
	}

	return true
}

func (m *model) historyPop() {
	if len(m.history) == 0 {
		return
	}

	m.history = m.history[:len(m.history)-1]
}

func (m *model) historyDown() {
	m.historyIter -= 1
	if m.historyIter <= -1 {
		m.input.SetValue("")
		m.historyIter = -1
		return
	}

	index := len(m.history) - m.historyIter - 1
	m.input.SetValue(m.history[index])
}

func (m *model) historyUp() {
	if len(m.history) == 0 {
		return
	}

	m.historyIter += 1
	if m.historyIter >= len(m.history) {
		m.historyIter = len(m.history) - 1
	}

	index := len(m.history) - m.historyIter - 1
	m.input.SetValue(m.history[index])
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyUp:
			m.historyUp()
		case tea.KeyDown:
			m.historyDown()
		case tea.KeyLeft:
			m.input.Reset()
		}
	}

	old := m.input.Value()
	old_result := updateResult(old)

	m.input, cmd = m.input.Update(msg)

	in := m.input.Value()
	m.result = updateResult(in)

	if in == "" {
		m.fresh = false
	}

	if m.result == "" {
		m.historyIter = -1
	} else if in != old {
		if m.fresh && old_result != "" && strings.HasPrefix(in, old) {
			m.historyPop()
		}
		m.fresh = m.historyPush()
	}

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"E A D B [0-21]\n%v\n%v",
		m.input.View(),
		m.result)
}
