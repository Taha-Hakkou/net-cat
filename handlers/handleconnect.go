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
	messageLog []string
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

	sendHistory(conn)

	joinMsg := fmt.Sprintf("✅ %s has joined our chat...", name)
	broadcastToOthers(groupName, joinMsg, conn)
	logs(groupName, joinMsg+"\n")
	addToHistory(joinMsg)

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
					removeLogFile(groupName)
				}
			}
			groupsMu.Unlock()

			leaveMsg := fmt.Sprintf("🔴 %s has left our chat...", name)
			broadcastToOthers(groupName, leaveMsg, conn)
			logs(groupName, leaveMsg+"\n")
			addToHistory(leaveMsg)
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
					addToHistory(changeMsg)
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

		addToHistory(formatted)
		// Broadcast to all OTHER users
		broadcastToOthers(groupName, formatted, conn)
		logs(groupName, formatted+"\n")
		
		// Send prompt to THIS user only
		sendPrompt(conn, name)
	}
}

func CreateGroup(conn net.Conn) (string, string, error) {
	groupsMu.Lock()
	if len(Groups) >= MAX_GROUPS {
		groupsMu.Unlock()
		conn.Write([]byte("Server is full. Choose again!\n"))
		return "", "", nil
	}

	groupName := fmt.Sprintf("room%d", len(Groups)+1)
	Groups[groupName] = make(map[net.Conn]string)
	groupsMu.Unlock()

	name, err := getClientName(conn, groupName)
	if err != nil {
		fmt.Println("Invalid name. Disconnecting client.")
		return "", "", err
	}

	groupsMu.Lock()
	if _, ok := Groups[groupName]; ok {
		Groups[groupName][conn] = name
	}
	groupsMu.Unlock()

	return groupName, name, nil
}

func JoinGroup(conn net.Conn) (string, string, error) {
	reader := bufio.NewReader(conn)

	groupsMu.RLock()
	if len(Groups) == 0 {
		groupsMu.RUnlock()
		conn.Write([]byte("No groups found. Choose again!\n"))
		return "", "", nil
	}
	groupsMu.RUnlock()

	for {
		groupsMu.RLock()
		for grp := range Groups {
			conn.Write([]byte(fmt.Sprintf("* %s [%d users]\n", grp, len(Groups[grp]))))
		}
		groupsMu.RUnlock()

		conn.Write([]byte("[ENTER GROUP NAME]:"))
		groupName, _ := reader.ReadString('\n')
		groupName = strings.ToLower(strings.TrimSpace(groupName))

		groupsMu.RLock()
		_, ok := Groups[groupName]
		groupsMu.RUnlock()
		if !ok {
			conn.Write([]byte("Group not found. Choose again!\n"))
			continue
		}

		name, err := getClientName(conn, groupName)
		if err != nil {
			fmt.Println("Invalid name. Disconnecting client.")
			return "", "", err
		}

		groupsMu.RLock()
		if len(Groups[groupName]) >= MAX_CLIENTS {
			groupsMu.RUnlock()
			conn.Write([]byte("Group is full. Choose again!\n"))
			continue
		}
		groupsMu.RUnlock()

		groupsMu.Lock()
		Groups[groupName][conn] = name
		groupsMu.Unlock()

		return groupName, name, nil
	}
}



// Broadcast to all users EXCEPT the sender
