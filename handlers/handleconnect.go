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
	Groups   = make(map[string]map[net.Conn]string)
	groupsMu sync.Mutex
	// before groups: Clients    = make(map[net.Conn]string)
	clientsMu  sync.Mutex // clients for every group
	messageLog []string
	logMu      sync.Mutex // log for every group
)
// the main function to handle connections(Name,Limit,prompt,broadcast, connect and disconnect...)
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	isSystemMessage := false

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
		// else: continue asking
	}

	sendHistory(conn)

	joinMsg := fmt.Sprintf("✅ %s has joined our chat...", name)
	broadcast(groupName, joinMsg, conn, isSystemMessage)
	isSystemMessage = true
	logs(groupName, joinMsg+"\n")
	addToHistory(joinMsg)

	flag := true
	for {
		if flag {
			prompt(groupName)
			isSystemMessage = false
		}

		message, err := reader.ReadString('\n')
		if err != nil {
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
			broadcast(groupName, leaveMsg, conn, isSystemMessage)
			prompt(groupName)
			isSystemMessage = true
			logs(groupName, leaveMsg+"\n")
			addToHistory(leaveMsg)
			flag = true
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
					broadcast(groupName, changeMsg, conn, isSystemMessage)
					isSystemMessage = true
					logs(groupName, changeMsg+"\n")
					addToHistory(changeMsg)
					break
				} else {
					groupsMu.Unlock()
					conn.Write([]byte("Group no longer exists.\n"))
					return
				}
			}
			continue
		}

		if message == "" || !Isvalidmessage(message) {
			flag = false
			groupsMu.Lock()
			clientName, ok := Groups[groupName][conn]
			groupsMu.Unlock()
			if !ok {
				return
			}
			formatted1 := fmt.Sprintf("[%s][%s]:",
				time.Now().Format("2006-01-02 15:04:05"),
				clientName)
			conn.Write([]byte(formatted1))
			continue
		}

		formatted := fmt.Sprintf("[%s][%s]: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			name,
			message)

		addToHistory(formatted)
		broadcast(groupName, formatted, conn, isSystemMessage)
		isSystemMessage = false
		logs(groupName, formatted+"\n")
		flag = true
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

	groupsMu.Lock()
	if len(Groups) == 0 {
		groupsMu.Unlock()
		conn.Write([]byte("No groups found. Choose again!\n"))
		return "", "", nil
	}
	groupsMu.Unlock()

	for {
		groupsMu.Lock()
		for grp := range Groups {
			conn.Write([]byte(fmt.Sprintf("* %s [%d users]\n", grp, len(Groups[grp]))))
		}
		groupsMu.Unlock()

		conn.Write([]byte("[ENTER GROUP NAME]:"))
		groupName, _ := reader.ReadString('\n')
		groupName = strings.ToLower(strings.TrimSpace(groupName))

		groupsMu.Lock()
		_, ok := Groups[groupName]
		groupsMu.Unlock()
		if !ok {
			conn.Write([]byte("Group not found. Choose again!\n"))
			continue
		}

		name, err := getClientName(conn, groupName)
		if err != nil {
			fmt.Println("Invalid name. Disconnecting client.")
			return "", "", err
		}

		groupsMu.Lock()
		if len(Groups[groupName]) >= MAX_CLIENTS {
			groupsMu.Unlock()
			conn.Write([]byte("Server full. Choose again!\n"))
			continue
		}
		Groups[groupName][conn] = name
		groupsMu.Unlock()

		return groupName, name, nil
	}
}
