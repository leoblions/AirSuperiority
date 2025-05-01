package main

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PROJECTILE_SPEED         = 3
	PROJECTILE_H             = 20
	PROJECTILE_W             = 8
	PROJECTILES_MAX          = 10
	PROJECTILE_OFFSET_X      = 50
	PROJECTILE_OFFSET_Y      = 1
	PROJECTILE_BORDER        = 100
	PROJECTILE_MIN_INTERVAL  = 500
	PROJECTILE_PLAYER_DAMAGE = 25
)

var (
	projectileColorE = color.RGBA{0xff, 0xff, 0x30, 0xff}
	projectileColorP = color.RGBA{0xe0, 0xe0, 0x6f, 0xff}
)

type Projectile struct {
	game             *Game
	imageP, imageE   *ebiten.Image
	projectileUnitsP [PROJECTILES_MAX]ProjectileUnit
	projectileUnitsE [PROJECTILES_MAX]ProjectileUnit
	lastTimeMilli    int64
	testRect         Movable
}

/*
kinds:
0 = player
1 = enemy
*/
const (
	PROJ_E = iota
	PROJ_P
)

type ProjectileUnit struct {
	worldX, worldY, kind int
	velX, velY           int
	active               bool
}

func (punit *ProjectileUnit) Motion() {
	punit.worldX += punit.velX
	punit.worldY += punit.velY
}

func (punit *ProjectileUnit) Dimensions() (int, int, int, int) {
	return punit.worldX, punit.worldY, PROJECTILE_W, PROJECTILE_H
}

func NewProjectile(g *Game) *Projectile {
	c := &Projectile{}
	c.game = g
	c.lastTimeMilli = time.Now().UnixMilli()
	c.projectileUnitsE = [PROJECTILES_MAX]ProjectileUnit{}
	c.projectileUnitsP = [PROJECTILES_MAX]ProjectileUnit{}
	c.initImages()
	c.testRect = Movable{0, 0, 0, 0, PROJECTILE_W, PROJECTILE_H}
	//c.projectileUnitsE[0] = ProjectileUnit{200, 200, 1, 0, 3, true}
	return c
}

func (c *Projectile) initImagesF() {
	c.imageP = ebiten.NewImage(PROJECTILE_W, PROJECTILE_H)
	c.imageE = ebiten.NewImage(PROJECTILE_W, PROJECTILE_H)
	SPRITESHEET := "rocketY.png"
	path := filepath.Join(c.game.imageSubdir, SPRITESHEET)
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	c.imageP.Fill(projectileColorP)
	c.imageE.Fill(projectileColorE)
	c.imageP = ScaleImage(img, PROJECTILE_W, PROJECTILE_H)
}

func (c *Projectile) initImages() {

	img, _, err := image.Decode(bytes.NewReader(RocketR))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage := ebiten.NewImageFromImage(img)
	c.imageP = ScaleImage(ebitenImage, PROJECTILE_W, PROJECTILE_H)

	//c.images = SpriteCutter(ebitenImage, 100, 100, 5, 1)
	//c.image = SubImage(ebitenImage, 100, 0, 100, 100)

	img, _, err = image.Decode(bytes.NewReader(RocketY))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage = ebiten.NewImageFromImage(img)
	c.imageE = ScaleImage(ebitenImage, PROJECTILE_W, PROJECTILE_H)

}

func (c *Projectile) Draw(screen *ebiten.Image) {
	for i := range PROJECTILES_MAX {
		// player
		var projectile = c.projectileUnitsP[i]
		if projectile.active {
			screenX, screenY := c.game.WorldToScreen(projectile.worldX, projectile.worldY)

			c.drawProjectile(screen, screenX, screenY, PROJ_P)
		}
		// enemy
		projectile = c.projectileUnitsE[i]
		if projectile.active {
			screenX, screenY := c.game.WorldToScreen(projectile.worldX, projectile.worldY)
			c.drawProjectile(screen, screenX, screenY, PROJ_E)
		}
	}

}

