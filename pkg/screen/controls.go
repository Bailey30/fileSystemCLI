package screen

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

func controls(ui *Ui, dir *Dir) {
	// var searchBarContent = ""

	// Poll event
	ev := ui.screen.PollEvent()
	// fmt.Println("ev", &ev)
	// Process event
	switch ev := ev.(type) {
	case *tcell.EventResize:
		ui.screen.Sync()
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyUp:
			dir.UpdateDirectoryPosition(Up, ui)
		case tcell.KeyDown:
			dir.UpdateDirectoryPosition(Down, ui)
		case tcell.KeyEscape:
			ui.screen.Fini()
			os.Exit(0)
		case tcell.KeyLeft:
			dir.UpdateDirectoryContents(Backwards, ui)
		case tcell.KeyRight:
			dir.UpdateDirectoryContents(Forwards, ui)
			// case tcell.KeyRune:
			// 	// If a printable character is pressed, update the search bar content
			// 	searchBarContent += string(ev.Rune())
			// 	return searchBarContent
			// case tcell.KeyBackspace, tcell.KeyBackspace2:
			// 	// If the Backspace key is pressed, remove the last character from the search bar content
			// 	if len(searchBarContent) > 0 {
			// 		searchBarContent = searchBarContent[:len(searchBarContent)-1]
			// 	}
			// 	return searchBarContent

		}
	}

}

func SearchEvent(ui *Ui, searchBarContent string) string {
	for {
		ev := ui.screen.PollEvent()

		switch ev.(type) {
		case *tcell.EventKey:
			// Handle keyboard events here
			keyEvent := ev.(*tcell.EventKey)
			if keyEvent.Key() == tcell.KeyRune {
				// If a printable character is pressed, update the search bar content
				searchBarContent += string(keyEvent.Rune())

				// fmt.Println("ev", keyEvent.Rune())
			} else if keyEvent.Key() == tcell.KeyBackspace || keyEvent.Key() == tcell.KeyBackspace2 {
				// If the Backspace key is pressed, remove the last character from the search bar content
				if len(searchBarContent) > 0 {
					searchBarContent = searchBarContent[:len(searchBarContent)-1]
				}
			}

		case *tcell.EventResize:
			ui.screen.Sync()
		}
		return searchBarContent
	}
}
