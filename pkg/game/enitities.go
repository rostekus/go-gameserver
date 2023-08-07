package game

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

const (
	snakeBody      = '*'
	snakeFgColor   = termbox.ColorRed
	snakeBgColor   = termbox.ColorDefault
	fruitBody      = '@'
	fruitFgColor   = termbox.ColorGreen
	fruitBgColor   = termbox.ColorWhite
	openentBody    = 'X'
	oponentFgColor = termbox.ColorBlue
	oponentBgColor = termbox.ColorYellow
)

type DrawableElemenet struct {
	x, y   int
	body   rune
	fg, bg termbox.Attribute
}

type Drawable interface {
	Draw() []DrawableElemenet
}

type Snake struct {
	pos []DrawableElemenet
	v   velocity
}

type Fruit struct {
	pos []DrawableElemenet
}

type OtherPlayer struct {
	pos []DrawableElemenet
}

func CreateOtherPlayerFromSlice(arr []int32) *OtherPlayer {
	lenght := len(arr)
	if lenght <= 2 {
		lenght = 1
	} else {
		lenght = len(arr) / 2
	}

	pos := make([]DrawableElemenet, lenght)
	for i := 0; i < lenght; i++ {
		pos[i] = DrawableElemenet{
			int(arr[i]), int(arr[i+1]), openentBody, oponentFgColor, oponentBgColor,
		}
	}
	return &OtherPlayer{pos: pos}
}

func NewRandomSnake(maxX, maxY int) *Snake {
	head := DrawableElemenet{
		rand.Intn(maxX), rand.Intn(maxY),
		snakeBody, snakeFgColor, snakeBgColor,
	}
	snk := &Snake{}
	snk.pos = append([]DrawableElemenet{}, head)
	return snk
}
func NewRandomFruit() *Fruit {

	fr := &Fruit{}
	fr.pos = make([]DrawableElemenet, 2)
	return fr
}

func (f *Fruit) AddFruit(maxX, maxY int) {

	fruit := DrawableElemenet{
		rand.Intn(maxX), rand.Intn(maxY),
		fruitBody, fruitFgColor, fruitBgColor,
	}
	f.pos = append(f.pos, fruit)
}

func (o *OtherPlayer) Draw() []DrawableElemenet {
	return o.pos
}
func (s *Snake) Draw() []DrawableElemenet {
	return s.pos
}

func (f *Fruit) Draw() []DrawableElemenet {
	return f.pos
}

func (f *Fruit) EatFruit(x, y int) bool {
	for i, fruit := range f.pos {
		if fruit.x == x && fruit.y == y {
			f.pos = append(f.pos[:i], f.pos[:+1]...)
			return true
		}
	}
	return false
}

// we move only head
// all other elements should move by 1
func (s *Snake) Move() {
	for i := len(s.pos) - 1; i >= 1; i-- {
		s.pos[i].x = s.pos[i-1].x
		s.pos[i].y = s.pos[i-1].y
	}
	// snake always has first element
	s.pos[0].x += s.v.x
	s.pos[0].y += +s.v.y
}
func (s *Snake) Eat() {
	head := s.pos[0]

	newX := head.x + s.v.x
	newY := head.y + s.v.y
	newSegment := DrawableElemenet{
		x:    newX,
		y:    newY,
		body: head.body,
		fg:   head.fg,
		bg:   head.bg,
	}

	s.pos = append(s.pos, newSegment)
}
