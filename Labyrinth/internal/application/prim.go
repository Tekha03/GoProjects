package application

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"labyrinths/internal/domain"
)

type Prim struct {}

func RandomInt(n int) (int, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errors.New("failed to generate random int"), err)
	}

	return int(number.Int64()), nil
}

func IsCorrectCoord(height, width int, coord domain.Point) bool {
	return coord.Y >= 0 && coord.Y < height && coord.X >= 0 && coord.X < width
}

func (p *Prim) getNearestPoints(coord domain.Point, labyrinth *domain.Maze) []domain.Point {
	directions := []domain.Point {
		{Y: coord.Y + 1, X: coord.X},
		{Y: coord.Y - 1, X: coord.X},
		{Y: coord.Y, X: coord.X + 1},
		{Y: coord.Y, X: coord.X - 1},
	}

	var nearby []domain.Point
	for _, diredirection := range directions {
		neighbor := domain.Point{Y: diredirection.Y, X: diredirection.X}
		if IsCorrectCoord(labyrinth.Height, labyrinth.Width, neighbor) {
			nearby = append(nearby, neighbor)
		}
	}

	return nearby
}

func (p *Prim) Generate(height, width int) (domain.Maze, error) {
	labyrinth := p.newLabyrinth(height, width)
	start, err := p.setStartPoint(height, width, &labyrinth)

	if err != nil {
		return domain.Maze{}, fmt.Errorf("%w: %v", errors.New("failed to set starting point"), err)
	}

	err = p.buildWalls(height, width, start, &labyrinth)
	if err != nil {
		return  domain.Maze{}, fmt.Errorf("%w: %v", errors.New("failed to build walls"), err)
	}

	return labyrinth, nil
}

func (p *Prim) newLabyrinth(height, width int) domain.Maze {
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

func (p *Prim) setStartPoint(height, width int, labyrinth *domain.Maze) (domain.Point, error) {
	row, err := RandomInt(height)
	if err != nil {
		return domain.Point{}, fmt.Errorf("%w: %v", errors.New("failed to set row of starting point"), err)
	}

	column, err := RandomInt(width)
	if err != nil {
		return domain.Point{}, fmt.Errorf("%w: %v", errors.New("failed to set column of starting point"), err)
	}

	labyrinth.Grid[row][column].Type = domain.Start
	return domain.Point{Y: row, X: column, Type: domain.Start}, nil
}

func (p *Prim) buildWalls(height, width int, start domain.Point, labyrinth *domain.Maze) error {
	walls := []domain.Point{
		{Y: start.Y - 1, X: start.X},
		{Y: start.Y + 1, X: start.X},
		{Y: start.Y, X: start.X - 1},
		{Y: start.Y, X: start.X + 1},
	}

	for len(walls) > 0 {
		index, err := RandomInt(len(walls))
		if err != nil {
			return fmt.Errorf("%w: %v", errors.New("failed to set wall"), err)
		}

		wall := walls[index]
		walls = append(walls[:index], walls[index + 1:]...)

		if !IsCorrectCoord(height, width, wall) || labyrinth.Grid[wall.Y][wall.X].Type != domain.Wall {
			continue
		}

		nearby := p.getNearestPoints(wall, labyrinth)
		if len(nearby) == 1 {
			labyrinth.Grid[wall.Y][wall.X].Type = domain.Space
			labyrinth.Grid[nearby[0].Y][nearby[0].X].Type = domain.Space

			newWalls := []domain.Point{
				{Y: wall.Y - 1, X: wall.X},
				{Y: wall.Y + 1, X: wall.X},
				{Y: wall.Y, X: wall.X + 1},
				{Y: wall.Y, X: wall.X - 1},
			}

			for _, newWall := range newWalls {
				if IsCorrectCoord(height, width, newWall) && labyrinth.Grid[newWall.Y][newWall.X].Type == domain.Wall {
					walls = append(walls, newWall)
				}
			}
		}
	}

	return nil
}
