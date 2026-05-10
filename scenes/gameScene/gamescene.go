package gameScene

import (
	"fmt"
	"image/color"
	"time"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/puzzle"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
)

type GameScene struct {
	puzzle            *puzzle.Puzzle
	headerHeight      float32
	footerHeight      float32
	pictureName       string
	moves             int
	headerButtons     []*common.Button
	footerButtons     []*common.Button
	startTime         int64
	endTime           int64
	isPuzzleCompleted bool
	image             *ebiten.Image
	showGhost         bool
	showImage         bool
	previewImage      *common.PreviewImage

	frameCache      *ebiten.Image
	frameCacheValid bool

	lastInputTime  time.Time
	lastMouseX     int
	lastMouseY     int
	lastTimeUpdate int64
	elapsedTimeStr string

	text *common.TextRenderer
}

func NewGameScene(gameImage *common.GameImage) *GameScene {
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetVsyncEnabled(false)

	text := common.NewTextRenderer(common.RobotoBoldFontName, common.BodyTextColor, 40, etxt.Center)
	buttonOptions := []common.ButtonOptFunc{
		common.ButtonOption.WithColor(common.HeaderButtonColor),
		common.ButtonOption.WithHoverColor(common.HeaderButtonHoverColor),
		common.ButtonOption.WithFontColor(common.BodyTextColor),
		common.ButtonOption.WithFontSize(18),
	}

	margin := 10.
	const footerHeight = 56
	image := gameImage.GetImage()

	previewImage := common.NewPreviewImage(
		image,
		margin,
		0,
		0.5,
		common.PreviewImageOption.WithBGColor(color.Black),
		common.PreviewImageOption.WithBorderColor(common.PrimaryColor))

	previewImage.Y = float64(common.ScreenHeight) - float64(previewImage.ScaledH) - float64(footerHeight) - margin

	s := &GameScene{
		text:              text,
		headerHeight:      64,
		footerHeight:      footerHeight,
		pictureName:       gameImage.GetName(),
		startTime:         time.Now().Unix(),
		endTime:           0,
		lastInputTime:     time.Now(),
		frameCache:        ebiten.NewImage(common.ScreenWidth, common.ScreenHeight),
		moves:             0,
		showGhost:         false,
		showImage:         false,
		isPuzzleCompleted: false,
		image:             image,
		previewImage:      previewImage,
		headerButtons: []*common.Button{
			common.NewButton(
				1025, 12,
				80, 40,
				"Restart",
				buttonOptions...,
			),
			common.NewButton(
				1150, 12,
				80, 40,
				"Home",
				buttonOptions...,
			),
		},
		footerButtons: []*common.Button{
			common.NewButton(
				20, 0,
				80, 36,
				"Image",
				append(buttonOptions,
					common.ButtonOption.WithToggle(true),
					common.ButtonOption.WithActiveColor(common.PrimaryColor),
					common.ButtonOption.WithColor(common.SurfaceColor),
					common.ButtonOption.WithFontSize(16),
				)...,
			),
			common.NewButton(
				140, 0,
				80, 36,
				"Ghost",
				append(buttonOptions,
					common.ButtonOption.WithToggle(true),
					common.ButtonOption.WithActiveColor(common.PrimaryColor),
					common.ButtonOption.WithColor(common.SurfaceColor),
					common.ButtonOption.WithFontSize(16),
				)...,
			),
		},
	}

	pp := puzzle.NewPuzzlePicture(s.image)
	s.puzzle = pp.CreatePuzzle(gameImage.GetPieceCount())

	return s
}

