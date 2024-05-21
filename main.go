package main

import (
	"fmt"
	"log"
	"slices"

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
	history     []string
	historyIter int
}

func initialModel() model {
	var input = textinput.New()
	input.Focus()
	input.CharLimit = 5

	return model{
		input:       input,
		result:      "",
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

func (m *model) historyPush() {
	in := m.input.Value()
	if !slices.Contains(m.history, in) {
		m.history = append(m.history, in)
		m.historyIter = 0
	}

	if len(m.history) > 50 {
		_, m.history = m.history[0], m.history[1:]
	}
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

	m.input, cmd = m.input.Update(msg)
	m.result = updateResult(m.input.Value())

	if m.result == "" {
		m.historyIter = -1
	} else {
		m.historyPush()
	}

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"E A D B [0-21]\n%v\n%v",
		m.input.View(),
		m.result)
}
