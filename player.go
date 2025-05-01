package main

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	PLAYER_BG_COLOR = color.RGBA{0xff, 0x10, 0x00, 0xff}
)

const (
	PLAYER_SIZE                      = 100
	PLAYER_DEFAULT_SPEED             = 3
	PLAYER_XMAX                      = WINDOW_WIDTH - PLAYER_SIZE
	PLAYER_YMAX                      = WINDOW_HEIGHT - PLAYER_SIZE
	PLAYER_HIT_ENEMY_DAMAGE          = 30
	PLAYER_COLLIDE_PROJECTILE_DAMAGE = 20
	PLAYER_RESPAWN_COUNT             = 60
	HEALTH_MAX                       = 100
	FUEL_MAX                         = 100
)

type Player struct {
	game         *Game
	drawPulser   func() bool
	images       []*ebiten.Image
	image        *ebiten.Image
	respawnCount int
	speed        int
	imageID      int
	sprint       bool
	active       bool
	motionFlags  [4]bool
	Movable
}

func NewPlayer(g *Game) *Player {
	c := &Player{}
	c.game = g
	c.imageID = 2
	c.image = ebiten.NewImage(PLAYER_SIZE, PLAYER_SIZE)
	c.image.Fill(PLAYER_BG_COLOR)
	c.setPositionBottomMiddle()
	c.speed = PLAYER_DEFAULT_SPEED
	c.sprint = false
	c.motionFlags = [...]bool{false, false, false, false}
	c.width = PLAYER_SIZE
	c.height = PLAYER_SIZE
	c.active = true
	c.drawPulser = Pulser(10)

	//screenY := (float64)(c.worldY - c.game.screenLocY)
	//fmt.Println(" player screen y ", screenY)
	c.initImages()
	return c
}

func (c *Player) Dimensions() (int, int, int, int) {
	return c.worldX, c.worldY, PLAYER_SIZE, PLAYER_SIZE
}

func (c *Player) initImagesF() {
	SPRITESHEET := "airplanePlayer.png"
	path := filepath.Join(c.game.imageSubdir, SPRITESHEET)
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	c.images = SpriteCutter(img, 100, 100, 5, 1)

	c.image = SubImage(img, 100, 0, 100, 100)
}

func (c *Player) initImages() {

	img, _, err := image.Decode(bytes.NewReader(AirplanePlayer))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage := ebiten.NewImageFromImage(img)
	c.images = SpriteCutter(ebitenImage, 100, 100, 5, 1)
	c.image = SubImage(ebitenImage, 100, 0, 100, 100)

}

func (c *Player) heal(healthAmount int) {
	newHealth := c.game.health + healthAmount
	if newHealth >= 0 && newHealth <= HEALTH_MAX {
		c.game.health = newHealth

	} else {
		c.game.health = HEALTH_MAX
	}
	c.game.hud.recalculateBarImages()
}

func (c *Player) refuel(fuelAmount int) {
	newFuel := c.game.fuel + fuelAmount
	if newFuel >= 0 && newFuel <= FUEL_MAX {
		c.game.fuel = newFuel

	} else {
		c.game.fuel = FUEL_MAX
	}
	c.game.hud.recalculateBarImages()
}

func (c *Player) fireProjectile() {
	c.game.projectile.addPlayerProjectile(c.worldX, c.worldY)
}

func (c *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screenX := (float64)(c.worldX - c.game.screenLocX)
	screenY := (float64)(c.worldY - c.game.screenLocY)
	_ = screenX
	_ = screenY
	op.GeoM.Translate(screenX, screenY)
	if c.game.mode == GAMEOVER || (c.respawnCount > 0 && c.drawPulser()) {

	} else {
		screen.DrawImage(c.images[c.imageID], op)
	}

}

func (c *Player) setPositionBottomMiddle() {
	//screen coords
	x := (WINDOW_WIDTH / 2) - (PLAYER_SIZE / 2)
	y := (WINDOW_HEIGHT) - PLAYER_SIZE
	c.worldX = x + c.game.screenLocX
	c.worldY = y + c.game.screenLocY

}

func (c *Player) Update() error {
	var err error
	c.imageID = 0
	if c.game.mode == PLAY {
		c.setPlayerImage()
		c.playerMotion()
		c.checkPlayerCollideEntity()
	}

	c.motionFlags = [...]bool{false, false, false, false}
	c.sprint = false
	if c.respawnCount > 0 {
		c.respawnCount -= 1
	}
	return err
}

func (c *Player) checkPlayerCollideEntity() {

	for i, entityUnit := range c.game.entity.entityUnits {
		collided := Intersect(c, &entityUnit)
		if collided && entityUnit.active {
			c.game.entity.entityUnits[i].active = false

			wx, wy, _, _ := entityUnit.Dimensions()
			kind := 0
			if entityUnit.kind > 7 {
				kind = 1
			}
			c.game.explosion.addExplosion(wx, wy, kind)
			c.game.incrementScore()
			c.takeDamage(PLAYER_HIT_ENEMY_DAMAGE)
			//fmt.Println("projectile hit entity")
			return
		}

	}

}

func (c *Player) takeDamage(damageAmount int) {
	newHealth := c.game.health - damageAmount
	if newHealth > 0 {
		c.game.health = newHealth
	} else {
		c.game.health = 0
		c.die()
	}
	c.game.hud.recalculateBarImages()

}

func (c *Player) die() {
	c.respawnCount = PLAYER_RESPAWN_COUNT
	c.game.decrementLives()
	c.setPositionBottomMiddle()
	c.game.health = GAME_START_HEALTH
	c.game.fuel = GAME_START_FUEL
	if c.game.lives < 0 {
		c.game.mode = GAMEOVER
		c.game.setStatusStringToMode()
		c.game.middleRSU.visible = true
	}

}

func (c *Player) playerMotion() {
	c.velX, c.velY = 0, 0
	if c.sprint {
		c.speed = PLAYER_DEFAULT_SPEED + 2
	} else {
		c.speed = PLAYER_DEFAULT_SPEED
	}
	if c.motionFlags[0] {
		c.velY = -c.speed
	}
	if c.motionFlags[1] {
		c.velY = c.speed
	}
	if c.motionFlags[2] {
		c.velX = -c.speed
	}
	if c.motionFlags[3] {
		c.velX = c.speed
	}

	var propWorldX = c.worldX + c.velX
	var propWorldY = c.worldY + c.velY
	c.worldX = Clamp(0, PLAYER_XMAX, propWorldX)
	c.worldY = Clamp(0, PLAYER_YMAX, propWorldY)
}

func (c *Player) setPlayerImage() {
	c.imageID = 2

	if c.motionFlags[2] {
		c.imageID = 1
		if c.sprint {
			c.imageID = 0
		}
	}
	if c.motionFlags[3] {
		c.imageID = 3
		if c.sprint {
			c.imageID = 4
		}
	}

}
