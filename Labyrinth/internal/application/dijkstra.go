package application

import (
	"container/heap"
	"fmt"

	"labyrinths/internal/domain"
)

type Dijkstra struct {}

func (d *Dijkstra) Search(labyrinth domain.Maze, start, end domain.Point) ([]domain.Point, error) {
	if !IsCorrectCoord(labyrinth.Height, labyrinth.Height, start) || !IsCorrectCoord(labyrinth.Height, labyrinth.Width, end) {
		return nil, fmt.Errorf("invalid start or end coord")
	}

	priorityQueue := &PriorityQueue{}
	heap.Init(priorityQueue)
	heap.Push(priorityQueue, &Item{coord: start, priority: 0})

	distance := make(map[domain.Point]int)
	previous := make(map[domain.Point]*domain.Point)
	visited  := make(map[domain.Point]bool)
	distance[start] = 0

	for priorityQueue.Len() > 0 {
		currentCell := heap.Pop(priorityQueue).(*Item)
		current := currentCell.coord

		if visited[current] {
			continue
		}
		visited[current] = true

		if current == end {
			return d.restorePath(previous, end), nil
		}

		directions := []domain.Point{
			{Y: current.Y + 1, X: current.X},
			{Y: current.Y - 1, X: current.X},
			{Y: current.Y, X: current.X + 1},
			{Y: current.Y, X: current.X - 1},
		}

		for _, direction := range directions {
			if IsCorrectCoord(labyrinth.Height, labyrinth.Width, direction) &&
			   labyrinth.Grid[direction.Y][direction.X].Type != domain.Wall {
				newDistance := distance[current] + 1

				oldDistance, exists := distance[direction]

				if !exists || newDistance < oldDistance {
					distance[direction] = newDistance
					heap.Push(priorityQueue, &Item{
						coord: direction,
						priority: newDistance,
					})

					parent := current
					previous[direction] = &parent
				}
			}
		}
	}

	return nil, fmt.Errorf("couldn't find the path")
}

func (d *Dijkstra) restorePath(previous map[domain.Point]*domain.Point, end domain.Point) []domain.Point {
	path := []domain.Point{}
	for point := &end; point != nil; point = previous[*point] {
		path = append([]domain.Point{*point}, path...)
	}

	return  path
}
