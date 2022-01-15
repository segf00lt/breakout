package main

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	WIN_W = 1200
	WIN_H = 800

	BRICK_W = 60
	BRICK_H = 20
	SPACE = 10

	PADDLE_W = 60
	PADDLE_H = 10
	PADDLE_X = WIN_W / 2 - PADDLE_W / 2
	PADDLE_Y = 100

	BALL_R = 8
	BALL_X = WIN_W / 2
	BALL_Y = 400
)

type Wall struct {
	line pixel.Line
}

func (wll *Wall) collision(b *Ball) bool {
	p := wll.IntersectCircle(b.circle)

	if p == pixel.ZV {
		return false
	}

	// calculate appropriate redirection
}

type Brick struct {
	rect pixel.Rect
	color color.Color
	id int
}

func (brk *Brick) draw(imd *imdraw.IMDraw) {
	imd.Color = brk.color
	imd.Push(brk.rect.Min, brk.rect.Max)
	imd.Rectangle(0)
}

type Paddle struct {
	rect pixel.Rect
	color color.Color
}

func (pdl *Paddle) draw(imd *imdraw.IMDraw) {
	imd.Color = pdl.color
	imd.Push(pdl.rect.Min, pdl.rect.Max)
	imd.Rectangle(0)
}

func (pdl *Paddle) move(w *pixelgl.Window) {
	if !w.MouseInsideWindow() {
		return
	}

	x := w.MousePosition().X 
	if x >= (WIN_W - PADDLE_W) {
		x = WIN_W - PADDLE_W
	}
	pdl.rect.Min.X = x
	pdl.rect.Max.X = x + PADDLE_W
}

type Ball struct {
	circle pixel.Circle
	color color.Color
	speed float64
	angle float64
}

func (bll *Ball) draw(imd *imdraw.IMDraw) {
	imd.Color = bll.color
	imd.Push(bll.circle.Center)
	imd.Circle(bll.circle.Radius, 0)
}

func (bll *Ball) move() {
	step_x := math.Cos(bll.angle)
	step_y := math.Sin(bll.angle)

	bll.circle.Center.X += step_x * bll.speed
	bll.circle.Center.Y += step_y * bll.speed
}

func (bll *Ball) accelerate() {
	if bll.speed <= 1.0 {
		bll.speed += 0.005
	}
}

func (bll *Ball) redirect(v pixel.Vec) {
}

func run() {
	// window configuration
	cfg := pixelgl.WindowConfig {
		Title: "Breakout",
		Bounds: pixel.R(0, 0, WIN_W, WIN_H),
		Resizable: false,
		Undecorated: true,
		VSync: false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetCursorVisible(false)

	// create brick rows
	n_row := 7
	n_col := WIN_W / (BRICK_W + SPACE)
	n_bricks := n_row * n_col

	bricks := make([]Brick, n_bricks, n_bricks) 

	initial_x := float64(5)
	initial_y := float64(WIN_H - 210 - SPACE)
	x_pos := initial_x
	y_pos := initial_y

	for i := 0; i < n_bricks; i++ {
		b := Brick {
			rect: pixel.R(x_pos, y_pos, x_pos + BRICK_W, y_pos + BRICK_H),
			color: pixel.RGB(1, 0, 0),
			id: i,
		}

		bricks[i] = b

		x_pos += BRICK_W + SPACE

		if (i + 1) % n_col == 0 {
			y_pos += BRICK_H + SPACE
			x_pos = initial_x
		}
	}
	
	// create paddle
	paddle := Paddle {
		rect: pixel.R(PADDLE_X, PADDLE_Y, PADDLE_X + PADDLE_W, PADDLE_Y + PADDLE_H),
		color: pixel.RGB(0, 0, 0),
	}

	// create ball
	ball := Ball {
		circle: pixel.C(pixel.V(BALL_X, BALL_Y), BALL_R),
		color: pixel.RGB(0.5, 1, 0.5),
		speed: 0.2,
		angle: -2.3,
	}

	// imd will draw all our shapes
	imd := imdraw.New(nil)

	for !win.Closed() {
		win.Clear(colornames.Skyblue)
		imd.Clear()

		paddle.move(win)
		ball.move()

		for _,b := range bricks {
			b.draw(imd)
		}

		paddle.draw(imd)
		ball.draw(imd)

		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
