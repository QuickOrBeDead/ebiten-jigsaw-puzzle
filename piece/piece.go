package piece

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/edge"
)

type PieceId int

type Piece struct {
	Id             PieceId
	GroupId        GroupId
	Neighbors      [4]PieceId
	hasNoNeighbors bool
	Pos            image.Point
	Geo            *PieceGeometry
	image          *ebiten.Image
	hitMask        *image.Alpha
	drawOpts       *ebiten.DrawImageOptions
}

func NewPiece(id PieceId, groupId GroupId, pos image.Point, geo *PieceGeometry, neighbors [4]PieceId, picture *ebiten.Image) *Piece {
	p := &Piece{
		Id:             id,
		GroupId:        groupId,
		Pos:            pos,
		Geo:            geo,
		Neighbors:      neighbors,
		hasNoNeighbors: false,
	}

	p.drawOpts = &ebiten.DrawImageOptions{}
	p.drawOpts.GeoM.Translate(float64(p.Pos.X), float64(p.Pos.Y))

	boxSize := geo.BoxSize
	mask := common.CreateMask(int(boxSize.W), int(boxSize.H), geo.Path, false)

	p.hitMask = common.CreateImageAlpha(mask)
	p.image = createPieceImage(boxSize, mask, picture)

	mask.Deallocate()

	return p
}

func (p *Piece) SetEdgeType(e edge.Edge, t edge.EdgeType) {
	p.Geo.Edges[e] = t
}

func (p *Piece) GetEdgeType(e edge.Edge) edge.EdgeType {
	return p.Geo.Edges[e]
}

func (p *Piece) SetNeighbor(e edge.Edge, pieceId PieceId) {
	p.Neighbors[e] = pieceId
	p.hasNoNeighbors = p.Neighbors[edge.Top] == 0 && p.Neighbors[edge.Right] == 0 && p.Neighbors[edge.Bottom] == 0 && p.Neighbors[edge.Left] == 0
}

func (p *Piece) GetNeighbor(e edge.Edge) PieceId {
	return p.Neighbors[e]
}

func (p *Piece) Contains(x, y int) bool {
	if x < p.Pos.X || x >= p.Pos.X+p.hitMask.Rect.Dx() ||
		y < p.Pos.Y || y >= p.Pos.Y+p.hitMask.Rect.Dy() {
		return false
	}

	return p.hitMask.Pix[(y-p.Pos.Y)*p.hitMask.Stride+x-p.Pos.X] > 0
}

func (p *Piece) Move(dx, dy int) {
	p.Pos.X += dx
	p.Pos.Y += dy
	p.drawOpts.GeoM.Reset()
	p.drawOpts.GeoM.Translate(float64(p.Pos.X), float64(p.Pos.Y))
}

func (p *Piece) GetCenterBoxCoordinates() (int, int) {
	x, y := p.Pos.X, p.Pos.Y
	if p.Geo.Edges[edge.Left] > 0 {
		x += int(p.Geo.EdgeGeos[edge.Left].TabHeight)
	}

	if p.Geo.Edges[edge.Top] > 0 {
		y += int(p.Geo.EdgeGeos[edge.Top].TabHeight)
	}

	return x, y
}

func (p *Piece) Draw(screen *ebiten.Image) {
	screen.DrawImage(p.image, p.drawOpts)
}

func createPieceImage(boxSize common.Size, mask *ebiten.Image, picture *ebiten.Image) *ebiten.Image {
	intermediate := ebiten.NewImage(int(boxSize.W), int(boxSize.H))

	// Draw mask to intermediate
	intermediate.DrawImage(mask, nil)

	// Apply composite operation
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceIn
	intermediate.DrawImage(picture, op)
	return intermediate
}
