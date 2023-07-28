package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell"
)

type Direction int

const (
	Up Direction = iota
	Down
	Forwards
	Backwards
	Jump
	Refresh
)

func IsRunningCommand(ui *Ui) bool {
	splitInput := strings.Split(ui.searchInput, " ")
	isCommand := false

	if ui.searchInput[0] == '/' && len(splitInput) > 0 || splitInput[0] == "/n" || splitInput[0] == "/d" || !ui.confirm {
		isCommand = true
	}

	return isCommand
}

func controls(ui *Ui, dir *Dir, filteredDir *Dir) {
	// splitInput := strings.Split(ui.searchInput, " ")
	var activeDir = dir
	if len(ui.searchInput) > 0 && !IsRunningCommand((ui)) {
		activeDir = filteredDir
	}

	// Poll event
	ev := ui.screen.PollEvent()

	switch ui.confirm {
	case true:

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			ui.screen.Sync()
			w, h := ui.screen.Size()
			ui.xmax = w
			ui.ymax = h
			ui.dirHeight = h - 1
		case *tcell.EventKey:

			switch ev.Key() {
			case tcell.KeyUp:
				activeDir.UpDown(Up, ui)
			case tcell.KeyDown:
				activeDir.UpDown(Down, ui)
			case tcell.KeyEscape:
				ui.screen.Fini()
				os.Exit(0)
			case tcell.KeyLeft:
				activeDir.Traverse(Backwards, ui, dir)
			case tcell.KeyRight:
				activeDir.Traverse(Forwards, ui, dir)
			case tcell.KeyRune:

				fmt.Println(ev.Rune())

				// If a printable character is pressed, update the search bar content
				ui.searchInput += string(ev.Rune())

				filteredDir.Filter(dir, ui.searchInput)

			case tcell.KeyBackspace, tcell.KeyBackspace2:

				// If the Backspace key is pressed, remove the last character from the search bar content
				if len(ui.searchInput) > 0 {
					ui.searchInput = ui.searchInput[:len(ui.searchInput)-1]
				}

				filteredDir.Filter(dir, ui.searchInput)

			case tcell.KeyEnter:
				split := strings.Split(ui.searchInput, " ")
				fmt.Println(split)

				if len(ui.searchInput) > 0 {
					if split[0] == "/" {
						JumpToDirectory(ui, dir)
					} else if split[0] == "/n" {
						dir.create(split[1], ui)
					} else if split[0] == "/d" {
						dir.Delete(ui)
					}

				} else {
					activeDir.Open("new")
				}
			case tcell.KeyHome:
				{
					activeDir.Open("same")

				}
			}
		}

	case false:
		switch ev := ev.(type) {
		case *tcell.EventResize:
			ui.screen.Sync()
			// w, h := ui.screen.Size()
			// ui.xmax = w
			// ui.ymax = h
			// ui.dirHeight = h - 1
		case *tcell.EventKey:

			switch ev.Rune() {
			case 'n':
				fmt.Println("You pressed 'n'")
				dir.CancelDelete(ui)
			case 'y':
				fmt.Println("You pressed 'y'")
				dir.ConfirmDelete(ui)
			}

			switch ev.Key() {

			case tcell.KeyEscape:
				ui.screen.Fini()
				os.Exit(0)

			}
		}
	}

}