func (g *GameScene) Update(context *common.SceneContext) error {
	now := time.Now()
	hasInput := false

	mx, my := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		hasInput = true
		g.puzzle.SetPieceBeingDragged(mx, my)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		hasInput = true
		g.puzzle.HandleDraggedPieceSnapping(mx, my)
		if g.puzzle.DropPuzzlePieces() {
			g.incrementMoves()
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		hasInput = true
		g.puzzle.MoveDraggedPiece(mx, my)
	}

	if mx != g.lastMouseX || my != g.lastMouseY {
		hasInput = true
		g.lastMouseX, g.lastMouseY = mx, my
	}

	if hasInput {
		g.lastInputTime = now
		g.frameCacheValid = false
		ebiten.SetTPS(60)
	} else if now.Sub(g.lastInputTime) > 2*time.Second {
		ebiten.SetTPS(1)
	}

	for _, button := range g.headerButtons {
		button.Update()
		if button.Clicked {
			g.frameCacheValid = false
			switch button.Label {
			case "Home":
				context.SceneManager.SetScene("Home")
			case "Restart":
				context.SceneManager.SetScene("Game")
			}
		}
	}

	for _, button := range g.footerButtons {
		button.Update()

		if button.Clicked {
			g.frameCacheValid = false
			switch button.Label {
			case "Image":
				g.showImage = !g.showImage
			case "Ghost":
				g.showGhost = !g.showGhost
			}
		}
	}

	if !g.isPuzzleCompleted {
		if g.puzzle.IsSolved() {
			g.isPuzzleCompleted = true
			g.endTime = now.Unix()
			g.frameCacheValid = false
		} else if now.Unix() != g.lastTimeUpdate {
			g.lastTimeUpdate = now.Unix()
			g.elapsedTimeStr = formatDuration(now.Unix() - g.startTime)
			g.frameCacheValid = false
		}
	}

	return nil
}

func (g *GameScene) incrementMoves() {
	if !g.isPuzzleCompleted {
		g.moves++
	}
}

func (g *GameScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	if g.frameCacheValid {
		return
	}
	g.renderFrame(g.frameCache)
	screen.DrawImage(g.frameCache, nil)
	g.frameCacheValid = true
}

func (g *GameScene) renderFrame(dst *ebiten.Image) {
	dst.Fill(common.BackgroundColor)

	if g.showGhost {
		opt := &ebiten.DrawImageOptions{}

		opt.ColorScale.Scale(0.5, 0.5, 0.5, 1)

		opt.GeoM.Translate(
			(float64(common.ScreenWidth)-float64(g.image.Bounds().Dx()))/2,
			(float64(common.ScreenHeight)-float64(g.image.Bounds().Dy()))/2,
		)
		dst.DrawImage(g.image, opt)
	}

	g.puzzle.Draw(dst)

	g.drawHeader(dst)
	g.drawFooter(dst)
}

func (g *GameScene) drawHeader(screen *ebiten.Image) {
	h := g.headerHeight

	common.DrawPanel(screen, 0, 0, float32(common.ScreenWidth), h, common.HeaderColor, false, nil)

	vector.FillRect(screen, 0, 0, float32(common.ScreenWidth), 1, color.RGBA{255, 255, 255, 14}, false)

	common.DrawPanel(screen, 0, h-3, float32(common.ScreenWidth), 3, common.PrimaryColor, false, nil)

	for i := 0; i < 4; i++ {
		alpha := uint8(35 - 7*i)
		vector.FillRect(screen, 0, h+float32(i), float32(common.ScreenWidth), 1, color.RGBA{0, 0, 0, alpha}, false)
	}

	g.text.SetColor(common.TitleColor)
	g.text.SetSize(32)
	g.text.SetAlign(etxt.Left)
	g.text.DrawEmbossedAutoWithShadow(screen, fmt.Sprintf("Jigsaw Puzzle - %s", g.pictureName), 32, 32, color.RGBA{0, 0, 0, 100}, 1, 1)

	for _, button := range g.headerButtons {
		button.Draw(screen)
	}
}

