package main

import (
	"fmt"
)

func convertSize(size int) string {

	switch {
	case size < 1024:
		return fmt.Sprintf("%dB", size)
	case size < 1048576:
		return fmt.Sprintf("%.2fKB", float64(size)/1024)
	case size < 1073741824:
		return fmt.Sprintf("%.2fMB", float64(size)/1048576)
	case size < 1099511627776:
		return fmt.Sprintf("%.2fGB", float64(size)/1073741824)
	default:
		return fmt.Sprintf("%.2fTB", float64(size)/1099511627776)
	}
}
