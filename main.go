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
	input       textinput.Model
	result      string
	fresh       bool
	history     *deque.Deque[string]
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
		// TODO: Write tests for history.
		history:     deque.New[string](),
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
	if m.history.Len() > 0 {
		if in == m.history.Back() {
			return false
		}
	}
	m.history.PushBack(in)
	m.historyIter = 0

	if m.history.Len() > 50 {
		m.history.PopFront()
	}

	return true
}

func (m *model) historyPop() {
	if m.history.Len() == 0 {
		return
	}

	m.history.PopBack()
}

func (m *model) historyDown() {
	m.historyIter -= 1
	if m.historyIter <= -1 {
		m.input.Reset()
		m.historyIter = -1
		return
	}

	m.setInput()
}

func (m *model) historyUp() {
	if m.history.Len() == 0 {
		return
	}

	m.historyIter += 1
	if m.historyIter >= m.history.Len() {
		m.historyIter = m.history.Len() - 1
	}

	m.setInput()
}

// Set textinput based on the historyIter.
func (m *model) setInput() {
	index := m.history.Len() - m.historyIter - 1
	m.input.SetValue(m.history.At(index))
	m.input.CursorEnd()
}

func (m *model) UpdateInput(msg tea.Msg) {
	old := m.input.Value()
	old_result := updateResult(old)

	m.input, _ = m.input.Update(msg)

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

	m.UpdateInput(msg)

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"E A D B [0-21]\n%v\n%v",
		m.input.View(),
		m.result)
}
