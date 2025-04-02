package main

import (
	"fmt"
	"strconv"
)

func convertSize(size string) string {
	num, err := strconv.Atoi(size)
	if err != nil {
		panic(err)
	}

	switch {
	case num < 1024:
		return fmt.Sprintf("%dB", num)
	case num < 1048576:
		return fmt.Sprintf("%.2fKB", float64(num)/1024)
	case num < 1073741824:
		return fmt.Sprintf("%.2fMB", float64(num)/1048576)
	case num < 1099511627776:
		return fmt.Sprintf("%.2fGB", float64(num)/1073741824)
	default:
		return fmt.Sprintf("%.2fTB", float64(num)/1099511627776)
	}
}
