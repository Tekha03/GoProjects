package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println(`Usage: maze-app [-hV] [COMMAND]
Maze generator and solver CLI application.
  -h, --help      Show this help message and exit.
  -V, --version   Print version information and exit.
Commands:
  generate  Generate a maze with specified algorithm and dimensions.
  solve     Solve a maze with specified algorithm and points.`)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("expected command")
		return
	}

	Launch()
}
