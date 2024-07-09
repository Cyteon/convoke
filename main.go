package main

import (
	"convoke/utils"
)

func main() {
	utils.LogInfo("Starting convoke", "cyan")
	utils.LoadDB()
}
