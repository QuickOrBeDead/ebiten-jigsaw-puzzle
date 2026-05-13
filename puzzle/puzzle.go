package puzzle

import (
	"math"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/piece"

	"github.com/hajimehoshi/ebiten/v2"
)

type Puzzle struct {
	PieceCount         int
	pieceMap           map[piece.PieceId]*piece.Piece
	groupsMap          map[piece.GroupId]*piece.PieceGroup
	groupSlice         []*piece.PieceGroup
	activeGroup        piece.GroupId
	cachedImage        *ebiten.Image
	cacheValid         bool
	totalSlotsCount    int
	completePercentage int
}

func NewPuzzle(pieces []*piece.Piece, groups []*piece.PieceGroup, rows, cols int) *Puzzle {
	groupsMap := make(map[piece.GroupId]*piece.PieceGroup, len(groups))
	for _, g := range groups {
		groupsMap[g.Id] = g
	}

	pieceMap := make(map[piece.PieceId]*piece.Piece, len(pieces))
	for _, p := range pieces {
		pieceMap[p.Id] = p
	}

	return &Puzzle{
		PieceCount:         len(pieces),
		pieceMap:           pieceMap,
		groupsMap:          groupsMap,
		groupSlice:         groups,
		cacheValid:         false,
		completePercentage: 0,
		totalSlotsCount:    len(pieces)*4 - ((rows + cols) * 2),
	}
}

func (puz *Puzzle) SetPieceBeingDragged(mx, my int) {
	for i := len(puz.groupSlice) - 1; i >= 0; i-- {
		g := puz.groupSlice[i]
		if g.Contains(mx, my) {
			puz.activeGroup = g.Id

			copy(puz.groupSlice[i:], puz.groupSlice[i+1:])

			puz.groupSlice[len(puz.groupSlice)-1] = g
			puz.cacheValid = false
			return
		}
	}
}

func (puz *Puzzle) MoveDraggedPiece(mx, my int) {
	if puz.activeGroup > 0 {
		ag := puz.groupsMap[puz.activeGroup]
		ag.ChangePosition(mx, my)
	}
}

func (puz *Puzzle) HandleDraggedPieceSnapping(mx, my int) {
	if puz.activeGroup > 0 {
		ag := puz.groupsMap[puz.activeGroup]

		connectGroupId, dx, dy := ag.CheckPieceMerge(puz.pieceMap)
		if connectGroupId > 0 {
			puz.mergeGroups(connectGroupId, ag.Id, dx, dy)
		}
	}
}

func (puz *Puzzle) DropPuzzlePieces() bool {
	if puz.activeGroup > 0 {
		puz.activeGroup = 0
		puz.cacheValid = false
		return true
	}

	return false
}

func (puz *Puzzle) GroupCount() int {
	return len(puz.groupsMap)
}

func (puz *Puzzle) IsSolved() bool {
	return puz.GroupCount() == 1
}

func (p *Puzzle) GetCompletePercentage() int {
	return p.completePercentage
}

func (puz *Puzzle) Draw(screen *ebiten.Image) {
	if puz.cachedImage == nil {
		puz.cachedImage = ebiten.NewImage(common.ScreenWidth, common.ScreenHeight)
	}

	if !puz.cacheValid {
		puz.cachedImage.Clear()
		for _, g := range puz.groupSlice {
			if puz.activeGroup == 0 || g.Id != puz.activeGroup {
				g.Draw(puz.cachedImage)
			}
		}
		puz.cacheValid = true
	}

	if puz.activeGroup > 0 {
		screen.DrawImage(puz.cachedImage, nil)

		ag := puz.groupsMap[puz.activeGroup]
		ag.DrawShadow(screen)
		ag.Draw(screen)
		return
	}

	screen.DrawImage(puz.cachedImage, nil)
}

func (puz *Puzzle) mergeGroups(sourceId, targetId piece.GroupId, translateX, translateY int) {
	s := puz.groupsMap[sourceId]
	t := puz.groupsMap[targetId]

	for _, p := range s.Pieces {
		p.GroupId = targetId
		p.Move(translateX, translateY)

		t.MergePath(p)
	}

	t.ResetNeighbors(puz.pieceMap, s.Pieces)
	t.Pieces = append(t.Pieces, s.Pieces...)
	t.UpdateConnectedSlotsCount()
	s.Pieces = nil

	puz.deleteGroup(sourceId)
	puz.updatePercentage()

	puz.cacheValid = false
}

func (puz *Puzzle) deleteGroup(sourceId piece.GroupId) {
	delete(puz.groupsMap, sourceId)
	w := 0
	for _, g := range puz.groupSlice {
		if g.Id != sourceId {
			puz.groupSlice[w] = g
			w++
		}
	}
	puz.groupSlice = puz.groupSlice[:w]
}

func (puz *Puzzle) updatePercentage() {
	connectedSlotsCount := 0
	for _, g := range puz.groupSlice {
		connectedSlotsCount += g.ConnectedSlotsCount
	}

	puz.completePercentage = int(math.Floor(float64(connectedSlotsCount) / float64(puz.totalSlotsCount) * 100))
}
