package agent

import (
	"github.com/BattlesnakeOfficial/rules"
	_ "log"
)

type CellKind int

const (
	CellEmpty CellKind = iota
	CellFood
	CellSnakeHead
	CellSnakeBody
	CellSnakeTail
)

type Cell interface {
	Kind() CellKind
	IsPassable() bool
	Coordinates() rules.Point
	Neighbours(board *Board) []Cell
	PassableNeighbours(board *Board) []Cell
	VoronoiOwner() string
	VoronoiDistance() int
}

type EmptyCell struct {
	coordinates     rules.Point
	voronoiOwner   string
	voronoiDist    int
}

func (e EmptyCell) VoronoiOwner() string {
	return e.voronoiOwner
}

func (e EmptyCell) VoronoiDistance() int {
	return e.voronoiDist
}

func (e EmptyCell) Kind() CellKind {
	return CellEmpty
}
func (e EmptyCell) IsPassable() bool {
	return true
}
func (e EmptyCell) Coordinates() rules.Point {
	return e.coordinates
}

type FoodCell struct {
	coordinates     rules.Point
	voronoiOwner   string
	voronoiDist    int
}

func (f FoodCell) VoronoiOwner() string {
	return f.voronoiOwner
}

func (f FoodCell) VoronoiDistance() int {
	return f.voronoiDist
}

func (f FoodCell) Kind() CellKind {
	return CellFood
}
func (f FoodCell) IsPassable() bool {
	return true
}
func (f FoodCell) Coordinates() rules.Point {
	return f.coordinates
}

// We can differentiate snake parts by a named type:
type SnakePartType int

const (
	SnakePartHead SnakePartType = iota
	SnakePartBody
	SnakePartTail
)

type SnakePartCell struct {
	coordinates        rules.Point
	SnakeID            string
	PartType           SnakePartType
	WillVanishNextTurn bool
	voronoiOwner      string
	voronoiDist       int
}

func (s SnakePartCell) VoronoiOwner() string {
	return s.voronoiOwner
}

func (s SnakePartCell) VoronoiDistance() int {
	return s.voronoiDist
}

func (s SnakePartCell) Coordinates() rules.Point {
	return s.coordinates
}
func (s SnakePartCell) Kind() CellKind {
	switch s.PartType {
	case SnakePartHead:
		return CellSnakeHead
	case SnakePartBody:
		return CellSnakeBody
	case SnakePartTail:
		return CellSnakeTail
	default:
		return CellEmpty // should never happen
	}
}

func (s SnakePartCell) IsPassable() bool {
	return s.PartType == SnakePartTail && s.WillVanishNextTurn
}

func getNeighbours(board *Board, pos rules.Point) []Cell {
	deltas := []rules.Point{{X: 0, Y: 1}, {X: 0, Y: -1}, {X: 1, Y: 0}, {X: -1, Y: 0}}
	neighbours := make([]Cell, 0, 4)

	for _, d := range deltas {
		newX, newY := pos.X+d.X, pos.Y+d.Y
		if newX >= 0 && newX < board.Width && newY >= 0 && newY < board.Height {
			neighbours = append(neighbours, board.Cells[newY][newX])
		}
	}
	return neighbours
}

func getPassableNeighbours(board *Board, pos rules.Point) []Cell {
	neighbours := getNeighbours(board, pos)
	passable := make([]Cell, 0, len(neighbours))
	for _, n := range neighbours {
		if n.IsPassable() {
			passable = append(passable, n)
		}
	}
	return passable
}

func (e EmptyCell) Neighbours(board *Board) []Cell {
	return getNeighbours(board, e.coordinates)
}

func (e EmptyCell) PassableNeighbours(board *Board) []Cell {
	return getPassableNeighbours(board, e.coordinates)
}

func (f FoodCell) Neighbours(board *Board) []Cell {
	return getNeighbours(board, f.coordinates)
}

func (f FoodCell) PassableNeighbours(board *Board) []Cell {
	return getPassableNeighbours(board, f.coordinates)
}

func (s SnakePartCell) Neighbours(board *Board) []Cell {
	return getNeighbours(board, s.coordinates)
}

func (s SnakePartCell) PassableNeighbours(board *Board) []Cell {
	return getPassableNeighbours(board, s.coordinates)
}

type Board struct {
	Width, Height int
	Cells         [][]Cell
}

func NewBoard(g GameSnapshot) *Board {
	board := &Board{
		Width:  g.Width(),
		Height: g.Height(),
		Cells:  make([][]Cell, g.Height()),
	}

	// Initialize cell slices
	for y := 0; y < g.Height(); y++ {
		board.Cells[y] = make([]Cell, g.Width())
	}

	// Place food
	for _, food := range g.Food() {
		if food.Y < g.Height() && food.X < g.Width() {
			board.Cells[food.Y][food.X] = FoodCell{coordinates: food}
		}
	}

	// Place snakes
	for _, snake := range g.AliveSnakes() {
		body := snake.Body()
		if len(body) == 0 {
			continue
		}

		// Head
		if body[0].Y < g.Height() && body[0].X < g.Width() {
			cell := SnakePartCell{
				coordinates:        body[0],
				SnakeID:            snake.ID(),
				PartType:           SnakePartHead,
				WillVanishNextTurn: false,
			}
			board.Cells[body[0].Y][body[0].X] = cell
		}

		// Body
		for i := 1; i < len(body)-1; i++ {
			if body[i].Y < g.Height() && body[i].X < g.Width() {
				cell := SnakePartCell{
					coordinates:        body[i],
					SnakeID:            snake.ID(),
					PartType:           SnakePartBody,
					WillVanishNextTurn: false,
				}
				board.Cells[body[i].Y][body[i].X] = cell
			}
		}

		// Tail
		if len(body) > 1 {
			tail := body[len(body)-1]
			if tail.Y < g.Height() && tail.X < g.Width() {
				cell := SnakePartCell{
					coordinates:        tail,
					SnakeID:            snake.ID(),
					PartType:           SnakePartTail,
					WillVanishNextTurn: snake.Health() < 100 && g.Turn() >= 3,
				}
				board.Cells[tail.Y][tail.X] = cell
			}
		}
	}

	// Fill remaining cells with EmptyCell
	for y := 0; y < g.Height(); y++ {
		for x := 0; x < g.Width(); x++ {
			if board.Cells[y][x] == nil {
				board.Cells[y][x] = EmptyCell{coordinates: rules.Point{X: x, Y: y}}
			}
		}
	}

	// Calculate Voronoi partitioning
	for y := 0; y < g.Height(); y++ {
		for x := 0; x < g.Width(); x++ {
			minDist := g.Height() * g.Width()
			var closestSnakeID string
			pos := rules.Point{X: x, Y: y}

			// Find closest snake head
			for _, snake := range g.AliveSnakes() {
				head := snake.Head()
				dist := abs(head.X - x) + abs(head.Y - y)
				if dist < minDist {
					minDist = dist
					closestSnakeID = snake.ID()
				}
			}

			// Update cell's Voronoi information
			switch cell := board.Cells[y][x].(type) {
			case EmptyCell:
				cell.voronoiOwner = closestSnakeID
				cell.voronoiDist = minDist
				board.Cells[y][x] = cell
			case FoodCell:
				cell.voronoiOwner = closestSnakeID
				cell.voronoiDist = minDist
				board.Cells[y][x] = cell
			case SnakePartCell:
				cell.voronoiOwner = closestSnakeID
				cell.voronoiDist = minDist
				board.Cells[y][x] = cell
			}
		}
	}

	return board
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
