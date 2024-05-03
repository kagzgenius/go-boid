package main

import (
	//"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	screenWidth  = 2000
	screenHeight = 2000

	boidWidth  = 20
	boidHeight = 20

	boidNum = 80 // Number of boids

	canvasImage *ebiten.Image
	boidImage   *ebiten.Image

	boids []boid // holds each boid

)

type Game struct {
	canvasImage *ebiten.Image
}

func init() {

	canvasImage = ebiten.NewImage(screenWidth, screenHeight)

	canvasImage.Fill(color.Black)

	boidImage = ebiten.NewImage(boidWidth, boidHeight)
	//
	// var R uint8 = uint8(rand.Intn(255))
	// var G uint8 = uint8(rand.Intn(255))
	// var B uint8 = uint8(rand.Intn(255))
	// var A uint8 = 1
	// boidImage.Fill(color.RGBA{R, G, B, A})

	boids = make([]boid, boidNum) // number of boids

	ebiten.SetFullscreen(true)

	// We will intailize the boids here

	for j := 0; j < len(boids)-1; j++ {

		boids[j].pos.x = rand.Intn(
			canvasImage.Bounds().Dx() + 1000,
		) // intailizea boid at random x-coordinate
		boids[j].pos.y = rand.Intn(
			canvasImage.Bounds().Dy() + 1000,
		) // intailizea boid at random y-coordinate

	}

}

type pos struct {
	x int //direction, pos
	y int //direction, pos
}

type boid struct {
	velocity pos // velcoity vector

	pos pos // position vector

	mag float64 // TBD

	steer float64 // used to steer angle of boid - TODO

}

func (g *Game) drawBoid(screen *ebiten.Image) {
	for j := 0; j < len(boids); j++ {

		var xPos = boids[j].pos.x
		var yPos = boids[j].pos.y

		ops := &ebiten.DrawImageOptions{}
		ops.GeoM.Translate(float64(xPos), float64(yPos))
		//	ops.GeoM.Rotate(float64(boids[0].steer))

		screen.DrawImage(boidImage, ops)

	}

}

func (g *Game) boidMove() {

	//var boidOne boid

	for j := 0; j < len(boids); j++ {

		var accreqThree pos = getFlockCentering(boids[j])
		var accreqOne pos = getCollisionAvoidance(boids[j])
		var accreqTwo pos = getVelocityMatching(boids[j])
		var accreqFour pos = boundPosition(boids[j])

		boids[j].velocity.x = boids[j].velocity.x + accreqOne.x + accreqTwo.x + accreqThree.x + accreqFour.x
		boids[j].velocity.y = boids[j].velocity.y + accreqOne.y + accreqTwo.y + accreqThree.y + accreqFour.y
		limitVelocity(boids[j])
		boids[j].pos.x = boids[j].pos.x + boids[j].velocity.x
		boids[j].pos.y = boids[j].pos.y + boids[j].velocity.y

	}
}

func boundPosition(b boid) pos {

	var vector pos
	if b.pos.x < 0 {
		vector.x = 10
	} else if b.pos.x > canvasImage.Bounds().Dx() {
		vector.x = -10
	}
	if b.pos.y < 0 {
		vector.y = 10
	} else if b.pos.y > canvasImage.Bounds().Dy() {
		vector.y = -10
	}
	return vector
}

func limitVelocity(b boid) {
	var vlim int = 20
	/* var vector pos */

	var dist = int(math.Sqrt(math.Pow(float64(b.pos.x), 2) + math.Pow(float64(b.pos.y), 2)))

	if dist > vlim {
		b.velocity.x = (b.velocity.x / dist) * vlim
		b.velocity.y = (b.velocity.y / dist) * vlim
	}
}

func getDistance(b1 boid, b2 boid) int {
	var val1 = b1.pos.x - b2.pos.x
	var val2 = b1.pos.y - b2.pos.y

	var dist = math.Sqrt(math.Pow(float64(val1), 2) + math.Pow(float64(val2), 2))
	return int(dist)
}

func getCollisionAvoidance(b boid) pos { // collision avoidance procedure

	var c pos
	c.x = 0
	c.y = 0

	for j := 0; j < len(boids); j++ {
		if boids[j] != b {
			var bj = boids[j]
			if v := getDistance(bj, b); v < 15 {

				var val1 int = boids[j].pos.x - b.pos.x
				var val2 int = boids[j].pos.y - b.pos.y

				c.x = c.x - val1
				c.y = c.y - val2
			}
		}
	}
	return c
}

func getVelocityMatching(b boid) pos { // velocity matching procedure
	var vector pos
	for j := 0; j < len(boids); j++ {
		if boids[j] != b {
			vector.x = vector.x + boids[j].velocity.x
			vector.y = vector.y + boids[j].velocity.y
		}
	}

	vector.x = vector.x/len(boids) - 1
	vector.y = vector.y/len(boids) - 1

	vector.x = (vector.x - b.velocity.x) / 8
	vector.y = (vector.y - b.velocity.y) / 8

	return vector
}

func getFlockCentering(b boid) pos { //Flock centering procedure

	var vector pos

	for j := 0; j < len(boids); j++ {
		if boids[j] != b {
			vector.x = vector.x + boids[j].pos.x
			vector.y = vector.y + boids[j].pos.y

		}
	}
	vector.x = vector.x / (len(boids) - 1)
	vector.y = vector.y / (len(boids) - 1)

	vector.x = (vector.x - b.pos.x) / 100
	vector.y = (vector.y - b.pos.y) / 100

	return vector

}

func (g *Game) Update() error {

	var R uint8 = uint8(150)
	var G uint8 = uint8(150)
	var B uint8 = uint8(255)
	var A uint8 = 1
	boidImage.Fill(color.RGBA{R, G, B, A})
	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {

	g.drawBoid(screen)
	g.boidMove()

	//screen.DrawImage(boidImage, nil)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	return screenWidth, screenHeight
}

func main() {
	game := &Game{}

	/* 	ebiten.SetTPS(30) */
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids - A bird simulation")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
