package common

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

func (ButtonOptionBuilder) WithSize(size ButtonSize) ButtonOptFunc {
	return func(b *Button) {
		b.BtnSize = size
	}
}

func (ButtonOptionBuilder) WithColorEnum(color ButtonColor) ButtonOptFunc {
	return func(b *Button) {
		b.BtnColor = color
	}
}

func (ButtonOptionBuilder) WithType(typ ButtonType) ButtonOptFunc {
	return func(b *Button) {
		b.BtnType = typ
	}
}

func (ButtonOptionBuilder) WithToggle(isToggle bool) ButtonOptFunc {
	return func(b *Button) {
		if isToggle {
			b.BtnType = ButtonTypeToggle
		} else {
			b.BtnType = ButtonTypeNormal
		}
	}
}

func (ButtonOptionBuilder) WithCornerRadius(radius float32) ButtonOptFunc {
	return func(b *Button) {
		b.CornerRadius = radius
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
	CornerRadius float32
	BtnSize      ButtonSize
	BtnColor     ButtonColor
	BtnType      ButtonType
	text         *TextRenderer
	hitMask      *image.Alpha
	cacheItem    *buttonCacheItem
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
		BtnSize:      resolveButtonSize(height),
		BtnColor:     ButtonColorPrimary,
		BtnType:      ButtonTypeNormal,
	}

	for _, opt := range opts {
		opt(btn)
	}

	btn.text = NewTextRenderer(btn.FontName, btn.FontColor, btn.FontSize, etxt.Center)
	btn.cacheItem = buttonCacheManager.addImage(btn.BtnSize, btn.BtnColor, btn.BtnType)
	btn.hitMask = btn.cacheItem.hitMask

	return btn
}

func (b *Button) currentImg(isActive, hovered, pressed bool) *ebiten.Image {
	if isActive {
		if pressed {
			return b.cacheItem.pressedActiveImg
		}
		if hovered {
			return b.cacheItem.hoveredActiveImg
		}
		return b.cacheItem.normalActiveImg
	}
	if pressed {
		return b.cacheItem.pressedImg
	}
	if hovered {
		return b.cacheItem.hoveredImg
	}
	return b.cacheItem.normalImg
}

func (b *Button) Draw(screen *ebiten.Image) {
	isActive := b.BtnType == ButtonTypeToggle && b.IsActive
	img := b.currentImg(isActive, b.Hovered, b.Pressed)
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.X), float64(b.Y))
	screen.DrawImage(img, op)

	cx := int(b.X) + int(b.Width/2) + int(b.CornerRadius)
	cy := int(b.Y) + int(b.Height/2)
	if b.Pressed {
		cy += 2
	}
	b.text.DrawEmbossedAutoWithShadow(screen, b.Label, cx, cy, color.RGBA{0, 0, 0, 70}, 1, 1)
}

func (b *Button) DrawTest(screen *ebiten.Image, isActive, hovered, pressed bool, x, y float64) {
	img := b.currentImg(isActive, hovered, pressed)
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(img, op)

	cx := int(x) + int(b.Width/2) + int(b.CornerRadius)
	cy := int(y) + int(b.Height/2)
	if pressed {
		cy += 2
	}
	b.text.DrawEmbossedAutoWithShadow(screen, b.Label, cx, cy, color.RGBA{0, 0, 0, 70}, 1, 1)
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
			if b.BtnType == ButtonTypeToggle {
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
