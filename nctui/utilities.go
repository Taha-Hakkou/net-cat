package nctui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var Title = []string{
	"░▒▓███████▓▒░░▒▓████████▓▒░▒▓████████▓▒░ ",
	"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     ",
	"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     ",
	"░▒▓█▓▒░░▒▓█▓▒░▒▓██████▓▒░    ░▒▓█▓▒░     ",
	"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     ",
	"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     ",
	"░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░  ░▒▓█▓▒░     ",
	"                                         ",
	"                                         ",
	" ░▒▓██████▓▒░ ░▒▓██████▓▒░▒▓████████▓▒░  ",
	"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░      ",
	"░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░      ",
	"░▒▓█▓▒░      ░▒▓████████▓▒░ ░▒▓█▓▒░      ",
	"░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░      ",
	"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░      ",
	" ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░      ",
	"                                         ",
	"                                         ",
}

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	tx, ty := len([]rune(Title[0])), len(Title)

	// Title
	if v, err := g.SetView("title", 0, 0, tx, ty); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		for _, line := range Title {
			fmt.Fprintln(v, line)
		}
	}

	// Logs
	if v, err := g.SetView("logs", 0, ty, tx, max(maxY-1, ty+1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " LOGS "
		v.Wrap = true
		// fmt.Fprintln(v, "")
	}

	vGap, hGap := 1, 2
	pWidth := 25

	// Frame
	fx1 := tx + 2
	fy1 := 0
	fx2 := max(maxX-1, fx1+1)
	fy2 := max(maxY-1, 1)
	if _, err := g.SetView("frame", fx1, fy1, fx2, fy2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	midHeight := (fy2 - fy1 - 3*vGap) / 2

	// Groups
	baseX := fx1 + hGap
	baseY := fy1 + vGap
	if v, err := g.SetView("groups", baseX, baseY, baseX+pWidth, baseY+max(midHeight, 1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " GROUPS "
		v.Wrap = true
		// fmt.Fprintln(v, "")
	}

	// Clients
	baseY += max(midHeight, 1) + vGap
	if v, err := g.SetView("clients", baseX, baseY, baseX+pWidth, max(fy2-vGap, baseY+1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " CLIENTS "
		v.Wrap = true
		// fmt.Fprintln(v, "")
	}

	// Chat
	baseX += pWidth + hGap
	baseY = fy1 + vGap
	if v, err := g.SetView("chat", baseX, baseY, max(fx2-hGap, baseX+hGap), max(fy2-vGap, baseY+1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " CHAT "
		v.Wrap = true
		v.Autoscroll = false
		// fmt.Fprintln(v, "")
	}

	return nil
}

func SetKeybindings(g *gocui.Gui) error {
	// Switch views
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, SwitchView); err != nil {
		return err
	}

	// View scrolling
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, ScrollUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, ScrollDown); err != nil {
		return err
	}

	// Group Selection
	if err := g.SetKeybinding("groups", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("groups", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	// Quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}

	return nil
}
