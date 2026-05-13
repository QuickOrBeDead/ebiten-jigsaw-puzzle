package common

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
)

type ButtonOptFunc func(*Button)

type ButtonOptionBuilder struct{}

var ButtonOption = ButtonOptionBuilder{}

func (ButtonOptionBuilder) WithFontSize(fontSize float64) ButtonOptFunc {
	return func(b *Button) {
		b.FontSize = fontSize
	}
}

func (ButtonOptionBuilder) WithFontColor(fontColor color.Color) ButtonOptFunc {
	return func(b *Button) {
		b.FontColor = fontColor
	}
}

func (ButtonOptionBuilder) WithFontName(fontName string) ButtonOptFunc {
	return func(b *Button) {
		b.FontName = fontName
	}
}

func (ButtonOptionBuilder) WithColor(color color.RGBA) ButtonOptFunc {
	return func(b *Button) {
		b.Color = color
	}
}

func (ButtonOptionBuilder) WithHoverColor(hoverColor color.RGBA) ButtonOptFunc {
	return func(b *Button) {
		b.HoverColor = hoverColor
	}
}

func (ButtonOptionBuilder) WithActiveColor(activeColor color.RGBA) ButtonOptFunc {
	return func(b *Button) {
		b.ActiveColor = activeColor
	}
}

func (ButtonOptionBuilder) WithShadowColor(shadowColor color.RGBA) ButtonOptFunc {
	return func(b *Button) {
		b.ShadowColor = shadowColor
	}
}

func (ButtonOptionBuilder) WithToggle(isToggle bool) ButtonOptFunc {
	return func(b *Button) {
		b.IsToggle = isToggle
	}
}

func (ButtonOptionBuilder) WithBorder(borderColor color.RGBA, borderWidth float32) ButtonOptFunc {
	return func(b *Button) {
		b.BorderColor = borderColor
		b.BorderWidth = borderWidth
	}
}

func (ButtonOptionBuilder) WithCornerRadius(radius float32) ButtonOptFunc {
	return func(b *Button) {
		b.CornerRadius = radius
	}
}

func (ButtonOptionBuilder) WithIcon(icon *ebiten.Image) ButtonOptFunc {
	return func(b *Button) {
		b.Icon = icon
	}
}

type Button struct {
	X              float32
	Y              float32
	Width          float32
	Height         float32
	Label          string
	Color          color.RGBA
	HoverColor     color.RGBA
	ActiveColor    color.RGBA
	ShadowColor    color.RGBA
	Pressed        bool
	Hovered        bool
	Clicked        bool
	OnClick        func()
	FontSize       float64
	FontColor      color.Color
	FontName       string
	IsToggle       bool
	IsActive       bool
	BorderColor    color.RGBA
	BorderWidth    float32
	Icon           *ebiten.Image
	CornerRadius   float32
	text           *TextRenderer
	path           *vector.Path
	hitMask        *image.Alpha
	hoverProgress  float32
	activeProgress float32
}

func NewButton(x, y, width, height float32, label string, opts ...ButtonOptFunc) *Button {
	fontSize := 16.0
	fontColor := color.Black
	fontName := RobotoBoldFontName
	cornerRadius := height / 2

	btn := &Button{
		X:              x,
		Y:              y,
		Width:          width,
		Height:         height,
		Label:          label,
		Color:          color.RGBA{R: 54, G: 153, B: 255, A: 255},
		HoverColor:     color.RGBA{R: 72, G: 176, B: 255, A: 255},
		ActiveColor:    color.RGBA{R: 40, G: 130, B: 240, A: 255},
		ShadowColor:    color.RGBA{R: 0, G: 0, B: 0, A: 50},
		FontName:       fontName,
		FontSize:       fontSize,
		FontColor:      fontColor,
		CornerRadius:   cornerRadius,
		BorderWidth:    0,
		IsToggle:       false,
		IsActive:       false,
		hoverProgress:  0,
		activeProgress: 0,
	}

	for _, opt := range opts {
		opt(btn)
	}

	r := btn.CornerRadius
	path := vector.Path{}
	path.MoveTo(r, height)
	path.Arc(r, height/2, r, -math.Pi*1.5, -math.Pi*0.5, vector.Clockwise)
	path.LineTo(r, 0)
	path.LineTo(width+r, 0)
	path.Arc(width+r, height/2, r, -math.Pi*0.5, math.Pi*0.5, vector.Clockwise)
	path.LineTo(width+r, height)
	path.Close()

	btn.path = &path
	btn.hitMask = CreateHitMask(int(width+2*r), int(height), &path, false)
	btn.text = NewTextRenderer(btn.FontName, btn.FontColor, btn.FontSize, etxt.Center)

	return btn
}

