package puzzle

import (
	"fmt"
	"image"
	"math"
	"math/rand/v2"
	"runtime"
	"sync"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/edge"
	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/piece"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
)

type PuzzlePicture struct {
	image *ebiten.Image
}

func NewPuzzlePicture(image *ebiten.Image) *PuzzlePicture {
	w, h := image.Bounds().Dx(), image.Bounds().Dy()
	scale := 1.0
	if w > 600 {
		scale = 600.0 / float64(w)
		w = int(float64(w) * scale)
		h = int(float64(h) * scale)
	}

	puzzleImage := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	puzzleImage.DrawImage(image, op)

	return &PuzzlePicture{image: puzzleImage}
}

func (p *PuzzlePicture) CreatePuzzle(pieceCount int) *Puzzle {
	cols, rows := findClosestDivisors(pieceCount)
	pieceDataArr := calculatePieceData(p, cols, rows, pieceCount)

	// Create pieces in parallel using a worker pool (limits goroutines for scalability)
	pieces := make([]*piece.Piece, pieceCount)
	groups := make([]*piece.PieceGroup, pieceCount)

	numWorkers := runtime.NumCPU()
	if numWorkers > pieceCount {
		numWorkers = pieceCount
	}
	jobs := make(chan int, pieceCount)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start worker goroutines
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				pieceId := piece.PieceId(idx + 1)
				pieceGroupId := piece.GroupId(pieceId)
				d := pieceDataArr[idx]

				neighbors := [4]piece.PieceId{}
				left := piece.PieceId(-1)
				if pieceId%piece.PieceId(cols) != 1 {
					left = pieceId - 1
				}

				right := piece.PieceId(-1)
				if pieceId%piece.PieceId(cols) != 0 {
					right = pieceId + 1
				}

				top := piece.PieceId(-1)
				if pieceId-piece.PieceId(cols) > 0 {
					top = pieceId - piece.PieceId(cols)
				}

				bottom := piece.PieceId(-1)
				if pieceId+piece.PieceId(cols) <= piece.PieceId(pieceCount) {
					bottom = pieceId + piece.PieceId(cols)
				}

				neighbors[edge.Left] = left
				neighbors[edge.Right] = right
				neighbors[edge.Top] = top
				neighbors[edge.Bottom] = bottom

				pic := p.image.SubImage(image.Rect(int(d.picX), int(d.picY), int(d.picX+d.geo.BoxSize.W), int(d.picY+d.geo.BoxSize.H))).(*ebiten.Image)
				p := piece.NewPiece(pieceId, pieceGroupId, image.Pt(int(d.randomX), int(d.randomY)), d.geo, neighbors, pic)

				pieces[idx] = p
				groups[idx] = piece.NewPieceGroup(pieceGroupId, p)
			}
		}()
	}

	// Send jobs to workers
	for i := 0; i < pieceCount; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return NewPuzzle(pieces, groups, rows, cols)
}

type pieceData struct {
	top, right, bottom, left float32
	randomX, randomY         float32
	geo                      *piece.PieceGeometry
	picX, picY               int
}

