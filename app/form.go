package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type form struct {
	fields        []textinput.Model
	style         lip.Style
	selectedInput int
}

var defaultFormStyle = lip.NewStyle().
	Align(lip.Center, lip.Center)

type invalidDomain struct {
    name string
}

func (err invalidDomain) Error() string {
    return err.name
}

func validateDomain(domain string) error {
	if !strings.HasPrefix(domain, "https://eclass.") {
		return invalidDomain{
            name: "Invalid domain '" + domain + "'!",
        }
	}

	return nil
}

func NewForm() form {
	username := textinput.New()
	username.Prompt = "Username: "
	password := textinput.New()
	password.Prompt = "Password: "
	password.EchoCharacter = '*'
	password.EchoMode = textinput.EchoPassword
	domain := textinput.New()
    // domain.Validate = validateDomain
	domain.Prompt = "Domain:   "
	domain.Placeholder = "https://eclass.<domain>.<xyz>"

	return form{
		fields: []textinput.Model{
			username,
			password,
			domain,
		},
		style:         defaultFormStyle,
		selectedInput: 0,
	}
}

func (f form) Init() tea.Cmd {
	return nil
}

func (f form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "tab":
			if f.selectedInput < 2 {
				f.selectedInput++
			}
		case "shift+tab":
			if f.selectedInput > 0 {
				f.selectedInput--
			}
		case "enter":
			log.Println("quit?")
			return f, tea.Quit
		}
	case tea.WindowSizeMsg:
		f.style.Height(msg.Height)
		f.style.Width(msg.Width)
	}

	cmd = f.updateFields(msg)
	return f, cmd
}

func (f form) updateFields(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.fields))

	for i, field := range f.fields {
		if i == f.selectedInput {
			field.Focus()
			field.SetValue(field.Value())
		} else {
			field.Blur()
		}
		f.fields[i], cmds[i] = field.Update(msg)
	}

	return tea.Batch(cmds...)
}

var (
	selectedInputStyle = lip.NewStyle().Foreground(lip.Color("5"))
	inputStyle         = lip.NewStyle()
	boxStyle           = lip.NewStyle().BorderStyle(lip.RoundedBorder()).Width(64)
)

func (f form) View() string {
	var rendered []string

	for i, field := range f.fields {
		if i == f.selectedInput {
			rendered = append(rendered, selectedInputStyle.Render(field.View()))
		} else {
			rendered = append(rendered, inputStyle.Render(field.View()))
		}
	}

	str := lip.JoinVertical(lip.Left, rendered...)

	return f.style.Render(boxStyle.Render(str))
}
