package domain

type Generator interface {
	Generate(height, width int) (Maze, error)
}

type Solver interface {
	Search(maze Maze, start, end Point) ([]Point, error)
}


