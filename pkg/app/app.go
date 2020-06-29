package app

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/yfedoruck/abc/pkg/fail"
	"math/rand"
	"time"
)

const (
	ScreenWidth  = 720.0  //590
	ScreenHeight = 1125.0 //960
	Scale        = 1      //960
	CubeWidth    = 50
	Delta        = 600
	SmallDelta   = 30
	SideDelta    = 100
)

type Game struct {
	fps        <-chan time.Time
	last       float64
	count      int
	frameNum   int
	Background *ebiten.Image
	Square     *Square
	screen     *ebiten.Image
	scale      float64
	field      Field
	tick       bool
	tx         int
	ty         int
	figure     TFig
	nextFig    TFig
	delta      float64
}

func NewGame() *Game {
	f := NewField()
	tf := NewFig(f, RandomNum())
	g := &Game{
		frameNum:   8,
		Background: LoadSprite("background.png"),
		Square:     NewSquare(), //LoadSprite("R.png"),
		scale:      Scale,       //0.379,
		field:      f,
		last:       tick(),
		tx:         f.area.Bounds().Min.X,
		ty:         f.area.Bounds().Min.Y,
		figure:     tf,
		delta:      Delta,
		nextFig:    NewFig(f, RandomNum()),
	}
	g.Fps()
	return g
}

func (r *Game) Fps() {
	r.fps = time.Tick(time.Second / 60)
}

func (r *Game) Update(screen *ebiten.Image) error {
	return nil
}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (r *Game) Draw(screen *ebiten.Image) {
	r.screen = screen
	r.DrawBg()
	r.DrawNextFigure()
	r.tickTack()
	if r.field.IsGameEnd() {
		r.Restart()
		return
	}
	r.DrawSquare()

	<-r.fps
}

func (r *Game) tickTack() {
	now := tick()
	if now-r.last > r.delta {
		if !r.tick {
			r.tick = true
		}
		r.last = now
	} else if r.tick {
		r.tick = false
	}
}

func (r *Game) tact(fn func(), delta float64) {
	now := tick()
	if now-r.last > delta {
		fn()
		r.last = now
	}
}

func (r *Game) DrawBg() {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(r.scale, r.scale)
	op.GeoM.Translate(0, 0)
	err := r.screen.DrawImage(r.Background, op)
	fail.Check(err)
}

func (r *Game) DrawSquare() {
	if r.figure.NotStopped() {
		r.listenMoving()
		r.FallDown()
		//r.listenRotate()
		//r.listenFall()
	} else {
		r.ResetDelta()
		r.SetNewFigure()
	}

	r.DrawFigure(r.figure)

	r.DrawWall()
}

func (r *Game) DrawWall() {
	for i := 0; i < r.field.NumX; i++ {
		for j := 0; j < r.field.NumY; j++ {
			if r.field.matrix[i][j] == true {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*CubeWidth+r.tx), float64(j*CubeWidth+r.ty))
				err := r.screen.DrawImage(r.Square.sprite, op)
				fail.Check(err)
			}
		}
	}
}

func (r Game) DrawFigure(figure TFig) {
	for _, point := range figure.a {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(point.x*CubeWidth+r.tx), float64(point.y*CubeWidth+r.ty))
		err := r.screen.DrawImage(r.Square.sprite, op)
		fail.Check(err)
	}
}

func (r *Game) SetNewFigure() {
	r.figure = r.nextFig
	r.nextFig = NewFig(r.field, RandomNum())
}

func (r Game) DrawNextFigure() {
	x := r.tx + 545.0
	y := r.ty + 115.0
	if r.nextFig.Type.IsI() {
		x = r.tx + 520.0
		y = r.ty + 125.0
	}
	for _, point := range r.nextFig.a {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(point.x*CubeWidth+x), float64(point.y*CubeWidth+y))
		err := r.screen.DrawImage(r.Square.sprite, op)
		fail.Check(err)
	}
}

func (r *Game) Restart() {
	r.field.Clear()
	r.SetNewFigure()
}

func RandomNum() Tetromino {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Tetromino(generator.Intn(7))
}

func (r *Game) listenMoving() {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		r.tact(r.figure.MoveRight, SideDelta)
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		r.tact(r.figure.MoveLeft, SideDelta)
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		r.figure.Rotate()
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		r.delta = SmallDelta
	default:
		r.ResetDelta()
	}
}

func (r *Game) listenRotate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		r.figure.Rotate()
	}
}

func (r *Game) listenFall() {
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		r.delta = SmallDelta
	} else {
		r.ResetDelta()
	}
}

func (r *Game) ResetDelta() {
	r.delta = Delta
}

func (r *Game) MoveRight() {
	r.figure.MoveRight()
}
func (r *Game) MoveLeft() {
	r.figure.MoveLeft()
}
func (r *Game) FallDown() {
	if !r.tick {
		return
	}
	r.figure.FallDown(r.field)
}

type Point struct {
	x, y int
}

func tick() float64 {
	return float64(time.Now().UnixNano() / 1e6)
}
