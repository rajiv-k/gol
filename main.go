package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rajiv-k/gol/internal/gol"
)

const (
	generationShowTimeMillis = 300
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: ./gol <path/to/game.txt>")
		os.Exit(1)
	}

	gameFile := os.Args[1]

	world, err := gol.NewWorld(20, 20)
	if err != nil {
		log.Fatalf("could not create new game: %v", err)
	}

	if err := world.LoadFromFile(gameFile); err != nil {
		log.Fatalf("could not spawn: %v", err)
	}

	for i := 0; i < 50; i++ {
		log.Printf("Game of Life | Generation: %v\n", i)
		world.Show()
		time.Sleep(time.Millisecond * generationShowTimeMillis)
		for j := 0; j < world.Height()+1; j++ {
			fmt.Printf("\033[2K\033M")
		}
	}
}
