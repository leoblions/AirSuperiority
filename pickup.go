package main

import (
	"bytes"
	"image"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PICKUP_H           = 30
	PICKUP_W           = 30
	PICKUP_DROP_OFFSET = 50
	PICKUPS_MAX        = 10
	PICKUP_KINDS       = 4
	PICKUP_DURATION    = 500
	PICKUP_DROP_FREQ   = 4
)

const (
	PICKUP_HEALTH1 = iota
	PICKUP_HEALTH2
	PICKUP_FUEL1
	PICKUP_FUEL2
)

type Pickup struct {
	game         *Game
	pickupImages [PICKUP_KINDS]*ebiten.Image
	pickupUnits  [PICKUPS_MAX]*PickupUnit

	lastTimeMilli int64
}

type PickupUnit struct {
	worldX, worldY, kind, life int
	active                     bool
}

func (punit *PickupUnit) Dimensions() (int, int, int, int) {
	return punit.worldX, punit.worldY, PICKUP_W, PICKUP_H
}

func NewPickup(g *Game) *Pickup {
	c := &Pickup{}
	c.game = g
	c.pickupUnits = [PICKUPS_MAX]*PickupUnit{}
	c.pickupImages = [PICKUP_KINDS]*ebiten.Image{}
	c.initImages()

	return c
}

func (c *Pickup) initImages() {

	img, _, err := image.Decode(bytes.NewReader(Icons2Png))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage := ebiten.NewImageFromImage(img)
	subImageA := SubImage(ebitenImage, 0, 0, 100, 100)
	subImageB := SubImage(ebitenImage, 100, 0, 100, 100)
	subImageC := SubImage(ebitenImage, 0, 200, 100, 100)
	subImageD := SubImage(ebitenImage, 100, 200, 100, 100)
	c.pickupImages[0] = ScaleImage(subImageA, PICKUP_W, PICKUP_H)
	c.pickupImages[1] = ScaleImage(subImageB, PICKUP_W, PICKUP_H)
	c.pickupImages[2] = ScaleImage(subImageC, PICKUP_W, PICKUP_H)
	c.pickupImages[3] = ScaleImage(subImageD, PICKUP_W, PICKUP_H)

}

func (c *Pickup) Draw(screen *ebiten.Image) {
	for i := range PICKUPS_MAX {

		var pickup = c.pickupUnits[i]
		if nil != pickup && pickup.active {
			screenX, screenY := c.game.WorldToScreen(pickup.worldX, pickup.worldY)

			c.drawPickup(screen, screenX, screenY, pickup.kind)
		}

	}

}

func (c *Pickup) drawPickup(screen *ebiten.Image, screenX, screenY, kind int) {
	var image *ebiten.Image
	image = c.pickupImages[kind]
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate((float64)(screenX), (float64)(screenY))
	screen.DrawImage(image, op)

}

func (c *Pickup) InBounds(punit *PickupUnit) bool {
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

func (c *Pickup) addPickup(worldX, worldY, kind int) *PickupUnit {

	for i := range PROJECTILES_MAX {
		if nil == c.pickupUnits[i] || !c.pickupUnits[i].active {
			c.pickupUnits[i] = &PickupUnit{worldX, worldY, kind, PICKUP_DURATION, true}

			return c.pickupUnits[i]

		}
	}
	return nil
}

func (c *Pickup) dropLoot(eunit EntityUnit) *PickupUnit {
	chance := rand.IntN(PICKUP_DROP_FREQ) == 0
	//chance := true
	for i := range PROJECTILES_MAX {
		if chance && (nil == c.pickupUnits[i] || !c.pickupUnits[i].active) {
			kind := rand.IntN(PICKUP_KINDS)
			//fmt.Println("drop loot ", kind)
			c.pickupUnits[i] = &PickupUnit{eunit.worldX + PICKUP_DROP_OFFSET,
				eunit.worldY + PICKUP_DROP_OFFSET, kind, PICKUP_DURATION, true}

			return c.pickupUnits[i]

		}
	}
	return nil
}

func (c *Pickup) checkUnitCollidePlayer(punit *PickupUnit) {
	if !punit.active {
		return
	}
	collided := Intersect(punit, c.game.player)
	if collided && c.game.player.active {
		c.playerTouchPickupAction(punit.kind)

		punit.active = false

		c.game.incrementScore()

	}

}

func (c *Pickup) playerTouchPickupAction(kind int) {
	switch kind {
	case 0:
		c.game.player.heal(25)
	case 1:
		c.game.player.heal(35)
	case 2:
		c.game.player.refuel(25)
	case 3:
		c.game.player.refuel(55)
	default:
		c.game.player.takeDamage(PROJECTILE_PLAYER_DAMAGE)
		wx, wy, _, _ := c.game.player.Dimensions()
		explosionKind := 2

		c.game.explosion.addExplosion(wx, wy, explosionKind)

	}
}

func (c *Pickup) loopPickups() {
	for i := range PICKUPS_MAX {
		if eunit := c.pickupUnits[i]; nil != eunit {
			if !c.InBounds(eunit) || eunit.life <= 0 {
				eunit.active = false
			} else {
				if eunit.life > 0 {
					eunit.life -= 1
				}
				c.checkUnitCollidePlayer(eunit)
			}
		}

	}

}

func (c *Pickup) Update() error {

	c.loopPickups()
	var err error
	return err
}
