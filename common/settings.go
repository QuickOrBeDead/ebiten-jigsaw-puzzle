package common

import "image/color"

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	// Spawn area bounds for puzzle pieces
	SpawnAreaLeftMaxX     = 320 // Max X for left spawn column (320px from left)
	SpawnAreaRightMinX    = 960 // Min X for right spawn column (ScreenWidth - SpawnAreaLeftMaxX)
	SpawnAreaMinY         = 68  // Min Y for spawn area
	SpawnAreaBottomOffset = 60  // Offset from ScreenHeight for spawn area max Y
)

var (
	// Primary background colors - deep dark blue-gray theme
	BackgroundColor = color.RGBA{R: 18, G: 18, B: 24, A: 255} // #121218
	HeaderColor     = color.RGBA{R: 26, G: 29, B: 40, A: 255} // #1A1D28
	FooterColor     = color.RGBA{R: 26, G: 29, B: 40, A: 255} // #1A1D28

	// Surface colors for panels, cards, and UI elements
	SurfaceColor       = color.RGBA{R: 32, G: 32, B: 42, A: 255} // #20202A
	SurfaceHoverColor  = color.RGBA{R: 40, G: 40, B: 52, A: 255} // #282834
	SurfaceActiveColor = color.RGBA{R: 48, G: 48, B: 60, A: 255} // #30303C

	// Text colors
	TitleColor        = color.RGBA{R: 212, G: 168, B: 75, A: 255}  // #D4A84B - Warm amber/gold
	BodyTextColor     = color.RGBA{R: 226, G: 228, B: 235, A: 255} // #E2E4EB - Soft white
	MutedTextColor    = color.RGBA{R: 142, G: 146, B: 163, A: 255} // #8E92A3 - Secondary text
	DisabledTextColor = color.RGBA{R: 90, G: 90, B: 100, A: 255}   // #5A5A64 - Disabled state

	// Accent colors - professional blue
	PrimaryColor       = color.RGBA{R: 70, G: 140, B: 245, A: 255} // #468CF5
	PrimaryHoverColor  = color.RGBA{R: 92, G: 160, B: 255, A: 255} // #5CA0FF
	PrimaryActiveColor = color.RGBA{R: 55, G: 120, B: 230, A: 255} // #3778E6

	// Semantic colors
	SuccessColor = color.RGBA{R: 76, G: 175, B: 132, A: 255} // #4CAF84 - Teal-green
	WarningColor = color.RGBA{R: 255, G: 180, B: 60, A: 255} // #FFB43C
	ErrorColor   = color.RGBA{R: 225, G: 85, B: 85, A: 255}  // #E15555

	// Header button colors - slightly more contrast against header
	HeaderButtonColor       = color.RGBA{R: 42, G: 44, B: 56, A: 255} // #2A2C38
	HeaderButtonHoverColor  = color.RGBA{R: 58, G: 60, B: 72, A: 255} // #3A3C48
	HeaderButtonActiveColor = color.RGBA{R: 68, G: 70, B: 84, A: 255} // #444654

	// Shadow colors
	ShadowColor     = color.RGBA{R: 0, G: 0, B: 0, A: 40} // Subtle shadow
	ShadowColorDark = color.RGBA{R: 0, G: 0, B: 0, A: 80} // Darker shadow
)
