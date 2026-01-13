package infrastructure

import (
	"bufio"
	"fmt"
	"os"

	"labyrinths/internal/domain"
)

func SaveLabyrinthToFile(maze domain.Maze, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error in creating file: %s", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, row := range maze.Grid {
		for _, cell := range row {
			fmt.Fprintf(writer, "%s", string(cell.Type))
		}
		fmt.Fprintf(writer, "\n")
	}

	return writer.Flush()
}

func SetPath(labyrinth *domain.Maze, path []domain.Point, start, end domain.Point) {
	for _, cell := range path {
		labyrinth.Grid[cell.Y][cell.X].Type = domain.Path
	}

	labyrinth.Grid[start.Y][start.X].Type = domain.Start
	labyrinth.Grid[end.Y][end.X].Type = domain.End
}
