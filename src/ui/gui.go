package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	ColorBlue          = color.NRGBA{R: 0x03, G: 0x67, B: 0xFD, A: 0xFF} // #0367FD
	ColorBlack         = color.NRGBA{R: 0x16, G: 0x16, B: 0x16, A: 0xFF} // #161616
	ColorWhite         = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // #FFFFFF
	ColorNeonGreen     = color.NRGBA{R: 0xD1, G: 0xF8, B: 0x01, A: 0xFF} // #D1F801
	ColorSidebarActive = color.NRGBA{R: 0x3E, G: 0x55, B: 0x78, A: 0xFF} // #3E5578
	ColorErrorRed      = color.NRGBA{R: 0xE8, G: 0x3E, B: 0x3E, A: 0xFF}

	ColorTitleLeft   = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	ColorTitleMiddle = color.NRGBA{R: 0x2F, G: 0x61, B: 0x92, A: 0xFF}
	ColorTitleRight  = color.NRGBA{R: 0xDF, G: 0xDF, B: 0xDF, A: 0xFF}

	RegionColors = map[rune]color.NRGBA{
		'A': {R: 0x0A, G: 0x84, B: 0xFF, A: 0xFF}, //Blue
		'B': {R: 0xFF, G: 0x45, B: 0x3A, A: 0xFF}, //Red
		'C': {R: 0x30, G: 0xD1, B: 0x58, A: 0xFF}, //Green
		'D': {R: 0xFF, G: 0x9F, B: 0x0A, A: 0xFF}, //Orange
		'E': {R: 0xBF, G: 0x5A, B: 0xF2, A: 0xFF}, //Purple
		'F': {R: 0x5A, G: 0xC8, B: 0xFA, A: 0xFF}, //Cyan
		'G': {R: 0xFF, G: 0xD6, B: 0x0A, A: 0xFF}, //Yellow
		'H': {R: 0xFF, G: 0x2D, B: 0x55, A: 0xFF}, //Pink
		'I': {R: 0x5E, G: 0x5C, B: 0xE6, A: 0xFF}, //Indigo
		'J': {R: 0x63, G: 0xE6, B: 0xE2, A: 0xFF}, //Mint
		'K': {R: 0xAC, G: 0x8E, B: 0x68, A: 0xFF}, //Brown
		'L': {R: 0x63, G: 0x63, B: 0x66, A: 0xFF}, //Gray
		'M': {R: 0xFF, G: 0x6B, B: 0x00, A: 0xFF}, //Deep Orange (dark mode)
		'N': {R: 0x00, G: 0xC7, B: 0xBE, A: 0xFF}, //Teal
		'O': {R: 0x34, G: 0xAA, B: 0xDC, A: 0xFF}, //Light Blue
		'P': {R: 0xFF, G: 0xCC, B: 0x00, A: 0xFF}, //Gold
		'Q': {R: 0x8E, G: 0x8E, B: 0x93, A: 0xFF}, //Mid Gray
		'R': {R: 0xAE, G: 0xEA, B: 0xD6, A: 0xFF}, //Seafoam (light)
		'S': {R: 0xFF, G: 0x3B, B: 0x30, A: 0xFF}, //System Red
		'T': {R: 0x00, G: 0x7A, B: 0xFF, A: 0xFF}, //System Blue
		'U': {R: 0x34, G: 0xC7, B: 0x59, A: 0xFF}, //System Green
		'V': {R: 0xAF, G: 0x52, B: 0xDE, A: 0xFF}, //System Purple
		'W': {R: 0xFF, G: 0xC8, B: 0x00, A: 0xFF}, //System Yellow
		'X': {R: 0xFF, G: 0x3A, B: 0x30, A: 0xFF}, //Coral / Hot Red
		'Y': {R: 0x1C, G: 0x1C, B: 0x1E, A: 0xFF}, //Label (near black)
		'Z': {R: 0xE5, G: 0xE5, B: 0xEA, A: 0xFF}, //Fill Gray (light)
	}
)

func GetRegionColor(letter rune) color.NRGBA {
	if c, ok := RegionColors[letter]; ok {
		return c
	}
	return ColorBlue
}

type App struct {
	fyneApp fyne.App
	window  fyne.Window
	content *fyne.Container
	//sidebar
	navHomeText *canvas.Text
	navPlayText *canvas.Text
	//pages
	homePage fyne.CanvasObject
	playPage fyne.CanvasObject
}

