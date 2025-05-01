package main

import (
	"bytes"
	"image"
	"log"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ENTITY_SPEED                = 2
	ENTITY_KINDS                = 18
	ENTITY_H                    = 100
	ENTITY_W                    = 100
	ENTITY_SCALE                = 0.5
	ENTITYS_MAX                 = 10
	ENTITY_OFFSET_X             = 50
	ENTITY_OFFSET_Y             = 1
	ENTITY_BORDER               = 300
	ENTITY_MIN_INTERVAL         = 2000
	ENTITY_RAND_INTERVAL_MAX    = 2000
	ENTITY_START_Y              = -300
	ENTITY_LOOT_DROP_MAX        = 3
	ENTITY_START_X_MAX          = WINDOW_WIDTH - ENTITY_W
	DIFFICULTY_SPAWN_SPEED_STEP = 200
	//ENEMY_PROJECTILE_SPEED   = 2
)

var (
// projectileColorE = color.RGBA{0xff, 0xff, 0x30, 0xff}
// projectileColorP = color.RGBA{0xe0, 0xe0, 0x6f, 0xff}
)

type Entity struct {
	game                *Game
	images              [ENTITY_KINDS]ebiten.Image
	entityUnits         [ENTITYS_MAX]EntityUnit
	lastTimeMilli       int64
	entitySpawnInterval int64
	enemyFirePositionY  int
}

type EntityUnit struct {
	worldX, worldY, kind int
	velX, velY           int
	active, fired        bool
	Movable
}

func (punit *EntityUnit) Motion() {
	punit.worldX += punit.velX
	punit.worldY += punit.velY
}

func (c *Entity) FireProjectile(eunit *EntityUnit) {
	if !eunit.fired && c.enemyFirePositionY-eunit.worldY < 3 {
		// if difficulty is low, abort more often
		randNum := rand.IntN(9)
		if randNum > c.game.difficulty {
			eunit.fired = true
			return
		}
		projectileY := eunit.worldY + eunit.height
		var projectileUnit = c.game.projectile.addEnemyProjectile(eunit.worldX, projectileY)
		if nil != projectileUnit {

			dx := c.game.player.worldX - eunit.worldX
			// dy := c.game.player.worldY - eunit.worldY
			// if dy == 0 {
			// 	dy = 1
			// }
			projectileUnit.velY = eunit.velY + 1
			if dx > 50 {

				projectileUnit.velX = projectileUnit.velY
			} else if dx < -50 {
				projectileUnit.velX = -projectileUnit.velY
			} else {
				projectileUnit.velX = 0
			}
			eunit.fired = true
		}

	}
}

func (c *EntityUnit) Dimensions() (int, int, int, int) {
	return c.worldX, c.worldY, c.width, c.height
	//return 5, 5, 5, 5
}

func NewEntity(g *Game) *Entity {
	c := &Entity{}
	c.game = g
	c.enemyFirePositionY = 100
	c.lastTimeMilli = time.Now().UnixMilli()
	c.entitySpawnInterval = ENTITY_MIN_INTERVAL + int64(DIFFICULTY_SPAWN_SPEED_STEP*c.game.difficulty)
	c.entityUnits = [ENTITYS_MAX]EntityUnit{}
	c.initImages()
	//c.entityUnits[0] = EntityUnit{200, 200, 4, 0, 1, true}
	c.addEntity(100, 100, 2)
	return c
}

func (c *Entity) initImages() {
	var imgI image.Image
	var img *ebiten.Image

	// SPRITESHEET1 := "jets.png"
	// SPRITESHEET2 := "airplanes2.png"
	// SPRITESHEET3 := "airplanes3.png"
	// path := filepath.Join(c.game.imageSubdir, SPRITESHEET1)
	// img, _, err := ebitenutil.NewImageFromFile(path)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	imgI, _, err := image.Decode(bytes.NewReader(Airplanes1))
	if err != nil {
		log.Fatal(err)
	}
	img = ebiten.NewImageFromImage(imgI)
	// jets
	//first row
	c.images[0] = *SubImage(img, 0, 0, 200, 200)
	c.images[1] = *SubImage(img, 200, 0, 200, 200)
	c.images[2] = *SubImage(img, 400, 0, 200, 250)
	//second row
	c.images[3] = *SubImage(img, 0, 250, 200, 250)
	c.images[4] = *SubImage(img, 200, 250, 200, 250)
	c.images[5] = *SubImage(img, 400, 250, 200, 300)

	// white airplanes
	// path = filepath.Join(c.game.imageSubdir, SPRITESHEET2)
	// img, _, err = ebitenutil.NewImageFromFile(path)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	imgI, _, err = image.Decode(bytes.NewReader(Airplanes2))
	if err != nil {
		log.Fatal(err)
	}
	img = ebiten.NewImageFromImage(imgI)
	//first row
	c.images[6] = *SubImage(img, 0, 0, 200, 250)
	c.images[7] = *SubImage(img, 200, 0, 200, 250)
	c.images[8] = *SubImage(img, 400, 0, 200, 250)
	//second row
	c.images[9] = *SubImage(img, 0, 250, 200, 250)
	c.images[10] = *SubImage(img, 200, 250, 200, 250)
	c.images[11] = *SubImage(img, 400, 250, 200, 250)

	// military airplanes
	// path = filepath.Join(c.game.imageSubdir, SPRITESHEET3)
	// img, _, err = ebitenutil.NewImageFromFile(path)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	imgI, _, err = image.Decode(bytes.NewReader(Airplanes3))
	if err != nil {
		log.Fatal(err)
	}
	img = ebiten.NewImageFromImage(imgI)
	//first row
	c.images[12] = *SubImage(img, 0, 0, 200, 200)
	c.images[13] = *SubImage(img, 200, 0, 200, 200)
	c.images[14] = *SubImage(img, 400, 0, 200, 200)
	//second row
	c.images[15] = *SubImage(img, 0, 200, 200, 200)
	c.images[16] = *SubImage(img, 200, 200, 200, 200)
	c.images[17] = *SubImage(img, 400, 200, 200, 200)

	for i := range c.images {
		var height = c.images[i].Bounds().Dy()
		var width = c.images[i].Bounds().Dx()
		newImage := ebiten.NewImage(width, height)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, (float64)(-height))
		op.GeoM.Scale(ENTITY_SCALE, -ENTITY_SCALE)
		newImage.DrawImage(&c.images[i], op)
		c.images[i] = *newImage

	}
}

