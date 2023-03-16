package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"testing"

	"github.com/Huray-hub/eclass-utils/assignment/config"
	"github.com/Huray-hub/eclass-utils/auth"
	"github.com/Huray-hub/eclass-utils/course"
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
	Credentials: auth.Credentials{
		Username: "asdf",
		Password: "asdf",
	},

	Options: config.Options{
		Options: course.Options{
			BaseDomain: "eclass.uniwa.gr",
		},
		PlainText:      false,
		IncludeExpired: true,
		ExportICS:      false,
	},
}

type test struct {
	inputs   string
	expected string
}

var tests = []test{
	{inputs: "q", expected: ""},
	{inputs: "iq", expected: "false"},
	{inputs: "iiq", expected: "falsetrue"},
}

func TestIncludeExpired(t *testing.T) {

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Errorf("error creating cookie jar: %v", err)
	}
	m := NewCourseList(conf, &http.Client{
		Jar: jar,
	})
	m.testing = true

	for _, test := range tests {

		in := strings.NewReader(test.inputs)

		p := tea.NewProgram(m, tea.WithoutRenderer(), tea.WithInput(in))

		result := withRedirectStdout(func() {
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		})

		if result != test.expected {
			t.Errorf("Expected '%s' but got '%s'", test.expected, result)
		}
	}
}
