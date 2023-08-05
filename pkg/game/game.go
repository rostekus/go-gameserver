package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

func writeText(x, y int, s string, fg, bg termbox.Attribute) {
	for i, ch := range s {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

type Game struct {
	sn    snake
	score int
	v     coord
	frt   fruit
	// Game field dimensions.
	fieldWidth, fieldHeight int
}

func NewGame(w, h int) Game {
	return Game{
		fieldWidth:  w,
		fieldHeight: h,
		sn:          newSnake(w, h),
		frt:         newFruit(w, h),
		v:           coord{1, 0},
		score:       0,
	}
}

func drawSnakePosition(g *Game) {
	str := fmt.Sprintf("(score : %d, %d) (%d, %d)", g.score, g.frt.pos.y, g.v.x, g.v.y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
}

func drawSnake(sn snake) {
	termbox.SetCell(sn.pos.x, sn.pos.y, snakeBody, snakeFgColor, snakeBgColor)
}

func drawFruit(frt fruit) {
	termbox.SetCell(frt.pos.x, frt.pos.y, fruitBody, fruitFgColor, snakeBgColor)
}

func (g *Game) draw() {
	termbox.Clear(snakeFgColor, snakeBgColor)
	drawSnakePosition(g)
	drawSnake(g.sn)
	if g.sn.pos.x == g.frt.pos.x && g.sn.pos.y == g.frt.pos.y {
		g.frt.pos.y = rand.Intn(10)
		g.frt.pos.x = rand.Intn(10)
		g.score++
	}
	drawFruit(g.frt)
	termbox.Flush()
}

func moveSnake(s snake, v coord, fw, fh int) snake {
	s.pos.x = s.pos.x + v.x
	s.pos.y = s.pos.y + v.y
	return s
}

func (g *Game) step() {
	if g.sn.pos.x+g.v.x < g.fieldWidth && g.sn.pos.x+g.v.x >= 0 {
		if g.sn.pos.y+g.v.y < g.fieldHeight && g.sn.pos.y+g.v.y > 0 {
			g.sn = moveSnake(g.sn, g.v, g.fieldWidth, g.fieldHeight)
		}
	}

	g.draw()
}

func (g *Game) moveLeft()  { g.v = coord{-1, 0} }
func (g *Game) moveRight() { g.v = coord{1, 0} }
func (g *Game) moveUp()    { g.v = coord{0, -1} }
func (g *Game) moveDown()  { g.v = coord{0, 1} }

func (g *Game) GetPosSnake() (int, int) {
	return g.sn.pos.x, g.sn.pos.y
}

func (g *Game) Start() {

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(70 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyArrowDown:
					g.moveDown()
				case termbox.KeyArrowUp:
					g.moveUp()
				case termbox.KeyArrowLeft:
					g.moveLeft()
				case termbox.KeyArrowRight:
					g.moveRight()
				case termbox.KeyEsc:
					return
				}
			}
		case <-ticker.C:
			g.step()
		}
	}
}
