package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

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
	input.Placeholder = "String and fret e.g. A5 or e9"
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

func parse(s string) (string, int, error) {
	var (
		note string
		fret int
	)
	re := regexp.MustCompile(`^\s*([A-Ga-g](b|#)?)\s*([0-9]+)\s*$`)
	matches := re.FindStringSubmatch(s)

	if len(matches) == 0 {
		return note, fret, errors.New("Invalid input")
	}

	note = matches[1]
	fret, _ = strconv.Atoi(matches[3])
	return note, fret, nil
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

	if notes, fret, err := parse(m.input.Value()); err == nil {
		m.result = half_step_plus(notes, fret)
	} else {
		m.result = ""
	}

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s",
		m.input.View(),
		m.result)
}