func (b *Button) Draw(screen *ebiten.Image) {
	baseColor := b.Color
	if b.IsToggle && b.IsActive {
		baseColor = b.ActiveColor
	}

	hoverColor := b.HoverColor
	r := float32(baseColor.R) + (float32(hoverColor.R)-float32(baseColor.R))*b.hoverProgress
	g := float32(baseColor.G) + (float32(hoverColor.G)-float32(baseColor.G))*b.hoverProgress
	bl := float32(baseColor.B) + (float32(hoverColor.B)-float32(baseColor.B))*b.hoverProgress
	a := float32(baseColor.A) + (float32(hoverColor.A)-float32(baseColor.A))*b.hoverProgress
	currentColor := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(bl), A: uint8(a)}

	b.drawShadow(screen)
	drawButtonBevel(screen, b.X, b.Y, b.Width, b.Height, b.CornerRadius, currentColor, b.Pressed, DefaultButtonShaderConfig)

	if b.BorderWidth > 0 {
		StrokePathWithColor(screen, b.path, b.X, b.Y, b.BorderWidth, b.BorderColor, true)
	}

	if b.IsToggle && b.IsActive {
		highlightColor := color.RGBA{R: 100, G: 200, B: 255, A: 100}
		StrokePathWithColor(screen, b.path, b.X, b.Y, 2, highlightColor, true)
	}

	b.drawText(screen)

	if b.Icon != nil {
		iconHeight := b.Height * 0.6
		aspect := float64(b.Icon.Bounds().Dx()) / float64(b.Icon.Bounds().Dy())
		iconWidth := iconHeight * float32(aspect)
		_ = iconWidth
		iconX := b.X + 10
		iconY := b.Y + (b.Height-iconHeight)/2
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(float64(iconHeight)/float64(b.Icon.Bounds().Dy()), float64(iconHeight)/float64(b.Icon.Bounds().Dy()))
		opt.GeoM.Translate(float64(iconX), float64(iconY))
		screen.DrawImage(b.Icon, opt)
	}
}

func (b *Button) Update() {
	b.Clicked = false
	b.Hovered = b.isHovered()
	b.updateHoverAnimation()

	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if pressed && b.Hovered {
		b.Pressed = true
	}

	if b.Pressed && !pressed {
		if b.Hovered {
			if b.IsToggle {
				b.IsActive = !b.IsActive
			}
			b.Clicked = true
			if AudioManager != nil {
				AudioManager.PlayClick()
			}
			if b.OnClick != nil {
				b.OnClick()
			}
		}
		b.Pressed = false
	}
}

func (b *Button) updateHoverAnimation() {
	target := float32(0.0)
	if b.Hovered {
		target = 1.0
	}
	b.hoverProgress += (target - b.hoverProgress) * 0.15
	if math.Abs(float64(target-b.hoverProgress)) < 0.01 {
		b.hoverProgress = target
	}

	if b.IsToggle {
		activeTarget := float32(0.0)
		if b.IsActive {
			activeTarget = 1.0
		}
		b.activeProgress += (activeTarget - b.activeProgress) * 0.15
		if math.Abs(float64(activeTarget-b.activeProgress)) < 0.01 {
			b.activeProgress = activeTarget
		}
	}
}

func (b *Button) isHovered() bool {
	mx, my := ebiten.CursorPosition()
	return b.hitMask.AlphaAt(mx-int(b.X), my-int(b.Y)).A > 0
}

func (b *Button) drawShadow(screen *ebiten.Image) {
	maxOffset := 5
	baseAlpha := 35.0

	if b.Hovered {
		maxOffset = 7
		baseAlpha = 45
	}

	if b.Pressed {
		maxOffset = 2
		baseAlpha = 15
	}

	for i := 1; i <= maxOffset; i++ {
		alpha := uint8(baseAlpha / (1.0 + float64(i)/float64(maxOffset)))
		offset := float32(i)
		shadowColor := color.RGBA{b.ShadowColor.R, b.ShadowColor.G, b.ShadowColor.B, alpha}
		StrokePathWithColor(screen, b.path, b.X+offset, b.Y+offset, 2, shadowColor, true)
	}
}

func (b *Button) drawBackground(screen *ebiten.Image, bColor color.Color, offsetX, offsetY float32) {
	FillPathWithColor(screen, b.path, b.X+offsetX, b.Y+offsetY, bColor, true)
}

func (b *Button) drawText(screen *ebiten.Image) {
	textX, textY := b.getCenter()
	if b.Pressed {
		textY += 2
	}
	b.text.SetFont(b.FontName)
	b.text.SetSize(b.FontSize)
	b.text.SetAlign(etxt.Center)

	b.text.SetColor(b.FontColor)
	b.text.DrawEmbossedAutoWithShadow(screen, b.Label, textX, textY, color.RGBA{0, 0, 0, 70}, 1, 1)
}

func (b *Button) getCenter() (int, int) {
	r := int(b.CornerRadius)
	return int(b.X) + int(b.Width/2) + r, int(b.Y) + int(b.Height/2)
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
