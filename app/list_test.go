package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/Huray-hub/eclass-utils/assignments/config"
	tea "github.com/charmbracelet/bubbletea"
)

func withRedirectStdout(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	res, _ := ioutil.ReadAll(r)
	os.Stdout = old

	return string(res)
}

var conf = config.Config{
	Credentials: config.Credentials{
		Username: "asdf",
		Password: "asdf",
	},

	Options: config.Options{
		BaseDomain:     "eclass.uniwa.gr",
		PlainText:      false,
		IncludeExpired: true,
		ExportICS:      false,
	},
}

func TestIncludeExpiredKey(t *testing.T) {

	m := NewList(conf)
	m.testing = true

	in := strings.NewReader("iq")

    p := tea.NewProgram(m, tea.WithoutRenderer(), tea.WithInput(in))

	result := withRedirectStdout(func() {
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	})

	expected := "false"

	if result != expected {
		t.Errorf("Expected '%s' but got '%s'", result, expected)
	}
}
