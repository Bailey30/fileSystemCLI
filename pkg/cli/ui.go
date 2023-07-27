package cli

import (
	"log"
	"os"
	"strings"

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
	screen      tcell.Screen
	xmax, ymax  int
	defStyle    tcell.Style
	boxStyle    tcell.Style
	textStyle   tcell.Style
	highlight   tcell.Style
	pathStyle   tcell.Style
	dirHeight   int
	searchInput string
	confirm     bool
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
	highlight := tcell.StyleDefault.Background(tcell.ColorLimeGreen)
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorLimeGreen)
	pathStyle := tcell.StyleDefault.Background(tcell.ColorPink)

	s.SetStyle(defStyle)

	s.Clear()

	xmax, ymax := s.Size()

	return &Ui{
		screen:      s,
		xmax:        xmax,
		ymax:        ymax,
		defStyle:    defStyle,
		boxStyle:    boxStyle,
		textStyle:   textStyle,
		highlight:   highlight,
		pathStyle:   pathStyle,
		dirHeight:   ymax - 3,
		searchInput: "",
		confirm:     true,
	}
}

type SearchBarContent struct {
	Search string
}

func InitUi(dir *Dir, filteredDir *Dir, err2 *Error) error {
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
		drawDir(ui, dir, filteredDir, err2)
		drawSearch(ui)
		ui.screen.Show()
		controls(ui, dir, filteredDir) // needs to be after Show()
	}
}

func drawText(ui *Ui, dir *Dir, err2 *Error) {

	// directory path
	for row := 0; row < 1; row++ {
		drawLine(ui, row, dir.path, ui.pathStyle)
	}

	// list of files and folders
	for row := dir.topOffset; row <= ui.dirHeight-dir.topOffset; row++ {

		style := ui.textStyle
		if row == dir.highlightIndex {
			style = ui.highlight
		}

		index := row + dir.drawBeginning - dir.topOffset // enables the scroll

		if index < len(dir.files) {
			drawLine(ui, row, formatName(dir.files[index].name), style)
		}

	}

}
func drawLine(ui *Ui, row int, text string, style tcell.Style) {
	for col, rune := range []rune(text) {
		ui.screen.SetContent(col, row, rune, nil, style)
	}
}

func formatName(name string) string {
	split := strings.Split(name, ".")

	if len(split) == 1 {
		return name + "/"
	} else {
		return name
	}
}

func drawDir(ui *Ui, dir *Dir, filteredDir *Dir, err2 *Error) {
	// xmax, ymax := ui.xmax, ui.ymax

	// XXX: manual clean without flush to avoid flicker on Windows
	wtot, htot := ui.screen.Size()

	// Fill background
	for row := 0; row < htot; row++ {
		for col := 0; col < wtot; col++ {

			ui.screen.SetContent(col, row, ' ', nil, ui.boxStyle)
		}
	}

	directoryToBeDisplayed := dir

	if len(ui.searchInput) > 0 && !IsRunningCommand((ui)) {
		directoryToBeDisplayed = filteredDir
	}

	drawText(ui, directoryToBeDisplayed, err2)
}

func drawSearch(ui *Ui) {
	wtot, htot := ui.screen.Size()
	// Fill background
	for i := 0; i < htot; i++ {
		for j := wtot - 1; j < htot; j++ {
			ui.screen.SetContent(i, j, ' ', nil, ui.boxStyle)
		}
	}

	// Draw borders
	for col := 0; col <= ui.xmax-2; col++ {
		// ui.screen.SetContent(col+1, 0, tcell.RuneHLine, nil, ui.boxStyle)
		ui.screen.SetContent(col, htot-2, tcell.RuneHLine, nil, ui.textStyle)
	}

	for i, r := range ui.searchInput {
		ui.screen.SetContent(i, htot-1, r, nil, ui.textStyle)
	}
	// drawLine(ui, htot-1, searchBarContent, ui.textStyle)

}
