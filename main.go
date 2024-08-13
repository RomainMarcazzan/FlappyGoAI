package main

import (
	"image/color"
	"log"
	"strconv"

	"github.com/RomainMarcazzan/flappy_go_ai/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth     = 640
	screenHeight    = 480
	birdXPosition   = 50
	birdSize        = 20
	gravity         = 0.3
	jumpImpulse     = -5
	maxFallingSpeed = 10
	pipeWidth       = birdSize * 2
	openingGapMin   = 100
	openingGapMax   = 200
	gapBetweenPipes = 200
	pipeSpeed       = 2
)

type Sprite struct {
	X, Y, W, H float32
}

type Bird struct {
	Sprite
	velocity float32
}

type Pipe struct {
	top    Sprite
	bottom Sprite
}

type Game struct {
	bird      Bird
	pipes     []Pipe
	score     int
	lastPipeX float32
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Score: "+strconv.Itoa(g.score))
	vector.DrawFilledRect(screen, float32(g.bird.X), float32(g.bird.Y), float32(g.bird.W), float32(g.bird.H), color.White, false)

	for _, pipe := range g.pipes {
		vector.DrawFilledRect(screen, float32(pipe.top.X), float32(pipe.top.Y), float32(pipe.top.W), float32(pipe.top.H), color.White, false)
		vector.DrawFilledRect(screen, float32(pipe.bottom.X), float32(pipe.bottom.Y), float32(pipe.bottom.W), float32(pipe.bottom.H), color.White, false)
	}
}

func (g *Game) Update() error {
	g.bird.ApplyPhysics()
	g.bird.HandleInput()
	g.bird.Contain(g)

	for i := range g.pipes {
		g.pipes[i].Slide()
		if g.pipes[i].top.X+pipeWidth == birdXPosition {
			g.score++
		}
	}

	if g.bird.CollidesWith(g.pipes) {
		g.Reset()
	}

	g.removeOffScreenPipes()

	if len(g.pipes) == 0 || g.pipes[len(g.pipes)-1].top.X < screenWidth-gapBetweenPipes {
		g.addPipe()
	}

	return nil
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Flappy GoAi")

	bird := Bird{
		Sprite: Sprite{
			X: birdXPosition,
			Y: (screenHeight / 2) - (birdSize / 2),
			W: birdSize,
			H: birdSize,
		},
		velocity: 0,
	}

	game := &Game{
		bird:      bird,
		pipes:     []Pipe{},
		lastPipeX: screenWidth,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) removeOffScreenPipes() {
	var newPipes []Pipe
	for _, pipe := range g.pipes {
		if pipe.top.X+pipeWidth > 0 {
			newPipes = append(newPipes, pipe)
		}
	}
	g.pipes = newPipes
}

func (g *Game) addPipe() {
	g.lastPipeX += pipeWidth

	openingGap := utils.RandomRange(openingGapMin, openingGapMax)
	topPipeHeight := utils.RandomRange(0, screenHeight-openingGap)
	bottomPipeHeight := screenHeight - topPipeHeight - openingGap

	g.pipes = append(g.pipes, Pipe{
		top: Sprite{
			X: g.lastPipeX,
			Y: 0,
			W: pipeWidth,
			H: topPipeHeight,
		},
		bottom: Sprite{
			X: g.lastPipeX,
			Y: screenHeight - bottomPipeHeight,
			W: pipeWidth,
			H: bottomPipeHeight,
		},
	})
}

func (b *Bird) HandleInput() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		b.velocity = jumpImpulse
	}
}

func (b *Bird) ApplyPhysics() {
	b.velocity += gravity
	if b.velocity > maxFallingSpeed {
		b.velocity = maxFallingSpeed
	}

	b.Y += b.velocity
}

func (b *Bird) Contain(game *Game) {
	if b.Y <= 0 {
		b.Y = 0
		b.velocity = 0
	}

	if b.Y+birdSize >= screenHeight {
		/* b.Y = screenHeight - birdSize
		b.velocity = 0 */
		game.Reset()
	}
}

func (b *Bird) CollidesWith(pipes []Pipe) bool {
	for _, pipe := range pipes {
		if b.X < pipe.top.X+pipe.top.W && b.X+b.W > pipe.top.X &&
			(b.Y < pipe.top.Y+pipe.top.H || b.Y+b.H > pipe.bottom.Y) {
			return true
		}
	}
	return false
}

func (g *Game) Reset() {
	g.bird.Y = (screenHeight / 2) - (birdSize / 2)
	g.bird.velocity = 0
	g.pipes = []Pipe{}
	g.lastPipeX = screenWidth
	g.score = 0
}

func (p *Pipe) Slide() {
	p.top.X -= pipeSpeed
	p.bottom.X -= pipeSpeed
}
