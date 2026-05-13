package home

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
)

type HomeScene struct {
	playBtn       *common.Button
	howToPlayBtn  *common.Button
	settingsBtn   *common.Button
	creditsBtn    *common.Button
	text          *common.TextRenderer
}

func NewHomeScene() *HomeScene {
	text := common.NewTextRenderer(common.RobotoBoldFontName, common.TitleColor, 40, etxt.Center)

	btnW := float32(260)
	btnH := float32(60)
	gapY := float32(20)

	startX := (float32(common.ScreenWidth) - btnW) / 2

	topY := float32(170)

	btnOpts := []common.ButtonOptFunc{
		common.ButtonOption.WithFontSize(24),
		common.ButtonOption.WithFontColor(common.BodyTextColor),
		common.ButtonOption.WithColor(common.SurfaceColor),
		common.ButtonOption.WithHoverColor(common.SurfaceHoverColor),
		common.ButtonOption.WithShadowColor(common.ShadowColor),
	}

	return &HomeScene{
		text: text,
		playBtn: common.NewButton(
			startX, topY,
			btnW, btnH,
			"Play",
			append(btnOpts,
				common.ButtonOption.WithColor(common.PrimaryColor),
				common.ButtonOption.WithHoverColor(common.PrimaryHoverColor),
			)...,
		),
		howToPlayBtn: common.NewButton(
			startX, topY+(btnH+gapY)*1,
			btnW, btnH,
			"How to Play",
			btnOpts...,
		),
		settingsBtn: common.NewButton(
			startX, topY+(btnH+gapY)*2,
			btnW, btnH,
			"Settings",
			btnOpts...,
		),
		creditsBtn: common.NewButton(
			startX, topY+(btnH+gapY)*3,
			btnW, btnH,
			"Credits",
			btnOpts...,
		),
	}
}

func (h *HomeScene) Update(context *common.SceneContext) error {
	h.playBtn.Update()
	if h.playBtn.Clicked {
		context.SceneManager.SetScene("startGame")
		return nil
	}

	h.howToPlayBtn.Update()
	if h.howToPlayBtn.Clicked {
		context.SceneManager.SetScene("howToPlay")
		return nil
	}

	h.settingsBtn.Update()
	if h.settingsBtn.Clicked {
		context.SceneManager.SetScene("settings")
		return nil
	}

	h.creditsBtn.Update()
	if h.creditsBtn.Clicked {
		context.SceneManager.SetScene("credits")
		return nil
	}

	return nil
}

func (h *HomeScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	screen.Fill(common.BackgroundColor)

	const headerH float32 = 120
	common.DrawPanel(screen, 0, 0, float32(common.ScreenWidth), headerH, common.HeaderColor, false, nil)

	vector.FillRect(screen, 0, 0, float32(common.ScreenWidth), 1, color.RGBA{255, 255, 255, 14}, false)

	common.DrawPanel(screen, 0, headerH-3, float32(common.ScreenWidth), 3, common.PrimaryColor, false, nil)

	for i := 0; i < 4; i++ {
		alpha := uint8(35 - 7*i)
		vector.FillRect(screen, 0, headerH+float32(i), float32(common.ScreenWidth), 1, color.RGBA{0, 0, 0, alpha}, false)
	}

	h.text.SetColor(common.TitleColor)
	h.text.SetSize(42)
	h.text.DrawEmbossedAutoWithShadow(screen, "Jigsaw Puzzle", common.ScreenWidth/2, 50, color.RGBA{0, 0, 0, 100}, 2, 2)

	h.text.SetColor(common.MutedTextColor)
	h.text.SetSize(20)
	h.text.DrawEmbossedHozCenterWithShadow(screen, "Choose an option to get started", 100, color.RGBA{0, 0, 0, 80}, 1, 1)

	isHovered := false
	for _, btn := range []*common.Button{h.playBtn, h.howToPlayBtn, h.settingsBtn, h.creditsBtn} {
		btn.Draw(screen)
		if btn.Hovered {
			isHovered = true
		}
	}

	if isHovered {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
}
