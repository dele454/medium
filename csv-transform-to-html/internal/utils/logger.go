package utils

import (
	"fmt"
	"log"
)

type Color string

const (
	ColorError = "\u001b[31m"
	ColorOK    = "\u001b[32m"
	ColorWarn  = "\u001b[33m"
	ColorReset = "\u001b[0m"
)

// Log logs message to stdout with color codes.
func Log(color Color, message any) {
	switch color {
	case ColorError:
		message = fmt.Sprintf("ERROR: %+v", message)
	case ColorWarn:
		message = fmt.Sprintf("WARN: %+v", message)
	default:
		message = fmt.Sprintf("INFO: %+v", message)
	}

	log.Println(string(color), message, string(ColorReset))
}
