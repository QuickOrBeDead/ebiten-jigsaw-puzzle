package common

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameImage struct {
	image      *ebiten.Image
	name       string
	pieceCount int
}

func NewGameImage() *GameImage {
	return &GameImage{}
}

func (g *GameImage) SetImage(name string, image *ebiten.Image) {
	g.image = image
	g.name = name
}

func (g *GameImage) GetImage() *ebiten.Image {
	if g.image == nil {
		panic("game image is not set")
	}

	return g.image
}

func (g *GameImage) GetName() string {
	return g.name
}

func (g *GameImage) SetPieceCount(count int) {
	g.pieceCount = count
}

func (g *GameImage) GetPieceCount() int {
	return g.pieceCount
}
