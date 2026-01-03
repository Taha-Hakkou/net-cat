package nctui

import (
	"github.com/jroimartin/gocui"
)

// Switches between GUI views
func SwitchView(g *gocui.Gui, v *gocui.View) error {
	cur := g.CurrentView()
	views := g.Views()
	for i, view := range views {
		if cur == nil || view == cur {
			var next *gocui.View
			var j int
			if cur != nil {
				j = (i + 1) % len(views)
				cur.FgColor = gocui.ColorDefault
			}
			for views[j].Name() == "title" || views[j].Name() == "frame" {
				j = (j + 1) % len(views)
			}
			next = views[j]
			next.FgColor = gocui.ColorCyan
			g.SetCurrentView(next.Name())
			return nil
		}
	}
	return nil
}

// Scrolling View's content

var scrollableViews = []string{"logs", "chat"}

func ScrollUp(g *gocui.Gui, v *gocui.View) error { // v = nil
	v = g.CurrentView()
	if v == nil {
		return nil
	}
	for _, sv := range scrollableViews {
		if sv == v.Name() {
			ox, oy := v.Origin()
			if oy > 0 {
				return v.SetOrigin(ox, oy-1)
			}
			return nil
		}
	}
	return nil
}

func ScrollDown(g *gocui.Gui, v *gocui.View) error { // v = nil
	v = g.CurrentView()
	if v == nil {
		return nil
	}
	for _, sv := range scrollableViews {
		if sv == v.Name() {
			ox, oy := v.Origin()
			return v.SetOrigin(ox, oy+1)
		}
	}
	return nil
}

// Group Selection

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if SelectedIndex < len(groups)-1 {
		SelectedIndex++
		SelectedGroup = groups[SelectedIndex]
	}
	UpdateGroups(g)
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if SelectedIndex > 0 {
		SelectedIndex--
		SelectedGroup = groups[SelectedIndex]
	}
	UpdateGroups(g)
	return nil
}

// Quit
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
