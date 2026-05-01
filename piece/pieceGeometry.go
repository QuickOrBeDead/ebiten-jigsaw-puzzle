package piece

import (
	"math"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/edge"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PieceGeometry struct {
	Size                     common.Size
	Top, Bottom, Right, Left float32
	Edges                    [4]edge.EdgeType
	EdgeGeos                 [4]edge.EdgeGeometry
	BoxSize                  common.Size
	Path                     *vector.Path
}

func NewPieceGeometry(size common.Size, top, right, bottom, left float32) *PieceGeometry {
	p := &PieceGeometry{
		Size:   size,
		Top:    top,
		Bottom: bottom,
		Right:  right,
		Left:   left,
	}

	edges := createEdges(top, bottom, right, left)
	edgeGeos := createEdgeGeos(size, top, bottom, right, left)
	boxSize := calculateBoxSize(size, edges, edgeGeos)
	path := generatePath(size, edges, edgeGeos)

	p.BoxSize = boxSize
	p.Path = path
	p.Edges = edges
	p.EdgeGeos = edgeGeos

	return p
}

func createEdges(top float32, bottom float32, right float32, left float32) [4]edge.EdgeType {
	edges := [4]edge.EdgeType{}
	edges[edge.Top] = edge.GetEdgeType(top)
	edges[edge.Bottom] = edge.GetEdgeType(bottom)
	edges[edge.Right] = edge.GetEdgeType(right)
	edges[edge.Left] = edge.GetEdgeType(left)
	return edges
}

func createEdgeGeos(size common.Size, top float32, bottom float32, right float32, left float32) [4]edge.EdgeGeometry {
	edgeGeos := [4]edge.EdgeGeometry{}
	edgeGeos[edge.Top] = edge.EdgeGeometry{
		TabWidth:  int(.2 * float32(size.W)),
		TabHeight: int(.2 * float32(size.W)),
		NeckWidth: int(.1 * float32(size.W)),
		Shoulder:  int(top * float32(size.W)),
	}

	edgeGeos[edge.Bottom] = edge.EdgeGeometry{
		TabWidth:  int(.2 * float32(size.W)),
		TabHeight: int(.2 * float32(size.W)),
		NeckWidth: int(.1 * float32(size.W)),
		Shoulder:  int(bottom * float32(size.W)),
	}

	edgeGeos[edge.Right] = edge.EdgeGeometry{
		TabWidth:  int(.2 * float32(size.H)),
		TabHeight: int(.2 * float32(size.H)),
		NeckWidth: int(.1 * float32(size.H)),
		Shoulder:  int(right * float32(size.H)),
	}

	edgeGeos[edge.Left] = edge.EdgeGeometry{
		TabWidth:  int(.2 * float32(size.H)),
		TabHeight: int(.2 * float32(size.H)),
		NeckWidth: int(.1 * float32(size.H)),
		Shoulder:  int(left * float32(size.H)),
	}
	return edgeGeos
}

