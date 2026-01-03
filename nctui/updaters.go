package nctui

import (
	"fmt"
	"os"

	zone "zone/handlers"

	"github.com/jroimartin/gocui"
)

func UpdateLog(g *gocui.Gui) error {
	v, _ := g.View("logs")
	v.Clear()
	fmt.Fprintln(v, "")
	bytes, err := os.ReadFile("nc.log")
	if err == nil {
		fmt.Fprintln(v, string(bytes))
	}
	return nil
}

var (
	SelectedGroup string
	SelectedIndex int
	groups        []string
)

func UpdateGroups(g *gocui.Gui) error {
	v, _ := g.View("groups")
	v.Clear()
	fmt.Fprintln(v, "")

	groups = make([]string, 0, len(zone.Groups))
	for k := range zone.Groups {
		groups = append(groups, k)
	}
	for i := 1; i < len(groups); i++ {
		j := i //////////// check if sort is correct
		for j > 0 && groups[j-1] > groups[j] {
			groups[j-1], groups[j] = groups[j], groups[j-1]
			j--
		}
	}

	for i, group := range groups {
		if SelectedGroup == "" {
			SelectedGroup = group
		}
		// if group == SelectedGroup {
		if i == SelectedIndex {
			fmt.Fprintf(v, ">> %s\n", group)
		} else {
			fmt.Fprintf(v, "  %s\n", group)
		}
		// fmt.Fprintln(v, group)
	}
	return nil
}

func UpdateClients(g *gocui.Gui) error {
	v, _ := g.View("clients")
	v.Clear()
	fmt.Fprintln(v, "")
	for _, client := range zone.Groups[SelectedGroup] {
		fmt.Fprintln(v, client)
	}
	return nil
}

func UpdateChat(g *gocui.Gui) error {
	v, _ := g.View("chat")
	v.Clear()
	fmt.Fprintln(v, "")
	bytes, _ := os.ReadFile(fmt.Sprintf("logs/%s.txt", SelectedGroup))
	fmt.Fprintln(v, string(bytes))
	return nil
}
