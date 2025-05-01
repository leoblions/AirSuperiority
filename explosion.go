package main

import (
	"bytes"
	"image"
	"log"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	EXPLOSION_SPEED          = 3
	EXPLOSION_FRAMES_MAX     = 5
	EXPLOSION_IMAGES         = 6
	EXPLOSION_H              = 100
	EXPLOSION_W              = 100
	EXPLOSIONS_MAX           = 10
	EXPLOSION_OFFSET_X       = 50
	EXPLOSION_OFFSET_Y       = 1
	EXPLOSION_BORDER         = 100
	EXPLOSION_MIN_INTERVAL   = 500
	EXPLOSION_FRAME_INTERVAL = 70
)

var ()

type Explosion struct {
	game                         *Game
	imagesE1, imagesE2, imagesE3 []*ebiten.Image
	explosionUnits               [EXPLOSIONS_MAX]ExplosionUnit
	lastTimeMilli                int64
	lastTimeFrameMilli           int64
	testRect                     Movable
}

type ExplosionUnit struct {
	worldX, worldY int
	velX, velY     int
	frame, kind    int
	active         bool
}

func (punit *ExplosionUnit) Motion() {
	punit.worldX += punit.velX
	punit.worldY += punit.velY
}

func (punit *ExplosionUnit) Dimensions() (int, int, int, int) {
	return punit.worldX, punit.worldY, EXPLOSION_W, EXPLOSION_H
}

func NewExplosion(g *Game) *Explosion {
	c := &Explosion{}
	c.game = g
	c.lastTimeMilli = time.Now().UnixMilli()
	c.lastTimeFrameMilli = time.Now().UnixMilli()
	c.explosionUnits = [EXPLOSIONS_MAX]ExplosionUnit{}
	c.initImages()
	c.testRect = Movable{0, 0, 0, 0, EXPLOSION_W, EXPLOSION_H}

	return c
}

func (c *Explosion) initImages() {

	img, _, err := image.Decode(bytes.NewReader(Explosions1))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage := ebiten.NewImageFromImage(img)
	c.imagesE1 = SpriteCutter(ebitenImage, 150, 150, 4, 2)

	img, _, err = image.Decode(bytes.NewReader(Explosions2))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage = ebiten.NewImageFromImage(img)

	c.imagesE2 = SpriteCutter(ebitenImage, 150, 150, 4, 2)

	img, _, err = image.Decode(bytes.NewReader(Explosions3))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage = ebiten.NewImageFromImage(img)

	c.imagesE3 = SpriteCutter(ebitenImage, 100, 100, 6, 1)

}

func (c *Explosion) initImagesF() {
	SPRITESHEET1 := "explosions1.png"
	SPRITESHEET2 := "explosions2.png"
	path := filepath.Join(c.game.imageSubdir, SPRITESHEET1)
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	c.imagesE1 = SpriteCutter(img, 150, 150, 4, 2)

	path = filepath.Join(c.game.imageSubdir, SPRITESHEET2)
	img, _, err = ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	c.imagesE2 = SpriteCutter(img, 150, 150, 4, 2)

}

func (c *Explosion) Draw(screen *ebiten.Image) {
	for i := range EXPLOSIONS_MAX {

		var exp = c.explosionUnits[i]
		if exp.active {
			screenX, screenY := c.game.WorldToScreen(exp.worldX, exp.worldY)

			c.drawExplosion(screen, screenX, screenY, exp.frame, exp.kind)
		}

	}

}

func (c *Explosion) drawExplosion(screen *ebiten.Image, screenX, screenY, frame, kind int) {
	var image *ebiten.Image
	switch kind {
	case 0:
		image = c.imagesE1[frame]
	case 1:
		image = c.imagesE2[frame]
	case 2:
		image = c.imagesE3[frame]
	default:
		image = c.imagesE2[frame]
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate((float64)(screenX), (float64)(screenY))
	screen.DrawImage(image, op)

}

func (c *Explosion) projectileInBounds(punit *ExplosionUnit) bool {
	// true if unit in bounds
	if punit.worldX < -EXPLOSION_BORDER ||
		punit.worldX > EXPLOSION_BORDER+WINDOW_WIDTH ||
		punit.worldY < -EXPLOSION_BORDER ||
		punit.worldY > EXPLOSION_BORDER+WINDOW_HEIGHT {
		return false
	} else {
		return true
	}
}

func (c *Explosion) addExplosion(worldX, worldY, kind int) {
	var puArray = &[EXPLOSIONS_MAX]ExplosionUnit{}
	var velX, velY = 0, 0
	var worldXC, worldYC = worldX, worldY
	_ = velX
	_ = velY
	puArray = &c.explosionUnits
	var nowMilli = time.Now().UnixMilli()
	for i := range EXPLOSIONS_MAX {
		var limitReached = (nowMilli-c.lastTimeMilli > EXPLOSION_MIN_INTERVAL)
		if !puArray[i].active && limitReached {
			eunit := ExplosionUnit{}
			eunit.worldX = worldXC
			eunit.worldY = worldYC
			eunit.active = true
			eunit.frame = 0
			eunit.kind = kind
			puArray[i] = eunit
			c.lastTimeMilli = nowMilli
			c.game.sound.PlaySFX(kind)
			c.game.sound.StopSFX(4)
			break

		}
	}
}

func (c *Explosion) updateFrame(punit *ExplosionUnit) {
	if punit.frame < EXPLOSION_FRAMES_MAX {
		punit.frame += 1
	} else {
		punit.active = false
	}

}

func (c *Explosion) loopExplosions() {
	//check if time elapsed to change frame
	changeFrame := false
	nowMilli := time.Now().UnixMilli()
	limitReached := (nowMilli-c.lastTimeFrameMilli > EXPLOSION_FRAME_INTERVAL)
	if limitReached {
		changeFrame = true
		c.lastTimeFrameMilli = nowMilli

	}
	for i := range EXPLOSIONS_MAX {

		var punit = &c.explosionUnits[i]
		punit.Motion()
		if changeFrame {

			c.updateFrame(punit)
		}

	}

}

func (c *Explosion) Update() error {

	c.loopExplosions()
	var err error
	return err
}