func LaunchGUI() {
	a := &App{}
	a.fyneApp = app.NewWithID("com.queens.solver")
	a.window = a.fyneApp.NewWindow("Queens")
	a.window.Resize(fyne.NewSize(1024, 800))
	a.window.SetFixedSize(true)
	a.window.CenterOnScreen()

	sidebar := a.buildSidebar()
	a.content = container.NewStack()

	a.homePage = a.buildHomePage()
	a.content.Objects = []fyne.CanvasObject{a.homePage}
	a.setActiveNav("home")

	sidebarBg := canvas.NewRectangle(ColorBlack)
	sidebarBg.SetMinSize(fyne.NewSize(180, 0))
	sidebarWithBg := container.NewStack(sidebarBg, sidebar)

	mainLayout := container.NewBorder(nil, nil, sidebarWithBg, nil, a.content)
	a.window.SetContent(mainLayout)
	a.window.ShowAndRun()
}

func (a *App) buildSidebar() fyne.CanvasObject {
	a.navHomeText = canvas.NewText("Home", ColorWhite)
	a.navHomeText.TextSize = 16
	a.navHomeText.TextStyle = fyne.TextStyle{Bold: true}

	a.navPlayText = canvas.NewText("Play", ColorWhite)
	a.navPlayText.TextSize = 16
	a.navPlayText.TextStyle = fyne.TextStyle{Bold: true}

	homeNav := newTappableStack(func() { a.navigateTo("home") }, a.navHomeText)
	playNav := newTappableStack(func() { a.navigateTo("play") }, a.navPlayText)
	navSection := container.NewVBox(homeNav, playNav)

	devLabel := canvas.NewText("developed by", ColorWhite)
	devLabel.TextSize = 11
	devName := canvas.NewText("Jordhy Branenda", ColorWhite)
	devName.TextSize = 12
	devName.TextStyle = fyne.TextStyle{Bold: true}
	devSection := container.NewVBox(devLabel, devName)

	return container.NewBorder(
		nil,
		container.New(layout.NewCustomPaddedLayout(0, 20, 16, 0), devSection),
		nil, nil,
		container.New(layout.NewCustomPaddedLayout(100, 0, 16, 0), navSection),
	)
}

func (a *App) navigateTo(page string) {
	switch page {
	case "home":
		if a.homePage==nil {
			a.homePage = a.buildHomePage()
		}
		a.content.Objects = []fyne.CanvasObject{a.homePage}
		a.setActiveNav("home")
	case "play":
		a.openFileAndSolve()
		return
	}
	a.content.Refresh()
}

func (a *App) showPlayPage(page fyne.CanvasObject) {
	a.playPage = page
	a.content.Objects = []fyne.CanvasObject{a.playPage}
	a.setActiveNav("play")
	a.content.Refresh()
	a.window.SetFixedSize(false)
}

func (a *App) setActiveNav(page string) {
	switch page {
	case "home":
		a.navHomeText.Color =ColorSidebarActive
		a.navPlayText.Color =ColorWhite
	case "play":
		a.navHomeText.Color =ColorWhite
		a.navPlayText.Color =ColorSidebarActive
	}
	a.navHomeText.Refresh()
	a.navPlayText.Refresh()
}

func makeColoredText(text string, size float32, c color.Color, bold bool) *canvas.Text {
	t := canvas.NewText(text, c)
	t.TextSize = size
	t.TextStyle = fyne.TextStyle{Bold: bold}
	t.Alignment = fyne.TextAlignCenter
	return t
}

func makeTitle(size float32) fyne.CanvasObject {
	title := canvas.NewText("Queens", ColorBlack)
	title.TextSize = size
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	return container.NewCenter(title)
}

func (a *App) makeNeonButton(label string, onTap func()) *tappableStack {
	bg := canvas.NewRectangle(ColorNeonGreen)
	bg.SetMinSize(fyne.NewSize(140, 40))
	bg.CornerRadius = 20
	text := makeColoredText(label, 14, ColorBlack, true)
	return newTappableStack(onTap, bg, container.NewCenter(text))
}

type tappableStack struct {
	widget.BaseWidget
	onTap   func()
	content *fyne.Container
}

func newTappableStack(onTap func(), objects ...fyne.CanvasObject) *tappableStack {
	t:=&tappableStack{
		onTap:   onTap,
		content: container.NewStack(objects...),
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *tappableStack) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.content)
}

func (t *tappableStack) Tapped(_ *fyne.PointEvent) {
	if t.onTap!=nil {
		t.onTap()
	}
}

func (t *tappableStack) TappedSecondary(_ *fyne.PointEvent) {}