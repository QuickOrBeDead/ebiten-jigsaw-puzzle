package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/scenes/gameScene"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/scenes/homeScene"
)

type Game struct {
	SceneManager *common.SceneManager
	initialized  bool
}

func (g *Game) init() {
	gameImage := common.NewGameImage()
	g.SceneManager.AddScene("Home", func() common.Scene { return homeScene.NewHomeScene(gameImage) })
	g.SceneManager.AddScene("Game", func() common.Scene { return gameScene.NewGameScene(gameImage) })

	// Set the initial scene to Home
	g.SceneManager.SetScene("Home")

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
