package home

import (
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/hajimehoshi/ebiten/v2"
)

type testScene struct {
	btn *common.Button
}

func NewTestScene() *testScene {
	btnW := float32(260)
	btnH := float32(60)

	btnOpts := []common.ButtonOptFunc{
		common.ButtonOption.WithFontSize(24),
		common.ButtonOption.WithFontColor(common.BodyTextColor),
		common.ButtonOption.WithColor(common.SurfaceColor),
		common.ButtonOption.WithHoverColor(common.SurfaceHoverColor),
		common.ButtonOption.WithShadowColor(common.ShadowColor),
	}

	return &testScene{
		btn: common.NewButton(
			0, 0,
			btnW, btnH,
			"Play",
			btnOpts...,
		),
	}
}

func (h *testScene) Update(context *common.SceneContext) error {
	return nil
}

func (h *testScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	h.btn.DrawTest(screen, false, false, false, 5, 5)
	h.btn.DrawTest(screen, false, true, false, 5, 75)
	h.btn.DrawTest(screen, false, true, true, 5, 145)
}
