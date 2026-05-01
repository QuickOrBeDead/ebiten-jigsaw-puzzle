package edge

type Edge int

const (
	Top Edge = iota
	Right
	Bottom
	Left
)

type EdgeType int

const (
	EdgeFlat  EdgeType = 0
	EdgeTab   EdgeType = 1
	EdgeBlank EdgeType = -1
)

type EdgeGeometry struct {
	TabWidth  int
	TabHeight int
	NeckWidth int
	Shoulder  int
}

func GetEdgeType(v float32) EdgeType {
	if v == 0 {
		return EdgeFlat
	}

	if v > 0 {
		return EdgeTab
	}

	return EdgeBlank
}
