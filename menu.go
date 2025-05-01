package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	OPTIONSMENU = iota
	MAINMENU
	PAUSEMENU
)

const (
	BUTTON_HEIGHT     = 50
	BUTTON_WIDTH      = 200
	BUTTON_SPACING_Y  = 30
	BUTTON_AMOUNT     = 4
	BUTTON_LETTER_W   = 10
	BUTTON_LETTER_H   = 10
	MENU_TOP_SPACER   = 80
	MENU_DEBOUNCE_MS  = 400
	MENU_TITLE_WIDTH  = 200
	MENU_TITLE_HEIGHT = 100
	NUMBER_MAX        = 9
)

type Button struct {
	screenX int
	screenY int
	width   int
	height  int
	id      int
	active  bool
}

type Label struct {
	x     int
	y     int
	image *ebiten.Image
}

type Menu struct {
	game                      *Game
	modeChangeDelayToggle     func() bool
	buttonImage1              *ebiten.Image
	buttonImage2              *ebiten.Image
	buttonImageMP             *ebiten.Image
	titleImage, optionsImage  *ebiten.Image
	plusImage, minusImage     *ebiten.Image
	numberImages              [10]*ebiten.Image
	buttonSlice               []*Button
	labelSliceO               []*Label
	labelSliceM               []*Label
	labelSlice                []*Label
	labelStringsO             []string
	labelStringsM             []string
	selectedButton            int
	music, sfx                int
	musicY, difficultyY, sfxY int
	numberLabelX              int
	minusX, plusX             int
	titleImageX               int
	titleImageY               int
	screenX                   int
	screenY                   int
	menuMode                  int
}

func NewMenu(game *Game) *Menu {
	c := &Menu{}
	c.selectedButton = 0
	c.music = 5
	c.sfx = 5
	//c.difficulty = 5
	c.game = game
	c.menuMode = MAINMENU
	c.modeChangeDelayToggle = CreateDelayToggle(MENU_DEBOUNCE_MS)
	c.screenX = (WINDOW_WIDTH / 2) - (BUTTON_WIDTH / 2)
	c.screenY = MENU_TOP_SPACER + (WINDOW_HEIGHT / 2) - (((BUTTON_HEIGHT + BUTTON_SPACING_Y) * BUTTON_AMOUNT) / 2)
	c.initImages()
	c.createButtonMP()
	c.initNumberImages()
	c.initNumberLabelPositions()
	c.labelStringsM = []string{"NEW GAME", "CONTINUE", "OPTIONS", "EXIT"}
	c.labelStringsO = []string{"MUSIC VOL", "SFX VOL", "DIFFICULTY", "BACK"}
	c.initLabels()
	c.initButtons()
	return c
}

func (c *Menu) initButtons() {
	c.buttonSlice = []*Button{}

	for buttonID := range BUTTON_AMOUNT {
		screenX := c.screenX
		screenY := c.getButtonY(buttonID)
		width := BUTTON_WIDTH
		height := BUTTON_HEIGHT
		btn := &Button{screenX, screenY, width, height, buttonID, true}
		c.buttonSlice = append(c.buttonSlice, btn)

	}

}
func (c *Menu) initLabels() {

	c.labelSliceM = []*Label{}
	// main menu labels
	for _, labelString := range c.labelStringsM {
		img := c.game.rasterstring.StringToImage(labelString)
		// calculate offset of label from TLC of button
		x := (BUTTON_WIDTH / 2) - (BUTTON_LETTER_W*len(labelString))/2
		y := (BUTTON_HEIGHT / 2) - (BUTTON_LETTER_H)/2
		label := &Label{x, y, img}
		c.labelSliceM = append(c.labelSliceM, label)

	}

	// options menu labels
	for _, labelString := range c.labelStringsO {
		img := c.game.rasterstring.StringToImage(labelString)
		// calculate offset of label from TLC of button
		x := (BUTTON_WIDTH / 2) - (BUTTON_LETTER_W*len(labelString))/2
		y := (BUTTON_HEIGHT / 2) - (BUTTON_LETTER_H)/2
		label := &Label{x, y, img}
		c.labelSliceO = append(c.labelSliceO, label)

	}

}

