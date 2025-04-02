package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	filename     string
	size         string
	created      string
	Lastmodified string
	permission   string
	owner        string
	description  string
	width        int
	height       int
	path         string

	confirmDelete bool
	showRename    bool
	showDelete    bool
	confirmPrompt bool
	confirmRename bool

	err error

	RenameInput textinput.Model
}

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("no arguments provided.")
		return
	}

	m := &model{
		showDelete:    false,
		confirmDelete: false,
		confirmPrompt: false,
		showRename:    false,
		confirmRename: false,
		RenameInput:   textinput.New(),
	}

	m.getStat(args[0])
	m.getFileType(args[0])
	m.getPath(args[0])
	m.RenameInput.Focus()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
