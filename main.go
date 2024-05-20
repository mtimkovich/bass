package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	input  textinput.Model
	result string
}

func initialModel() model {
	var input = textinput.New()
	input.Focus()
	input.CharLimit = 5

	return model{
		input:  input,
		result: "",
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.input, cmd = m.input.Update(msg)
	m.result = updateResult(m.input.Value())

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"String and fret (e.g. A5 or e9)\n%s\n%s",
		m.input.View(),
		m.result)
}
