package gol

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type World struct {
	height  int
	width   int
	cells   [][]int
	nextGen [][]int
}

var (
	BadWorld = errors.New("it's a small world")
)

func NewWorld(height int, width int) (*World, error) {
	if height == 0 || width == 0 {
		return nil, BadWorld
	}

	cells := make([][]int, height)
	nextGen := make([][]int, height)
	for row := range cells {
		cells[row] = make([]int, width)
		nextGen[row] = make([]int, width)
	}

	return &World{height: height, width: width, cells: cells, nextGen: nextGen}, nil
}

func (w *World) Height() int {
	return w.height
}

func (w *World) Width() int {
	return w.width
}

func (w *World) isAlive(row, col int) bool {
	return w.cells[row][col] == 1
}

func (w *World) neighbours(row, col int) int {
	if row < 0 || row > w.height || col < 0 || col > w.width {
		return 0
	}

	n := 0

	// top
	if row-1 >= 0 {
		n += w.cells[row-1][col]
	}

	// top-right
	if row-1 >= 0 && col+1 < w.width {
		n += w.cells[row-1][col+1]
	}

	// right
	if col+1 < w.width {
		n += w.cells[row][col+1]
	}

	// bottom-right
	if row+1 < w.height && col+1 < w.width {
		n += w.cells[row+1][col+1]
	}

	// bottom
	if row+1 < w.height {
		n += w.cells[row+1][col]
	}

	// bottom-left
	if row+1 < w.height && col-1 >= 0 {
		n += w.cells[row+1][col-1]
	}

	// left
	if col-1 >= 0 {
		n += w.cells[row][col-1]
	}

	// top-left
	if row-1 >= 0 && col-1 >= 0 {
		n += w.cells[row-1][col-1]
	}
	return n
}

func (w *World) mutate() {
	for row := 0; row < w.height; row++ {
		for col := 0; col < w.width; col++ {
			n := w.neighbours(row, col)
			if w.isAlive(row, col) {
				if n > 3 || n < 2 {
					// overpopulation or underpopulation
					w.nextGen[row][col] = 0
				} else if n == 2 || n == 3 {
					// ideal
					w.nextGen[row][col] = 1
				}
			} else {
				if n == 3 {
					// repopulation
					w.nextGen[row][col] = 1
				}
			}
		}
	}
}

func (w *World) Show() {
	for row := 0; row < w.height; row++ {
		for col := 0; col < w.width; col++ {
			if w.cells[row][col] == 1 {
				fmt.Printf("■")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	w.mutate()
	for row := 0; row < w.height; row++ {
		for col := 0; col < w.width; col++ {
			w.cells[row][col] = w.nextGen[row][col]
		}
	}
}

func (w *World) Load(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	validChars := 0
	for offset := 0; offset < len(data); offset++ {
		if data[offset] == '#' || data[offset] == '.' {
			validChars++
		}
	}
	expectedValidChars := w.height * w.width
	if validChars != expectedValidChars {
		return fmt.Errorf("invalid size of data (expected %v valid bytes, got %v bytes)", expectedValidChars, validChars)
	}

	positions := make([][]int, w.height)
	for row := range positions {
		positions[row] = make([]int, w.width)
	}
	var row, col int
	for offset := 0; offset < len(data); offset++ {
		if col > 0 && col%w.width == 0 {
			row++
			col = 0
		}
		switch data[offset] {
		case '#':
			positions[row][col] = 1
			col++
			validChars++
		case '.':
			positions[row][col] = 0
			col++
			validChars++
		case ' ', '\t', '\n':
			// skip and process the next line
		default:
			return fmt.Errorf("invalid byte at %v: '%v'", offset, data[offset])
		}

	}

	for row := range len(positions) {
		for col := range positions[row] {
			w.cells[row][col] = positions[row][col]
		}
	}
	return nil
}

func (w *World) LoadFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return w.Load(f)
}
