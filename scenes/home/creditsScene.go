package home

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
)

type creditLineType int

const (
	creditBlank creditLineType = iota
	creditTitle
	creditSubtitle
	creditInfo
	creditSection
	creditItem
)

type creditLine struct {
	kind creditLineType
	text string
}

type CreditsScene struct {
	backBtn      *common.Button
	title        *common.TextRenderer
	body         *common.TextRenderer
	section      *common.TextRenderer
	item         *common.TextRenderer
	scrollY      float64
	contentLines []creditLine
	lineHeight   int
	isDragging   bool
}

const creditsLineHeight = 36

func newLine(kind creditLineType, text string) creditLine {
	return creditLine{kind: kind, text: text}
}

func NewCreditsScene(context *common.SceneContext) *CreditsScene {
	lines := []creditLine{
		newLine(creditTitle, "Jigsaw Puzzle"),
		newLine(creditBlank, ""),
		newLine(creditSubtitle, "A puzzle game built with Ebitengine."),
		newLine(creditBlank, ""),
		newLine(creditInfo, "Developer: Bora Akgün"),
		newLine(creditInfo, "Source Code: github.com/QuickOrBeDead/ebiten-jigsaw-puzzle"),
		newLine(creditBlank, ""),
		newLine(creditSection, "Libraries"),
		newLine(creditItem, "github.com/tinne26/etxt"),
		newLine(creditItem, "github.com/sqweek/dialog"),
		newLine(creditBlank, ""),
		newLine(creditSection, "Puzzle Images"),
		newLine(creditItem, "pixabay.com/photos/mountainous-mountain-landscape-5942962 – Image by Pixabay"),
		newLine(creditItem, "pixabay.com/photos/landscape-rice-terrace-5104510 – Image by CongVuphotographer"),
		newLine(creditItem, "pixabay.com/photos/house-4028391 – Image by Peggychoucair"),
		newLine(creditItem, "pixabay.com/photos/muhlviertel-7544316 – Image by Leonhard_Niederwimmer"),
		newLine(creditItem, "pixabay.com/photos/desert-4388204 – Image by grebmot"),
		newLine(creditItem, "pixabay.com/photos/sunrise-7591335 – Image by Nordseher"),
		newLine(creditBlank, ""),
		newLine(creditSection, "Fonts"),
		newLine(creditItem, "Roboto by Google"),
		newLine(creditBlank, ""),
		newLine(creditSection, "Sound Effects"),
		newLine(creditItem, "Interface Sounds Starter Pack by p0ss (opengameart.org) – CC-BY-SA 3.0"),
		newLine(creditItem, "Well Done by qubodup (opengameart.org) – CC0"),
		newLine(creditBlank, ""),
		newLine(creditSection, "Music"),
		newLine(creditItem, "Nocturne Op. 9 No. 2 by Frédéric Chopin – Performed by Peter Johnston – CC0 (musopen.org)"),
	}

	return &CreditsScene{
		title:        common.NewTextRenderer(common.RobotoBoldFontName, common.TitleColor, 32, etxt.Center),
		body:         common.NewTextRenderer(common.RobotoRegularFontName, common.BodyTextColor, 20, etxt.Center),
		section:      common.NewTextRenderer(common.RobotoBoldFontName, common.TitleColor, 20, etxt.Center),
		item:         common.NewTextRenderer(common.RobotoRegularFontName, common.BodyTextColor, 18, etxt.Center),
		contentLines: lines,
		lineHeight:   creditsLineHeight,
		backBtn: common.NewButton(
			common.ButtonTypeNormal, common.ButtonColorSecondary, common.ButtonSizeSmall,
			20, 12,
			"Back",
			common.ButtonOption.WithFontSize(18),
			common.ButtonOption.WithFontColor(common.BodyTextColor),
			common.ButtonOption.WithOnClick(func() {
				context.SceneManager.SetScene("home")
			}),
		),
	}
}

func (s *CreditsScene) totalContentHeight() int {
	return len(s.contentLines) * s.lineHeight
}

const creditsContentStartY = 120

func (s *CreditsScene) scrollBarBounds() (barX, barY, barW, barH float32) {
	return common.ScreenWidth - 12, 80, 6, common.ScreenHeight - 80
}

func (s *CreditsScene) maxScroll() float64 {
	viewportH := int(common.ScreenHeight - creditsContentStartY)
	maxS := s.totalContentHeight() - viewportH
	if maxS < 0 {
		return 0
	}
	return float64(maxS)
}