func (c *Menu) initNumberLabelPositions() {
	// set positions of number labeltext for options screen
	screenX := (WINDOW_WIDTH / 2) - ((BUTTON_LETTER_W * 3) / 2)
	c.musicY = c.getButtonY(0) + BUTTON_HEIGHT
	c.sfxY = c.getButtonY(1) + BUTTON_HEIGHT
	c.difficultyY = c.getButtonY(2) + BUTTON_HEIGHT
	c.numberLabelX = screenX

}

func (c *Menu) clickLeftOrRightOfButton() int {
	cx := c.game.input.mousePosition.x
	if cx <= (WINDOW_WIDTH / 2) {
		return -1
	} else {
		return 1
	}
}

func (c *Menu) initNumberImages() {

	//c.numberImages = [10]*ebiten.Image
	// main menu labels
	for i := range c.numberImages {
		//labelString := string(i)
		labelString := fmt.Sprintf(" %v ", i)
		img := c.game.rasterstring.StringToImage(labelString)
		width, height := img.Bounds().Dx(), img.Bounds().Dy()
		background := ebiten.NewImage(width, height)
		background.DrawImage(img, nil)
		c.numberImages[i] = background
		//c.labelSliceM = append(c.labelSliceM, label)

	}

}

func (c *Menu) getButtonY(buttonID int) int {
	buttonY := c.screenY + (BUTTON_HEIGHT+BUTTON_SPACING_Y)*buttonID
	return buttonY
}

func (c *Menu) createButtonMP() {

	img, _, err := image.Decode(bytes.NewReader(IconsPng))
	if err != nil {
		log.Fatal(err)
	}
	sheetEI := ebiten.NewImageFromImage(img)
	//c.images = SpriteCutter(ebitenImage, 100, 100, 5, 1)
	backgroundEI := SubImage(sheetEI, 0, 200, 400, 100)
	bgwidth := backgroundEI.Bounds().Dx()
	minusEI := SubImage(sheetEI, 400, 200, 100, 100)
	plusEI := SubImage(sheetEI, 500, 200, 100, 100)
	DrawImageAt(minusEI, backgroundEI, 0, 0)
	DrawImageAt(plusEI, backgroundEI, bgwidth-100, 0)
	c.buttonImageMP = ScaleImage(backgroundEI, BUTTON_WIDTH, BUTTON_HEIGHT)
}

func (c *Menu) initImages() {

	img, _, err := image.Decode(bytes.NewReader(IconsPng))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage := ebiten.NewImageFromImage(img)
	//c.images = SpriteCutter(ebitenImage, 100, 100, 5, 1)
	buttonImage1C := SubImage(ebitenImage, 0, 200, 400, 100)
	buttonImage2C := SubImage(ebitenImage, 0, 300, 400, 100)
	c.buttonImage1 = ScaleImage(buttonImage1C, BUTTON_WIDTH, BUTTON_HEIGHT)
	c.buttonImage2 = ScaleImage(buttonImage2C, BUTTON_WIDTH, BUTTON_HEIGHT)
	// title image
	img, _, err = image.Decode(bytes.NewReader(TitlePng))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage = ebiten.NewImageFromImage(img)
	c.titleImage = ScaleImage(ebitenImage, MENU_TITLE_WIDTH, MENU_TITLE_HEIGHT)

	c.titleImageX = (WINDOW_WIDTH / 2) - (c.titleImage.Bounds().Dx() / 2)
	c.titleImageY = (MENU_TOP_SPACER) - (c.titleImage.Bounds().Dy() / 2)
	// options image
	img, _, err = image.Decode(bytes.NewReader(OptionsPng))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage = ebiten.NewImageFromImage(img)
	c.optionsImage = ScaleImage(ebitenImage, MENU_TITLE_WIDTH, MENU_TITLE_HEIGHT)
	// plus image, minus Image
	img, _, err = image.Decode(bytes.NewReader(IconsPng))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage = ebiten.NewImageFromImage(img)
	c.plusImage = SubImage(ebitenImage, 500, 300, 100, 100)
	c.plusImage = ScaleImage(ebitenImage, MENU_TITLE_WIDTH, MENU_TITLE_HEIGHT)
	c.minusImage = SubImage(ebitenImage, 400, 300, 100, 100)
	c.minusImage = ScaleImage(ebitenImage, MENU_TITLE_WIDTH, MENU_TITLE_HEIGHT)

}