func (c *Projectile) drawProjectile(screen *ebiten.Image, screenX, screenY, kind int) {
	var image *ebiten.Image
	if kind == PROJ_P {
		image = c.imageP
	} else {
		image = c.imageE
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate((float64)(screenX), (float64)(screenY))
	screen.DrawImage(image, op)

}

func (c *Projectile) projectileInBounds(punit *ProjectileUnit) bool {
	// true if unit in bounds
	if punit.worldX < -PROJECTILE_BORDER ||
		punit.worldX > PROJECTILE_BORDER+WINDOW_WIDTH ||
		punit.worldY < -PROJECTILE_BORDER ||
		punit.worldY > PROJECTILE_BORDER+WINDOW_HEIGHT {
		return false
	} else {
		return true
	}
}

func (c *Projectile) addProjectile0(worldX, worldY, kind int) *ProjectileUnit {
	var puArray = &[PROJECTILES_MAX]ProjectileUnit{}
	var velX, velY = 0, 0
	var worldXC, worldYC = worldX, worldY
	_ = velX
	_ = velY
	if kind == PROJ_P {
		velY = -PROJECTILE_SPEED
		velX = 0
		worldXC += PROJECTILE_OFFSET_X
		worldYC += PROJECTILE_OFFSET_Y
		puArray = &c.projectileUnitsP
	} else {
		velY = PROJECTILE_SPEED
		velX = 0
		puArray = &c.projectileUnitsE
	}
	var nowMilli = time.Now().UnixMilli()
	for i := range PROJECTILES_MAX {
		var limitReached = (nowMilli-c.lastTimeMilli > PROJECTILE_MIN_INTERVAL)
		if !puArray[i].active && limitReached {
			puArray[i] = ProjectileUnit{worldXC, worldYC, kind, velX, velY, true}
			//fmt.Println("add projectil ", i)
			c.lastTimeMilli = nowMilli
			return &puArray[i]
			//puArray[i] = ProjectileUnit{worldX, worldY, kind, velX, velY, true}

		}
	}
	return nil
}

func (c *Projectile) addPlayerProjectile(worldX, worldY int) *ProjectileUnit {
	var puArray = &[PROJECTILES_MAX]ProjectileUnit{}
	var velX, velY = 0, 0
	var worldXC, worldYC = worldX, worldY

	velY = -PROJECTILE_SPEED
	velX = 0
	worldXC += PROJECTILE_OFFSET_X
	worldYC += PROJECTILE_OFFSET_Y
	puArray = &c.projectileUnitsP

	var nowMilli = time.Now().UnixMilli()
	for i := range PROJECTILES_MAX {
		var limitReached = (nowMilli-c.lastTimeMilli > PROJECTILE_MIN_INTERVAL)
		if !puArray[i].active && limitReached {
			puArray[i] = ProjectileUnit{worldXC, worldYC, PROJ_P, velX, velY, true}
			//fmt.Println("add projectil ", i)
			c.lastTimeMilli = nowMilli
			return &puArray[i]
			//puArray[i] = ProjectileUnit{worldX, worldY, kind, velX, velY, true}

		}
	}
	return nil
}

func (c *Projectile) addEnemyProjectile(worldX, worldY int) *ProjectileUnit {
	//var puArray = &[PROJECTILES_MAX]ProjectileUnit{}
	var velX, velY = 0, 0
	var worldXC, worldYC = worldX, worldY

	velY = PROJECTILE_SPEED
	velX = 0
	kind := PROJ_E
	//puArray = &c.projectileUnitsE
	for i := range PROJECTILES_MAX {
		if nil == &c.projectileUnitsE[i] || !c.projectileUnitsE[i].active {
			c.projectileUnitsE[i] = ProjectileUnit{worldXC, worldYC, kind, velX, velY, true}
			//fmt.Println("add projectil ", i)
			return &c.projectileUnitsE[i]

		}
	}
	return nil
}

func (c *Projectile) checkUnitCollideEntity(punit *ProjectileUnit) int {
	if !punit.active {
		return -1
	}
	for i, entityUnit := range c.game.entity.entityUnits {
		collided := Intersect(punit, &entityUnit)
		if collided && entityUnit.active {
			c.game.entity.entityUnits[i].active = false
			c.game.pickup.dropLoot(entityUnit)
			punit.active = false
			wx, wy, _, _ := entityUnit.Dimensions()
			explosionKind := 0
			if entityUnit.kind > 7 {
				explosionKind = 1
			}
			c.game.explosion.addExplosion(wx, wy, explosionKind)
			c.game.incrementScore()
			//fmt.Println("projectile hit entity")
			return i
		}

	}
	return -1

}

func (c *Projectile) checkUnitCollidePlayer(punit *ProjectileUnit) {
	if !punit.active || punit.kind != PROJ_E {
		return
	}
	collided := Intersect(punit, c.game.player)
	if collided && c.game.player.active {
		c.game.player.takeDamage(PROJECTILE_PLAYER_DAMAGE)

		punit.active = false
		wx, wy, _, _ := c.game.player.Dimensions()
		explosionKind := 2

		c.game.explosion.addExplosion(wx, wy, explosionKind)
		c.game.incrementScore()
		//fmt.Println("projectile hit entity")

	}

}

func (c *Projectile) loopProjectiles() {
	for i := range PROJECTILES_MAX {
		// player
		var punit = &c.projectileUnitsP[i]
		if !c.projectileInBounds(punit) {
			c.projectileUnitsP[i].active = false
		} else {
			punit.Motion()
			c.checkUnitCollideEntity(punit)
		}
		// enemy
		var eunit = &c.projectileUnitsE[i]
		if !c.projectileInBounds(eunit) {
			c.projectileUnitsE[i].active = false
		} else {
			eunit.Motion()
			c.checkUnitCollidePlayer(eunit)
		}
	}

}

func (c *Projectile) Update() error {

	c.loopProjectiles()
	var err error
	return err
}
