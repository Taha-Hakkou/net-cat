package nctui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func SwitchView(g *gocui.Gui, v *gocui.View) error {
	cur := g.CurrentView()
	views := g.Views()
	for i, view := range views {
		if cur == nil || view == cur {
			var next *gocui.View
			if cur == nil {
				next = views[0]
			} else {
				next = views[(i+1)%len(views)]
				cur.FgColor = gocui.ColorDefault
			}
			next.FgColor = gocui.ColorCyan
			cv, _ := g.SetCurrentView(next.Name())
			h, _ := g.View("settings")
			fmt.Fprintln(h, cv.Name())
			return nil
		}
	}
	return nil
}

func ScrollUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	if oy > 0 {
		// fmt.Fprintln(v, strconv.Itoa(ox))
		return v.SetOrigin(ox, oy-1)
	}
	return nil
}

func ScrollDown(g *gocui.Gui, v *gocui.View) error {
	cv := g.CurrentView()
	ox, oy := cv.Origin()
	return cv.SetOrigin(ox, oy+1)
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
