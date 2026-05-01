package common

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func CreateHitMask(width, height int, path *vector.Path, antiAlias bool) *image.Alpha {
	mask := CreateMask(width, height, path, antiAlias)
	h := CreateImageAlpha(mask)
	mask.Deallocate()
	return h
}

func CreateMask(width, height int, path *vector.Path, antiAlias bool) *ebiten.Image {
	mask := ebiten.NewImage(width, height)
	mask.Fill(color.Transparent)

	fillOpts := &vector.FillOptions{
		FillRule: vector.FillRuleEvenOdd,
	}
	drawOpts := &vector.DrawPathOptions{}
	drawOpts.AntiAlias = antiAlias
	drawOpts.ColorScale.ScaleWithColor(color.White)
	vector.FillPath(mask, path, fillOpts, drawOpts)
	return mask
}

func CreateImageAlpha(img *ebiten.Image) *image.Alpha {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	pixels := make([]byte, w*h*4)
	img.ReadPixels(pixels)

	a := image.NewAlpha(b)
	for i := 0; i < w*h; i++ {
		a.Pix[i] = pixels[i*4+3]
	}

	return a
}

func FillPathWithColor(screen *ebiten.Image, path *vector.Path, x, y float32, c color.Color, antiAlias bool) {
	var newPath vector.Path
	op := &vector.AddPathOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	newPath.AddPath(path, op)

	fillOpts := &vector.FillOptions{
		FillRule: vector.FillRuleEvenOdd,
	}
	drawOpts := &vector.DrawPathOptions{}
	drawOpts.AntiAlias = antiAlias
	drawOpts.ColorScale.ScaleWithColor(c)
	vector.FillPath(screen, &newPath, fillOpts, drawOpts)
}

func StrokePathWithColor(screen *ebiten.Image, path *vector.Path, x, y, strokeWidth float32, c color.Color, antiAlias bool) {
	var newPath vector.Path
	op := &vector.AddPathOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	newPath.AddPath(path, op)

	strokeOpts := &vector.StrokeOptions{Width: strokeWidth}
	drawOpts := &vector.DrawPathOptions{}
	drawOpts.AntiAlias = antiAlias
	drawOpts.ColorScale.ScaleWithColor(c)
	vector.StrokePath(screen, &newPath, strokeOpts, drawOpts)
}
