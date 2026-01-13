package application

import (
	"fmt"
	"math"
	"container/heap"

	"labyrinths/internal/domain"
)

func heuristic(x, y domain.Point) int {
	return int(math.Abs(float64(x.Y - y.Y)) + math.Abs(float64(x.X - y.X)))
}

type AStar struct {}

func (a *AStar) Search(labyrinth domain.Maze, start, end domain.Point) ([]domain.Point, error) {
	if !IsCorrectCoord(labyrinth.Height, labyrinth.Width, start) || !IsCorrectCoord(labyrinth.Height, labyrinth.Width, end) {
		return nil, fmt.Errorf("invalid start or end coord")
	}

	priorityQueue := &PriorityQueue{}
	heap.Init(priorityQueue)

	heap.Push(priorityQueue, &Item{coord: start, priority: 0})

	visited := make(map[domain.Point]bool)
	price := make(map[domain.Point]int)
	previous := make(map[domain.Point]*domain.Point)
	price[start] = 0

	for priorityQueue.Len() > 0 {
		currentCell := heap.Pop(priorityQueue).(*Item)
		current := currentCell.coord

		if current == end {
			return a.restorePath(previous, end), nil
		}

		visited[current] = true
		directions := []domain.Point{
			{Y: current.Y + 1, X: current.X},
			{Y: current.Y - 1, X: current.X},
			{Y: current.Y, X: current.X + 1},
			{Y: current.Y, X: current.X - 1},
		}

		for _, direction := range directions {
			if IsCorrectCoord(labyrinth.Height, labyrinth.Width, direction) &&
			   labyrinth.Grid[direction.Y][direction.X].Type != domain.Wall &&
			   !visited[direction] {

				newPrice := price[current] + 1
				currentPrice, ok := price[direction]

				if newPrice < currentPrice || !ok {
					price[direction] = newPrice
					heap.Push(priorityQueue, &Item{
						coord: 	  direction,
						priority: newPrice + heuristic(direction, end),
					})

					previous[direction] = &current
				}
			}
		}
	}

	return nil, fmt.Errorf("couldn't find the path")
}

func (a *AStar) restorePath(previous map[domain.Point]*domain.Point, end domain.Point) []domain.Point {
	path := []domain.Point{}

	for pos := &end; pos != nil; pos = previous[*pos] {
		path = append([]domain.Point{*pos}, path...)
	}

	return path
}
