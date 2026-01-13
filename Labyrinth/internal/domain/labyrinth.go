package domain

type CellType string

const (
	Wall  CellType = "#"
	Space CellType = " "
	Start CellType = "O"
	End   CellType = "X"
	Path  CellType = "."
)

type Point struct {
	Y int
	X int
	Type CellType
}

type Maze struct {
	Width 	int
	Height 	int
	Grid 	[][]Point
}
