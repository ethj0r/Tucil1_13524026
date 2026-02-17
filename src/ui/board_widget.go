package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"Tucil1_13524026/src/core"
)

type BoardWidget struct {
	widget.BaseWidget
	board  *core.Board
	queens map[core.Cell]bool
	size   int
}

func NewBoardWidget(board *core.Board, queens []core.Cell) *BoardWidget {
	qmap := make(map[core.Cell]bool)
	for _, q := range queens {
		qmap[q] = true
	}
	bw := &BoardWidget{
		board:  board,
		queens: qmap,
		size:   board.Size,
	}
	bw.ExtendBaseWidget(bw)
	return bw
}

func (bw *BoardWidget) CreateRenderer() fyne.WidgetRenderer {
	return &boardRenderer{bw: bw}
}

type boardRenderer struct {
	bw      *BoardWidget
	objects []fyne.CanvasObject
	built   bool
}

func (r *boardRenderer) buildGrid() {
	r.objects = nil

	n := r.bw.size
	cellSize := r.calcCellSize()
	gap := float32(4)
	cornerR := float32(8)
	boardPad := float32(16)

	totalInner := float32(n)*cellSize+float32(n-1)*gap
	outerSize := totalInner + boardPad*2

	outerRect := canvas.NewRectangle(ColorBlack)
	outerRect.CornerRadius = 20
	outerRect.Move(fyne.NewPos(0, 0))
	outerRect.Resize(fyne.NewSize(outerSize, outerSize))
	r.objects = append(r.objects, outerRect)

	// Cells
	for row:=0; row<n; row++ {
		for col:=0; col<n; col++ {
			x := boardPad+float32(col)*(cellSize+gap)
			y := boardPad+float32(row)*(cellSize+gap)

			letter := r.bw.board.Grid[row][col]
			cellColor := GetRegionColor(letter)

			cellRect := canvas.NewRectangle(cellColor)
			cellRect.CornerRadius = cornerR
			cellRect.Move(fyne.NewPos(x, y))
			cellRect.Resize(fyne.NewSize(cellSize, cellSize))
			r.objects = append(r.objects, cellRect)

			cell := core.Cell{Row: row, Col: col}
			if r.bw.queens[cell] {
				queenText :=canvas.NewText("♛", color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xE0})
				queenText.TextSize = cellSize * 0.55
				queenText.TextStyle = fyne.TextStyle{Bold: true}
				queenText.Alignment = fyne.TextAlignCenter
				queenText.Move(fyne.NewPos(x, y+(cellSize-queenText.TextSize)*0.5))
				queenText.Resize(fyne.NewSize(cellSize, cellSize))
				r.objects = append(r.objects, queenText)
			}
		}
	}
	r.built = true
}

func (r *boardRenderer) calcCellSize() float32 {
	n:=r.bw.size
	switch {
	case n <= 5:
		return 64
	case n <= 8:
		return 52
	case n <= 12:
		return 40
	default:
		return 32
	}
}

func (r *boardRenderer) Layout(_ fyne.Size) {
	if !r.built {
		r.buildGrid()
	}
}

func (r *boardRenderer) MinSize() fyne.Size {
	n := r.bw.size
	cellSize := r.calcCellSize()
	gap := float32(4)
	boardPad := float32(16)
	totalInner := float32(n)*cellSize + float32(n-1)*gap
	outer := totalInner + boardPad*2
	return fyne.NewSize(outer, outer)
}

func (r *boardRenderer) Refresh() {
	r.built = false
	r.buildGrid()
	canvas.Refresh(r.bw)
}

func (r *boardRenderer) Objects() []fyne.CanvasObject {
	if !r.built {
		r.buildGrid()
	}
	return r.objects
}

func (r *boardRenderer) Destroy() {}
