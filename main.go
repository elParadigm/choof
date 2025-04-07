package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	filename     string
	size         int64
	created      string
	Lastmodified string
	permission   string
	owner        string
	description  string
	width        int
	height       int
	path         string

	isFile        bool
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
		fmt.Println("please provide a file name or a path to file.")
		return
	}

	m := &model{
		isFile:        true,
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
	m.checkFile()

	if m.err != nil {
		fmt.Printf("Error: %v\n", m.err)
		os.Exit(0)
	}
	m.RenameInput.Focus()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)

		os.Exit(1)
	}

}
