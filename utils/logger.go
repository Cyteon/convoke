package utils

import (
	"log"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func GetColor(color string) string {
	switch color {
	case "red":
		color = Red
	case "green":
		color = Green
	case "yellow":
		color = Yellow
	case "blue":
		color = Blue
	case "magenta":
		color = Magenta
	case "cyan":
		color = Cyan
	case "gray":
		color = Gray
	case "white":
		color = White
	default:
		color = Reset
	}

	return color
}

func Log(message string, color string) {
	color = GetColor(color)

	log.Printf("| %s%s%s\n", color, message, Reset)
}

func LogFatal(message string, color string) {
	color = GetColor(color)

	log.Fatalf("| %s%s%s\n", color, message, Reset)
}
