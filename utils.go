package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	UTILS_NEWLINE = "\n"
	UTILS_COMMENT = '#'
)

func DrawImageAt(src *ebiten.Image, dest *ebiten.Image, screenX, screenY int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate((float64)(screenX), (float64)(screenY))
	dest.DrawImage(src, op)
}

func DrawImageAtF(src *ebiten.Image, dest *ebiten.Image, screenX, screenY float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate((screenX), (screenY))
	dest.DrawImage(src, op)
}

func (game *Game) WorldToScreen(worldX, worldY int) (int, int) {
	screenX := (worldX - game.screenLocX)
	screenY := (worldY - game.screenLocY)
	return screenX, screenY
}

func CropImage(src *ebiten.Image, width, height int) *ebiten.Image {
	output := ebiten.NewImage(width, height)
	output.DrawImage(src, nil)
	return output
}

func SubImage(src *ebiten.Image, startX, startY, width, height int) *ebiten.Image {
	output := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-(float64)(startX), -(float64)(startY))
	output.DrawImage(src, op)
	return output
}

func SpriteCutter(src *ebiten.Image, width, height, cols, rows int) []*ebiten.Image {
	output := []*ebiten.Image{}
	for y := range rows {
		for x := range cols {
			startX := x * width
			startY := y * height
			tempImage := SubImage(src, startX, startY, width, height)
			//fmt.Printf("Add image %v %v %v %v \n", startX, startY, width, height)
			output = append(output, tempImage)
		}
	}
	return output
}

func FlipImageVertical(src *ebiten.Image) *ebiten.Image {
	height := src.Bounds().Dy()
	width := src.Bounds().Dx()
	output := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, -(float64)(height))
	op.GeoM.Scale(0, -1)
	output.DrawImage(src, op)
	return output
}

func ScaleImage(src *ebiten.Image, width, height int) *ebiten.Image {
	heightO := (float64)(src.Bounds().Dy())
	widthO := (float64)(src.Bounds().Dx())
	heightN := (float64)(height)
	widthN := (float64)(width)
	heightF := (float64)(heightN / heightO)
	widthF := (float64)(widthN / widthO)
	//fmt.Println("Width ", widthF)
	output := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(widthF, heightF)
	output.DrawImage(src, op)
	return output
}

// https://go.dev/doc/tutorial/generics
func Clamp[T int | float64](min, max, test T) T {
	if test > max {
		return max
	} else if test < min {
		return min
	} else {
		return test
	}
}

func getListOfLinesFromFile(filePath string) []*string {
	//fmt.Printf("GLOLFF: \n")
	var linesList = []*string{}
	// open file handle
	file, err := os.Open(filePath)
	if nil != err {
		log.Fatal(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)

	// Section 2
	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
			//fmt.Printf("ReadLine: %q\n", line)
			lineAsString := string(line)
			linesList = append(linesList, &lineAsString)
		} else {
			break
		}
		if err != nil {
			break
		}
	}

	_ = err

	return linesList
}

func getListOfLinesFromBytes(bytesArray []byte) []*string {

	bytesAsString := string(bytesArray)
	myStrings := strings.Split(bytesAsString, UTILS_NEWLINE)
	var linesList = []*string{}

	for i, line := range myStrings {

		if len(line) > 0 {
			linesList = append(linesList, &myStrings[i])
		} else {
			continue
		}
	}

	return linesList
}

type Collider interface {
	// worldx, worldy width, height
	Dimensions() (int, int, int, int)
}

func Intersect(objA, objB Collider) bool {
	//true if colliding
	ax1, ay1, aw, ah := objA.Dimensions()
	bx1, by1, bw, bh := objB.Dimensions()
	ax2, ay2 := ax1+aw, ay1+ah
	bx2, by2 := bx1+bw, by1+bh
	if ax1 > bx2 || bx1 > ax2 || ay1 > by2 || by1 > ay2 {
		return false
	} else {

		return true
	}

}

func Pulser(tickPeriod int) func() bool {
	//returns true for a series of ticks, then false ,then repeat
	var currentTick = 0
	var tickMax = tickPeriod
	var state = false

	var checkTimeExpired = func() bool {

		if (currentTick) > tickMax {
			currentTick = 0
			state = !state
		} else {
			currentTick += 1
		}
		return state

	}

	return checkTimeExpired
}
