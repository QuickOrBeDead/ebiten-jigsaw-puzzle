package homeScene

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
)

type startGameDialog struct {
	open         bool
	previewImage *common.PreviewImage
	imageName    string

	panelX, panelY float32
	panelW, panelH float32

	slider       *common.Slider
	startButton  *common.Button
	cancelButton *common.Button

	text *common.TextRenderer
}

func newStartGameDialog() *startGameDialog {
	panelW := float32(480)
	panelH := float32(420)
	panelX := (float32(common.ScreenWidth) - panelW) / 2
	panelY := (float32(common.ScreenHeight) - panelH) / 2

	sliderW := float32(380)
	sliderX := panelX + (panelW-sliderW)/2
	sliderY := panelY + 215

	btnY := panelY + 300
	btnH := float32(44)

	d := &startGameDialog{
		panelX: panelX,
		panelY: panelY,
		panelW: panelW,
		panelH: panelH,
		slider: common.NewSlider(sliderX, sliderY, sliderW, 50, 12, 300, 12),
		startButton: common.NewButton(
			panelX+panelW-240, btnY,
			160, btnH,
			"Start Puzzle",
			common.ButtonOption.WithFontSize(20),
			common.ButtonOption.WithFontColor(common.BodyTextColor),
			common.ButtonOption.WithColor(common.PrimaryColor),
			common.ButtonOption.WithHoverColor(common.PrimaryHoverColor),
			common.ButtonOption.WithShadowColor(common.ShadowColor),
		),
		cancelButton: common.NewButton(
			panelX+50, btnY,
			120, btnH,
			"Cancel",
			common.ButtonOption.WithFontSize(20),
			common.ButtonOption.WithFontColor(common.BodyTextColor),
			common.ButtonOption.WithColor(common.SurfaceColor),
			common.ButtonOption.WithHoverColor(common.SurfaceHoverColor),
			common.ButtonOption.WithShadowColor(common.ShadowColor),
		),
		text: common.NewTextRenderer(common.RobotoBoldFontName, common.BodyTextColor, 24, etxt.Center),
	}

	return d
}

func (d *startGameDialog) Open(img *ebiten.Image, name string) {
	cx := float64(d.panelX) + float64(d.panelW)/2
	cy := float64(d.panelY) + 65

	d.previewImage = common.NewPreviewImageWithMax(
		img, cx, cy, 160, 120,
		common.PreviewImageOption.WithBGColor(color.Black),
		common.PreviewImageOption.WithBorderColor(common.MutedTextColor))
	d.previewImage.X -= float64(d.previewImage.ScaledW / 2)

	d.imageName = name
	d.slider.Value = 60
	if d.slider.Value < d.slider.Min {
		d.slider.Value = d.slider.Min
	}
	if d.slider.Value > d.slider.Max {
		d.slider.Value = d.slider.Max
	}
	d.open = true
}

func (d *startGameDialog) Close() {
	d.open = false
	d.previewImage = nil
}

func (d *startGameDialog) IsOpen() bool {
	return d.open
}

func (d *startGameDialog) PieceCount() int {
	return d.slider.Value
}

func (d *startGameDialog) PreviewImage() *ebiten.Image {
	return d.previewImage.Image
}

func (d *startGameDialog) ImageName() string {
	return d.imageName
}

func (d *startGameDialog) startClicked() bool {
	return d.startButton.Clicked
}

func (d *startGameDialog) cancelClicked() bool {
	return d.cancelButton.Clicked
}

func (d *startGameDialog) Update() {
	if !d.open {
		return
	}

	d.slider.Update()
	d.startButton.Update()
	d.cancelButton.Update()
}

func (d *startGameDialog) Draw(screen *ebiten.Image) {
	if !d.open {
		return
	}

	vector.FillRect(screen, 0, 0, common.ScreenWidth, common.ScreenHeight, color.RGBA{0, 0, 0, 200}, true)

	common.DrawPanel(screen, d.panelX, d.panelY, d.panelW, d.panelH, common.SurfaceColor, true, common.PrimaryColor)

	common.DrawSoftShadow(screen, d.panelX, d.panelY-2, d.panelW, d.panelH+4, color.RGBA{0, 0, 0, 100}, 6, 6)

	d.text.SetColor(common.TitleColor)
	d.text.SetSize(26)
	d.text.SetAlign(etxt.Center)
	d.text.DrawHorizontalCenter(screen, "Choose Piece Count", int(d.panelY)+30)

	d.drawPreview(screen)

	d.text.SetColor(common.BodyTextColor)
	d.text.SetSize(22)
	d.text.SetAlign(etxt.Center)
	d.text.DrawHorizontalCenter(screen, fmt.Sprintf("%d pieces", d.slider.Value), int(d.panelY)+195)

	d.slider.Draw(screen)
	d.startButton.Draw(screen)
	d.cancelButton.Draw(screen)
}

func (d *startGameDialog) drawPreview(screen *ebiten.Image) {
	if d.previewImage == nil {
		return
	}

	d.previewImage.Draw(screen)
}
