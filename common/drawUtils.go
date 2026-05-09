package common

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawShadowForPath(screen *ebiten.Image, x, y float64, path *vector.Path) {
	maxOffset := 2
	baseAlpha := 85
	for i := 1; i <= maxOffset; i++ {
		alpha := uint8(baseAlpha / i)
		offset := float64(i)
		shadowColor := color.RGBA{0, 0, 0, alpha}
		var newPath vector.Path
		op := &vector.AddPathOptions{}
		op.GeoM.Translate(x+offset, y+offset)
		newPath.AddPath(path, op)
		fillOpts := &vector.FillOptions{
			FillRule: vector.FillRuleEvenOdd,
		}
		drawOpts := &vector.DrawPathOptions{}
		drawOpts.AntiAlias = true
		drawOpts.ColorScale.ScaleWithColor(shadowColor)
		vector.FillPath(screen, &newPath, fillOpts, drawOpts)
	}
}

// DrawPanel draws a rounded rectangle panel with optional border
func DrawPanel(screen *ebiten.Image, x, y, w, h float32, bgColor color.Color, hasBorder bool, borderColor color.Color) {
	vector.FillRect(screen, x, y, w, h, bgColor, true)
	if hasBorder && borderColor != nil {
		vector.StrokeRect(screen, x, y, w, h, 1, borderColor, true)
	}
}

// DrawProgressBar draws a modern progress bar with rounded ends
func DrawProgressBar(screen *ebiten.Image, x, y, w, h float32, progress float32, bgColor, fillColor, borderColor color.Color) {
	vector.FillRect(screen, x, y, w, h, bgColor, true)
	fillWidth := w * progress
	if fillWidth > 0 {
		vector.FillRect(screen, x, y, fillWidth, h, fillColor, true)
	}
	if borderColor != nil {
		vector.StrokeRect(screen, x, y, w, h, 1, borderColor, true)
	}
}

// DrawSoftShadow draws a softer, more diffuse shadow behind a rectangle
func DrawSoftShadow(screen *ebiten.Image, x, y, w, h float32, shadowColor color.Color, offsetX, offsetY float32) {
	shadow := getRGBA(shadowColor)
	for i := 0; i < 4; i++ {
		alpha := uint8(float32(shadow.A) * (1.0 - float32(i)/4.0))
		color := color.RGBA{R: shadow.R, G: shadow.G, B: shadow.B, A: alpha}
		offset := float32(i + 1)
		vector.FillRect(screen, x+offsetX+offset, y+offsetY+offset, w, h, color, true)
	}
}

// Helper to extract RGBA values from color.Color
func getRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}
