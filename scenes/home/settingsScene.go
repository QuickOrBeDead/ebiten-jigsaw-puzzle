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
	backBtn        *common.Button
	title          *common.TextRenderer
	body           *common.TextRenderer
	sfxSlider      *common.Slider
	musicSlider    *common.Slider
	prevSFXValue   int
	prevMusicValue int
}

func NewSettingsScene(context *common.SceneContext) *SettingsScene {
	settings := common.GetSettings()
	sliderX := float32(common.ScreenWidth-300) / 2
	sfxSlider := common.NewSlider(sliderX, 140, 300, 40, 0, 100, 1)
	sfxSlider.Value = int(settings.SFXVolume * 100)
	musicSlider := common.NewSlider(sliderX, 220, 300, 40, 0, 100, 1)
	musicSlider.Value = int(settings.MusicVolume * 100)

	return &SettingsScene{
		title: common.NewTextRenderer(common.RobotoBoldFontName, common.TitleColor, 32, etxt.Center),
		body:  common.NewTextRenderer(common.RobotoRegularFontName, common.BodyTextColor, 22, etxt.Left),
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
		sfxSlider:      sfxSlider,
		musicSlider:    musicSlider,
		prevSFXValue:   sfxSlider.Value,
		prevMusicValue: musicSlider.Value,
	}
}

func (s *SettingsScene) Update(context *common.SceneContext) error {
	s.backBtn.Update(context)
	s.sfxSlider.Update()
	if s.sfxSlider.Value != s.prevSFXValue {
		s.prevSFXValue = s.sfxSlider.Value
		vol := float64(s.sfxSlider.Value) / 100.0
		if common.AudioManager != nil {
			common.AudioManager.SetSFXVolume(vol)
		}
		common.GetSettings().SFXVolume = vol
		common.SaveSettings()
	}

	s.musicSlider.Update()
	if s.musicSlider.Value != s.prevMusicValue {
		s.prevMusicValue = s.musicSlider.Value
		vol := float64(s.musicSlider.Value) / 100.0
		if common.AudioManager != nil {
			common.AudioManager.SetMusicVolume(vol)
		}
		common.GetSettings().MusicVolume = vol
		common.SaveSettings()
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

	sectionY := 90
	sliderX := (common.ScreenWidth - 300) / 2

	s.body.SetColor(common.TitleColor)
	s.body.SetSize(20)
	s.body.SetAlign(etxt.Center)
	s.body.DrawEmbossedAutoWithShadow(screen, "Audio", common.ScreenWidth/2, sectionY, color.RGBA{0, 0, 0, 80}, 1, 1)

	s.body.SetAlign(etxt.Left)
	s.body.SetColor(common.BodyTextColor)
	s.body.SetSize(18)
	s.body.DrawEmbossedAutoWithShadow(screen, "SFX Volume", sliderX, sectionY+35, color.RGBA{0, 0, 0, 80}, 1, 1)
	s.sfxSlider.Draw(screen)
	s.body.DrawEmbossedAutoWithShadow(screen, fmt.Sprintf("%d%%", s.sfxSlider.Value), sliderX+320, sectionY+35, color.RGBA{0, 0, 0, 80}, 1, 1)

	s.body.DrawEmbossedAutoWithShadow(screen, "Music Volume", sliderX, sectionY+115, color.RGBA{0, 0, 0, 80}, 1, 1)
	s.musicSlider.Draw(screen)
	s.body.DrawEmbossedAutoWithShadow(screen, fmt.Sprintf("%d%%", s.musicSlider.Value), sliderX+320, sectionY+115, color.RGBA{0, 0, 0, 80}, 1, 1)
}
