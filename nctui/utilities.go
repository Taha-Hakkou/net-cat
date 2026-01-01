package nctui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func SetKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, SwitchView); err != nil {
		return err
	}
	if err := g.SetKeybinding("chat", gocui.KeyArrowUp, gocui.ModNone, ScrollUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("chat", gocui.KeyArrowDown, gocui.ModNone, ScrollDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}

	// Group Selection
	if err := g.SetKeybinding("groups", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("groups", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	// if err := g.SetKeybinding("groups", gocui.KeyEnter, gocui.ModNone, choose); err != nil {
	// 	return err
	// }
	return nil
}

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Settings
	if v, err := g.SetView("settings", 0, 0, maxX/4-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " SETTINGS "
		v.Wrap = true
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "settings")
	}

	// Groups
	if v, err := g.SetView("groups", 0, maxY/2+1, maxX/4-1, max(maxY/2+2, maxY-1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " GROUPS "
		v.Wrap = true
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "groups")
	}

	// Clients
	if v, err := g.SetView("clients", maxX/4, 0, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " CLIENTS "
		v.Wrap = true
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "clients")
	}

	// Chat
	if v, err := g.SetView("chat", maxX/2+1, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " CHAT "
		v.Wrap = true
		v.Autoscroll = false
	}

	return nil
}
