package common

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type PreviewImageOptFunc func(*PreviewImage)

type PreviewImageOptionBuilder struct{}

var PreviewImageOption = PreviewImageOptionBuilder{}

func (PreviewImageOptionBuilder) WithBGColor(c color.Color) PreviewImageOptFunc {
	return func(p *PreviewImage) {
		p.BGColor = c
	}
}

func (PreviewImageOptionBuilder) WithBorderColor(c color.Color) PreviewImageOptFunc {
	return func(p *PreviewImage) {
		p.BorderColor = c
	}
}

func (PreviewImageOptionBuilder) WithCaption(text *TextRenderer, captionText string, color color.Color) PreviewImageOptFunc {
	return func(p *PreviewImage) {
		p.text = text
		p.captionText = captionText
		p.captionColor = color
	}
}

type PreviewImage struct {
	Image            *ebiten.Image
	X, Y             float64
	Scale            float64
	ScaledW, ScaledH float32
	BGColor          color.Color
	BorderColor      color.Color
	text             *TextRenderer
	captionText      string
	captionColor     color.Color
}

func NewPreviewImageWithMax(img *ebiten.Image, x, y, maxW, maxH float64, opts ...PreviewImageOptFunc) *PreviewImage {
	imgW := float64(img.Bounds().Dx())
	imgH := float64(img.Bounds().Dy())

	scale := math.Min(maxW/imgW, maxH/imgH)
	if scale > 1 {
		scale = 1
	}

	return NewPreviewImage(img, x, y, scale, opts...)
}

func NewPreviewImage(img *ebiten.Image, x, y, scale float64, opts ...PreviewImageOptFunc) *PreviewImage {
	bounds := img.Bounds()
	pi := &PreviewImage{
		Image:       img,
		X:           x,
		Y:           y,
		Scale:       scale,
		ScaledW:     float32(float64(bounds.Dx()) * scale),
		ScaledH:     float32(float64(bounds.Dy()) * scale),
		BGColor:     SurfaceColor,
		BorderColor: PrimaryColor,
	}

	for _, opt := range opts {
		opt(pi)
	}

	return pi
}

func (p *PreviewImage) Draw(screen *ebiten.Image) {
	if p.Image == nil {
		return
	}

	bounds := p.Image.Bounds()
	scaledW := float32(float64(bounds.Dx()) * p.Scale)
	scaledH := float32(float64(bounds.Dy()) * p.Scale)

	DrawSoftShadow(screen, float32(p.X), float32(p.Y), scaledW, scaledH, color.RGBA{0, 0, 0, 180}, 5, 5)
	DrawPanel(screen,
		float32(p.X)-4, float32(p.Y)-4,
		scaledW+8, scaledH+8,
		p.BGColor, true, p.BorderColor)

	geoM := ebiten.GeoM{}
	geoM.Scale(p.Scale, p.Scale)
	geoM.Translate(p.X, p.Y)
	opt := &ebiten.DrawImageOptions{GeoM: geoM, Filter: ebiten.FilterLinear}
	screen.DrawImage(p.Image, opt)

	if p.text != nil && len(p.captionText) > 0 {
		p.text.SetColor(p.captionColor)
		p.text.SetSize(14)
		p.text.SetAlign(etxt.Center)
		p.text.Draw(screen, p.captionText, int(p.X+float64(p.ScaledW)/2), int(p.Y+float64(p.ScaledH))+20)
	}
}

func (p *PreviewImage) IsPointInImage(x, y float64) bool {
	return x >= p.X && x <= p.X+float64(p.ScaledW) &&
		y >= p.Y && y <= p.Y+float64(p.ScaledH)
}
