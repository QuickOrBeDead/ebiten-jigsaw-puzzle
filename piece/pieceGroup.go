package piece

import (
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/edge"

	"github.com/hajimehoshi/ebiten/v2"
)

type GroupId int

type PieceGroup struct {
	Id                  GroupId
	Pieces              []*Piece
	ConnectedSlotsCount int
	prevX, prevY        int
}

var inverseEdge = [4]edge.Edge{edge.Bottom, edge.Left, edge.Top, edge.Right}

const MergeSnapThreshold = 5 // Max pixel distance to trigger piece merge

func NewPieceGroup(id GroupId, piece *Piece) *PieceGroup {
	return &PieceGroup{
		Id:                  id,
		Pieces:              []*Piece{piece},
		ConnectedSlotsCount: 0,
	}
}

func (g *PieceGroup) Contains(mx, my int) bool {
	for _, p := range g.Pieces {
		if p.Contains(mx, my) {
			g.prevX = mx
			g.prevY = my
			return true
		}
	}

	return false
}

func (g *PieceGroup) ChangePosition(mx, my int) {
	dx := mx - g.prevX
	dy := my - g.prevY

	for _, p := range g.Pieces {
		p.Move(dx, dy)
	}

	g.prevX = mx
	g.prevY = my
}

func (g *PieceGroup) Draw(screen *ebiten.Image) {
	for _, p := range g.Pieces {
		p.Draw(screen)
	}
}

func (g *PieceGroup) CheckPieceMerge(pieceMap map[PieceId]*Piece) (GroupId, int, int) {
	for _, p := range g.Pieces {
		if p.hasNoNeighbors {
			continue
		}
		px, py := p.GetCenterBoxCoordinates()

		for ni, n := range p.Neighbors {
			if n <= 0 {
				continue
			}

			np := pieceMap[n]
			ne := edge.Edge(ni)

			if g.Id == np.GroupId {
				p.SetNeighbor(ne, 0)

				npe := inverseEdge[ne]
				np.SetNeighbor(npe, 0)
				continue
			}

			nx, ny := np.GetCenterBoxCoordinates()

			switch ne {
			case edge.Top:
				dx := px - nx
				dy := py - (ny + int(np.Geo.Size.H))

				if common.Abs(dy) <= MergeSnapThreshold &&
					common.Abs(dx) <= MergeSnapThreshold {
					return np.GroupId, dx, dy
				}
			case edge.Right:
				dx := (px + int(p.Geo.Size.W)) - nx
				dy := py - ny

				if common.Abs(dx) <= MergeSnapThreshold &&
					common.Abs(dy) <= MergeSnapThreshold {
					return np.GroupId, dx, dy
				}
			case edge.Bottom:
				dx := px - nx
				dy := (py + int(p.Geo.Size.H)) - ny

				if common.Abs(dy) <= MergeSnapThreshold &&
					common.Abs(dx) <= MergeSnapThreshold {
					return np.GroupId, dx, dy
				}
			case edge.Left:
				dx := px - (nx + int(np.Geo.Size.W))
				dy := py - ny

				if common.Abs(dx) <= MergeSnapThreshold &&
					common.Abs(dy) <= MergeSnapThreshold {
					return np.GroupId, dx, dy
				}
			}
		}
	}
	return 0, 0, 0
}

func (g *PieceGroup) ResetNeighbors(pieceMap map[PieceId]*Piece, pieces []*Piece) {
	for _, p := range pieces {
		if p.hasNoNeighbors {
			continue
		}

		for ni, n := range p.Neighbors {
			if n <= 0 {
				continue
			}

			np := pieceMap[n]
			if np.GroupId == g.Id {
				ne := edge.Edge(ni)
				p.SetNeighbor(ne, 0)
				np.SetNeighbor(inverseEdge[ne], 0)
			}
		}
	}
}

func (g *PieceGroup) UpdateConnectedSlotsCount() {
	connectedSlotsCount := 0
	for _, p := range g.Pieces {
		for _, n := range p.Neighbors {
			if n == 0 {
				connectedSlotsCount++
			}
		}
	}

	g.ConnectedSlotsCount = connectedSlotsCount
}
