package main

import "fmt"

func printHelp() {
	fmt.Println(`Simple Storage Service.

**Usage:**
	triple-s [-port <N>] [-dir <S>]  
	triple-s --help
		    
**Options:**
	- --help     Show this screen.
	- --port N   Port number
	- --dir S    Path to the directory
	`)
}
