package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"

	"github.com/Huray-hub/eclass-utils/assignments/config"
	auth "github.com/Huray-hub/eclass-utils/authentication"
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	homeCache, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err.Error())
	}

	path := filepath.Join(homeCache, "eclass-tui")
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.OpenFile(
		filepath.Join(path, "assignments.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	// defer func() { NOTE: uhhhhhhhh, remove? ask
	// 	if err = file.Close(); err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// }()

	log.SetOutput(file)

}

func updateConfiguration(o config.Options) {
	log.Print("Not Implemented: updating the configuration file.")
}

func main() {
	opts, creds, err := config.Import()
	if err != nil {
		log.Fatal(err)
	}
	conf := config.Config{
		Credentials: *creds,
		Options:     *opts,
	}

	{
		auth_creds := auth.Credentials{
			Username: "asdf",
			Password: "asdf",
		}

		jar, err := cookiejar.New(nil)
		if err != nil {
			return
		}

		client := &http.Client{
			Jar: jar,
		}
		err = auth.Login(context.Background(), "https://"+conf.Options.BaseDomain, auth_creds, client)
		if err != nil {
			log.Fatal(err)
		}
	}
    logged_in := false // TODO: find out how to check if we logged in correctly
    var ret tea.Model

    if ! logged_in {
        p := tea.NewProgram(NewForm())
        if ret, err = p.Run(); err != nil {
            log.Fatalf("Alas, there's been an error: %v", err)
        }

        login_form := ret.(form) // NOTE: should always succeed
        opts.BaseDomain = login_form.fields[0].Value()
        conf.Credentials.Username = login_form.fields[1].Value()
        conf.Credentials.Password = login_form.fields[2].Value()
    }


	m := NewList(conf)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}

}
