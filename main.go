package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/scenes/game"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/scenes/home"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/scenes/startGame"
)

type Game struct {
	SceneManager *common.SceneManager
	initialized  bool
}

func (g *Game) init() {
	gameImage := common.NewGameImage()
	g.SceneManager.AddScene("home", func() common.Scene { return home.NewHomeScene() })
	g.SceneManager.AddScene("startGame", func() common.Scene { return startGame.NewStartGameScene(gameImage) })
	g.SceneManager.AddScene("game", func() common.Scene { return game.NewGameScene(gameImage) })
	g.SceneManager.AddScene("howToPlay", func() common.Scene { return home.NewHowToPlayScene() })
	g.SceneManager.AddScene("settings", func() common.Scene { return home.NewSettingsScene() })
	g.SceneManager.AddScene("credits", func() common.Scene { return home.NewCreditsScene() })

	g.SceneManager.SetScene("home")
}

func (g *Game) Update() error {
	if !g.initialized {
		g.init()
		g.initialized = true
	}

	return g.SceneManager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.SceneManager.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) { return common.ScreenWidth, common.ScreenHeight }

func main() {
	ebiten.SetWindowSize(common.ScreenWidth, common.ScreenHeight)
	ebiten.SetWindowTitle("Jigsaw Puzzle")

	game := &Game{
		SceneManager: common.NewSceneManager(),
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
