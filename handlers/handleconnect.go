package zone

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	MAX_GROUPS  int = 4
	MAX_CLIENTS int = 10
)

// Declaring global variables
var (
	Groups     = make(map[string]map[net.Conn]string)
	groupsMu   sync.RWMutex
	clientsMu  sync.RWMutex
	// Change: Store message history per group instead of globally
	messageLog = make(map[string][]string)
	logMu      sync.Mutex
)

// the main function to handle connections(Name,Limit,prompt,broadcast, connect and disconnect...)
func HandleConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write(peng())

	var name, groupName string

	conn.Write([]byte("[1] Join Existing Group Chat (default)\n"))
	conn.Write([]byte("[2] Create New Group Chat\n"))
	reader := bufio.NewReader(conn)

	for {
		conn.Write([]byte("[ENTER YOUR ANSWER]:"))
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		var g, n string
		var e error

		if answer == "2" {
			g, n, e = CreateGroup(conn)
		} else {
			g, n, e = JoinGroup(conn)
		}

		if e != nil {
			return
		}

		if g != "" && n != "" {
			groupName = g
			name = n
			break
		}
	}

	// Send history for THIS group only
	sendHistory(conn, groupName)

	joinMsg := fmt.Sprintf("✅ %s has joined our chat...", name)
	broadcastToOthers(groupName, joinMsg, conn)
	logs(groupName, joinMsg+"\n")
	addToHistory(groupName, joinMsg)

	// Send initial prompt to THIS user only
	sendPrompt(conn, name)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			// Handle disconnect
			groupsMu.Lock()
			if _, ok := Groups[groupName]; ok {
				delete(Groups[groupName], conn)
				if len(Groups[groupName]) == 0 {
					delete(Groups, groupName)
					// Clean up history for empty group
					logMu.Lock()
					delete(messageLog, groupName)
					logMu.Unlock()
					removeLogFile(groupName)
				}
			}
			groupsMu.Unlock()

			leaveMsg := fmt.Sprintf("🔴 %s has left our chat...", name)
			broadcastToOthers(groupName, leaveMsg, conn)
			logs(groupName, leaveMsg+"\n")
			addToHistory(groupName, leaveMsg)
			return
		}
		message = strings.TrimSpace(message)

		// Change name
		if message == "/name" {
			for {
				newName, err := getClientName(conn, groupName)
				if err != nil {
					conn.Write([]byte("Invalid name, try again.\n"))
					continue
				}

				groupsMu.Lock()
				if _, ok := Groups[groupName]; ok {
					oldName := name
					name = newName
					Groups[groupName][conn] = newName
					groupsMu.Unlock()

					changeMsg := fmt.Sprintf("🔁 %s changed name to %s", oldName, newName)
					broadcastToOthers(groupName, changeMsg, conn)
					logs(groupName, changeMsg+"\n")
					addToHistory(groupName, changeMsg)
					break
				} else {
					groupsMu.Unlock()
					conn.Write([]byte("Group no longer exists.\n"))
					return
				}
			}
			// Send prompt to THIS user only after name change
			sendPrompt(conn, name)
			continue
		}

		if message == "" || !Isvalidmessage(message) {
			// Send prompt back to THIS user only
			sendPrompt(conn, name)
			continue
		}

		formatted := fmt.Sprintf("[%s][%s]: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			name,
			message)

		addToHistory(groupName, formatted)
		// Broadcast to all OTHER users
		broadcastToOthers(groupName, formatted, conn)
		logs(groupName, formatted+"\n")
		
		// Send prompt to THIS user only
		sendPrompt(conn, name)
	}
}