func calculatePieceData(p *PuzzlePicture, cols int, rows int, pieceCount int) []pieceData {
	pictureW, pictureH := p.image.Bounds().Dx(), p.image.Bounds().Dy()

	// Distribute pixels evenly across columns and rows
	colWidths := make([]int, cols)
	baseW, extraW := pictureW/cols, pictureW%cols
	for j := 0; j < cols; j++ {
		colWidths[j] = baseW
		if j < extraW {
			colWidths[j]++
		}
	}

	rowHeights := make([]int, rows)
	baseH, extraH := pictureH/rows, pictureH%rows
	for i := 0; i < rows; i++ {
		rowHeights[i] = baseH
		if i < extraH {
			rowHeights[i]++
		}
	}

	// Precompute starting positions for each column and row
	colX := make([]int, cols)
	rowY := make([]int, rows)
	for j := 1; j < cols; j++ {
		colX[j] = colX[j-1] + colWidths[j-1]
	}
	for i := 1; i < rows; i++ {
		rowY[i] = rowY[i-1] + rowHeights[i-1]
	}
	maxPieceW, maxPieceH := colWidths[0], rowHeights[0]
	for j := 1; j < cols; j++ {
		if colWidths[j] > maxPieceW {
			maxPieceW = colWidths[j]
		}
	}
	for i := 1; i < rows; i++ {
		if rowHeights[i] > maxPieceH {
			maxPieceH = rowHeights[i]
		}
	}
	randomX1Min, randomX1Max, randomX2Min, randomX2Max, randomYMin, randomYMax := float32(0.), float32(common.SpawnAreaLeftMaxX)-float32(maxPieceW), float32(common.SpawnAreaRightMinX)-float32(maxPieceW), float32(common.ScreenWidth)-float32(maxPieceW), float32(common.SpawnAreaMinY), float32(common.ScreenHeight)-float32(maxPieceH)-common.SpawnAreaBottomOffset

	pieceDataArr := make([]pieceData, pieceCount)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			idx := i*cols + j

			top, right, bottom, left := randomTabPosition(), randomTabPosition(), randomTabPosition(), randomTabPosition()

			if i == 0 {
				top = 0
			}

			if j == 0 {
				left = 0
			}

			if i == rows-1 {
				bottom = 0
			} else {
				if rand.N(2) == 0 {
					bottom = -bottom
				}
			}

			if j == cols-1 {
				right = 0
			} else {
				if rand.N(2) == 0 {
					right = -right
				}
			}

			if j > 0 {
				left = -pieceDataArr[i*cols+j-1].geo.Right
			}

			if i > 0 {
				top = -pieceDataArr[(i-1)*cols+j].geo.Bottom
			}

			// Generate a random position for the piece within the bounds of the puzzle image
			var randomX float32
			if i%2 == 0 {
				randomX = randomX1Min + rand.Float32()*(randomX1Max-randomX1Min)
			} else {
				randomX = randomX2Min + rand.Float32()*(randomX2Max-randomX2Min)
			}

			randomY := randomYMin + rand.Float32()*(randomYMax-randomYMin)

			pieceSize := common.Size{W: colWidths[j], H: rowHeights[i]}

			geo := piece.NewPieceGeometry(pieceSize, top, right, bottom, left)

			picX := colX[j]
			if left > 0 {
				picX = picX - geo.EdgeGeos[edge.Left].TabHeight
			}

			picY := rowY[i]
			if top > 0 {
				picY = picY - geo.EdgeGeos[edge.Top].TabHeight
			}

			pieceDataArr[idx] = pieceData{
				top:     top,
				right:   right,
				bottom:  bottom,
				left:    left,
				randomX: randomX,
				randomY: randomY,
				geo:     geo,
				picX:    picX,
				picY:    picY,
			}
		}
	}
	return pieceDataArr
}

func randomTabPosition() float32 {
	return 0.3 + (0.7-0.3)*rand.Float32()
}

func findClosestDivisors(c int) (int, int) {
	d1, d2, minDiff := 0, 0, math.MaxInt

	for i := int(math.Sqrt(float64(c))); i > 0; i-- {
		if c%i == 0 {
			diff := int(math.Abs(float64(i - c/i)))
			if diff < minDiff {
				minDiff = diff
				d1 = int(math.Max(float64(i), float64(c/i)))
				d2 = int(math.Min(float64(i), float64(c/i)))
			}
		}
	}

	if minDiff == math.MaxInt {
		panic(fmt.Errorf("%d has no closest divisors", c))
	}

	if float64(d1)/float64(d2) > 3 {
		panic(fmt.Errorf("closest divisors' ratio of %d cannot be bigger than 3", c))
	}

	return d1, d2
}
