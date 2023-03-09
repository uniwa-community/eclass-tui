package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"

	"github.com/Huray-hub/eclass-utils/assignment/config"
	"github.com/Huray-hub/eclass-utils/auth"
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
	conf, err := config.ImportDefault()
	if err != nil {
        log.Fatal(fmt.Errorf("error importing user config: %v", err))
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}

	client := &http.Client{
		Jar: jar,
	}

	// attempt to login
    log.Println("Attempting login with", conf.Credentials.Username, "|", conf.Credentials.Password, "|", conf.Options.BaseDomain)
	client, err = auth.Login(context.Background(), "https://"+conf.Options.BaseDomain, conf.Credentials, client)
	// if we failed ask user for credentials
	if err != nil {
        log.Println(fmt.Errorf("login failed (asking user now): %v", err))
		p := tea.NewProgram(NewForm(err, *conf))
		f, err := p.Run()
		if err != nil {
			log.Fatalf("Alas, there's been an error: %v", err)
		}

		login_form := f.(form) // NOTE: should always succeed
		conf.Options.BaseDomain = login_form.fields[Domain].Value()
        conf.Credentials.Username = login_form.fields[Username].Value()
		conf.Credentials.Password = login_form.fields[Password].Value()

        log.Println("Attempting login with", conf.Credentials.Username, conf.Credentials.Password, conf.Options.BaseDomain)
        client, err = auth.Login(context.Background(), "https://"+conf.Options.BaseDomain, conf.Credentials, client)
        if err != nil || client == nil {
            log.Fatal(fmt.Errorf("unexpected login failure: %v", err))
        }
	}

	m := NewList(*conf, client)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}

}