func (c *Menu) drawNumberLabels(screen *ebiten.Image) {
	// music
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.numberLabelX), float64(c.musicY))
	screen.DrawImage(c.numberImages[c.music], op)
	// sfx
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.numberLabelX), float64(c.sfxY))
	screen.DrawImage(c.numberImages[c.sfx], op)
	// difficulty
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.numberLabelX), float64(c.difficultyY))
	screen.DrawImage(c.numberImages[c.game.difficulty], op)

}

func (c *Menu) Draw(screen *ebiten.Image) {

	if c.menuMode == MAINMENU {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(c.titleImageX), float64(c.titleImageY))

		screen.DrawImage(c.titleImage, op)

	} else if c.menuMode == OPTIONSMENU {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(c.titleImageX), float64(c.titleImageY))

		screen.DrawImage(c.optionsImage, op)
		c.drawNumberLabels(screen)
	}

	for i, btn := range c.buttonSlice {
		if btn.active {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(btn.screenX), float64(btn.screenY))
			if i != 3 && c.menuMode == OPTIONSMENU {
				screen.DrawImage(c.buttonImageMP, op)
			} else {
				screen.DrawImage(c.buttonImage1, op)
			}

			if i < len(c.labelSlice) {
				op = &ebiten.DrawImageOptions{}
				label := c.labelSlice[i]
				op.GeoM.Translate(float64(btn.screenX+label.x), float64(btn.screenY+label.y))
				screen.DrawImage(label.image, op)
			}

		}

	}

}

func (c *Menu) Update() error {
	if c.menuMode == OPTIONSMENU {
		c.labelSlice = c.labelSliceO
	} else if c.menuMode == MAINMENU {
		c.labelSlice = c.labelSliceM
	}
	if c.game.input.mouseL {
		buttonID := c.clickIntersectButton()
		//fmt.Println("left click ", buttonID)

		if c.modeChangeDelayToggle() {
			c.buttonPressAction(buttonID)
			//fmt.Println("accepted click ", buttonID)
		}

	}
	return nil

}

func (c *Menu) keyActivateButton() {
	c.buttonPressAction(c.selectedButton)
}

func (c *Menu) keyChangeButton(isUp bool) {

	nextButton := c.selectedButton
	if isUp {
		nextButton -= 1
	} else {
		nextButton += 1
	}
	c.selectedButton = Clamp(0, int(BUTTON_AMOUNT-1), nextButton)
}

func (c *Menu) buttonPressAction(buttonID int) {
	switch buttonID {
	case 0:
		if c.menuMode == MAINMENU {

			c.game.mode = PLAY
			c.game.setStatusStringToMode()
			c.game.resetGame()
		} else if c.menuMode == OPTIONSMENU {
			change := c.clickLeftOrRightOfButton()
			c.music = Clamp(0, NUMBER_MAX, change+c.music)
			c.game.sound.SetMusicVolume(c.music)
		}

	case 1:

		if c.menuMode == MAINMENU {

			c.game.mode = PLAY
			c.game.setStatusStringToMode()
		} else if c.menuMode == OPTIONSMENU {
			change := c.clickLeftOrRightOfButton()
			c.sfx = Clamp(0, NUMBER_MAX, change+c.sfx)
			c.game.sound.SetSFXVolume(c.sfx)
		}
	case 2:
		if c.menuMode == MAINMENU {

			c.menuMode = OPTIONSMENU
		} else if c.menuMode == OPTIONSMENU {
			change := c.clickLeftOrRightOfButton()
			c.game.difficulty = Clamp(0, NUMBER_MAX, change+c.game.difficulty)
		}
	case 3:
		if c.menuMode == MAINMENU {

			os.Exit(0)
		} else if c.menuMode == OPTIONSMENU {
			c.menuMode = MAINMENU
		}
	}
}

func (c *Menu) clickIntersectButton() int {
	cx, cy := c.game.input.mousePosition.x, c.game.input.mousePosition.y
	for _, btn := range c.buttonSlice {
		bx2 := btn.screenX + btn.width
		by2 := btn.screenY + btn.height
		if (btn.screenX <= cx && cx <= bx2) && (btn.screenY <= cy && cy <= by2) {
			return btn.id
		}
	}
	return -1
}
