package screen

import (
	"log"
	"os"

	"github.com/gdamore/tcell"
)

type Error struct {
	message string
	isError bool
}

func NewError(err string) *Error {
	if err == "" {
		return &Error{
			message: "",
			isError: false,
		}
	} else {
		return &Error{
			message: err,
			isError: true,
		}
	}
}

type Ui struct {
	screen     tcell.Screen
	xmax, ymax int
	defStyle   tcell.Style
	boxStyle   tcell.Style
	textStyle  tcell.Style
	highlight  tcell.Style
	dirHeight  int
}

func NewUi() *Ui {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorSlateGray).Foreground(tcell.ColorIndianRed)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorRed)
	highlight := tcell.StyleDefault.Background(tcell.ColorWhite)
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	s.SetStyle(defStyle)

	s.Clear()

	xmax, ymax := s.Size()

	return &Ui{
		screen:    s,
		xmax:      xmax,
		ymax:      ymax,
		defStyle:  defStyle,
		boxStyle:  boxStyle,
		textStyle: textStyle,
		highlight: highlight,
		dirHeight: ymax - 1,
	}
}

type SearchBarContent string

func GetUi(dir *Dir, err2 *Error) error {
	ui := NewUi()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		ui.screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	normalLogFile, err := os.Create("logs.log")
	if err != nil {
		panic(err)
	}
	os.Stdout = normalLogFile

	// Event loop
	for {
		ui.screen.Clear()
		drawDir(ui, dir, err2)
		ui.screen.Show()
		controls(ui, dir) // needs to be after Show()
		// drawSearch(ui, sbc)
		// sbc := SearchEvent(ui, searchBarContent)
		// Update screen
	}
}

func drawText(ui *Ui, dir *Dir, err2 *Error) {

	for row := 0; row < len(dir.files); row++ {

		style := ui.textStyle
		if row == dir.highlightIndex {
			style = ui.highlight
		}

		index := row + dir.drawBeginning // enables the scroll

		if index < len(dir.files) {
			drawLine(ui, row, dir.files[index].path, style)
		}

	}

}

func drawDir(ui *Ui, dir *Dir, err2 *Error) {
	// xmax, ymax := ui.xmax, ui.ymax

	// XXX: manual clean without flush to avoid flicker on Windows
	wtot, htot := ui.screen.Size()

	// Fill background
	for i := 0; i < wtot; i++ {
		for j := 0; j < htot; j++ {

			ui.screen.SetContent(i, j, ' ', nil, ui.boxStyle)
		}
	}

	// Draw borders
	// for col := 0; col <= xmax-3; col++ {
	// 	ui.screen.SetContent(col+1, 0, tcell.RuneHLine, nil, ui.boxStyle)
	// 	ui.screen.SetContent(col+1, ymax-1, tcell.RuneHLine, nil, ui.boxStyle)
	// }
	// for row := 0 + 1; row < ymax; row++ {
	// 	ui.screen.SetContent(0, row, tcell.RuneVLine, nil, ui.boxStyle)
	// 	ui.screen.SetContent(xmax-1, row, tcell.RuneVLine, nil, ui.boxStyle)
	// }

	drawText(ui, dir, err2)
}

func drawLine(ui *Ui, row int, text string, style tcell.Style) {
	for col, rune := range []rune(text) {
		ui.screen.SetContent(col, row, rune, nil, style)
	}

}

func drawSearch(ui *Ui, searchBarContent string) {
	wtot, htot := ui.screen.Size()
	// Fill background
	for i := 0; i < wtot; i++ {
		for j := htot - 1; j < htot; j++ {

			ui.screen.SetContent(i, j, ' ', nil, ui.boxStyle)
		}
	}

	// Draw borders
	for col := 0; col <= ui.xmax-3; col++ {
		// ui.screen.SetContent(col+1, 0, tcell.RuneHLine, nil, ui.boxStyle)
		ui.screen.SetContent(col, ui.ymax-2, tcell.RuneHLine, nil, ui.textStyle)
	}

	for i, r := range searchBarContent {

		ui.screen.SetContent(i, htot-1, r, nil, ui.textStyle)
	}
	// drawLine(ui, htot-1, searchBarContent, ui.textStyle)

}
