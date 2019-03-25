package textgeometry

import (
	"log"
	"strings"

	wordwrap "github.com/mitchellh/go-wordwrap"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
)

var (
	magicFactor float64
)

func init() {

	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	calcFont := truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	// Calculate the factor X in the equation FontSize * X * characters = pixelWidth
	longText := "Your crew boards the station, cautiously moving between corridors. Suddenly a man-sized arachnid bursts from a vent in the ceiling, followed by countless more. You fight your way back to the airlock and are forced to leave before accounting for all crew members. Not everybody made it back."
	_, advance := font.BoundString(calcFont, longText)
	magicFactor = float64(advance.Ceil()/calcFont.Metrics().Height.Ceil()) / float64(len(longText))

	// Todo: replace with magicFactor = 0.47645, or just 0.5?
	// Do we gain anything from calculating this if we do binary search after?
}

// PredictChars estimates the number of the max number of characters allowed to fill the width of maxPixels
func PredictChars(fnt font.Face, maxPixels int) int {
	return int(float64(maxPixels) / (float64(fnt.Metrics().Height.Ceil()) * magicFactor))
}

// LinesMaxWidthPixels calculates the max width in pixels for given lines and font
func LinesMaxWidthPixels(lines []string, fnt font.Face) float64 {
	max := -1.0
	for _, line := range lines {
		advance := float64(font.MeasureString(fnt, line).Ceil())
		if advance > max {
			max = advance
		}
	}
	return max
}

// BoundingBox calculates the width and height of the bounding box of a text with specified lines and font
// Uses a fixed line height from font.Metrics().Height
func BoundingBox(lines []string, fnt font.Face) (int, int) {
	height := int(fnt.Metrics().Height.Ceil()) * len(lines)
	width := -1
	for _, line := range lines {
		advance := font.MeasureString(fnt, line).Ceil()
		if advance > width {
			width = advance
		}
	}
	return width, height
}

func MaxWrapPosition(text string, fnt font.Face, maxPixels int) int {
	prediction := uint(PredictChars(fnt, maxPixels))
	for lim := prediction + 10; lim > prediction-10; lim-- {
		lines := strings.Split(wordwrap.WrapString(text, lim), "\n")
		w, _ := BoundingBox(lines, fnt)
		if w < maxPixels {
			return int(lim)
		}
	}
	return -1
}
