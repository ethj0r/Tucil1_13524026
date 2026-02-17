package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"

	"Tucil1_13524026/src/core"
)

func (a *App) buildHomePage() fyne.CanvasObject {
	welcome := makeColoredText("Hi, welcome", 28, ColorBlack, true)
	title := makeTitle(64)
	img := canvas.NewImageFromFile("src/ui/assets/queen.png")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(280, 360))

	playBtn := a.makeNeonButton("Play", func() {
		a.openFileAndSolve()
	})

	page:=container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(welcome),
		container.NewVBox(),
		title,
		layout.NewSpacer(),
		container.NewCenter(img),
		layout.NewSpacer(),
		container.NewCenter(playBtn),
		layout.NewSpacer(),
	)

	bg := canvas.NewRectangle(ColorWhite)
	return container.NewStack(bg, container.NewPadded(page))
}

func (a *App) openFileAndSolve() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err!=nil {
			a.showErrorPage("Failed to open file picker:\n" + err.Error())
			return
		}
		if reader==nil {
			return 
		}
		defer reader.Close()

		filepath := reader.URI().Path()

		board, err :=core.ParseBoard(filepath)
		if err != nil {
			a.showErrorPage("Error parsing board:\n" + err.Error())
			return
		}

		if err:=board.Validate(); err != nil {
			a.showErrorPage("Invalid board configuration:\n" + err.Error())
			return
		}

		solver := core.NewSolver(board)
		result := solver.Solve()

		if !result.Found {
			a.showErrorPage("No solution found!\nThe puzzle has no valid queen placement.")
			return
		}

		playPage := a.buildPlaySuccessPage(board, &result)
		a.showPlayPage(playPage)

	}, a.window)
}
