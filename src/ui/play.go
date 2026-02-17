package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"
	"time"

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
	saveTxtBtn := a.makeNeonButton("Save TXT", func() {
		a.saveSolution(board, result.Solution)
	})
	saveImgBtn := a.makeNeonButton("Save Image", func() {
		a.saveSolutionAsImage(boardWidget)
	})
	btnRow := container.NewHBox(layout.NewSpacer(), playAgainBtn, saveTxtBtn, saveImgBtn, layout.NewSpacer())

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

func (a *App) showSolvingPage(board *core.Board) {
	title := makeTitle(52)
	boardWidget := NewBoardWidget(board, []core.Cell{})
	boardContainer := container.NewCenter(boardWidget)

	statusText := makeColoredText("Solving...", 16, ColorBlack, true)
	iterText := makeColoredText("Iterations: 0", 14, ColorBlack, false)
	
	progressSection := container.NewVBox(
		container.NewCenter(statusText),
		container.NewCenter(iterText),
	)

	page := container.NewVBox(
		layout.NewSpacer(),
		title,
		layout.NewSpacer(),
		boardContainer,
		layout.NewSpacer(),
		container.NewPadded(progressSection),
		layout.NewSpacer(),
	)

	bg := canvas.NewRectangle(ColorWhite)
	pageObj := container.NewStack(bg, container.NewPadded(page))
	a.showPlayPage(pageObj)

	type UpdateMsg struct {
		queens []core.Cell
		iter   int
		done   bool
		result *core.SolverResult
	}
	updateChan := make(chan UpdateMsg, 100)

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		
		var lastUpdate UpdateMsg
		for {
			select {
			case msg := <-updateChan:
				if msg.done {
					if !msg.result.Found {
						a.showErrorPage("No solution found!\nThe puzzle has no valid queen placement.")
					} else {
						playPage := a.buildPlaySuccessPage(board, msg.result)
						a.showPlayPage(playPage)
					}
					return
				}
				lastUpdate = msg
			case <-ticker.C:
				if lastUpdate.iter > 0 {
					boardWidget.UpdateQueens(lastUpdate.queens)
					iterText.Text = fmt.Sprintf("Iterations: %s", formatNumber(lastUpdate.iter))
					iterText.Refresh()
				}
			}
		}
	}()

	go func() {
		solver := core.NewSolver(board)
		solver.SetUpdateCallback(func(currentQueens []core.Cell, iteration int) {
			select {
			case updateChan <- UpdateMsg{queens: currentQueens, iter: iteration}:
			default:
			}
		})

		result := solver.Solve()
		
		time.Sleep(100 * time.Millisecond)
		updateChan <- UpdateMsg{done: true, result: &result}
	}()
}

func (a *App) saveSolutionAsImage(boardWidget *BoardWidget) {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		img := a.renderBoardToImage(boardWidget)
		if img == nil {
			dialog.ShowError(fmt.Errorf("failed to render board"), a.window)
			return
		}
		//encode ke png
		if err := png.Encode(writer, img); err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		dialog.ShowInformation("Saved", "Image saved successfully!", a.window)
	}, a.window)
}

func (a *App) renderBoardToImage(bw *BoardWidget) image.Image {
	if bw == nil || bw.board == nil {
		return nil
	}

	n := bw.size
	cellSize := 80
	gap := 4
	boardPad := 16
	cornerR := 8
	totalInner := n*cellSize + (n-1)*gap
	imgSize := totalInner + boardPad*2
	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	fillRect(img, 0, 0, imgSize, imgSize, ColorBlack, cornerR)

	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			x := boardPad + col*(cellSize+gap)
			y := boardPad + row*(cellSize+gap)

			letter := bw.board.Grid[row][col]
			cellColor := GetRegionColor(letter)
			fillRect(img, x, y, cellSize, cellSize, cellColor, cornerR)
			cell := core.Cell{Row: row, Col: col}
			if bw.queens[cell] {
				drawQueenSymbol(img, x, y, cellSize)
			}
		}
	}

	return img
}

func fillRect(img *image.RGBA, x, y, w, h int, col color.NRGBA, cornerRadius int) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			dx := 0
			dy := 0
			if px < x+cornerRadius {
				dx = x + cornerRadius - px
			} else if px >= x+w-cornerRadius {
				dx = px - (x + w - cornerRadius - 1)
			}
			if py < y+cornerRadius {
				dy = y + cornerRadius - py
			} else if py >= y+h-cornerRadius {
				dy = py - (y + h - cornerRadius - 1)
			}
			
			if dx*dx+dy*dy > cornerRadius*cornerRadius {
				continue
			}

			if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
				img.Set(px, py, col)
			}
		}
	}
}

func drawQueenSymbol(img *image.RGBA, x, y, size int) {
	centerX := x + size/2
	centerY := y + size/2
	radius := size / 3
	
	queenColor := color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	for py := centerY - radius; py <= centerY+radius; py++ {
		for px := centerX - radius; px <= centerX+radius; px++ {
			dx := px - centerX
			dy := py - centerY
			if dx*dx+dy*dy <= radius*radius {
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, queenColor)
				}
			}
		}
	}
}
