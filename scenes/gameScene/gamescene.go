package gameScene

import (
	"fmt"
	"image/color"
	"time"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/puzzle"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	text *common.TextRenderer
}

func NewGameScene(gameImage *common.GameImage) *GameScene {
	text := common.NewTextRenderer(common.RobotoBoldFontName, common.BodyTextColor, 40, etxt.Center)
	buttonOptions := []common.ButtonOptFunc{
		common.ButtonOption.WithColor(common.HeaderButtonColor),
		common.ButtonOption.WithHoverColor(common.HeaderButtonHoverColor),
		common.ButtonOption.WithFontColor(common.BodyTextColor),
		common.ButtonOption.WithFontSize(18),
	}

	s := &GameScene{
		text:              text,
		headerHeight:      64,
		footerHeight:      56,
		pictureName:       gameImage.GetName(),
		startTime:         time.Now().Unix(),
		endTime:           0,
		moves:             0,
		showGhost:         false,
		showImage:         false,
		isPuzzleCompleted: false,
		image:             gameImage.GetImage(),
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
			common.NewButton(
				260, 0,
				80, 36,
				"Arrange",
				append(buttonOptions,
					common.ButtonOption.WithColor(common.SurfaceColor),
					common.ButtonOption.WithFontSize(16),
				)...,
			),
		},
	}

	pp := puzzle.NewPuzzlePicture(s.image)
	s.puzzle = pp.CreatePuzzle(24)

	return s
}

func (g *GameScene) Update(context *common.SceneContext) error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		g.puzzle.SetPieceBeingDragged(mx, my)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		g.puzzle.HandleDraggedPieceSnapping(mx, my)
		if g.puzzle.DropPuzzlePieces() {
			g.incrementMoves()
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		g.puzzle.MoveDraggedPiece(mx, my)
	}

	for _, button := range g.headerButtons {
		button.Update()
		if button.Clicked {
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
			switch button.Label {
			case "Image":
				g.showImage = !g.showImage
			case "Ghost":
				g.showGhost = !g.showGhost
			}
		}
	}

	if !g.isPuzzleCompleted {
		g.isPuzzleCompleted = g.puzzle.IsSolved()
		g.endTime = time.Now().Unix()
	}

	return nil
}

func (g *GameScene) incrementMoves() {
	if !g.isPuzzleCompleted {
		g.moves++
	}
}

func (g *GameScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	screen.Fill(common.BackgroundColor)

	if g.showGhost {
		opt := &ebiten.DrawImageOptions{}

		opt.ColorScale.Scale(0.5, 0.5, 0.5, 1) // make ghost image semi-transparent

		// translate the ghost image to center
		opt.GeoM.Translate(
			(float64(common.ScreenWidth)-float64(g.image.Bounds().Dx()))/2,
			(float64(common.ScreenHeight)-float64(g.image.Bounds().Dy()))/2,
		)
		screen.DrawImage(g.image, opt)
	}

	g.puzzle.Draw(screen)

	g.drawHeader(screen)
	g.drawFooter(screen)
}

func (g *GameScene) drawHeader(screen *ebiten.Image) {
	common.DrawPanel(screen, 0, 0, float32(common.ScreenWidth), g.headerHeight, common.HeaderColor, false, nil)

	common.DrawPanel(screen, 0, g.headerHeight-2, float32(common.ScreenWidth), 2, common.PrimaryColor, false, nil)

	g.text.SetColor(common.TitleColor)
	g.text.SetSize(32)
	g.text.SetAlign(etxt.Left)
	g.text.DrawWithShadow(screen, fmt.Sprintf("Jigsaw Puzzle - %s", g.pictureName), 32, 32, color.RGBA{0, 0, 0, 100}, 1, 1)

	for _, button := range g.headerButtons {
		button.Draw(screen)
	}
}

func (g *GameScene) drawFooter(screen *ebiten.Image) {
	footerY := float32(common.ScreenHeight) - g.footerHeight

	common.DrawPanel(screen, 0, footerY, float32(common.ScreenWidth), g.footerHeight, common.FooterColor, false, nil)

	common.DrawPanel(screen, 0, footerY, float32(common.ScreenWidth), 2, common.PrimaryColor, false, nil)

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
	g.text.Draw(
		screen,
		fmt.Sprintf("%d%%", g.puzzle.GetCompletePercentage()),
		int(barX+barWidth/2),
		int(barY+barHeight/2),
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
	g.text.Draw(
		screen,
		fmt.Sprintf("Moves: %d  |  Time: %s", g.moves, g.getElapsedTime()),
		int(float32(common.ScreenWidth)-20),
		int(footerY+g.footerHeight/2),
	)

	if g.isPuzzleCompleted {
		g.drawCompletionBanner(screen)
	}

	if g.showImage {
		scale := 0.5
		margin := 10.
		opt := &ebiten.DrawImageOptions{}

		opt.GeoM.Scale(scale, scale)
		opt.GeoM.Translate(
			margin,
			float64(common.ScreenHeight)-float64(g.image.Bounds().Dy())*scale-float64(g.footerHeight)-margin,
		)

		common.DrawPanel(screen,
			float32(margin-2), float32(common.ScreenHeight)-float32(g.image.Bounds().Dy())*float32(scale)-g.footerHeight-2,
			float32(g.image.Bounds().Dx())*float32(scale)+4, float32(g.image.Bounds().Dy())*float32(scale)+4,
			color.RGBA{0, 0, 0, 0}, true, common.PrimaryColor)

		screen.DrawImage(g.image, opt)
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
	g.text.DrawHorizontalCenter(screen, "Puzzle Completed!", int(bannerY)+40)

	elapsed := g.getElapsedTime()
	g.text.SetColor(common.BodyTextColor)
	g.text.SetSize(20)
	g.text.DrawHorizontalCenter(screen, fmt.Sprintf("Time: %s  |  Moves: %d", elapsed, g.moves), int(bannerY)+75)
}

func (g *GameScene) getElapsedTime() string {
	var elapsed int64
	if g.isPuzzleCompleted {
		elapsed = g.endTime - g.startTime
	} else {
		elapsed = time.Now().Unix() - g.startTime
	}

	days := elapsed / (24 * 3600)
	hours := (elapsed % (24 * 3600)) / 3600
	minutes := (elapsed % 3600) / 60
	seconds := elapsed % 60

	var timeStr string
	if days > 0 {
		timeStr = fmt.Sprintf("%d.%02d:%02d:%02d", days, hours, minutes, seconds)
	} else if hours > 0 {
		timeStr = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	} else {
		timeStr = fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
	return timeStr
}
