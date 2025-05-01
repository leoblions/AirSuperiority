package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	_ "os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PLAY = iota
	PAUSED
	MENU
	GAMEOVER
)

const (
	WINDOW_HEIGHT            = 480
	WINDOW_WIDTH             = 640
	WINDOW_TITLE             = "AIR SUPERIORITY"
	GAME_GAMEOVER_STRING     = "GAME OVER"
	GAME_LIVES_TS            = "LIVES %v"
	GAME_SCORE_TS            = "SCORE %v"
	GAME_LIVES_X             = WINDOW_WIDTH - 100
	GAME_SCORE_X             = WINDOW_WIDTH - 100
	GAME_STATUS_X            = WINDOW_WIDTH/2 - 20
	GAME_LIVES_Y             = 10
	GAME_SCORE_Y             = 30
	GAME_MIDDLE_Y            = WINDOW_HEIGHT / 2
	GAME_POINTS_PER_NEW_LIFE = 30
	GAME_GODMODE             = false
	GAME_START_HEALTH        = 100
	GAME_START_FUEL          = 100
	GAME_START_MODE          = MENU
	GAME_START_LIVES         = 3
	GAME_START_VOLUME        = 0.5
)

type Component interface {
	Update() error
	Draw(*ebiten.Image)
}

type Movable struct {
	worldX int
	worldY int
	velX   int
	velY   int
	width  int
	height int
}

type Game struct {
	components   []Component
	input        *Input
	background   *Background
	projectile   *Projectile
	pickup       *Pickup
	rasterstring *Rasterstring
	player       *Player
	explosion    *Explosion
	entity       *Entity
	hud          *HUD
	sound        *Sound
	menu         *Menu
	mode         int
	screenLocX   int
	screenLocY   int
	lives        int
	health       int
	fuel         int
	difficulty   int
	score        int
	loaded       bool
	imageSubdir  string
	soundSubdir  string
	statusString string
	audioContext *audio.Context
	scoreRSU     *RasterstringUnit
	livesRSU     *RasterstringUnit
	statusRSU    *RasterstringUnit
	middleRSU    *RasterstringUnit
	//input
	touchIDs   []ebiten.TouchID
	gamepadIDs []ebiten.GamepadID
	keyIDs     []ebiten.Key
	didDraw    bool
	godmode    bool
}

func NewGame() *Game {
	var g = &Game{}

	g.loaded = false
	g.imageSubdir = "images"
	g.soundSubdir = "sound"
	g.mode = GAME_START_MODE
	g.godmode = false
	g.lives = 3
	g.difficulty = 5
	g.health = GAME_START_HEALTH
	g.fuel = GAME_START_FUEL
	g.input = NewInput(g)
	g.components = []Component{}

	g.statusString = "PLAY"

	g.background = NewBackground(g)
	g.components = append(g.components, g.background)

	g.rasterstring = NewRasterString(g)
	g.middleRSU = g.rasterstring.AddRasterStringUnit(GAME_GAMEOVER_STRING, GAME_STATUS_X, GAME_MIDDLE_Y)
	g.scoreRSU = g.rasterstring.AddRasterStringUnit(fmt.Sprintf(GAME_SCORE_TS, g.score), GAME_SCORE_X, GAME_SCORE_Y)
	g.livesRSU = g.rasterstring.AddRasterStringUnit(fmt.Sprintf(GAME_LIVES_TS, g.lives), GAME_LIVES_X, GAME_LIVES_Y)
	g.statusRSU = g.rasterstring.AddRasterStringUnit(g.statusString, GAME_STATUS_X, GAME_LIVES_Y)
	g.components = append(g.components, g.rasterstring)
	g.setStatusStringToMode()
	g.middleRSU.visible = false

	g.player = NewPlayer(g)
	g.components = append(g.components, g.player)

	g.hud = NewHUD(g)
	g.components = append(g.components, g.hud)

	g.explosion = NewExplosion(g)
	g.components = append(g.components, g.explosion)

	g.projectile = NewProjectile(g)
	g.components = append(g.components, g.projectile)

	g.pickup = NewPickup(g)
	g.components = append(g.components, g.pickup)

	g.entity = NewEntity(g)
	g.components = append(g.components, g.entity)

	g.sound = NewSound(g)
	g.sound.musicVolume = GAME_START_VOLUME
	g.sound.sfxVolume = GAME_START_VOLUME

	g.menu = NewMenu(g)
	//g.components = append(g.components, g.menu)

	exe, err := os.Executable()
	path := filepath.Join(exe, "images")
	_ = path
	if err != nil {

	}
	//fmt.Println("cwd = " + path)
	g.loaded = true
	return g
}

func (g *Game) Update() error {
	g.isKeyJustPressed()
	g.input.MouseHandler()
	if g.mode == PLAY {
		// if   done loading, and play mode
		for _, v := range g.components {

			v.Update()
		}
		g.didDraw = false
	}
	if g.mode == MENU {
		g.background.Update()
		g.menu.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, WINDOW_TITLE)

	if g.mode != MENU {
		for _, v := range g.components {

			v.Draw(screen)
		}
	}
	if g.mode == MENU {
		g.background.Draw(screen)
		g.menu.Draw(screen)
	}

}

func (g *Game) incrementScore() {
	g.score += 1
	g.scoreRSU.SetText(fmt.Sprintf(GAME_SCORE_TS, g.score))
	if g.score%GAME_POINTS_PER_NEW_LIFE == 0 {
		g.incrementLives()
	}
}

func (g *Game) resetScore() {
	g.score = 0
	g.scoreRSU.SetText(fmt.Sprintf(GAME_SCORE_TS, g.score))
}

func (g *Game) resetGame() {
	g.resetScore()
	if g.mode == GAMEOVER {
		g.mode = PLAY
		g.middleRSU.visible = false
	}
	g.lives = GAME_START_LIVES
	g.health = GAME_START_HEALTH
	g.fuel = GAME_START_FUEL
	g.entity.removeAll()
	g.livesRSU.SetText(fmt.Sprintf(GAME_LIVES_TS, g.lives))
}

func (g *Game) incrementLives() {
	g.lives += 1
	g.livesRSU.SetText(fmt.Sprintf(GAME_LIVES_TS, g.lives))
}

func (g *Game) setStatusString(newStatus string) {

	g.statusRSU.SetText(newStatus)
}

func (g *Game) setStatusStringToMode() {
	var statusString = ""
	switch g.mode {
	case PLAY:
		statusString = "PLAY"
	case PAUSED:
		statusString = "PAUSED"
	case MENU:
		statusString = "MENU"
	}

	g.statusRSU.SetText(statusString)
}

func (g *Game) decrementLives() {
	g.lives -= 1
	g.livesRSU.SetText(fmt.Sprintf(GAME_LIVES_TS, g.lives))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WINDOW_WIDTH, WINDOW_HEIGHT
}

func main() {
	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	ebiten.SetWindowTitle(WINDOW_TITLE)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
