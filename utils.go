package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func colorPermissions(per string) string {
	user := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6f61")).Render(per[1:4])
	group := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5733")).Render(per[4:7])
	other := lipgloss.NewStyle().Foreground(lipgloss.Color("#d500f9")).Render(per[7:])
	return user + group + other
}
func detailedPermission(per string) string {
	s := "Permission: \n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6f61")).Render(" User   : ")
	for i, char := range per {
		if char == 'r' {
			s += "read "
		} else if char == 'w' {
			s += "write "
		} else if char == 'x' {
			s += "execute "
		}
		if i == 3 {
			s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5733")).Render(" Group  : ")
		}
		if i == 6 {
			s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#d500f9")).Render(" Others : ")
		}
	}
	return s

}

func convertSizeBinary(size int) string {

	switch {
	case size < 1024:
		return fmt.Sprintf("%diB", size)
	case size < 1048576:
		return fmt.Sprintf("%.2fKiB", float64(size)/1024)
	case size < 1073741824:
		return fmt.Sprintf("%.2fMiB", float64(size)/1048576)
	case size < 1099511627776:
		return fmt.Sprintf("%.2fGiB", float64(size)/1073741824)
	default:
		return fmt.Sprintf("%.2fTiB", float64(size)/1099511627776)
	}
}

func convertSizeDecimal(size int) string {
	const (
		KB = 1_000
		MB = 1_000_000
		GB = 1_000_000_000
		TB = 1_000_000_000_000
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	case size < GB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size < TB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	default:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	}
}
