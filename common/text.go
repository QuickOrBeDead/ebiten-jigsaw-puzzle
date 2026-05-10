package common

import (
	"embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"github.com/tinne26/etxt/cache"
	"github.com/tinne26/etxt/font"
	"golang.org/x/image/font/sfnt"
)

const (
	RobotoRegularFontName = "Roboto"
	RobotoBoldFontName    = "Roboto Bold"
)

type TextRenderer struct {
	renderer *etxt.Renderer
	font     *sfnt.Font
	color    color.Color
	align    etxt.Align
	size     float64
}

var (
	//go:embed fonts/*
	fonts        embed.FS
	textRenderer *etxt.Renderer
	fontLibrary  *font.Library
)

func init() {
	fontLibrary = font.NewLibrary()
	loaded, skipped, err := fontLibrary.ParseAllFromFS(fonts, "fonts")
	if err != nil {
		log.Fatalf("Error while loading fonts: %s", err.Error())
	}

	fontLibrary.EachFont(func(s string, f *etxt.Font) error {
		log.Printf("%s font is loaded\n", s)

		return nil
	})

	log.Printf("Loaded fonts: %d, Skipped fonts: %d\n", loaded, skipped)

	textRenderer = etxt.NewRenderer()
	glyphsCache := cache.NewDefaultCache(16 * 1024 * 1024) // 16MiB cache
	textRenderer.SetCacheHandler(glyphsCache.NewHandler())
	textRenderer.SetColor(color.RGBA{239, 91, 91, 255})
	textRenderer.SetAlign(etxt.Center)
	textRenderer.SetSize(32)
}

func NewTextRenderer(fontName string, color color.Color, size float64, align etxt.Align) *TextRenderer {
	font := getFont(fontName)

	r := &TextRenderer{
		renderer: textRenderer,
		font:     font,
		color:    color,
		align:    align,
		size:     size,
	}

	return r
}

func getFont(fontName string) *sfnt.Font {
	font := fontLibrary.GetFont(fontName)
	if font == nil {
		log.Fatalf("font '%s' not found\n", fontName)
	}
	return font
}

func (t *TextRenderer) DrawHorizontalCenter(target *ebiten.Image, text string, y int) {
	t.SetAlign(etxt.HorzCenter)
	t.Draw(target, text, target.Bounds().Dx()/2, y)
}

func (t *TextRenderer) DrawCenter(target *ebiten.Image, text string) {
	t.SetAlign(etxt.Center)
	t.Draw(target, text, target.Bounds().Dx()/2, target.Bounds().Dy()/2)
}

func (t *TextRenderer) Draw(target *ebiten.Image, text string, x, y int) {
	t.renderer.SetFont(t.font)
	t.renderer.SetColor(t.color)
	t.renderer.SetSize(t.size)
	t.renderer.SetAlign(t.align)
	t.renderer.Draw(target, text, x, y)
}

func (t *TextRenderer) SetColor(color color.Color) {
	t.color = color
}

func (t *TextRenderer) SetAlign(align etxt.Align) {
	t.align = align
}

func (t *TextRenderer) SetSize(size float64) {
	t.size = size
}

func (t *TextRenderer) SetFont(fontName string) {
	t.font = getFont(fontName)
}

// DrawEmbossedAutoWithShadow draws 3D embossed text with a drop shadow to ground it
func (t *TextRenderer) DrawEmbossedAutoWithShadow(target *ebiten.Image, text string, x, y int, dropColor color.Color, dropOffX, dropOffY int) {
	r, g, b, a := t.color.RGBA()
	fc := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	highlight := lighterColor(fc, 0.7)
	shadow := darkerColor(fc, 0.5)

	t.renderer.SetFont(t.font)
	t.renderer.SetSize(t.size)
	t.renderer.SetAlign(t.align)
	t.renderer.SetColor(dropColor)
	t.renderer.Draw(target, text, x+dropOffX, y+dropOffY)

	t.DrawEmbossed(target, text, x, y, highlight, shadow)
}

// DrawEmbossedHozCenterWithShadow draws horizontally centered 3D embossed text with a drop shadow
func (t *TextRenderer) DrawEmbossedHozCenterWithShadow(target *ebiten.Image, text string, y int, dropColor color.Color, dropOffX, dropOffY int) {
	t.SetAlign(etxt.HorzCenter)
	t.DrawEmbossedAutoWithShadow(target, text, target.Bounds().Dx()/2, y, dropColor, dropOffX, dropOffY)
}

// DrawEmbossed draws text with a 3D embossed effect: highlight top-left, shadow bottom-right
func (t *TextRenderer) DrawEmbossed(target *ebiten.Image, text string, x, y int, highlightColor, shadowColor color.Color) {
	t.renderer.SetFont(t.font)
	t.renderer.SetSize(t.size)
	t.renderer.SetAlign(t.align)

	t.renderer.SetColor(highlightColor)
	t.renderer.Draw(target, text, x-1, y-1)

	t.renderer.SetColor(shadowColor)
	t.renderer.Draw(target, text, x+1, y+1)

	t.renderer.SetColor(t.color)
	t.renderer.Draw(target, text, x, y)
}

// DrawEmbossedAuto draws 3D embossed text, auto-computing highlight/shadow from the text color
func (t *TextRenderer) DrawEmbossedAuto(target *ebiten.Image, text string, x, y int) {
	r, g, b, a := t.color.RGBA()
	fc := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	highlight := lighterColor(fc, 0.7)
	shadow := darkerColor(fc, 0.5)
	t.DrawEmbossed(target, text, x, y, highlight, shadow)
}

// DrawEmbossedHorizontalCenter draws horizontally centered 3D embossed text
func (t *TextRenderer) DrawEmbossedHorizontalCenter(target *ebiten.Image, text string, y int) {
	t.SetAlign(etxt.HorzCenter)
	t.DrawEmbossedAuto(target, text, target.Bounds().Dx()/2, y)
}

// DrawWithShadow draws text with a shadow effect
func (t *TextRenderer) DrawWithShadow(target *ebiten.Image, text string, x, y int, shadowColor color.Color, offsetX, offsetY int) {
	t.renderer.SetFont(t.font)
	t.renderer.SetColor(shadowColor)
	t.renderer.SetSize(t.size)
	t.renderer.SetAlign(t.align)
	t.renderer.Draw(target, text, x+offsetX, y+offsetY)

	t.renderer.SetColor(t.color)
	t.renderer.Draw(target, text, x, y)
}

// DrawWithOutline draws text with an outline effect
func (t *TextRenderer) DrawWithOutline(target *ebiten.Image, text string, x, y int, outlineColor color.Color, outlineWidth int) {
	t.renderer.SetFont(t.font)
	t.renderer.SetColor(outlineColor)
	t.renderer.SetSize(t.size)
	t.renderer.SetAlign(t.align)

	for dx := -outlineWidth; dx <= outlineWidth; dx++ {
		for dy := -outlineWidth; dy <= outlineWidth; dy++ {
			if dx != 0 || dy != 0 {
				t.renderer.Draw(target, text, x+dx, y+dy)
			}
		}
	}

	t.renderer.SetColor(t.color)
	t.renderer.Draw(target, text, x, y)
}
