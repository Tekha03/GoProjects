package application

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"labyrinths/internal/domain"
	"labyrinths/internal/infrastructure"
)

type File struct {}

func (f *File) GenerateToFile(algorithm string, height, width int, output string) error {
	var generator domain.Generator

	switch algorithm {
	case "prim":
		generator = &Prim{}
	case "dfs":
		generator = &DFS{}
	default:
		return fmt.Errorf("unknown generator: %s", algorithm)
	}

	labyrinth, err := generator.Generate(height, width)
	if err != nil {
		return  err
	}

	return infrastructure.SaveLabyrinthToFile(labyrinth, output)
}

func (f *File) SolveToFile(algorithm, file string, start, end domain.Point, output string) error {
	labyrinth, err := f.loadLabyrinthFromFile(file)
	if err != nil {
		return err
	}

	var solver domain.Solver
	switch strings.ToLower(algorithm) {
	case "astar":
		solver = &AStar{}
	case "dijkstra":
		solver = &Dijkstra{}
	default:
		return fmt.Errorf("unknown solver: %s", algorithm)
	}

	path, err := solver.Search(labyrinth, start, end)
	if err != nil {
		return err
	}

	infrastructure.SetPath(&labyrinth, path, start, end)
	return infrastructure.SaveLabyrinthToFile(labyrinth, output)
}

func (f *File) loadLabyrinthFromFile(filename string) (domain.Maze, error) {
	file, err := os.Open(filename)
	if err != nil {
		return domain.Maze{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid [][]domain.Point

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]domain.Point, len(line))

		for i, char := range line {
			row[i] = domain.Point{Type: domain.CellType(char)}
		}
		grid = append(grid, row)
	}

	if len(grid) == 0 {
		return domain.Maze{}, errors.New("empty labyrinth file")
	}

	return domain.Maze{
		Height: len(grid),
		Width:  len(grid[0]),
		Grid:   grid,
	}, nil
}
