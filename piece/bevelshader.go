package piece

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type BevelConfig struct {
	SampleDist         float32
	LightDirX          float32
	LightDirY          float32
	HighlightIntensity float32
	ShadowIntensity    float32
}

var DefaultBevelConfig = BevelConfig{
	SampleDist:         1.8,
	LightDirX:          -0.707,
	LightDirY:          -0.707,
	HighlightIntensity: 0.5,
	ShadowIntensity:    0.4,
}

const bevelShaderSrc string = `
//kage:unit pixels

package main

var LightDir vec2
var SampleDist float

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	src := imageSrc0At(srcPos)

	if src.a < 0.001 {
		return vec4(0, 0, 0, 0)
	}

	d := SampleDist
	aR := imageSrc0At(srcPos + vec2(d, 0.0)).a
	aL := imageSrc0At(srcPos - vec2(d, 0.0)).a
	aD := imageSrc0At(srcPos + vec2(0.0, d)).a
	aU := imageSrc0At(srcPos - vec2(0.0, d)).a
	gx := aR - aL
	gy := aD - aU
	best := gx*gx + gy*gy

	if best < 0.0001 {
		return src
	}

	gm := sqrt(best)

	nx := -gx / gm
	ny := -gy / gm

	facing := nx*LightDir.x + ny*LightDir.y

	strength := gm * gm * gm

	highlight := clamp(facing, 0.0, 1.0) * strength * 0.5
	shadow := clamp(-facing, 0.0, 1.0) * strength * 0.4

	r := src.r + highlight - shadow
	g := src.g + highlight - shadow
	b := src.b + highlight - shadow

	return vec4(clamp(r, 0.0, 1.0), clamp(g, 0.0, 1.0), clamp(b, 0.0, 1.0), src.a)
}
`

var (
	bevelShader     *ebiten.Shader
	bevelShaderOnce sync.Once
)

func initBevelShader() {
	bevelShaderOnce.Do(func() {
		var err error
		bevelShader, err = ebiten.NewShader([]byte(bevelShaderSrc))
		if err != nil {
			panic(err)
		}
	})
}

func applyBevelShader(img *ebiten.Image, cfg BevelConfig) {
	initBevelShader()

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if w == 0 || h == 0 {
		return
	}

	result := ebiten.NewImage(w, h)
	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = img
	op.Uniforms = map[string]any{
		"LightDir":   []float32{cfg.LightDirX, cfg.LightDirY},
		"SampleDist": cfg.SampleDist,
	}
	result.DrawRectShader(w, h, bevelShader, op)

	img.Clear()
	img.DrawImage(result, nil)
	result.Deallocate()
}
