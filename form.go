package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Huray-hub/eclass-utils/assignment/config"
	"github.com/Huray-hub/eclass-utils/auth"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type form struct {
	fields        []textinput.Model
	submit        button
	style         lip.Style
	selectedInput int
	session       *http.Client
	loginFailed   int
	reason        error
}

const (
	Username = iota
	Password
	Domain
	Submit
)

var defaultFormStyle = lip.NewStyle().
	Align(lip.Center, lip.Center)

type invalid struct {
	why string
}

func (err invalid) Error() string {
	return err.why
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (f form) validate() bool {
	errs := make([]error, 3)
	errs[Username] = validateUsername(f.fields[Username].Value())
	errs[Password] = validatePassword(f.fields[Password].Value())
	errs[Domain] = validateDomain(f.fields[Domain].Value())

	for _, v := range errs {
		if v != nil {
			return false
		}
	}
	return true
}

func validateDomain(domain string) error {
	err := validateLength(domain, "Domain")
	if err != nil {
		return err
	}

	prefix := "eclass."

	// if the domain so far does not have the prefix
	// where prefix is truncated to the size of the domain
	if !strings.HasPrefix(domain, prefix[0:min(len(domain), len(prefix))]) {
		return invalid{
			why: "Domain must start with: " + prefix + "",
		}
	}

	return nil
}

func validateUsername(name string) error {
	return validateLength(name, "Username")
}

func validatePassword(pass string) error {
	return validateLength(pass, "Password")
}

func validateLength(in string, what string) error {
	if len(in) == 0 {
		return invalid{
			why: what + " must be non empty!",
		}
	}
	return nil
}

func NewForm(e error, conf config.Config) form {
	username := textinput.New()
	username.Prompt = "Username: "
    username.SetValue(conf.Credentials.Username)
	username.Validate = validateUsername

	password := textinput.New()
	password.Prompt = "Password: "
    username.SetValue(conf.Credentials.Password)
	password.Validate = validatePassword
	password.EchoCharacter = '*'
	password.EchoMode = textinput.EchoPassword

	domain := textinput.New()
	domain.Validate = validateDomain
    domain.SetValue(conf.Options.BaseDomain)
	domain.Prompt = "Domain:   "
	domain.Placeholder = "eclass.<domain>.<xyz>"

	return form{
		fields: []textinput.Model{
			Username: username,
			Password: password,
			Domain:   domain,
		},
		submit:        NewButton("Submit"),
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
			if f.selectedInput != Submit {
				f.selectedInput++
			}
		case "shift+tab":
			if f.selectedInput != Username {
				f.selectedInput--
			}
		case "enter":
			if f.selectedInput != Submit {
				f.selectedInput++
			} else if f.validate() {
				log.Println("Attempting login..")
				return f, f.loginCmd
			}
		case "ctrl+c":
			log.Fatal("Keyboard SIGINT")
		}
	case tea.WindowSizeMsg:
		f.style.Height(msg.Height)
		f.style.Width(msg.Width)
	case loginFail:
		f.loginFailed++
		log.Println("Login failed!")
	}

	cmd = f.updateFields(msg)
	return f, cmd
}

type loginFail struct{ Err error }

func (l loginFail) Error() string {
	return l.Err.Error()
}

type startSpinner struct{}

func startSpinnerCmd() tea.Msg {
	return startSpinner{}
}

func (f *form) loginCmd() tea.Msg {
	creds := auth.Credentials{
		Username: f.fields[Username].Value(),
		Password: f.fields[Password].Value(),
	}

	log.Println("Attempting login with", creds.Username, creds.Password, f.fields[Domain].Value())

	client, err := auth.Login(
		context.Background(),
		"https://"+f.fields[Domain].Value(),
		creds,
		nil,
	)

	if err != nil || client == nil {
		return loginFail{fmt.Errorf("login failed: %v", err)}
	}
	time.Sleep(time.Second * 2)

	log.Println("Login success!")
	return tea.Quit()
}

func (f *form) updateFields(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.fields)+1)

	for i, field := range f.fields {
		if i == f.selectedInput {
			field.Focus()
			field.SetValue(field.Value())
		} else {
			field.Blur()
		}
		f.fields[i], cmds[i] = field.Update(msg)
	}

	if f.selectedInput == 3 {
		f.submit.Focus()
	} else {
		f.submit.Blur()
	}
	f.submit, cmds[Submit] = f.submit.Update(msg)

	return tea.Batch(cmds...)
}

var (
	selectedInputStyle = lip.NewStyle().Foreground(lip.Color("5"))
	warningStyle       = lip.NewStyle().Foreground(lip.Color("1"))
	inputStyle         = lip.NewStyle()
	boxStyle           = lip.NewStyle().BorderStyle(lip.RoundedBorder()).Width(64)
)

func (f form) View() string {
	var inputsBoxes []string

	var warnings []string
	for i, field := range f.fields {
		if i == f.selectedInput {
			inputsBoxes = append(inputsBoxes, selectedInputStyle.Render(field.View()))
		} else {
			inputsBoxes = append(inputsBoxes, inputStyle.Render(field.View()))
		}
		if field.Err != nil {
			warnings = append(warnings, warningStyle.Render(field.Err.Error()))
		}
	}

    var loginFailed_msg string
	if f.loginFailed == 1 {
        loginFailed_msg = "Failed to login!"
	} else if f.loginFailed > 1 {
        loginFailed_msg = fmt.Sprintf("Failed to login(%d)!", f.loginFailed - 1)

    }
		warnings = append(warnings, warningStyle.Render(loginFailed_msg))

	str := lip.JoinVertical(lip.Left,
		boxStyle.Render(
			lip.JoinVertical(lip.Left, inputsBoxes...),
		),
		f.submit.View(),
		boxStyle.Render(
			lip.JoinVertical(lip.Left, warnings...),
		),
	)

	return f.style.Render(str)
}
