package common

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Slider struct {
	X, Y, Width, Height   float32
	Min, Max, Step, Value int

	TrackColor       color.RGBA
	FillColor        color.RGBA
	ThumbColor       color.RGBA
	ThumbHoverColor  color.RGBA
	ThumbActiveColor color.RGBA

	trackHeight    float32
	thumbRadius    float32
	dragging       bool
	hovered        bool
	hoverProgress  float32
	activeProgress float32
}

func NewSlider(x, y, width, height float32, min, max, step int) *Slider {
	v := 60
	if v < min {
		v = min
	}
	if v > max {
		v = max
	}
	v = (v / step) * step

	return &Slider{
		X:                x,
		Y:                y,
		Width:            width,
		Height:           height,
		Min:              min,
		Max:              max,
		Step:             step,
		Value:            v,
		TrackColor:       color.RGBA{60, 60, 70, 255},
		FillColor:        PrimaryColor,
		ThumbColor:       color.RGBA{200, 200, 210, 255},
		ThumbHoverColor:  color.RGBA{230, 230, 240, 255},
		ThumbActiveColor: PrimaryColor,
		trackHeight:      8,
		thumbRadius:      12,
	}
}

func (s *Slider) centerY() float32 {
	return s.Y + s.Height/2
}

func (s *Slider) trackLeft() float32 {
	return s.X + s.thumbRadius
}

func (s *Slider) trackRight() float32 {
	return s.X + s.Width - s.thumbRadius
}

func (s *Slider) trackWidth() float32 {
	return s.trackRight() - s.trackLeft()
}

func (s *Slider) thumbCenterX() float32 {
	tw := s.trackWidth()
	if tw <= 0 {
		return s.trackLeft()
	}
	ratio := float32(s.Value-s.Min) / float32(s.Max-s.Min)
	return s.trackLeft() + ratio*tw
}

func (s *Slider) setValueFromX(mx float32) {
	tw := s.trackWidth()
	if tw <= 0 {
		s.Value = s.Min
		return
	}
	ratio := (mx - s.trackLeft()) / tw
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	rawValue := float32(s.Min) + ratio*float32(s.Max-s.Min)
	steppedValue := int(math.Round(float64(rawValue)/float64(s.Step))) * s.Step
	if steppedValue < s.Min {
		steppedValue = s.Min
	}
	if steppedValue > s.Max {
		steppedValue = s.Max
	}
	s.Value = steppedValue
}

func (s *Slider) Update() {
	mx, my := ebiten.CursorPosition()
	fmx, fmy := float32(mx), float32(my)

	thumbCX := s.thumbCenterX()
	thumbCY := s.centerY()
	dx := fmx - thumbCX
	dy := fmy - thumbCY
	dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	isOverThumb := dist <= s.thumbRadius+4

	isOverTrack := fmx >= s.X && fmx <= s.X+s.Width &&
		fmy >= s.Y && fmy <= s.Y+s.Height

	s.hovered = isOverThumb || (isOverTrack && !s.dragging)

	if !s.dragging && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if isOverThumb {
			s.dragging = true
		} else if isOverTrack {
			s.setValueFromX(fmx)
			newCX := s.thumbCenterX()
			ndx := fmx - newCX
			if float32(math.Sqrt(float64(ndx*ndx))) <= s.thumbRadius+4 {
				s.dragging = true
			}
		}
	}

	if s.dragging {
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.dragging = false
		} else {
			s.setValueFromX(fmx)
		}
	}

	target := float32(0.0)
	if s.hovered && !s.dragging {
		target = 1.0
	}
	s.hoverProgress += (target - s.hoverProgress) * 0.15
	if float64(math.Abs(float64(target-s.hoverProgress))) < 0.01 {
		s.hoverProgress = target
	}

	activeTarget := float32(0.0)
	if s.dragging {
		activeTarget = 1.0
	}
	s.activeProgress += (activeTarget - s.activeProgress) * 0.15
	if float64(math.Abs(float64(activeTarget-s.activeProgress))) < 0.01 {
		s.activeProgress = activeTarget
	}
}

func (s *Slider) Draw(screen *ebiten.Image) {
	cy := s.centerY()
	thumbCX := s.thumbCenterX()
	trackTop := cy - s.trackHeight/2

	drawTrack(screen, s.trackLeft(), trackTop, s.trackWidth(), s.trackHeight, s.TrackColor)

	if s.Value > s.Min {
		fillWidth := thumbCX - s.trackLeft()
		drawTrack(screen, s.trackLeft(), trackTop, fillWidth, s.trackHeight, s.FillColor)
	}

	thumbColor := s.ThumbColor
	if s.dragging {
		thumbColor = mixColors(s.ThumbHoverColor, s.ThumbActiveColor, s.activeProgress)
	} else if s.hovered {
		thumbColor = mixColors(s.ThumbColor, s.ThumbHoverColor, s.hoverProgress)
	}

	path := vector.Path{}
	path.MoveTo(thumbCX, cy)
	path.Arc(thumbCX, cy, s.thumbRadius, 0, 2*math.Pi, vector.Clockwise)
	path.Close()
	FillPathWithColor(screen, &path, 0, 0, thumbColor, true)

	StrokePathWithColor(screen, &path, 0, 0, 2, color.RGBA{255, 255, 255, 60}, true)
}

func drawTrack(screen *ebiten.Image, x, y, w, h float32, c color.RGBA) {
	vector.FillRect(screen, x, y, w, h, c, true)
}

func mixColors(a, b color.RGBA, t float32) color.RGBA {
	return color.RGBA{
		R: uint8(float32(a.R) + (float32(b.R)-float32(a.R))*t),
		G: uint8(float32(a.G) + (float32(b.G)-float32(a.G))*t),
		B: uint8(float32(a.B) + (float32(b.B)-float32(a.B))*t),
		A: uint8(float32(a.A) + (float32(b.A)-float32(a.A))*t),
	}
}
