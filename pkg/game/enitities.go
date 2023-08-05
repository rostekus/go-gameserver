package game

import (
	"log"
	"math/rand"

	"github.com/nsf/termbox-go"
)

const (
	snakeBody    = '*'
	snakeFgColor = termbox.ColorRed
	snakeBgColor = termbox.ColorDefault
	fruitBody    = '@'
	fruitFgColor = termbox.ColorRed
	fruitBgColor = termbox.ColorDefault
)

type coord struct {
	x, y int
}

type snake struct {
	pos coord
}

type fruit struct {
	pos coord
}

func newSnake(maxX, maxY int) snake {
	return snake{coord{rand.Intn(maxX), rand.Intn(maxY)}}
}

func newFruit(maxX, maxY int) fruit {
	log.Println(" fruit ", rand.Intn(maxX))
	return fruit{coord{rand.Intn(maxX), rand.Intn(maxY)}}
}
