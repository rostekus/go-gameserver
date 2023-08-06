package game

import (
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

type Drawable interface {
	Draw() (int, int, rune, termbox.Attribute, termbox.Attribute)
	Move(x, y int)
}

type Snake struct {
	x, y int
}

type Fruit struct {
	x, y int
}

func (s *Snake) Draw() (int, int, rune, termbox.Attribute, termbox.Attribute) {

	return s.x, s.y, snakeBody, snakeFgColor, snakeBgColor
}

func (f *Fruit) Draw() (int, int, rune, termbox.Attribute, termbox.Attribute) {
	return f.x, f.y, fruitBody, fruitFgColor, fruitBgColor
}

func (f *Fruit) Move(x, y int) {
	f.x = x
	f.y = y
}

func (s *Snake) Move(x, y int) {
	s.x = x
	s.y = y
}
