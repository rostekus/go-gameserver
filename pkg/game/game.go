package game

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

func writeText(x, y int, s string, fg, bg termbox.Attribute) {
	for i, ch := range s {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

type velocity struct {
	x, y int
}

type Game struct {
	sn            *Snake
	score         int
	v             velocity
	frt           *Fruit
	fieldWidth    int
	fieldHeight   int
	onlinePlayers map[string]Drawable
}

func NewGame(w, h int) *Game {
	game := Game{
		fieldWidth:    w,
		fieldHeight:   h,
		sn:            NewRandomSnake(w, h),
		frt:           NewRandomFruit(w, h),
		v:             velocity{1, 0},
		score:         0,
		onlinePlayers: make(map[string]Drawable),
	}
	game.onlinePlayers["player"] = NewRandomFruit(w, h)
	return &game
}

func NewRandomSnake(maxX, maxY int) *Snake {
	return &Snake{rand.Intn(maxX), rand.Intn(maxY)}
}

func NewRandomFruit(maxX, maxY int) *Fruit {
	return &Fruit{rand.Intn(maxX), rand.Intn(maxY)}
}
func (g *Game) isFruitEaten() bool {
	if g.sn.x == g.frt.x && g.sn.y == g.frt.y {
		return true
	}
	return false
}
func (g *Game) incScore() {
	g.score++
}
func (g *Game) newFruit() {
	g.frt.y = rand.Intn(g.fieldHeight)
	g.frt.x = rand.Intn(g.fieldWidth)
}
func (g *Game) draw() {
	termbox.Clear(snakeFgColor, snakeBgColor)
	if g.isFruitEaten() {
		g.incScore()
		g.newFruit()
	}

	g.sn.Draw()
	g.drawObject(g.sn)
	g.drawObject(g.frt)
	g.drawPlayers()
	termbox.Flush()
}
func (g *Game) drawObject(object Drawable) {
	termbox.SetCell(object.Draw())

}

func (g *Game) AddOnlinePlayer(pos []int32, username string) {
	g.onlinePlayers[username] = &Fruit{int(pos[0]), int(pos[1])}
}

func (g *Game) drawPlayers() {
	for _, v := range g.onlinePlayers {
		g.drawObject(v)
	}
}

func (g *Game) moveSnake() {
	if g.sn.x+g.v.x < g.fieldWidth && g.sn.x+g.v.x > 0 {
		if g.sn.y+g.v.y < g.fieldHeight && g.sn.y+g.v.y > 0 {
			x := g.sn.x + g.v.x
			y := g.sn.y + g.v.y
			g.sn.Move(x, y)
		}
	}
}

func (g *Game) step() {
	g.moveSnake()
	g.draw()
}

func (g *Game) GetPosSnake() (int, int) {
	return g.sn.x, g.sn.y
}

func (g *Game) moveLeft()  { g.v = velocity{-1, 0} }
func (g *Game) moveRight() { g.v = velocity{1, 0} }
func (g *Game) moveUp()    { g.v = velocity{0, -1} }
func (g *Game) moveDown()  { g.v = velocity{0, 1} }

func (g *Game) Start() {
	rand.Seed(time.Now().UnixNano())
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
