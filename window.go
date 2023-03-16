package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Huray-hub/eclass-utils/assignment/config"
	tea "github.com/charmbracelet/bubbletea"
)

type WindowState int

const (
	Login WindowState = iota + 1
	List
)

func (s WindowState) String() string {
	switch s {
	case 0:
		return "Uninitialized"
	case Login:
		return "Login"
	case List:
		return "List"
	default:
		return "Undefined"
	}
}

type window struct {
	login form
	list  courseList
	state WindowState
}

func NewWindow(state WindowState, conf config.Config, session *http.Client, err error) window {
	w := window{state: state}
	switch w.state {
	case Login:
		w.login = NewForm(err, conf)
	case List:
        log.Println("Window shows list")
		w.list = NewList(conf, session)
	}

	return w
}

type loginSuccess struct {
	conf    config.Config
	session *http.Client
}

func (window) Init() tea.Cmd {
	return nil
}
func (w window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case loginSuccess:
		w.state = List
		w.list = NewList(msg.conf, msg.session)
    case tea.KeyMsg:
        switch key := msg.String(); key {
        case "ctrl+c":
            log.Println("quit")
            return w, tea.Quit
        }
	}

	switch w.state {
	case Login:
		l, cmd := w.login.Update(msg)
		w.login = l.(form)
		return w, cmd
	case List:
		l, cmd := w.list.Update(msg)
		w.list = l.(courseList)
		return w, cmd
	default:
		log.Fatal(fmt.Errorf("window state %s", w.state))
		return w, nil // unreachable ?
	}
}

func (w window) View() string {
	switch w.state {
	case Login:
		return w.login.View()
	case List:
		return w.list.View()
	default:
		log.Fatal(fmt.Errorf("window state %s", w.state))
		return "" // unreachable
	}
}
