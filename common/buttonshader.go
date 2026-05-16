package common

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type ButtonShaderConfig struct {
	SampleDist         float32
	LightDirX          float32
	LightDirY          float32
	HighlightIntensity float32
	ShadowIntensity    float32
}

var DefaultButtonShaderConfig = ButtonShaderConfig{
	SampleDist:         1.8,
	LightDirX:          -0.707,
	LightDirY:          -0.707,
	HighlightIntensity: 0.5,
	ShadowIntensity:    0.4,
}

const buttonShaderSrc string = `
//kage:unit pixels

package main

var Center vec2
var Size vec2
var CornerRadius float
var ButtonColor vec4
var LightDir vec2
var SampleDist float
var HighlightIntensity float
var ShadowIntensity float

func sdRoundedRect(p vec2, size vec2, r float) float {
	q := abs(p) - size*0.5 + r
	return length(max(q, 0.0)) + min(max(q.x, q.y), 0.0) - r
}

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	p := dstPos.xy - Center
	d := sdRoundedRect(p, Size, CornerRadius)
	a := 1.0 - smoothstep(0.0, 1.2, d)
	if a < 0.005 {
		return vec4(0)
	}

	d2 := SampleDist

	nx := (sdRoundedRect(p + vec2(d2, 0.0), Size, CornerRadius) -
		sdRoundedRect(p - vec2(d2, 0.0), Size, CornerRadius)) / (2.0 * d2)
	ny := (sdRoundedRect(p + vec2(0.0, d2), Size, CornerRadius) -
		sdRoundedRect(p - vec2(0.0, d2), Size, CornerRadius)) / (2.0 * d2)

	gm := sqrt(nx*nx + ny*ny)
	if gm < 0.001 {
		return vec4(ButtonColor.rgb, a)
	}

	nx /= gm
	ny /= gm
	facing := nx*LightDir.x + ny*LightDir.y

	edgeDist := -d
	bevelRadius := SampleDist * 1.2
	edgeWeight := clamp(1.0 - edgeDist/bevelRadius, 0.0, 1.0)

	strength := gm * gm * gm * edgeWeight * edgeWeight

	highlight := clamp(facing, 0.0, 1.0) * strength * HighlightIntensity
	shadow := clamp(-facing, 0.0, 1.0) * strength * ShadowIntensity

	r := ButtonColor.r + highlight - shadow
	g := ButtonColor.g + highlight - shadow
	b := ButtonColor.b + highlight - shadow

	return vec4(clamp(r, 0.0, 1.0), clamp(g, 0.0, 1.0), clamp(b, 0.0, 1.0), a)
}
`

var (
	buttonShader     *ebiten.Shader
	buttonShaderOnce sync.Once
)

func initButtonShader() {
	buttonShaderOnce.Do(func() {
		var err error
		buttonShader, err = ebiten.NewShader([]byte(buttonShaderSrc))
		if err != nil {
			panic(err)
		}
	})
}

func drawButtonBevel(screen *ebiten.Image, x, y, w, h, r float32, col color.RGBA, pressed bool, cfg ButtonShaderConfig) {
	initButtonShader()

	bw := w + 2*r
	bh := h

	colR := float32(col.R) / 255
	colG := float32(col.G) / 255
	colB := float32(col.B) / 255
	colA := float32(col.A) / 255

	if pressed {
		colR *= 0.75
		colG *= 0.75
		colB *= 0.75
	}

	op := &ebiten.DrawRectShaderOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.Uniforms = map[string]any{
		"Center":             []float32{x + bw*0.5, y + bh*0.5},
		"Size":               []float32{bw, bh},
		"CornerRadius":       r,
		"ButtonColor":        []float32{colR, colG, colB, colA},
		"LightDir":           []float32{cfg.LightDirX, cfg.LightDirY},
		"SampleDist":         cfg.SampleDist,
		"HighlightIntensity": cfg.HighlightIntensity,
		"ShadowIntensity":    cfg.ShadowIntensity,
	}
	screen.DrawRectShader(int(bw), int(bh), buttonShader, op)
}
