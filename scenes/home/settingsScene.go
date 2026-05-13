package home

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
)

type SettingsScene struct {
	backBtn *common.Button
	title   *common.TextRenderer
	body    *common.TextRenderer
}

func NewSettingsScene() *SettingsScene {
	return &SettingsScene{
		title: common.NewTextRenderer(common.RobotoBoldFontName, common.TitleColor, 32, etxt.Center),
		body:  common.NewTextRenderer(common.RobotoRegularFontName, common.BodyTextColor, 22, etxt.Left),
		backBtn: common.NewButton(
			20, 12,
			80, 40,
			"Back",
			common.ButtonOption.WithFontSize(18),
			common.ButtonOption.WithFontColor(common.BodyTextColor),
			common.ButtonOption.WithColor(common.HeaderButtonColor),
			common.ButtonOption.WithHoverColor(common.HeaderButtonHoverColor),
		),
	}
}

func (s *SettingsScene) Update(context *common.SceneContext) error {
	s.backBtn.Update()
	if s.backBtn.Clicked {
		context.SceneManager.SetScene("home")
		return nil
	}
	return nil
}

func (s *SettingsScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
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
	s.title.DrawEmbossedAutoWithShadow(screen, "Settings", common.ScreenWidth/2, 32, color.RGBA{0, 0, 0, 100}, 1, 1)

	s.backBtn.Draw(screen)

	s.body.SetSize(22)
	s.body.SetAlign(etxt.Left)
	s.body.DrawEmbossedAutoWithShadow(screen, fmt.Sprintf("Screen: %dx%d", common.ScreenWidth, common.ScreenHeight), 60, 100, color.RGBA{0, 0, 0, 80}, 1, 1)
}
