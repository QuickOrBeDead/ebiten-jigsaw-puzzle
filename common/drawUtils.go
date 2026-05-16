package common

import (
	"image/color"
	"math"

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

// DrawProgressBar draws a 3D beveled progress bar
func DrawProgressBar(screen *ebiten.Image, x, y, w, h float32, progress float32, bgColor, fillColor, borderColor color.Color) {
	bg := getRGBA(bgColor)
	fg := getRGBA(fillColor)

	// Drop shadow
	DrawSoftShadow(screen, x, y, w, h, color.RGBA{0, 0, 0, 60}, 2, 2)

	// Track background
	vector.FillRect(screen, x, y, w, h, bgColor, true)

	// Track inset bevel (recessed groove: dark top, light bottom)
	vector.FillRect(screen, x, y, w, 2, darkerColor(bg, 0.15), true)
	vector.FillRect(screen, x, y+h-2, w, 2, lighterColor(bg, 0.15), true)

	// Fill
	fillWidth := w * progress
	if fillWidth > 0 {
		vector.FillRect(screen, x, y, fillWidth, h, fillColor, true)

		// Fill raised bevel (light top/left, dark bottom/right)
		vector.FillRect(screen, x, y, fillWidth, 2, lighterColor(fg, 0.2), true)
		vector.FillRect(screen, x, y+h-2, fillWidth, 2, darkerColor(fg, 0.15), true)
		vector.FillRect(screen, x, y, 2, h, lighterColor(fg, 0.1), true)
	}

	// Border
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
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

func lighterColor(c color.RGBA, factor float32) color.RGBA {
	return color.RGBA{
		R: uint8(math.Min(255, float64(c.R)+(255-float64(c.R))*float64(factor))),
		G: uint8(math.Min(255, float64(c.G)+(255-float64(c.G))*float64(factor))),
		B: uint8(math.Min(255, float64(c.B)+(255-float64(c.B))*float64(factor))),
		A: c.A,
	}
}

func darkerColor(c color.RGBA, factor float32) color.RGBA {
	return color.RGBA{
		R: uint8(float32(c.R) * (1 - factor)),
		G: uint8(float32(c.G) * (1 - factor)),
		B: uint8(float32(c.B) * (1 - factor)),
		A: c.A,
	}
}