func (s *CreditsScene) thumbSize() (float64, float64) {
	_, barY, _, barH := s.scrollBarBounds()
	maxS := s.maxScroll()
	thumbH := float64(barH) * float64(barH) / float64(s.totalContentHeight())
	thumbY := float64(barY)
	if maxS > 0 {
		thumbY = float64(barY) + (float64(barH)-thumbH)*s.scrollY/maxS
	}
	return thumbH, thumbY
}

func (s *CreditsScene) Update(context *common.SceneContext) error {
	s.backBtn.Update(context)

	_, wy := ebiten.Wheel()
	if wy != 0 {
		maxS := s.maxScroll()
		s.scrollY -= wy * 20
		s.scrollY = math.Max(0, math.Min(maxS, s.scrollY))
	}

	mx, my := ebiten.CursorPosition()
	barX, barY, barW, barH := s.scrollBarBounds()

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.isDragging = false
	}

	if s.isDragging {
		maxS := s.maxScroll()
		_, thumbH := s.thumbSize()
		trackH := float64(barH) - thumbH
		if trackH > 0 {
			relY := float64(my) - float64(barY) - thumbH/2
			s.scrollY = math.Max(0, math.Min(maxS, relY/trackH*maxS))
		}
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		fmx, fmy := float32(mx), float32(my)
		if fmx >= barX && fmx < barX+barW && fmy >= barY && fmy < barY+barH {
			thumbH, thumbY := s.thumbSize()
			s.isDragging = true
			if fmy < float32(thumbY) || fmy >= float32(thumbY+thumbH) {
				maxS := s.maxScroll()
				trackH := float64(barH) - thumbH
				if trackH > 0 {
					relY := float64(my) - float64(barY) - thumbH/2
					s.scrollY = math.Max(0, math.Min(maxS, relY/trackH*maxS))
				}
			}
		}
	}

	return nil
}

func (s *CreditsScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	screen.Fill(common.BackgroundColor)

	const headerH float32 = 64
	common.DrawPanel(screen, 0, 0, float32(common.ScreenWidth), headerH, common.HeaderColor, false, nil)
	vector.FillRect(screen, 0, 0, float32(common.ScreenWidth), 1, color.RGBA{255, 255, 255, 14}, false)
	common.DrawPanel(screen, 0, headerH-3, float32(common.ScreenWidth), 3, common.PrimaryColor, false, nil)
	for i := 0; i < 4; i++ {
		alpha := uint8(35 - 7*i)
		vector.FillRect(screen, 0, headerH+float32(i), float32(common.ScreenWidth), 1, color.RGBA{0, 0, 0, alpha}, false)
	}

	s.title.SetSize(32)
	s.title.DrawEmbossedAutoWithShadow(screen, "Credits", common.ScreenWidth/2, 32, color.RGBA{0, 0, 0, 100}, 1, 1)

	s.backBtn.Draw(screen)

	cx := common.ScreenWidth / 2
	startY := creditsContentStartY - int(s.scrollY)
	for i, cl := range s.contentLines {
		y := startY + i*s.lineHeight
		if y < 60 || y > common.ScreenHeight {
			continue
		}
		switch cl.kind {
		case creditTitle:
			s.title.SetSize(28)
			s.title.DrawEmbossedAutoWithShadow(screen, cl.text, cx, y, color.RGBA{0, 0, 0, 80}, 1, 1)
		case creditSection:
			s.section.SetSize(20)
			s.section.DrawEmbossedAutoWithShadow(screen, cl.text, cx, y, color.RGBA{0, 0, 0, 80}, 1, 1)
		case creditItem:
			s.item.SetSize(18)
			s.item.DrawEmbossedAutoWithShadow(screen, "• "+cl.text, cx, y, color.RGBA{0, 0, 0, 80}, 1, 1)
		default:
			s.body.SetSize(20)
			s.body.DrawEmbossedAutoWithShadow(screen, cl.text, cx, y, color.RGBA{0, 0, 0, 80}, 1, 1)
		}
	}

	if s.maxScroll() > 0 {
		barX, barY, barW, barH := s.scrollBarBounds()
		thumbH, thumbY := s.thumbSize()

		overScroll := false
		mx, my := ebiten.CursorPosition()
		fmx, fmy := float32(mx), float32(my)
		if fmx >= barX && fmx < barX+barW && fmy >= barY && fmy < barY+barH {
			overScroll = true
		}
		if overScroll {
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
		}

		vector.FillRect(screen, barX, barY, barW, barH, color.RGBA{255, 255, 255, 30}, true)
		thumbColor := common.PrimaryColor
		if s.isDragging {
			thumbColor = common.PrimaryHoverColor
		}
		vector.FillRect(screen, barX, float32(thumbY), barW, float32(thumbH), thumbColor, true)
	}
}
