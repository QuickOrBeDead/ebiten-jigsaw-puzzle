package common

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type buttonCache struct {
	items []*buttonCacheItem
	ids   map[buttonCacheKey]int
	index int
}

type buttonCacheItem struct {
	path                                                *vector.Path
	hitMask                                             *image.Alpha
	normalImg, hoveredImg, pressedImg                   *ebiten.Image
	normalActiveImg, hoveredActiveImg, pressedActiveImg *ebiten.Image
	width, height, cornerRadius                         float32
}

type buttonCacheKey struct {
	width, height                  int
	color, hoverColor, activeColor color.RGBA
	shadowColor, borderColor       color.RGBA
	borderWidth                    float32
}

var buttonCacheManager *buttonCache

func init() {
	buttonCacheManager = &buttonCache{
		items: []*buttonCacheItem{},
		ids:   make(map[buttonCacheKey]int),
		index: 0,
	}
}

func (b *buttonCache) addImage(width, height float32, clr, hoverClr, activeClr, shadowClr, borderClr color.RGBA, borderWidth float32) *buttonCacheItem {
	key := buttonCacheKey{
		width: int(width), height: int(height),
		color: clr, hoverColor: hoverClr, activeColor: activeClr,
		shadowColor: shadowClr, borderColor: borderClr,
		borderWidth: borderWidth,
	}
	if id, ok := b.ids[key]; ok {
		return b.items[id]
	}

	newId := b.index
	b.ids[key] = newId
	b.index++

	cornerRadius := height / 2
	r := cornerRadius

	path := &vector.Path{}
	path.MoveTo(r, height)
	path.Arc(r, height/2, r, -math.Pi*1.5, -math.Pi*0.5, vector.Clockwise)
	path.LineTo(r, 0)
	path.LineTo(width+r, 0)
	path.Arc(width+r, height/2, r, -math.Pi*0.5, math.Pi*0.5, vector.Clockwise)
	path.LineTo(width+r, height)
	path.Close()

	hitMask := CreateHitMask(int(width+2*r), int(height), path, false)

	imgW := int(width + 2*r + 7)
	imgH := int(height + 7)

	createImg := func(hovered, pressed bool, baseColor color.RGBA, isActive bool) *ebiten.Image {
		img := ebiten.NewImage(imgW, imgH)
		maxOff := btnShadowOffset(pressed)
		baseAlpha := 35.0
		if pressed {
			baseAlpha = 15
		}
		for j := 1; j <= maxOff; j++ {
			alpha := uint8(baseAlpha / (1.0 + float64(j)/float64(maxOff)))
			off := float32(j)
			c := color.RGBA{shadowClr.R, shadowClr.G, shadowClr.B, alpha}
			StrokePathWithColor(img, path, off, off, 2, c, true)
		}
		stateImg := ebiten.NewImage(imgW, imgH)
		stateImg.DrawImage(img, nil)
		col := getColor(baseColor, hoverClr, hovered)
		drawButtonBevel(stateImg, 0, 0, width, height, cornerRadius, col, pressed, DefaultButtonShaderConfig)
		resultImg := ebiten.NewImage(imgW, imgH)
		resultImg.DrawImage(stateImg, nil)
		if borderWidth > 0 {
			StrokePathWithColor(resultImg, path, 0, 0, borderWidth, borderClr, true)
		}
		if isActive {
			StrokePathWithColor(resultImg, path, 0, 0, 2, color.RGBA{100, 200, 255, 100}, true)
		}
		return resultImg
	}

	item := &buttonCacheItem{
		path:             path,
		hitMask:          hitMask,
		normalImg:        createImg(false, false, clr, false),
		hoveredImg:       createImg(true, false, clr, false),
		pressedImg:       createImg(true, true, clr, false),
		normalActiveImg:  createImg(false, false, activeClr, true),
		hoveredActiveImg: createImg(true, false, activeClr, true),
		pressedActiveImg: createImg(true, true, activeClr, true),
		cornerRadius:     cornerRadius,
		width:            width,
		height:           height,
	}

	b.items = append(b.items, item)
	return item
}

func getColor(buttonColor, hoverColor color.RGBA, isHovered bool) color.RGBA {
	baseColor := buttonColor

	hover := float32(0)
	if isHovered {
		hover = 1.
	}

	r := float32(baseColor.R) + (float32(hoverColor.R)-float32(baseColor.R))*hover
	g := float32(baseColor.G) + (float32(hoverColor.G)-float32(baseColor.G))*hover
	bl := float32(baseColor.B) + (float32(hoverColor.B)-float32(baseColor.B))*hover
	a := float32(baseColor.A) + (float32(hoverColor.A)-float32(baseColor.A))*hover
	currentColor := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(bl), A: uint8(a)}
	return currentColor
}
