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
		frt:           NewRandomFruit(),
		v:             velocity{1, 0},
		score:         0,
		onlinePlayers: make(map[string]Drawable),
	}
	game.frt.AddFruit(w, h)
	return &game
}

func (g *Game) incScore() {
	g.score++
}
func (g *Game) draw() {
	termbox.Clear(snakeFgColor, snakeBgColor)
	headx, heady := g.sn.pos[0].x, g.sn.pos[0].y
	if g.frt.EatFruit(headx, heady) {
		g.incScore()
		g.sn.Eat()
		g.frt.AddFruit(g.fieldWidth, g.fieldHeight)
	}

	g.drawObject(g.sn)
	g.drawObject(g.frt)
	g.drawPlayers()
	termbox.Flush()
}
func (g *Game) drawObject(object Drawable) {
	for _, obj := range object.Draw() {
		termbox.SetCell(obj.x, obj.y, obj.body, obj.fg, obj.bg)
	}

}

func (g *Game) AddOnlinePlayer(pos []int32, username string) {
	g.onlinePlayers[username] = CreateOtherPlayerFromSlice(pos)
}

func (g *Game) drawPlayers() {
	for _, v := range g.onlinePlayers {
		g.drawObject(v)
	}
}

func (g *Game) moveSnake() {
	g.sn.v = g.v
	headx, heady := g.sn.pos[0].x, g.sn.pos[0].y
	if headx+g.v.x < g.fieldWidth && headx+g.v.x > 0 {
		if heady+g.v.y < g.fieldHeight && heady+g.v.y > 0 {
			g.sn.Move()
		}
	}
}

func (g *Game) step() {
	g.moveSnake()
	g.draw()
}

func (g *Game) GetPosSnake() []int32 {
	snakePos := make([]int32, 2*len(g.sn.pos))
	for i := 0; i < len(g.sn.pos); i += 2 {
		snakePos[i] = int32(g.sn.pos[i].x)
		snakePos[i+1] = int32(g.sn.pos[i].y)
	}
	return snakePos
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
