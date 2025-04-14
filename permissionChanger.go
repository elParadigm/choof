package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderColor    = lipgloss.Color("#2AC3DE")
	titleBGColor   = lipgloss.Color("#2AC3DE")
	titleTextColor = lipgloss.Color("#1A1B26")

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Background(titleBGColor).
			Foreground(titleTextColor).
			Bold(true).
			Padding(0, 1)
)

type permBit struct {
	label   string
	mask    os.FileMode
	checked bool
	focused bool
}

type permodel struct {
	filename string
	perms    []permBit
	cursor   int
	err      error
	done     bool
	width    int
	height   int
}

func NewPermModel(filename string) tea.Model {
	info, err := os.Stat(filename)
	if err != nil {
		return model{err: err}
	}

	mode := info.Mode().Perm()

	bits := []permBit{
		{"Owner Read", 0400, mode&0400 != 0, true},
		{"Owner Write", 0200, mode&0200 != 0, false},
		{"Owner Exec", 0100, mode&0100 != 0, false},
		{"Group Read", 0040, mode&0040 != 0, false},
		{"Group Write", 0020, mode&0020 != 0, false},
		{"Group Exec", 0010, mode&0010 != 0, false},
		{"Other Read", 0004, mode&0004 != 0, false},
		{"Other Write", 0002, mode&0002 != 0, false},
		{"Other Exec", 0001, mode&0001 != 0, false},
	}

	return permodel{
		filename: filename,
		perms:    bits,
	}
}

func RunPermissionEditor(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("❌ File does not exist:", filename)
		return
	}
	p := tea.NewProgram(NewPermModel(filename), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running permission editor:", err)
	}
}

func (m permodel) Init() tea.Cmd {
	return nil
}

func (m permodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.perms[m.cursor].focused = false
				m.cursor--
				m.perms[m.cursor].focused = true
			}
		case "down", "j":
			if m.cursor < len(m.perms)-1 {
				m.perms[m.cursor].focused = false
				m.cursor++
				m.perms[m.cursor].focused = true
			}
		case " ":
			m.perms[m.cursor].checked = !m.perms[m.cursor].checked
		case "enter":
			var mode os.FileMode
			for _, p := range m.perms {
				if p.checked {
					mode |= p.mask
				}
			}
			err := os.Chmod(m.filename, mode)
			if err != nil {
				m.err = err
			} else {
				m.done = true
			}
		}
	}
	return m, nil
}

func (m permodel) View() string {
	if m.err != nil {
		return boxStyle.Render(fmt.Sprintf("❌ Error: %v", m.err))
	}
	if m.done {
		content := fmt.Sprintf("✅ Permissions updated for '%s'\nPress q to quit.\n", m.filename)
		return boxStyle.Render(content)
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf(" Change permissions for: %s ", m.filename)))
	b.WriteString("\n\n")

	for _, p := range m.perms {
		cursor := " "
		if p.focused {
			cursor = ">"
		}
		check := "[ ]"
		if p.checked {
			check = "[x]"
		}
		b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, check, p.label))
	}

	b.WriteString("\n↑/↓ to move, space to toggle, Enter to apply, q to quit.")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, boxStyle.Render(b.String()))
}
