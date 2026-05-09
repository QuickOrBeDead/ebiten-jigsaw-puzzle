package piece

import (
	"math"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/edge"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GroupId int

type PieceGroup struct {
	Id                    GroupId
	Pieces                []*Piece
	ConnectedSlotsCount   int
	prevX, prevY          int
	path                  *vector.Path
	x, y                  int
	shadowImage           *ebiten.Image
	shadowImageCacheValid bool
	boundsX, boundsY      int
}

var inverseEdge = [4]edge.Edge{edge.Bottom, edge.Left, edge.Top, edge.Right}

const MergeSnapThreshold = 5 // Max pixel distance to trigger piece merge

func NewPieceGroup(id GroupId, p *Piece) *PieceGroup {
	pg := &PieceGroup{
		Id:                  id,
		Pieces:              []*Piece{p},
		ConnectedSlotsCount: 0,
		path:                &vector.Path{},
		x:                   p.Pos.X,
		y:                   p.Pos.Y,
	}
	pg.path.AddPath(p.Geo.Path, &vector.AddPathOptions{})
	pg.cacheShadow()
	return pg
}

func (pg *PieceGroup) MergePath(p *Piece) {
	op := &vector.AddPathOptions{}
	op.GeoM.Translate(float64(p.Pos.X-pg.x), float64(p.Pos.Y-pg.y))
	pg.path.AddPath(p.Geo.Path, op)
	pg.shadowImageCacheValid = false
}

func (g *PieceGroup) Bounds() (x, y, w, h int) {
	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32
	for _, p := range g.Pieces {
		rx := p.Pos.X - g.x
		ry := p.Pos.Y - g.y
		rw := p.Geo.BoxSize.W
		rh := p.Geo.BoxSize.H
		if rx < minX {
			minX = rx
		}
		if ry < minY {
			minY = ry
		}
		if rx+rw > maxX {
			maxX = rx + rw
		}
		if ry+rh > maxY {
			maxY = ry + rh
		}
	}
	return minX, minY, maxX - minX, maxY - minY
}

func (g *PieceGroup) cacheShadow() {
	shadowPadding := 10
	bx, by, bw, bh := g.Bounds()
	g.boundsX, g.boundsY = bx, by
	g.shadowImage = ebiten.NewImage(bw+shadowPadding*2, bh+shadowPadding*2)
	common.DrawShadowForPath(g.shadowImage, float64(-bx+shadowPadding), float64(-by+shadowPadding), g.path)
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

	g.x += dx
	g.y += dy
}

func (g *PieceGroup) Draw(screen *ebiten.Image) {
	for _, p := range g.Pieces {
		p.Draw(screen)
	}
}

func (g *PieceGroup) DrawShadow(screen *ebiten.Image) {
	if !g.shadowImageCacheValid {
		g.cacheShadow()
		g.shadowImageCacheValid = true
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(g.x+g.boundsX), float64(g.y+g.boundsY))
	screen.DrawImage(g.shadowImage, opt)
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
