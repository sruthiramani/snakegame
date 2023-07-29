package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
)

type Coordinates struct {
	Row, Col int
}
type Snake struct {
	Locations []Coordinates
	Visited   map[Coordinates]bool
	mu        sync.Mutex
}

type SnakeGame struct {
	BoardSize  int
	TotalMoves int
	snake      *Snake
}

func (s *Snake) Size() int {
	return len(s.Locations)
}

func NewSnake(startCoords ...Coordinates) *Snake {
	snake := Snake{
		Locations: make([]Coordinates, 0),
		Visited:   make(map[Coordinates]bool),
	}
	if len(startCoords) > 0 {
		snake.Locations = append(snake.Locations, startCoords...)
		for _, c := range startCoords {
			snake.Visited[c] = true
		}
	}

	return &snake
}

func NewSnakeGame(boardSize int, start ...Coordinates) *SnakeGame {
	snake := NewSnake(start...)
	return &SnakeGame{
		BoardSize: boardSize,
		snake:     snake,
	}
}

type Direction int

const (
	Up    Direction = 0
	Down  Direction = 1
	Left  Direction = 2
	Right Direction = 3
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	}
	return ""
}
func (sg *SnakeGame) MoveSnake(d Direction) {

	var next Coordinates

	fmt.Println("Moving ", d.String())
	switch d {
	case Up:
		next = Coordinates{-1, 0}
	case Down:
		next = Coordinates{1, 0}
	case Right:
		next = Coordinates{0, 1}
	case Left:
		next = Coordinates{0, -1}

	}

	sg.MoveSnakeHead(next, false)
	sg.TotalMoves++
	if sg.TotalMoves%5 == 0 {
		sg.MoveSnakeHead(next, true)

	}

	fmt.Println("Current Score: ", sg.TotalMoves)

}

func (sg *SnakeGame) Print() {
	//fmt.Println("tail->", sg.snake.Locations, "->head")
	for i := 0; i < sg.BoardSize; i++ {
		for j := 0; j < sg.BoardSize; j++ {
			if sg.snake.Visited[Coordinates{Row: i, Col: j}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}

}

func (sg *SnakeGame) MoveSnakeHead(dir Coordinates, incrSize bool) {

	if incrSize {

		curHead := sg.snake.Locations[sg.snake.Size()-1]
		nextHead := Coordinates{
			Row: curHead.Row + dir.Row,
			Col: curHead.Col + dir.Col,
		}

		if sg.IsGameOver(nextHead) {
			fmt.Printf("Game Over! Your score is: %d", sg.TotalMoves)
			os.Exit(0)
		}
		sg.snake.Locations = append(sg.snake.Locations, nextHead)
		sg.snake.Visited[nextHead] = true
	} else {
		curHead := sg.snake.Locations[sg.snake.Size()-1]
		nextHead := Coordinates{
			Row: curHead.Row + dir.Row,
			Col: curHead.Col + dir.Col,
		}

		if sg.IsGameOver(nextHead) {
			fmt.Printf("Game Over! Your score is: %d", sg.TotalMoves)
			os.Exit(0)
		}

		tail := sg.snake.Locations[0]
		sg.snake.Locations = sg.snake.Locations[1:]
		sg.snake.Locations = append(sg.snake.Locations, nextHead)
		delete(sg.snake.Visited, tail)
		sg.snake.Visited[nextHead] = true

	}

}

func (sg *SnakeGame) IsGameOver(nextHead Coordinates) bool {

	if nextHead.Row >= sg.BoardSize || nextHead.Row < 0 || nextHead.Col >= sg.BoardSize || nextHead.Col < 0 {
		return true
	}

	if sg.snake.Visited[nextHead] {
		fmt.Println("loop")
		return true
	}

	return false

}

func getRandomDir() Direction {

	dir := rand.Intn(4)
	return Direction(dir)
}

func getRandomStart(size int) ([]Coordinates, Direction) {

	dir := getRandomDir()

	switch dir {
	case Up:
		startRow := rand.Intn(size-4) + 3
		startCol := rand.Intn(size)
		return []Coordinates{
			{
				Row: startRow,
				Col: startCol,
			},
			{
				Row: startRow - 1,
				Col: startCol,
			},
			{
				Row: startRow - 2,
				Col: startCol,
			},
		}, dir
	case Down:
		startRow := rand.Intn(size-4) + 3
		startCol := rand.Intn(size)
		return []Coordinates{
			{
				Row: startRow - 2,
				Col: startCol,
			},
			{
				Row: startRow - 1,
				Col: startCol,
			},
			{
				Row: startRow,
				Col: startCol,
			},
		}, dir
	case Right:
		startRow := rand.Intn(size)
		startCol := rand.Intn(size-4) + 3
		return []Coordinates{
			{
				Row: startRow,
				Col: startCol - 2,
			},
			{
				Row: startRow,
				Col: startCol - 1,
			},
			{
				Row: startRow,
				Col: startCol,
			},
		}, dir
	case Left:
		startRow := rand.Intn(size)
		startCol := rand.Intn(size-4) + 3
		return []Coordinates{
			{
				Row: startRow,
				Col: startCol,
			},
			{
				Row: startRow,
				Col: startCol - 1,
			},
			{
				Row: startRow,
				Col: startCol - 2,
			},
		}, dir

	}

	return []Coordinates{}, dir
}

func main() {

	start, dir := getRandomStart(10)

	fmt.Println("Choose a direction to move: Up(W), Down(S), Left(A), Right(D).\nRules:" +
		" Keep the Snake within the board. Snake grows in size with every 5 moves you make.")

	newGame := NewSnakeGame(10, start...)
	fmt.Println("Head facing ", dir)
	prev := dir
	newGame.Print()

	key := map[Direction]byte{
		Up:    'w',
		Down:  's',
		Left:  'a',
		Right: 'd',
	}

	direction := map[byte]Direction{
		'a': Left,
		'A': Left,
		'w': Up,
		'W': Up,
		'd': Right,
		'D': Right,
		's': Down,
		'S': Down,
	}

	next := key[dir]

	for {

		fmt.Scanf("%c", &next)
		nextDir, ok := direction[next]
		if !ok {
			nextDir = prev
		}
		switch nextDir {
		case Up:
			if prev == Down {
				fmt.Println("Can't move Down, while facing Up. Choose Up(W), Left(A) or Right(D)")
				continue
			}
			newGame.MoveSnake(nextDir)
			prev = nextDir
		case Down:
			if prev == Up {
				fmt.Println("Can't move Up, while facing Down. Choose Down(S), Left(A) or Right(D)")
				continue
			}
			newGame.MoveSnake(nextDir)
			prev = nextDir
		case Left:
			if prev == Right {
				fmt.Println("Can't move Left, while facing Right. Choose Right(D), Up(W) or Down(S)")
				continue
			}
			newGame.MoveSnake(nextDir)
			prev = nextDir
		case Right:
			if prev == Left {
				fmt.Println("Can't move Right, while facing Left. Choose Left(A), Up(W) or Down(S)")
				continue
			}
			newGame.MoveSnake(nextDir)
			prev = nextDir

		}
		newGame.Print()

	}

}