func (g *GameScene) drawFooter(screen *ebiten.Image) {
	footerY := float32(common.ScreenHeight) - g.footerHeight
	h := g.footerHeight

	common.DrawPanel(screen, 0, footerY, float32(common.ScreenWidth), h, common.FooterColor, false, nil)

	vector.FillRect(screen, 0, footerY, float32(common.ScreenWidth), 3, common.PrimaryColor, false)

	for i := 0; i < 4; i++ {
		alpha := uint8(35 - 7*i)
		vector.FillRect(screen, 0, footerY-float32(i+1), float32(common.ScreenWidth), 1, color.RGBA{0, 0, 0, alpha}, false)
	}

	// Progress bar in the center
	progress := float32(g.puzzle.GetCompletePercentage()) / 100.0
	barWidth := float32(280)
	barHeight := float32(30)
	barX := float32(common.ScreenWidth)/2 - barWidth/2
	barY := footerY + (g.footerHeight-barHeight)/2

	common.DrawProgressBar(screen, barX, barY, barWidth, barHeight,
		progress, common.SurfaceColor, common.SuccessColor, common.PrimaryColor)

	g.text.SetColor(color.RGBA{255, 255, 255, 255})
	g.text.SetSize(14)
	g.text.SetAlign(etxt.Center)
	g.text.DrawEmbossedAutoWithShadow(
		screen,
		fmt.Sprintf("%d%%", g.puzzle.GetCompletePercentage()),
		int(barX+barWidth/2),
		int(barY+barHeight/2),
		color.RGBA{0, 0, 0, 80}, 1, 1,
	)

	// Draw footer buttons - centered vertically
	for _, button := range g.footerButtons {
		button.Y = footerY + (g.footerHeight-button.Height)/2
		button.Draw(screen)
	}

	// Draw stats on the right - vertically centered
	g.text.SetColor(common.BodyTextColor)
	g.text.SetSize(18)
	g.text.SetAlign(etxt.Right)
	g.text.DrawEmbossedAutoWithShadow(
		screen,
		fmt.Sprintf("Pieces: %d  |  Moves: %d  |  Time: %s", g.puzzle.PieceCount, g.moves, g.getElapsedTime()),
		int(float32(common.ScreenWidth)-20),
		int(footerY+g.footerHeight/2),
		color.RGBA{0, 0, 0, 80}, 1, 1,
	)

	if g.isPuzzleCompleted {
		g.drawCompletionBanner(screen)
	}

	if g.showImage {
		g.previewImage.Draw(screen)
	}
}

func (g *GameScene) drawCompletionBanner(screen *ebiten.Image) {
	bannerY := float32(common.ScreenHeight)/2 - 60
	bannerH := float32(120)

	common.DrawPanel(screen, float32(common.ScreenWidth)/2-200, bannerY, 400, bannerH,
		color.RGBA{20, 20, 30, 220}, true, common.SuccessColor)

	g.text.SetColor(common.SuccessColor)
	g.text.SetSize(36)
	g.text.SetAlign(etxt.Center)
	g.text.DrawEmbossedHozCenterWithShadow(screen, "Puzzle Completed!", int(bannerY)+40, color.RGBA{0, 0, 0, 100}, 2, 2)

	elapsed := g.getElapsedTime()
	g.text.SetColor(common.BodyTextColor)
	g.text.SetSize(20)
	g.text.DrawEmbossedHozCenterWithShadow(screen, fmt.Sprintf("Time: %s  |  Moves: %d", elapsed, g.moves), int(bannerY)+75, color.RGBA{0, 0, 0, 80}, 1, 1)
}

func formatDuration(elapsed int64) string {
	days := elapsed / (24 * 3600)
	hours := (elapsed % (24 * 3600)) / 3600
	minutes := (elapsed % 3600) / 60
	seconds := elapsed % 60

	if days > 0 {
		return fmt.Sprintf("%d.%02d:%02d:%02d", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (g *GameScene) getElapsedTime() string {
	if g.isPuzzleCompleted {
		return formatDuration(g.endTime - g.startTime)
	}
	return g.elapsedTimeStr
}
