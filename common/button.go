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

func (ButtonOptionBuilder) WithOnClick(onClickFunc func()) ButtonOptFunc {
	return func(b *Button) {
		b.OnClick = onClickFunc
	}
}

type Button struct {
	X            float32
	Y            float32
	HoverCursor  ebiten.CursorShapeType
	OnClick      func()
	FontSize     float64
	FontColor    color.Color
	FontName     string
	width        float32
	height       float32
	Label        string
	pressed      bool
	hovered      bool
	clicked      bool
	isToggle     bool
	isActive     bool
	borderColor  color.RGBA
	borderWidth  float32
	cornerRadius float32
	btnSize      ButtonSize
	btnColor     ButtonColor
	btnType      ButtonType
	text         *TextRenderer
	hitMask      *image.Alpha
	cacheItem    *buttonCacheItem
}

func NewButton(typ ButtonType, clr ButtonColor, size ButtonSize, x, y float32, label string, opts ...ButtonOptFunc) *Button {
	fontSize := 16.0
	fontColor := color.Black
	fontName := RobotoBoldFontName

	btnItem := buttonCacheManager.addImage(size, clr, typ)
	btn := &Button{
		X:            x,
		Y:            y,
		width:        btnItem.width,
		height:       btnItem.height,
		Label:        label,
		HoverCursor:  ebiten.CursorShapePointer,
		FontName:     fontName,
		FontSize:     fontSize,
		FontColor:    fontColor,
		cornerRadius: btnItem.cornerRadius,
		borderWidth:  0,
		isToggle:     typ == ButtonTypeToggle,
		isActive:     false,
		btnSize:      size,
		btnColor:     ButtonColorPrimary,
		btnType:      typ,
	}

	for _, opt := range opts {
		opt(btn)
	}

	btn.text = NewTextRenderer(btn.FontName, btn.FontColor, btn.FontSize, etxt.Center)
	btn.cacheItem = btnItem
	btn.hitMask = btnItem.hitMask

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
	isActive := b.btnType == ButtonTypeToggle && b.isActive
	img := b.currentImg(isActive, b.hovered, b.pressed)
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.X), float64(b.Y))
	screen.DrawImage(img, op)

	cx := int(b.X) + int(b.width/2) + int(b.cornerRadius)
	cy := int(b.Y) + int(b.height/2)
	if b.pressed {
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

	cx := int(x) + int(b.width/2) + int(b.cornerRadius)
	cy := int(y) + int(b.height/2)
	if pressed {
		cy += 2
	}
	b.text.DrawEmbossedAutoWithShadow(screen, b.Label, cx, cy, color.RGBA{0, 0, 0, 70}, 1, 1)
}

func (b *Button) Height() float32 {
	return b.height
}

func (b *Button) Clicked() bool {
	return b.clicked
}

func (b *Button) Update(context *SceneContext) {
	b.clicked = false
	b.hovered = b.isHovered()

	if b.hovered {
		context.Cursor = b.HoverCursor
	}

	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if pressed && b.hovered {
		b.pressed = true
	}

	if b.pressed && !pressed {
		if b.hovered {
			if b.btnType == ButtonTypeToggle {
				b.isActive = !b.isActive
			}
			b.clicked = true
			if AudioManager != nil {
				AudioManager.PlayClick()
			}
			if b.OnClick != nil {
				b.OnClick()
			}
		}
		b.pressed = false
	}
}

func (b *Button) Reset() {
	b.isActive = false
}

func (b *Button) isHovered() bool {
	mx, my := ebiten.CursorPosition()
	return b.hitMask.AlphaAt(mx-int(b.X), my-int(b.Y)).A > 0
}
