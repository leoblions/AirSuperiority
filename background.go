package main

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math/rand/v2"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	BACKGROUND_BG_COLOR        = color.RGBA{0x10, 0x10, 0xef, 0xff}
	BACKGROUND_Y_MAX           = (int)(WINDOW_HEIGHT * 1.5)
	BACKGROUND_Y_HALF          = (int)(WINDOW_HEIGHT * 0.5)
	BACKGROUND_CLOUD_MAX_SPEED = 6
	BACKGROUND_SKY_SPEED       = 10
	BACKGROUND_TICK_MAX        = 60
)

const (
	BACKGROUND_CLOUD_OFFSET_MAX = WINDOW_WIDTH / 2
	BACKGROUND_CLOUD_AMOUNT     = 7
	BACKGROUND_SPAWN_BUFFER     = 100
)

type Sprite struct {
	image         *ebiten.Image
	x, y          float64
	width, height int
	speed         float64
}

type Background struct {
	game                 *Game
	imageOcean           *Sprite
	cloudSprites         [BACKGROUND_CLOUD_AMOUNT]*Sprite
	imageC               string
	imageW               string
	backgroundHeight     int
	backgroundHeightHalf int
	backgroundStart      int
	backgroundStartC     int
	backgroundEnd        int
	backgroundEndC       int
	backgroundY1         int
	backgroundY2         int
	//backgroundYC                          int
	backgroundX            int
	cloudStartY, cloudEndY float64
	cloudX                 int
	oceanSpeed             int
	ticks                  int
	//cloudSpeed  int

	Movable
}

func NewBackground(g *Game) *Background {
	c := &Background{}
	c.game = g
	c.imageC = "clouds.png"
	c.imageW = "oceanBlank.png"
	c.cloudSprites = [BACKGROUND_CLOUD_AMOUNT]*Sprite{}
	//c.cloudSpeed = BACKGROUND_SKY_SPEED
	c.oceanSpeed = 1
	c.initImagesClouds()
	c.initSprites()
	c.backgroundY1 = c.backgroundStart
	c.backgroundY2 = 0
	c.cloudEndY = WINDOW_HEIGHT + BACKGROUND_SPAWN_BUFFER
	c.cloudStartY = 0 - BACKGROUND_SPAWN_BUFFER
	//c.backgroundYC = 0
	//c.cloudSprites[0].y = 0
	return c
}

func (c *Background) initImagesF() {
	var err error
	path := filepath.Join(c.game.imageSubdir, c.imageC)
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	c.cloudSprites[0].image = img
	c.backgroundStartC = -c.cloudSprites[0].image.Bounds().Dy()
	c.backgroundEndC = c.cloudSprites[0].image.Bounds().Dy()

	path = filepath.Join(c.game.imageSubdir, c.imageW)
	img, _, err = ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	c.imageOcean.image = CropImage(img, WINDOW_WIDTH, WINDOW_HEIGHT)

	c.backgroundHeight = c.imageOcean.image.Bounds().Dy()

	c.backgroundHeightHalf = c.backgroundHeight / 2
	c.backgroundStart = -c.backgroundHeight
	c.backgroundEnd = c.backgroundHeight
}

func (c *Background) initImagesClouds() {
	var err error
	// get sprite sheet from bytes array
	img1, _, err := image.Decode(bytes.NewReader(ImageClouds1))
	if err != nil {
		log.Fatal(err)
	}
	// convert image to ebiten image format

	imageCloud01 := ebiten.NewImageFromImage(img1)

	img2, _, err := image.Decode(bytes.NewReader(ImageClouds2))
	if err != nil {
		log.Fatal(err)
	}
	// convert image to ebiten image format
	imageCloud02 := ebiten.NewImageFromImage(img2)
	if nil == &c.cloudSprites {
		c.cloudSprites = [BACKGROUND_CLOUD_AMOUNT]*Sprite{}
	}
	cloudSpeeds := []float64{0.1, 0.2, 0.3, 0.4}
	for i := range BACKGROUND_CLOUD_AMOUNT {
		c.cloudSprites[i] = &Sprite{}
		c.cloudSprites[i].x = float64(rand.IntN(WINDOW_WIDTH))
		c.cloudSprites[i].y = float64(rand.IntN(WINDOW_HEIGHT))
		c.cloudSprites[i].speed = cloudSpeeds[rand.IntN(4)]
	}
	// sheet 1
	c.cloudSprites[0].image = SubImage(imageCloud01, 0, 0, 150, 150)
	c.cloudSprites[1].image = SubImage(imageCloud01, 150, 0, 150, 150)
	c.cloudSprites[2].image = SubImage(imageCloud01, 0, 150, 50, 50)
	c.cloudSprites[3].image = SubImage(imageCloud01, 50, 150, 200, 150)
	// sheet 2
	c.cloudSprites[4].image = SubImage(imageCloud02, 0, 0, 100, 100)
	c.cloudSprites[5].image = SubImage(imageCloud02, 100, 0, 100, 100)
	c.cloudSprites[6].image = SubImage(imageCloud02, 0, 100, 200, 100)

}

