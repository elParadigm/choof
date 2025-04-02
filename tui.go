package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	fileNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#1A1B26")).
			Background(lipgloss.Color("#2AC3DE")).
			MarginLeft(1).
			MarginBottom(1).
			PaddingLeft(1).
			PaddingRight(1).
			Italic(true)

	otherItemsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).
			PaddingLeft(3).
			MarginRight(1)

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E0E0E0")).
				PaddingLeft(2)

	leftBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("#E0E0E0"))

	mainContainerStyle = lipgloss.NewStyle().
				PaddingRight(3).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#2AC3DE"))

	controlsStyle = lipgloss.NewStyle().
			Faint(true).
			MarginTop(1).
			MarginLeft(1)

	DeleteConfirmation = lipgloss.JoinHorizontal(lipgloss.Left,
		"Do you want to Delete this file?(",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#1A1B26")).Background(lipgloss.Color("#F44336")).Render(" y "),
		"/",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#1A1B26")).Background(lipgloss.Color("#2E7D32")).Render(" n "),
		")")
	RenameConfirmation = lipgloss.JoinHorizontal(lipgloss.Left,
		"Do you want to Rename this file?(",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#1A1B26")).Background(lipgloss.Color("#F44336")).Render(" y "),
		"/",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#1A1B26")).Background(lipgloss.Color("#2E7D32")).Render(" n "),
		")")
)

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.err != nil {
			m.err = nil
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "d", "D":
			if !m.showDelete && !m.showRename {
				m.showDelete = true
				return m, nil
			}

		case "ctrl+r":
			if m.showRename {
				m.RenameInput.SetValue("")
				m.showRename = false
				return m, nil
			} else {
				m.showRename = true
				return m, nil
			}
		case "enter":
			if m.showRename {
				m.confirmRename = !m.confirmRename
			}
		case "ctrl+o":
			return m, tea.ExecProcess(exec.Command("xdg-open", m.path), func(err error) tea.Msg { return err })

		case "y", "Y":
			if m.showDelete && !m.showRename {
				m.confirmDelete = true
				m.DeleteFile(m.path)
				return m, tea.Quit

			}
			if m.showRename && m.confirmRename {
				newPath := filepath.Join(filepath.Dir(m.path), m.RenameInput.Value())
				m.RenameFile(m.path, "/"+m.RenameInput.Value())
				m.getStat(newPath)
				m.getFileType(newPath)
				m.getPath(newPath)
				m.RenameInput.SetValue("")
				m.showRename = false
			}

		case "n", "N", "esc":
			if m.showDelete && !m.showRename {
				m.showDelete = false
				return m, nil
			}
			if m.showRename && m.confirmRename {
				m.RenameInput.SetValue("")
				println("canceled")
				m.showRename = false
				return m, nil
			}

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}
	if m.showRename {
		if m.confirmRename {
			m.RenameInput.Blur()
		} else {
			m.RenameInput.Focus()
		}
		var cmd tea.Cmd
		m.RenameInput, cmd = m.RenameInput.Update(msg)
		return m, cmd
	}
	// Forward messages to the table (for ↑/↓ navigation)
	var cmd tea.Cmd

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Background(lipgloss.Color("#1A1B26")).
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF0000"))

		return errorStyle.Render(fmt.Sprintf("Error: %v\n\nPress any key to continue", m.err))
	}

	if m.showDelete {
		return lipgloss.Place(25, 10, lipgloss.Left, lipgloss.Center, DeleteConfirmation, lipgloss.WithWhitespaceChars(" "))
	}
	if m.showRename {
		return lipgloss.Place(m.width, 10, lipgloss.Left, lipgloss.Center, RenameConfirmation+"\n"+m.RenameInput.View(), lipgloss.WithWhitespaceChars(" "))
	}
	title := fileNameStyle.Render("File: " + m.filename)

	fileInfo := leftBorderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			otherItemsStyle.Render("Size: "+m.size),
			otherItemsStyle.Render("Last modified: "+m.Lastmodified),
			otherItemsStyle.Render("Created: "+m.created),
			otherItemsStyle.Render("Owner: "+m.owner[:len(m.owner)-1]),
			otherItemsStyle.Render("Permission: "+m.permission),
		),
	)

	description := descriptionStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			"Description: ",
			m.description,
		),
	)

	metadata := otherItemsStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			"SHA256: "+m.calculateSHA256(m.path),
			"MD5: "+m.calculateMD5(m.path),
			"Path: "+m.path,
		),
	)

	metadata = lipgloss.NewStyle().
		MarginTop(1).
		Render(metadata)

	controls := controlsStyle.Render(
		"Ctrl+c/q - Close | d - Delete file | Ctrl+r - Rename file | Ctrl+o - Open file",
	)
	// Compose the main view
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		fileInfo,
		description,
	)

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		mainContent,
		metadata,
	)

	mainContainer := mainContainerStyle.Render(fullContent)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		mainContainer,
		controls,
	)
}
