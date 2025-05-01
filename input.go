package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type pos struct {
	x int
	y int
}

type Input struct {
	game                  *Game
	modeChangeDelayToggle func() bool
	run                   bool
	mouseL                bool
	mouseR                bool
	mouseM                bool

	mousePosition pos
}

func NewInput(g *Game) *Input {
	input := &Input{}
	input.game = g
	input.modeChangeDelayToggle = CreateDelayToggle(900)
	return input
}

func (input *Input) MouseHandler() {
	// Paint the brush by mouse dragging
	mx, my := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		input.mouseL = true
	} else {
		input.mouseL = false
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		input.mouseR = true
	} else {
		input.mouseR = false
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		input.mouseM = true
	} else {
		input.mouseM = false
	}
	input.mousePosition = pos{
		x: mx,
		y: my,
	}
}

func (g *Game) isKeyJustPressed() {
	// runs when update isnt being called
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])

	// make list of keyboard keys
	//g.keyIDs = inpututil.AppendJustPressedKeys(g.keyIDs[:0])

	g.keyIDs = inpututil.AppendPressedKeys(g.keyIDs[:0])

	if len(g.keyIDs) > 0 {
		for _, v := range g.keyIDs {
			switch v {
			case ebiten.KeySpace:
				g.player.sprint = true
			case ebiten.KeyEnter:
				if g.mode == MENU {
					g.menu.keyActivateButton()
				} else if g.mode == GAMEOVER {
					g.resetGame()
				}
			case ebiten.KeyShift:
				g.player.sprint = true
			case ebiten.KeyP:
				if g.input.modeChangeDelayToggle() {
					if g.mode == PLAY {
						g.mode = PAUSED
						g.setStatusStringToMode()
					} else if g.mode == PAUSED {
						g.mode = PLAY

						g.setStatusStringToMode()
					}
				}
			case ebiten.KeyEscape:
				if g.input.modeChangeDelayToggle() {
					if g.mode == PLAY || g.mode == PAUSED {
						g.mode = MENU
						g.setStatusStringToMode()
					} else if g.mode == MENU {
						g.mode = PLAY

						g.setStatusStringToMode()
					}
				}
			case ebiten.KeyW:
				g.player.motionFlags[0] = true
			case ebiten.KeyS:
				g.player.motionFlags[1] = true
			case ebiten.KeyA:
				g.player.motionFlags[2] = true
			case ebiten.KeyD:
				g.player.motionFlags[3] = true
			case ebiten.KeyF:
				g.player.fireProjectile()
				g.sound.PlaySFX(4)
				if g.mode == GAMEOVER {
					g.resetGame()
				}
			case ebiten.KeySemicolon:
				g.lives = 0
			case ebiten.KeyUp:
				g.player.motionFlags[0] = true
				if g.mode == MENU {
					g.menu.keyChangeButton(true)
				}
			case ebiten.KeyDown:
				g.player.motionFlags[1] = true
				if g.mode == MENU {
					g.menu.keyChangeButton(false)
				}
			case ebiten.KeyLeft:
				g.player.motionFlags[2] = true
			case ebiten.KeyRight:
				g.player.motionFlags[3] = true

			}
		}

	}

	g.gamepadIDs = ebiten.AppendGamepadIDs(g.gamepadIDs[:0])
	for _, g := range g.gamepadIDs {
		if ebiten.IsStandardGamepadLayoutAvailable(g) {
			if inpututil.IsStandardGamepadButtonJustPressed(g, ebiten.StandardGamepadButtonRightBottom) {

			}
			if inpututil.IsStandardGamepadButtonJustPressed(g, ebiten.StandardGamepadButtonRightRight) {

			}
		} else {
			// The button 0/1 might not be A/B buttons.
			if inpututil.IsGamepadButtonJustPressed(g, ebiten.GamepadButton0) {

			}
			if inpututil.IsGamepadButtonJustPressed(g, ebiten.GamepadButton1) {

			}
		}
	}

}