func (c *Background) initSprites() {
	var err error
	// get sprite sheet from bytes array
	// img, _, err := image.Decode(bytes.NewReader(ImageClouds))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // convert image to ebiten image format
	// c.cloudSprites[0] = &Sprite{}
	// c.cloudSprites[0].image = ebiten.NewImageFromImage(img)

	// c.backgroundStartC = -c.cloudSprites[0].image.Bounds().Dy()
	// c.backgroundEndC = c.cloudSprites[0].image.Bounds().Dy() * 2

	img, _, err := image.Decode(bytes.NewReader(ImageOcean))
	if err != nil {
		log.Fatal(err)
	}
	// convert image to ebiten image format
	c.imageOcean = &Sprite{}
	c.imageOcean.image = ebiten.NewImageFromImage(img)

	c.imageOcean.image = CropImage(c.imageOcean.image, WINDOW_WIDTH, WINDOW_HEIGHT)

	c.backgroundHeight = c.imageOcean.image.Bounds().Dy()

	c.backgroundHeightHalf = c.backgroundHeight / 2
	c.backgroundStart = -c.backgroundHeight
	c.backgroundEnd = c.backgroundHeight
}

func (c *Background) Draw(screen *ebiten.Image) {
	screen.Fill(BACKGROUND_BG_COLOR)
	DrawImageAt(c.imageOcean.image, screen, 0, 0)
	DrawImageAt(c.imageOcean.image, screen, 0, c.backgroundY1)
	DrawImageAt(c.imageOcean.image, screen, 0, c.backgroundY2)

	for i := range c.cloudSprites {
		DrawImageAtF(c.cloudSprites[i].image, screen, c.cloudSprites[i].x, c.cloudSprites[i].y)
	}

}

func (c *Background) Update() error {
	var err error
	// bgspeed := BACKGROUND_SKY_SPEED
	if c.ticks < BACKGROUND_TICK_MAX {
		c.ticks += 1
	} else {
		c.ticks = 0
	}

	// if c.backgroundY1 > c.backgroundEnd {
	// 	c.backgroundY1 = c.backgroundStart
	// } else if c.ticks%bgspeed == 0 {
	// 	c.backgroundY1 += 1
	// }

	// if c.backgroundY2 > c.backgroundEnd {
	// 	c.backgroundY2 = c.backgroundStart
	// } else {
	// 	c.backgroundY2 += 1
	// }

	for i := range c.cloudSprites {
		if c.cloudSprites[i].y > c.cloudEndY {
			c.cloudSprites[i].y = float64(c.cloudStartY)
			c.cloudSprites[i].x = float64(rand.IntN(WINDOW_WIDTH - 100))
			c.cloudSprites[i].speed = rand.Float64()
		} else {
			c.cloudSprites[i].y += c.cloudSprites[i].speed

		}
	}

	// if c.cloudSprites[0].y > c.backgroundEndC {
	// 	c.cloudSprites[0].y = c.backgroundStartC
	// 	c.cloudX = rand.IntN(BACKGROUND_CLOUD_OFFSET_MAX)
	// } else {
	// 	c.cloudSprites[0].y += c.cloudSpeed
	// }

	// if c.cloudSprites[1].y > c.backgroundEndC {
	// 	c.cloudSprites[1].y = c.backgroundStartC
	// 	c.cloudX = rand.IntN(BACKGROUND_CLOUD_OFFSET_MAX)
	// } else {
	// 	c.cloudSprites[1].y += c.cloudSpeed
	// }

	return err
}