func FlipArrayOfImagesVertical(inputImages []ebiten.Image) []ebiten.Image {
	newImages := []ebiten.Image{}
	for i := range inputImages {
		var height = inputImages[i].Bounds().Dy()
		var width = inputImages[i].Bounds().Dx()
		newImage := ebiten.NewImage(width, height)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, (float64)(-height))
		op.GeoM.Scale(ENTITY_SCALE, -ENTITY_SCALE)
		newImage.DrawImage(&inputImages[i], op)
		//inputImages[i] = newImage
		newImages = append(newImages, *newImage)

	}
	return newImages
}

func (c *Entity) Draw(screen *ebiten.Image) {
	for i := range ENTITYS_MAX {
		var entity = c.entityUnits[i]
		if entity.active {
			screenX, screenY := c.game.WorldToScreen(entity.worldX, entity.worldY)

			c.drawEntity(screen, screenX, screenY, entity.kind)
		}

	}

}

func (c *Entity) thirdOfScreen(screenX int) int {
	// -1 0 or 1 for left middle or right
	divider12 := WINDOW_WIDTH / 3
	divider23 := divider12 * 2
	if screenX < divider12 {
		return -1
	} else if screenX < divider23 {
		return 0
	} else {
		return 1
	}
}

func (c *Entity) drawEntity(screen *ebiten.Image, screenX, screenY, kind int) {
	var image = &c.images[kind]

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate((float64)(screenX), (float64)(screenY))
	screen.DrawImage(image, op)

}

func (c *Entity) projectileInBounds(punit *EntityUnit) bool {
	// true if unit in bounds
	if punit.worldX < -ENTITY_BORDER ||
		punit.worldX > ENTITY_BORDER+WINDOW_WIDTH ||
		punit.worldY < -ENTITY_BORDER ||
		punit.worldY > ENTITY_BORDER+WINDOW_HEIGHT {
		return false
	} else {
		return true
	}
}

func (c *Entity) addRandomEntity() {
	kind := rand.IntN(ENTITY_KINDS)
	worldX := rand.IntN(ENTITY_START_X_MAX)
	c.addEntity(worldX, ENTITY_START_Y, kind)

}

func (c *Entity) removeAll() {
	c.entityUnits = [ENTITYS_MAX]EntityUnit{}
}

func (c *Entity) addEntity(worldX, worldY, kind int) {
	var puArray = &[ENTITYS_MAX]EntityUnit{}
	var velX, velY = 0, 0
	var worldXC, worldYC = worldX, worldY
	_ = velX
	_ = velY
	if kind == 0 {
		velY = ENTITY_SPEED
		velX = 0
		worldXC += ENTITY_OFFSET_X
		worldYC += ENTITY_OFFSET_Y
		puArray = &c.entityUnits
	} else {
		velY = ENTITY_SPEED
		velX = 0
		puArray = &c.entityUnits
	}
	if kind >= 6 && kind < 12 {
		velX = c.thirdOfScreen(worldXC) * -1
	}
	var nowMilli = time.Now().UnixMilli()
	for i := range ENTITYS_MAX {
		var limitReached = (nowMilli-c.lastTimeMilli > c.entitySpawnInterval)
		if !puArray[i].active && limitReached {
			c.entitySpawnInterval = ENTITY_MIN_INTERVAL + rand.Int64N(ENTITY_RAND_INTERVAL_MAX)
			temp := EntityUnit{}
			temp.worldX, temp.worldY = worldXC, worldYC
			temp.velX, temp.velY = velX, velY
			temp.kind = kind
			temp.active = true
			temp.width = ENTITY_W
			temp.height = ENTITY_H
			puArray[i] = temp

			//fmt.Printf("add entity %v %v %v %v %v %v \n ", worldXC, worldYC, kind, velX, velY, true)
			c.lastTimeMilli = nowMilli
			break
			//puArray[i] = EntityUnit{worldX, worldY, kind, velX, velY, true}

		}
	}
}

func (c *Entity) loopEntitys() {
	for i := range ENTITYS_MAX {
		// player
		var punit = &c.entityUnits[i]
		if !c.projectileInBounds(punit) {
			c.entityUnits[i].active = false
		} else {
			punit.Motion()
			c.FireProjectile(punit)
		}

	}

}

func (c *Entity) Update() error {
	if c.game.mode == PLAY {
		c.loopEntitys()
		c.addRandomEntity()
	}
	var err error
	return err
}
