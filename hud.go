package main

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	HUD_BAR_HEIGHT = 10
	HUD_FUEL_MAX   = 100
	HUD_HEALTH_MAX = 100
	HUD_ICON_SIZE  = 15
)

var (
	healthColor = color.RGBA{0xcf, 0xff, 0x10, 0xef}
	fuelColor   = color.RGBA{0x10, 0x10, 0x10, 0xef}
)

type HUD struct {
	game                         *Game
	fuelBarImage, healthBarImage *ebiten.Image
	fuelIcon, healthIcon         *ebiten.Image
	//health                       int
	fuel  int
	barY1 int
	barY2 int
	barX  int
	iconX int
}

func NewHUD(g *Game) *HUD {
	c := &HUD{}
	c.game = g

	//c.health = HUD_HEALTH_MAX
	c.fuel = HUD_FUEL_MAX
	c.setPositions()
	c.initIconImages()
	c.recalculateBarImages()
	return c
}

// func (c *HUD) ReduceHealthBar(amount int) {
// 	c.health -= amount
// 	if c.health < 0 {
// 		c.health = 0
// 	}
// }

// func (c *HUD) RefillHealthBar() {
// 	c.health = HUD_HEALTH_MAX
// }

func (c *HUD) initIconImages() {

	img, _, err := image.Decode(bytes.NewReader(IconsPng))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage := ebiten.NewImageFromImage(img)
	//c.images = SpriteCutter(ebitenImage, 100, 100, 5, 1)
	fuelIconCut := SubImage(ebitenImage, 300, 0, 100, 100)
	healthIconCut := SubImage(ebitenImage, 500, 0, 100, 100)
	c.healthIcon = ScaleImage(healthIconCut, HUD_ICON_SIZE, HUD_ICON_SIZE)
	c.fuelIcon = ScaleImage(fuelIconCut, HUD_ICON_SIZE, HUD_ICON_SIZE)

}

func (c *HUD) initImages() {
	c.fuelBarImage = ebiten.NewImage(HUD_BAR_HEIGHT, HUD_BAR_HEIGHT)
	c.healthBarImage = ebiten.NewImage(HUD_BAR_HEIGHT, HUD_BAR_HEIGHT)

	c.healthBarImage.Fill(healthColor)
	c.fuelBarImage.Fill(fuelColor)
}

func (c *HUD) recalculateBarImages() {

	// health bar
	healthW := c.game.health
	if healthW < 1 {
		healthW = 1
	}
	// fuel
	c.fuel = c.game.fuel
	if c.fuel < 1 {
		c.fuel = 1
	}
	c.healthBarImage = ebiten.NewImage(healthW, HUD_BAR_HEIGHT)
	c.healthBarImage.Fill(healthColor)

	c.fuelBarImage = ebiten.NewImage(c.fuel, HUD_BAR_HEIGHT)
	c.fuelBarImage.Fill(fuelColor)

}

func (c *HUD) Draw(screen *ebiten.Image) {

	// health icon
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.iconX), float64(c.barY1))

	screen.DrawImage(c.healthIcon, op)

	// fuel icon
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.iconX), float64(c.barY2))

	screen.DrawImage(c.fuelIcon, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.barX), float64(c.barY1))

	screen.DrawImage(c.healthBarImage, op)
	// fuel bar

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.barX), float64(c.barY2))

	screen.DrawImage(c.fuelBarImage, op)

}

func (c *HUD) setPositions() {
	c.barY1 = HUD_BAR_HEIGHT * 3
	c.barY2 = HUD_BAR_HEIGHT * 5
	c.barX = HUD_BAR_HEIGHT * 3
	c.iconX = HUD_BAR_HEIGHT

}

func (c *HUD) Update() error {
	var err error
	return err
}
