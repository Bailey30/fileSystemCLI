package cli

import (
	"os"

	"github.com/gdamore/tcell"
)

type Direction int

const (
	Up Direction = iota
	Down
	Forwards
	Backwards
)

func controls(ui *Ui, dir *Dir, filteredDir *Dir) {

	var activeDir = dir
	if len(ui.searchInput) != 0 {
		activeDir = filteredDir
	}

	// Poll event
	ev := ui.screen.PollEvent()

	// Process event
	switch ev := ev.(type) {
	case *tcell.EventResize:
		ui.screen.Sync()
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyUp:
			activeDir.UpdateDirectoryPosition(Up, ui)
		case tcell.KeyDown:
			activeDir.UpdateDirectoryPosition(Down, ui)
		case tcell.KeyEscape:
			ui.screen.Fini()
			os.Exit(0)
		case tcell.KeyLeft:
			activeDir.UpdateDirectoryContents(Backwards, ui, dir)
		case tcell.KeyRight:
			activeDir.UpdateDirectoryContents(Forwards, ui, dir)
		case tcell.KeyRune:
			// If a printable character is pressed, update the search bar content
			ui.searchInput += string(ev.Rune())
			filteredDir.Filter(dir, ui.searchInput)

		case tcell.KeyBackspace, tcell.KeyBackspace2:
			// If the Backspace key is pressed, remove the last character from the search bar content
			if len(ui.searchInput) > 0 {
				ui.searchInput = ui.searchInput[:len(ui.searchInput)-1]
			}
			filteredDir.Filter(dir, ui.searchInput)

		}
	}

}
