package home

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
)

type HomeScene struct {
	playBtn      *common.Button
	howToPlayBtn *common.Button
	settingsBtn  *common.Button
	creditsBtn   *common.Button
	quitBtn      *common.Button
	text         *common.TextRenderer
	background   *ebiten.Image
}

func NewHomeScene(context *common.SceneContext) *HomeScene {
	text := common.NewTextRenderer(common.RobotoBoldFontName, common.TitleColor, 40, etxt.Center)

	gapY := float32(20)

	startX := (float32(common.ScreenWidth) - 260) / 2

	topY := float32(170)

	btnOpts := []common.ButtonOptFunc{
		common.ButtonOption.WithFontSize(24),
		common.ButtonOption.WithFontColor(common.BodyTextColor),
	}

	return &HomeScene{
		text:       text,
		background: drawBackground(text),
		playBtn: common.NewButton(
			common.ButtonTypeNormal, common.ButtonColorPrimary, common.ButtonSizeBig,
			startX, topY,
			"Play",
			append(btnOpts,
				common.ButtonOption.WithOnClick(func() {
					context.SceneManager.SetScene("startGame")
				}),
			)...,
		),
		howToPlayBtn: common.NewButton(
			common.ButtonTypeNormal, common.ButtonColorSecondary, common.ButtonSizeBig,
			startX, topY+(60+gapY)*1,
			"How to Play",
			append(btnOpts,
				common.ButtonOption.WithOnClick(func() {
					context.SceneManager.SetScene("howToPlay")
				}),
			)...,
		),
		settingsBtn: common.NewButton(
			common.ButtonTypeNormal, common.ButtonColorSecondary, common.ButtonSizeBig,
			startX, topY+(60+gapY)*2,
			"Settings",
			append(btnOpts,
				common.ButtonOption.WithOnClick(func() {
					context.SceneManager.SetScene("settings")
				}),
			)...,
		),
		creditsBtn: common.NewButton(
			common.ButtonTypeNormal, common.ButtonColorSecondary, common.ButtonSizeBig,
			startX, topY+(60+gapY)*3,
			"Credits",
			append(btnOpts,
				common.ButtonOption.WithOnClick(func() {
					context.SceneManager.SetScene("credits")
				}),
			)...,
		),
		quitBtn: common.NewButton(
			common.ButtonTypeNormal, common.ButtonColorSecondary, common.ButtonSizeBig,
			startX, topY+(60+gapY)*4,
			"Quit",
			btnOpts...,
		),
	}
}

func drawBackground(text *common.TextRenderer) *ebiten.Image {
	background := ebiten.NewImage(common.ScreenWidth, common.ScreenHeight)
	background.Fill(common.BackgroundColor)

	const headerH float32 = 120
	common.DrawPanel(background, 0, 0, float32(common.ScreenWidth), headerH, common.HeaderColor, false, nil)

	vector.FillRect(background, 0, 0, float32(common.ScreenWidth), 1, color.RGBA{255, 255, 255, 14}, false)

	common.DrawPanel(background, 0, headerH-3, float32(common.ScreenWidth), 3, common.PrimaryColor, false, nil)

	for i := 0; i < 4; i++ {
		alpha := uint8(35 - 7*i)
		vector.FillRect(background, 0, headerH+float32(i), float32(common.ScreenWidth), 1, color.RGBA{0, 0, 0, alpha}, false)
	}

	text.SetColor(common.TitleColor)
	text.SetSize(42)
	text.DrawEmbossedAutoWithShadow(background, "Jigsaw Puzzle", common.ScreenWidth/2, 50, color.RGBA{0, 0, 0, 100}, 2, 2)

	text.SetColor(common.MutedTextColor)
	text.SetSize(20)
	text.DrawEmbossedHozCenterWithShadow(background, "Choose an option to get started", 100, color.RGBA{0, 0, 0, 80}, 1, 1)
	return background
}

func (h *HomeScene) Update(context *common.SceneContext) error {
	h.playBtn.Update(context)
	h.howToPlayBtn.Update(context)
	h.settingsBtn.Update(context)
	h.creditsBtn.Update(context)
	h.quitBtn.Update(context)
	if h.quitBtn.Clicked() {
		return ebiten.Termination
	}

	return nil
}

func (h *HomeScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	screen.DrawImage(h.background, nil)
	h.playBtn.Draw(screen)
	h.howToPlayBtn.Draw(screen)
	h.settingsBtn.Draw(screen)
	h.creditsBtn.Draw(screen)
	h.quitBtn.Draw(screen)
}
