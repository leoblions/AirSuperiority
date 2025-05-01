package main

import (
	"bytes"
	"image"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var rasterStringStatic = struct {
}{}

const (
	CHAR_OFFSET_TO_INT          = 48
	RASTERSTRING_USE_BYTE_ARRAY = true
)

var (
	RASTERSTRING_BYTES_CHARMAP = CharmapLetters
)

type RasterstringUnit struct {
	screenX int
	screenY int

	stringContent string
	visible       bool
}

func (rsu *RasterstringUnit) GetText() string {
	return rsu.stringContent
}

func (rsu *RasterstringUnit) SetText(newText string) {
	rsu.stringContent = newText
}

type Rasterstring struct {
	game *Game

	letterHeight       int
	letterWidth        int
	letterKerning      int
	spritesheet        *ebiten.Image
	rasterstringUnits  []*RasterstringUnit
	currImage          *ebiten.Image
	runeImageMap       map[rune]*ebiten.Image
	dataFolder         string
	letterLocationFile string
	letterSpriteFile   string
	imageFolder        string
	initDone           bool
}

type characterRecord struct {
	col    int
	row    int
	letter rune
}

func NewRasterstringUnit(content string, startX, startY int) *RasterstringUnit {
	rsu := &RasterstringUnit{startX, startY, content, true}
	return rsu
}

func (p *Rasterstring) AddRasterStringUnit(content string, startX, startY int) *RasterstringUnit {
	//length := len(p.rasterstringUnits)
	newUnit := NewRasterstringUnit(content, startX, startY)
	p.rasterstringUnits = append(p.rasterstringUnits, newUnit)
	return newUnit
}

func NewRasterString(g *Game) *Rasterstring {
	p := &Rasterstring{}
	p.game = g
	p.rasterstringUnits = []*RasterstringUnit{}

	p.letterHeight = 10
	p.letterWidth = 10
	p.letterKerning = 10
	p.dataFolder = "data"
	p.imageFolder = "images"
	p.letterLocationFile = "charmap_letters.cfg"
	p.letterSpriteFile = "letterSpritesW.png"

	p.initImages()

	p.runeImageMap = p.initializeLetterSprites()
	p.initDone = true
	return p

}

func (c *Rasterstring) initImages() {

	img, _, err := image.Decode(bytes.NewReader(RasterLetters))
	if err != nil {
		log.Fatal(err)
	}
	c.spritesheet = ebiten.NewImageFromImage(img)

}

func (c *Rasterstring) initImagesF() {

	var err error
	imageDir := path.Join(c.imageFolder, c.letterSpriteFile)
	c.spritesheet, _, err = ebitenutil.NewImageFromFile(imageDir)
	c.currImage = SubImage(c.spritesheet, 0, 0, c.letterWidth, c.letterHeight)

	if err != nil {
		log.Fatal(nil)
	}
}

func (p *Rasterstring) Draw(screen *ebiten.Image) {
	//fmt.Println("draw rs")
	for _, rsu := range p.rasterstringUnits {
		if !rsu.visible {
			continue
		}
		xOffsetTotal := 0
		for _, letter := range rsu.stringContent {
			//fmt.Printf("draw rs %v %v %v \n", rsu.stringContent, rsu.screenX, rsu.screenY)
			letterImage := p.runeImageMap[letter]
			//fmt.Println(letter)

			if letterImage != nil {
				DrawImageAt(letterImage, screen, rsu.screenX+xOffsetTotal, rsu.screenY)

				xOffsetTotal += p.letterKerning
			} else if letter == ' ' {
				xOffsetTotal += p.letterKerning
			}

		}

	}

}

func (p *Rasterstring) GetRasterStringAsSingleImage(rsu RasterstringUnit) *ebiten.Image {
	width := len(rsu.stringContent) * p.letterWidth
	stringImage := ebiten.NewImage(width, p.letterHeight)
	xOffsetTotal := 0
	for _, letter := range rsu.stringContent {
		letterImage := p.runeImageMap[letter]

		if letterImage != nil {
			//fmt.Println("draw ", c)
			DrawImageAt(stringImage, letterImage, xOffsetTotal, 0)
			xOffsetTotal += p.letterKerning
		} else if letter == ' ' {
			xOffsetTotal += p.letterKerning
		}

	}

	return stringImage
}

func (p *Rasterstring) StringToImage(label string) *ebiten.Image {

	if nil == p.runeImageMap {
		log.Fatal("Rasterstring runeImageMap is null")
	}
	//fmt.Println("letter images ", len(p.runeImageMap))
	width := len(label) * p.letterWidth
	//fmt.Println("Label Width ", width)
	stringImage := ebiten.NewImage(width, p.letterHeight)
	//fmt.Println("Label Height ", p.letterHeight)
	xOffsetTotal := 0
	for _, letter := range label {
		letterImage := p.runeImageMap[letter]

		if letterImage != nil {
			//fmt.Println("draw ", c)
			//fmt.Println("StringToImage letter found ", rune(letter))
			DrawImageAt(letterImage, stringImage, xOffsetTotal, 0)
			xOffsetTotal += p.letterKerning
		} else if letter == ' ' {
			xOffsetTotal += p.letterKerning
		}

	}

	return stringImage
}

func (p *Rasterstring) Update() error {
	return nil
}

func (p *Rasterstring) readSpriteLocationTable() []*characterRecord {

	dataFilePath := path.Join(p.dataFolder, p.letterLocationFile)
	var linesList []*string
	if RASTERSTRING_USE_BYTE_ARRAY {

		linesList = getListOfLinesFromBytes(RASTERSTRING_BYTES_CHARMAP)
	} else {

		linesList = getListOfLinesFromFile(dataFilePath)
	}
	recordsList := []*characterRecord{}
	for _, line := range linesList {
		runeSlice := []rune(*line)
		if len(runeSlice) < 3 {
			// exclude invalid lines that are too short
			continue
		}
		row := int(runeSlice[1]) - CHAR_OFFSET_TO_INT
		col := int(runeSlice[0]) - CHAR_OFFSET_TO_INT
		cr := characterRecord{col, row, runeSlice[2]}
		recordsList = append(recordsList, &cr)
	}

	return recordsList

}

func (p *Rasterstring) initializeLetterSprites() map[rune]*ebiten.Image {
	recordsList := p.readSpriteLocationTable()
	// imageDir := path.Join(p.imageFolder, p.letterSpriteFile)
	// rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//letterSpriteList := []*ebiten.Image{}
	runeImageMap := map[rune]*ebiten.Image{}
	for _, record := range recordsList {
		x := record.col * p.letterWidth
		y := record.row * p.letterHeight
		w := p.letterWidth
		h := p.letterHeight
		thisRune := record.letter
		letterImage := SubImage(p.spritesheet, x, y, w, h)
		runeImageMap[thisRune] = letterImage
	}
	return runeImageMap
}