func generatePath(size common.Size, edges [4]edge.EdgeType, edgeGeos [4]edge.EdgeGeometry) *vector.Path {
	path := vector.Path{}

	w, h := float32(size.W), float32(size.H)
	topGeo, bottomGeo, leftGeo, rightGeo := edgeGeos[edge.Top], edgeGeos[edge.Bottom], edgeGeos[edge.Left], edgeGeos[edge.Right]

	var x, y float32
	if edges[edge.Top] <= edge.EdgeFlat {
		y = 0
	} else {
		y = float32(topGeo.TabHeight)
	}

	if edges[edge.Left] <= edge.EdgeFlat {
		x = 0
	} else {
		x = float32(leftGeo.TabHeight)
	}

	// Top
	path.MoveTo(x, y)

	if edges[edge.Top] != edge.EdgeFlat {
		shoulderW := float32(math.Abs(float64(topGeo.Shoulder)))

		path.LineTo(x+shoulderW-float32(topGeo.NeckWidth), y)

		path.CubicTo(
			x+shoulderW-float32(topGeo.NeckWidth), y-float32(topGeo.TabHeight)*0.2*common.Sign(float32(topGeo.Shoulder)),
			x+shoulderW-float32(topGeo.TabWidth), y-float32(topGeo.TabHeight)*common.Sign(float32(topGeo.Shoulder)),
			x+shoulderW, y-float32(topGeo.TabHeight)*common.Sign(float32(topGeo.Shoulder)),
		)
		path.CubicTo(
			x+shoulderW+float32(topGeo.TabWidth), y-float32(topGeo.TabHeight)*common.Sign(float32(topGeo.Shoulder)),
			x+shoulderW+float32(topGeo.NeckWidth), y-float32(topGeo.TabHeight)*0.2*common.Sign(float32(topGeo.Shoulder)),
			x+shoulderW+float32(topGeo.NeckWidth), y,
		)
	}

	// Right
	path.LineTo(x+w, y)

	if edges[edge.Right] != edge.EdgeFlat {
		shoulderW := float32(math.Abs(float64(rightGeo.Shoulder)))

		path.LineTo(x+w, y+shoulderW-float32(rightGeo.NeckWidth))

		path.CubicTo(
			x+w+float32(rightGeo.TabHeight)*0.2*common.Sign(float32(rightGeo.Shoulder)), y+shoulderW-float32(rightGeo.NeckWidth),
			x+w+float32(rightGeo.TabHeight)*common.Sign(float32(rightGeo.Shoulder)), y+shoulderW-float32(rightGeo.TabWidth),
			x+w+float32(rightGeo.TabHeight)*common.Sign(float32(rightGeo.Shoulder)), y+shoulderW,
		)
		path.CubicTo(
			x+w+float32(rightGeo.TabHeight)*common.Sign(float32(rightGeo.Shoulder)), y+shoulderW+float32(rightGeo.TabWidth),
			x+w+float32(rightGeo.TabHeight)*0.2*common.Sign(float32(rightGeo.Shoulder)), y+shoulderW+float32(rightGeo.NeckWidth),
			x+w, y+shoulderW+float32(rightGeo.NeckWidth),
		)
	}

	// Bottom
	path.LineTo(x+w, y+h)

	if edges[edge.Bottom] != edge.EdgeFlat {
		shoulderW := float32(math.Abs(float64(bottomGeo.Shoulder)))

		path.LineTo(x+shoulderW-float32(bottomGeo.NeckWidth), y+h)

		path.CubicTo(
			x+shoulderW-float32(bottomGeo.NeckWidth), y+h+float32(bottomGeo.TabHeight)*0.2*common.Sign(float32(bottomGeo.Shoulder)),
			x+shoulderW-float32(bottomGeo.TabWidth), y+h+float32(bottomGeo.TabHeight)*common.Sign(float32(bottomGeo.Shoulder)),
			x+shoulderW, y+h+float32(bottomGeo.TabHeight)*common.Sign(float32(bottomGeo.Shoulder)),
		)
		path.CubicTo(
			x+shoulderW+float32(bottomGeo.TabWidth), y+h+float32(bottomGeo.TabHeight)*common.Sign(float32(bottomGeo.Shoulder)),
			x+shoulderW+float32(bottomGeo.NeckWidth), y+h+float32(bottomGeo.TabHeight)*0.2*common.Sign(float32(bottomGeo.Shoulder)),
			x+shoulderW+float32(bottomGeo.NeckWidth), y+h,
		)
	}

	// Left
	path.LineTo(x, y+h)

	if edges[edge.Left] != edge.EdgeFlat {
		shoulderW := float32(math.Abs(float64(leftGeo.Shoulder)))

		path.LineTo(x, y+shoulderW-float32(leftGeo.NeckWidth))

		path.CubicTo(
			x-float32(leftGeo.TabHeight)*0.2*common.Sign(float32(leftGeo.Shoulder)), y+shoulderW-float32(leftGeo.NeckWidth),
			x-float32(leftGeo.TabHeight)*common.Sign(float32(leftGeo.Shoulder)), y+shoulderW-float32(leftGeo.TabWidth),
			x-float32(leftGeo.TabHeight)*common.Sign(float32(leftGeo.Shoulder)), y+shoulderW,
		)
		path.CubicTo(
			x-float32(leftGeo.TabHeight)*common.Sign(float32(leftGeo.Shoulder)), y+shoulderW+float32(leftGeo.TabWidth),
			x-float32(leftGeo.TabHeight)*0.2*common.Sign(float32(leftGeo.Shoulder)), y+shoulderW+float32(leftGeo.NeckWidth),
			x, y+shoulderW+float32(leftGeo.NeckWidth),
		)
	}

	path.Close()

	return &path
}

func calculateBoxSize(baseSize common.Size, edges [4]edge.EdgeType, edgeGeos [4]edge.EdgeGeometry) common.Size {
	tabWidth := 0
	if edges[edge.Left] != edge.EdgeFlat {
		tabWidth += edgeGeos[edge.Left].TabWidth
	}

	if edges[edge.Right] != edge.EdgeFlat {
		tabWidth += edgeGeos[edge.Right].TabWidth
	}

	tabHeight := 0
	if edges[edge.Top] != edge.EdgeFlat {
		tabHeight += edgeGeos[edge.Top].TabHeight
	}

	if edges[edge.Bottom] != edge.EdgeFlat {
		tabHeight += edgeGeos[edge.Bottom].TabHeight
	}

	return common.Size{W: baseSize.W + tabWidth, H: baseSize.H + tabHeight}
}
