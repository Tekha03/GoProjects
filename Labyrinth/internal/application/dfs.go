package application

import (
	"math/rand"

	"labyrinths/internal/domain"
)

type DFS struct {}

func (d *DFS) initializeLabyrinth(height, width int) domain.Maze {
	labyrinth := domain.Maze{
		Height: height,
		Width:  width,
		Grid:   make([][]domain.Point, height),
	}

	for i := range labyrinth.Grid {
		labyrinth.Grid[i] = make([]domain.Point, width)
		for j := range labyrinth.Grid[i] {
			labyrinth.Grid[i][j] = domain.Point{Y: i, X: j, Type: domain.Wall}
		}
	}

	return labyrinth
}

func (d *DFS) fillWalls(labyrinth *domain.Maze) {
	for i := 0; i < labyrinth.Height; i++ {
		labyrinth.Grid[i][0].Type = domain.Wall
		labyrinth.Grid[i][labyrinth.Width - 1].Type = domain.Wall
	}

	for i := 1; i < labyrinth.Height - 1; i++ {
		labyrinth.Grid[i][1].Type = domain.Wall
		labyrinth.Grid[i][labyrinth.Width - 2].Type = domain.Wall
	}

	for j := 0; j < labyrinth.Width; j++ {
		labyrinth.Grid[0][j].Type = domain.Wall
		labyrinth.Grid[labyrinth.Height - 1][j].Type = domain.Wall
	}

	for j := 1; j < labyrinth.Width - 1; j++ {
		labyrinth.Grid[1][j].Type = domain.Wall
		labyrinth.Grid[labyrinth.Height - 2][j].Type = domain.Wall
	}
}

func (d *DFS) dfs(labyrinth domain.Maze, currentCell domain.Point) {
	labyrinth.Grid[currentCell.Y][currentCell.X].Type = domain.Space

	directions := []domain.Point {
		{Y: -2, X: 0},
		{Y: +2, X: 0},
		{Y: 0, X: -2},
		{Y: 0, X: +2},
	}
	for i := range directions {
		j := rand.Intn(i + 1)
		directions[i], directions[j] = directions[j], directions[i]
	}

	for _, direction := range directions {
		nextCell := domain.Point{Y: currentCell.Y + direction.Y, X: currentCell.X + direction.X}

		if IsCorrectCoord(labyrinth.Height, labyrinth.Width, nextCell) && labyrinth.Grid[nextCell.Y][nextCell.X].Type == domain.Wall {
			labyrinth.Grid[currentCell.Y + direction.Y / 2][currentCell.X + direction.X / 2].Type = domain.Space
			d.dfs(labyrinth, nextCell)
		}
	}
}

func (d *DFS) Generate(height, width int) (domain.Maze, error) {
	labyrinth := d.initializeLabyrinth(height + 4, width + 4)
	start := domain.Point{Y: 2, X: 2}
	labyrinth.Grid[start.Y][start.X].Type = domain.Space

	d.dfs(labyrinth, start)
	d.fillWalls(&labyrinth)

	return labyrinth, nil
}
