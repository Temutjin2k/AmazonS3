package utils

import "fmt"

var Dir string // Storing Directory flag

func PrintHelp() {
	helpMessage := `Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`

	fmt.Println(helpMessage)
}
