package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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

	log.SetOutput(file)
}

func main() {
	conf, err := config.ImportDefault()
	if err != nil {
        log.Fatal(fmt.Errorf("error importing user config: %v", err))
	}

    var state WindowState
	// attempt to login
    ctx, cancel := context.WithCancel(context.Background())
    client, err := auth.Session(ctx, "https://"+conf.Options.BaseDomain, conf.Credentials, nil)
	// if we failed ask user for credentials
	if err != nil {
        cancel()
        state = Login
	} else {
        state = List
    }

	w := NewWindow(state, *conf, client, err)
	p := tea.NewProgram(w)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}

}
