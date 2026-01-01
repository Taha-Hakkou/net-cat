package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"zone/nctui"

	"github.com/jroimartin/gocui"

	zone "zone/handlers"
)

func main() {
	address := ":8989" // default port
	if len(os.Args) == 2 {
		newadress, err := strconv.Atoi(os.Args[1])
		if err != nil || newadress < 1024 || newadress > 65000 { // checking the validity of the port 1024>=port<=65000
			fmt.Println("check the validity of the port")
			return
		}
		address = ":" + strconv.Itoa(newadress)
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer ln.Close()

	//**********************************************************
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// Set GUI managers and key bindings
	g.SetManagerFunc(nctui.Layout)
	nctui.SetKeybindings(g)

	// fmt.Printf("Listening on port %s\n", address)

	go func() {
		for {
			time.Sleep(time.Second)
			g.Update(updateChat)
			g.Update(updateClients)
		}
	}()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}
			go zone.HandleConnection(conn)
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	// chat
	// groups
	// client names/ips

	// server status ?
	// send errors/messages with color
	// listener port

	//**********************************************************
}

func updateChat(g *gocui.Gui) error {
	v, _ := g.View("chat")
	v.Clear()
	fmt.Fprintln(v, "")
	bytes, _ := os.ReadFile("logs.txt")
	fmt.Fprintln(v, string(bytes))
	return nil
}

func updateClients(g *gocui.Gui) error {
	v, _ := g.View("clients")
	v.Clear()
	fmt.Fprintln(v, "")
	for _, client := range zone.Clients {
		fmt.Fprintln(v, client)
	}
	return nil
}