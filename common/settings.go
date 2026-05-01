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
	HeaderColor     = color.RGBA{R: 24, G: 24, B: 32, A: 255} // #181820
	FooterColor     = color.RGBA{R: 24, G: 24, B: 32, A: 255} // #181820

	// Surface colors for panels, cards, and UI elements
	SurfaceColor       = color.RGBA{R: 32, G: 32, B: 42, A: 255} // #20202A
	SurfaceHoverColor  = color.RGBA{R: 40, G: 40, B: 52, A: 255} // #282834
	SurfaceActiveColor = color.RGBA{R: 48, G: 48, B: 60, A: 255} // #30303C

	// Text colors
	TitleColor        = color.RGBA{R: 80, G: 200, B: 120, A: 255}  // #FFD764 - Warmer golden
	BodyTextColor     = color.RGBA{R: 220, G: 220, B: 230, A: 255} // #DCDCE6 - Soft white
	MutedTextColor    = color.RGBA{R: 140, G: 140, B: 155, A: 255} // #8C8C9B - Secondary text
	DisabledTextColor = color.RGBA{R: 90, G: 90, B: 100, A: 255}   // #5A5A64 - Disabled state

	// Accent colors - professional blue
	PrimaryColor       = color.RGBA{R: 66, G: 135, B: 245, A: 255} // #4287F5
	PrimaryHoverColor  = color.RGBA{R: 86, G: 155, B: 255, A: 255} // #569BFF
	PrimaryActiveColor = color.RGBA{R: 56, G: 115, B: 225, A: 255} // #3873E1

	// Semantic colors
	SuccessColor = color.RGBA{R: 80, G: 200, B: 120, A: 255} // #50C878
	WarningColor = color.RGBA{R: 255, G: 180, B: 60, A: 255} // #FFB43C
	ErrorColor   = color.RGBA{R: 220, G: 80, B: 80, A: 255}  // #DC5050

	// Header button colors
	HeaderButtonColor       = color.RGBA{R: 40, G: 40, B: 52, A: 255} // #282834
	HeaderButtonHoverColor  = color.RGBA{R: 56, G: 56, B: 68, A: 255} // #383844
	HeaderButtonActiveColor = color.RGBA{R: 66, G: 66, B: 80, A: 255} // #424250

	// Shadow colors
	ShadowColor     = color.RGBA{R: 0, G: 0, B: 0, A: 40} // Subtle shadow
	ShadowColorDark = color.RGBA{R: 0, G: 0, B: 0, A: 80} // Darker shadow
)
