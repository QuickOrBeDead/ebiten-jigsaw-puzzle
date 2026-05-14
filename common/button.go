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

func (ButtonOptionBuilder) WithOnClick(onClickFunc func()) ButtonOptFunc {
	return func(b *Button) {
		b.OnClick = onClickFunc
	}
}

type Button struct {
	X            float32
	Y            float32
	Width        float32
	Height       float32
	Label        string
	Color        color.RGBA
	HoverColor   color.RGBA
	ActiveColor  color.RGBA
	ShadowColor  color.RGBA
	Pressed      bool
	Hovered      bool
	HoverCursor  ebiten.CursorShapeType
	Clicked      bool
	OnClick      func()
	FontSize     float64
	FontColor    color.Color
	FontName     string
	IsToggle     bool
	IsActive     bool
	BorderColor  color.RGBA
	BorderWidth  float32
	Icon         *ebiten.Image
	CornerRadius float32
	text         *TextRenderer
	path         *vector.Path
	hitMask      *image.Alpha
	fgImg        [6]*ebiten.Image
}

const (
	stateNormal int = iota
	stateHovered
	statePressed
	stateActive
	stateActiveHovered
	stateActivePressed
)

func stateIndex(isActive, hovered, pressed bool) int {
	switch {
	case pressed:
		if isActive {
			return stateActivePressed
		}
		return statePressed
	case hovered:
		if isActive {
			return stateActiveHovered
		}
		return stateHovered
	default:
		if isActive {
			return stateActive
		}
		return stateNormal
	}
}

func btnShadowOffset(isPressed bool) int {
	if isPressed {
		return 2
	}
	return 7
}

func NewButton(x, y, width, height float32, label string, opts ...ButtonOptFunc) *Button {
	fontSize := 16.0
	fontColor := color.Black
	fontName := RobotoBoldFontName
	cornerRadius := height / 2

	btn := &Button{
		X:            x,
		Y:            y,
		Width:        width,
		Height:       height,
		Label:        label,
		HoverCursor:  ebiten.CursorShapePointer,
		Color:        color.RGBA{R: 54, G: 153, B: 255, A: 255},
		HoverColor:   color.RGBA{R: 72, G: 176, B: 255, A: 255},
		ActiveColor:  color.RGBA{R: 40, G: 130, B: 240, A: 255},
		ShadowColor:  color.RGBA{R: 0, G: 0, B: 0, A: 50},
		FontName:     fontName,
		FontSize:     fontSize,
		FontColor:    fontColor,
		CornerRadius: cornerRadius,
		BorderWidth:  0,
		IsToggle:     false,
		IsActive:     false,
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

	imgW := int(width + 2*r + 7)
	imgH := int(height + 7)

	type stateSpec struct {
		idx      int
		hovered  bool
		pressed  bool
		isActive bool
	}

	states := []stateSpec{
		{stateNormal, false, false, false},
		{stateHovered, true, false, false},
		{statePressed, true, true, false},
	}

	if btn.IsToggle {
		states = append(states,
			stateSpec{stateActive, false, false, true},
			stateSpec{stateActiveHovered, true, false, true},
			stateSpec{stateActivePressed, true, true, true},
		)
	}

	for _, s := range states {
		img := ebiten.NewImage(imgW, imgH)
		maxOff := btnShadowOffset(s.pressed)
		baseAlpha := 35.0
		if s.pressed {
			baseAlpha = 15
		}
		for i := 1; i <= maxOff; i++ {
			alpha := uint8(baseAlpha / (1.0 + float64(i)/float64(maxOff)))
			off := float32(i)
			c := color.RGBA{btn.ShadowColor.R, btn.ShadowColor.G, btn.ShadowColor.B, alpha}
			StrokePathWithColor(img, btn.path, off, off, 2, c, true)
		}

		nImg := ebiten.NewImage(imgW, imgH)
		nImg.DrawImage(img, nil)

		col := getColor(btn.Color, btn.ActiveColor, btn.HoverColor, btn.IsToggle, s.isActive, s.hovered)
		drawButtonBevel(nImg, 0, 0, btn.Width, btn.Height, btn.CornerRadius, col, s.pressed, DefaultButtonShaderConfig)

		if btn.BorderWidth > 0 {
			StrokePathWithColor(nImg, btn.path, 0, 0, btn.BorderWidth, btn.BorderColor, true)
		}
		if btn.IsToggle && s.isActive {
			StrokePathWithColor(nImg, btn.path, 0, 0, 2, color.RGBA{100, 200, 255, 100}, true)
		}

		cx := int(btn.Width/2) + int(btn.CornerRadius)
		cy := int(btn.Height / 2)
		if s.pressed {
			cy += 2
		}
		btn.text.DrawEmbossedAutoWithShadow(nImg, btn.Label, cx, cy, color.RGBA{0, 0, 0, 70}, 1, 1)

		if btn.Icon != nil {
			iconHeight := btn.Height * 0.6
			iconX := float32(10)
			iconY := (btn.Height - iconHeight) / 2
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(iconHeight)/float64(btn.Icon.Bounds().Dy()), float64(iconHeight)/float64(btn.Icon.Bounds().Dy()))
			op.GeoM.Translate(float64(iconX), float64(iconY))
			nImg.DrawImage(btn.Icon, op)
		}

		btn.fgImg[s.idx] = nImg
	}

	return btn
}

func (b *Button) Draw(screen *ebiten.Image) {
	isActive := b.IsToggle && b.IsActive
	idx := stateIndex(isActive, b.Hovered, b.Pressed)

	fgImg := b.fgImg[idx]

	if fgImg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(b.X), float64(b.Y))
		screen.DrawImage(fgImg, op)
	}
}

func getColor(buttonColor, activeColor, hoverColor color.RGBA, isToggle, isActive, isHovered bool) color.RGBA {
	baseColor := buttonColor
	if isToggle && isActive {
		baseColor = activeColor
	}

	hoverProgress := float32(0)
	if isHovered {
		hoverProgress = 1.
	}

	r := float32(baseColor.R) + (float32(hoverColor.R)-float32(baseColor.R))*hoverProgress
	g := float32(baseColor.G) + (float32(hoverColor.G)-float32(baseColor.G))*hoverProgress
	bl := float32(baseColor.B) + (float32(hoverColor.B)-float32(baseColor.B))*hoverProgress
	a := float32(baseColor.A) + (float32(hoverColor.A)-float32(baseColor.A))*hoverProgress
	currentColor := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(bl), A: uint8(a)}
	return currentColor
}

func (b *Button) Update(context *SceneContext) {
	b.Clicked = false
	b.Hovered = b.isHovered()

	if b.Hovered {
		context.Cursor = b.HoverCursor
	}

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
