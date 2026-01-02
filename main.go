package main

import (
	"errors"
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

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// Set GUI managers and key bindings
	g.SetManagerFunc(nctui.Layout)
	err = nctui.SetKeybindings(g)
	if err != nil {
		log.Panicln(err)
	}

	// fmt.Printf("Listening on port %s\n", address) // log in LOGS

	go func() {
		for {
			time.Sleep(2 * time.Second)
			g.Update(nctui.UpdateGroups)
			g.Update(nctui.UpdateClients)
			g.Update(nctui.UpdateChat)
		}
	}()

	// Accept loop
	go acceptLoop(ln)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func acceptLoop(ln net.Listener) { // maybe doesn't need go-routine ?!!
	for {
		conn, err := ln.Accept()
		if err != nil { // not understandable !!!!!
			if errors.Is(err, net.ErrClosed) {
				// listener closed → normal shutdown
				return
			}
			fmt.Println("Accept error:", err)
			continue // return only the routine ?? | OR continue
		}
		go zone.HandleConnection(conn)
	}
}

// TODO:
// client ips
// listener port
