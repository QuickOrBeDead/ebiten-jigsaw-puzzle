package home

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
)

type HowToPlayScene struct {
	backBtn *common.Button
	title   *common.TextRenderer
	body    *common.TextRenderer
}

func NewHowToPlayScene(context *common.SceneContext) *HowToPlayScene {
	return &HowToPlayScene{
		title: common.NewTextRenderer(common.RobotoBoldFontName, common.TitleColor, 36, etxt.Center),
		body:  common.NewTextRenderer(common.RobotoRegularFontName, common.BodyTextColor, 24, etxt.Center),
		backBtn: common.NewButton(
			20, 12,
			80, 40,
			"Back",
			common.ButtonOption.WithFontSize(18),
			common.ButtonOption.WithFontColor(common.BodyTextColor),
			common.ButtonOption.WithColor(common.HeaderButtonColor),
			common.ButtonOption.WithHoverColor(common.HeaderButtonHoverColor),
			common.ButtonOption.WithOnClick(func() {
				context.SceneManager.SetScene("home")
			}),
		),
	}
}

func (s *HowToPlayScene) Update(context *common.SceneContext) error {
	s.backBtn.Update(context)
	return nil
}

func (s *HowToPlayScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	screen.Fill(common.BackgroundColor)

	const headerH float32 = 64
	common.DrawPanel(screen, 0, 0, float32(common.ScreenWidth), headerH, common.HeaderColor, false, nil)
	vector.FillRect(screen, 0, 0, float32(common.ScreenWidth), 1, color.RGBA{255, 255, 255, 14}, false)
	common.DrawPanel(screen, 0, headerH-3, float32(common.ScreenWidth), 3, common.PrimaryColor, false, nil)
	for i := 0; i < 4; i++ {
		alpha := uint8(35 - 7*i)
		vector.FillRect(screen, 0, headerH+float32(i), float32(common.ScreenWidth), 1, color.RGBA{0, 0, 0, alpha}, false)
	}

	s.title.DrawEmbossedAutoWithShadow(screen, "How to Play", common.ScreenWidth/2, 32, color.RGBA{0, 0, 0, 100}, 1, 1)

	s.backBtn.Draw(screen)

	s.body.SetSize(24)
	s.body.SetAlign(etxt.Center)

	lines := []string{
		"1. Click \"Play\" on the main menu to browse puzzle images.",
		"2. Select an image from the gallery or upload your own JPG.",
		"3. Choose how many pieces you want (12-300).",
		"4. Drag puzzle pieces from the sides onto the board.",
		"5. Pieces snap together when placed correctly.",
		"",
		"Tips:",
		"- Use the \"Image\" toggle to preview the full picture.",
		"- Use the \"Ghost\" toggle to overlay a faded reference image.",
		"- The progress bar shows completion percentage.",
	}

	y := 100
	cx := common.ScreenWidth / 2
	for _, line := range lines {
		s.body.DrawEmbossedAutoWithShadow(screen, line, cx, y, color.RGBA{0, 0, 0, 80}, 1, 1)
		y += 36
	}
}
