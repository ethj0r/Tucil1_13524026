package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"

	"Tucil1_13524026/src/core"
)

func (a *App) buildPlaySuccessPage(board *core.Board, result *core.SolverResult) fyne.CanvasObject {
	title := makeTitle(52)
	boardWidget := NewBoardWidget(board, result.Solution)
	boardContainer := container.NewCenter(boardWidget)

	playAgainBtn := a.makeNeonButton("Play Again", func() {
		a.openFileAndSolve()
	})
	saveBtn := a.makeNeonButton("Save", func() {
		a.saveSolution(board, result.Solution)
	})
	btnRow := container.NewHBox(layout.NewSpacer(), playAgainBtn, saveBtn, layout.NewSpacer())

	timeMs := result.ExecutionTime.Milliseconds()
	timeText := fmt.Sprintf("Waktu pencarian: %d ms", timeMs)
	iterText := fmt.Sprintf("Banyak kasus yang ditinjau: %s kasus", formatNumber(result.Iterations))

	timeStat := makeColoredText(timeText, 15, ColorBlack, false)
	iterStat := makeColoredText(iterText, 15, ColorBlack, false)

	statsSection := container.NewVBox(
		container.NewCenter(timeStat),
		container.NewCenter(iterStat),
	)

	page := container.NewVBox(
		layout.NewSpacer(),
		title,
		layout.NewSpacer(),
		boardContainer,
		layout.NewSpacer(),
		btnRow,
		container.NewPadded(statsSection),
		layout.NewSpacer(),
	)

	bg := canvas.NewRectangle(ColorWhite)
	return container.NewStack(bg, container.NewPadded(page))
}

func (a *App) showErrorPage(errMsg string) {
	page := a.buildPlayErrorPage(errMsg)
	a.showPlayPage(page)
}

func (a *App) buildPlayErrorPage(errMsg string) fyne.CanvasObject {
	title := makeTitle(52)
	img := canvas.NewImageFromFile("src/ui/assets/angry.png")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(300, 300))

	playAgainBtn := a.makeNeonButton("Play Again", func() {
		a.openFileAndSolve()
	})
	btnRow := container.NewCenter(playAgainBtn)

	errTitle := makeColoredText("Error message", 18, ColorErrorRed, true)
	errLines := strings.Split(errMsg, "\n")
	errContainer := container.NewVBox(container.NewCenter(errTitle))
	for _, line:=range errLines {
		lineText := makeColoredText(line, 14, ColorErrorRed, false)
		errContainer.Add(container.NewCenter(lineText))
	}

	page := container.NewVBox(
		layout.NewSpacer(),
		title,
		layout.NewSpacer(),
		container.NewCenter(img),
		layout.NewSpacer(),
		btnRow,
		container.NewPadded(errContainer),
		layout.NewSpacer(),
	)

	bg := canvas.NewRectangle(ColorWhite)
	return container.NewStack(bg, container.NewPadded(page))
}

func (a *App) saveSolution(board *core.Board, queens []core.Cell) {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err!=nil {
			dialog.ShowError(err, a.window)
			return
		}
		if writer==nil {
			return
		}
		defer writer.Close()

		queenMap := make(map[core.Cell]bool)
		for _, q :=range queens {
			queenMap[q] = true
		}

		var sb strings.Builder
		for row:=0; row<board.Size; row++ {
			for col:=0; col<board.Size; col++ {
				cell := core.Cell{Row: row, Col: col}
				if queenMap[cell] {
					sb.WriteString("#")
				} else {
					sb.WriteRune(board.Grid[row][col])
				}
			}
			sb.WriteString("\n")
		}

		_, writeErr := writer.Write([]byte(sb.String()))
		if writeErr != nil {
			dialog.ShowError(writeErr, a.window)
			return
		}

		dialog.ShowInformation("Saved", "Solution saved successfully!", a.window)
	}, a.window)
}

//utils
func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s)<=3 {
		return s
	}
	var result strings.Builder
	for i, digit:=range s {
		if i>0 && (len(s)-i)%3==0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}
	return result.String()
}
