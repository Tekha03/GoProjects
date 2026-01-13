package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"labyrinths/internal/domain"
	"labyrinths/internal/application"
)

func Launch() {
	command := os.Args[1]
	f := &application.File{}
	switch command {
	case "generate":
		generateCommand := flag.NewFlagSet("generate", flag.ExitOnError)

		algorithm := generateCommand.String("algorithm", "prim", "prim|eller")
		height := generateCommand.Int("height", 10, "")
		width := generateCommand.Int("width", 10, "")
		out := generateCommand.String("output", "", "output file")

		generateCommand.StringVar(algorithm, "a", "prim", "algorithm alias")
		generateCommand.IntVar(height, "h", 10, "height alias")
		generateCommand.IntVar(width, "w", 10, "width alias")
		generateCommand.StringVar(out, "o", "", "output alias")

		generateCommand.Parse(os.Args[2:])

		err := f.GenerateToFile(*algorithm, *height, *width, *out)
		if err != nil {
			fmt.Println(err)
			return
		}

	case "solve":
		solveCommand := flag.NewFlagSet("solve", flag.ExitOnError)
		algorithm := solveCommand.String("algorithm", "astar", "dfs")
		file := solveCommand.String("file", "", "")
		start := solveCommand.String("start", "", "")
		end := solveCommand.String("end", "", "")
		out := solveCommand.String("output", "", "")

		solveCommand.Parse(os.Args[2:])
		startCoord, err := parseCoord(*start)
		if err != nil {
			fmt.Println(err)
			return
		}

		endCoord, err := parseCoord(*end)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = f.SolveToFile(
			*algorithm,
			*file,
			startCoord,
			endCoord,
			*out,
		)

		if err != nil {
			fmt.Println(err)
			return
		}

	default:
		fmt.Println("unknown command: ", command)
		return
	}
}

func parseCoord(s string) (domain.Point, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return domain.Point{}, fmt.Errorf("Invalid point format: %s, expected format: x,y", s)
	}

	var row, column int
	fmt.Sscanf(parts[0], "%d", &row)
	fmt.Sscanf(parts[1], "%d", &column)

	return domain.Point{Y: row + 1, X: column + 1}, nil
}
