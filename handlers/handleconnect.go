package zone

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const MAX_GROUPS int = 4

// Declaring global variables
var (
	// groups vars
	Groups   = make(map[string]map[net.Conn]string)
	groupsMu sync.Mutex
	////
	// Clients    = make(map[net.Conn]string)
	// substituted in:
	// handleconnect.go -> 4
	// sub-func.go -> 1
	// broadcast.go -> 3
	clientsMu  sync.Mutex
	messageLog []string
	logMu      sync.Mutex // log for every group
)

// the main function to handle connections(Name,Limit,prompt,broadcast, connect and disconnect...)
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	isSystemMessage := false

	conn.Write(peng())

	// name, err := getClientName(conn)
	// if err != nil {
	// 	fmt.Println("Invalid name. Disconnecting client.")
	// 	return
	// }
	var name string

	// groups logic /////////////////////
	reader := bufio.NewReader(conn) // reader exists in get client
	conn.Write([]byte("1. Join Existing Group Chat\n"))
	conn.Write([]byte("2. Create New Group Chat\n"))

	var groupName string
	for {
		conn.Write([]byte("[YOUR ANSWER][default->1]:"))
		answer, _ := reader.ReadString('\n') // error handling like in get client
		if strings.TrimSpace(answer) == "2" {
			// create new group
			groupsMu.Lock()
			if len(Groups) >= MAX_GROUPS {
				groupsMu.Unlock()
				/*_, err := */ conn.Write([]byte("Server is full. Choose again!\n"))
				// if err != nil {
				// 	fmt.Println("error writing to the client", err) // logs should be in terminal ui
				// }
				continue
			}

			n, err := getClientName(conn, groupName)
			if err != nil {
				fmt.Println("Invalid name. Disconnecting client.")
				return
			}
			name = n // because n work as local var to this scope

			groupName = fmt.Sprintf("room%d", len(Groups)+1) // handle group names when all users exit
			Groups[groupName] = make(map[net.Conn]string)
			Groups[groupName][conn] = name
			groupsMu.Unlock()
			break // break the infinite loop
		} else {
			// join existing group
			// groupsMu.Lock()
			if len(Groups) == 0 {
				// groupsMu.Unlock() // must be related to group creation logic !!!!!!!!!!!!!
				conn.Write([]byte("No groups found. Choose again!\n"))
				continue
			}
			for {
				for grp := range Groups {
					s := fmt.Sprintf("* %s [%d users]", grp, len(Groups[grp]))
					conn.Write([]byte(s))
				}
				conn.Write([]byte("[GROUP NAME]:"))
				groupName, _ = reader.ReadString('\n')
				groupName = strings.ToLower(strings.TrimSpace(groupName))
				groupsMu.Lock() // in case it got deleted, until client hoins group
				_, ok := Groups[groupName]
				if !ok {
					conn.Write([]byte("Group not found. Choose again!\n")) // locks + group delete logic
					continue                                               // choose another group
				}

				name, err := getClientName(conn, groupName)
				if err != nil {
					fmt.Println("Invalid name. Disconnecting client.")
					return
				}
				Groups[groupName][conn] = name
				groupsMu.Unlock()
				break
				// Clients = Groups[answer] // copy or refers to ?!
			}
		}
	}
	/////////////////////////

	clientsMu.Lock()
	if len(Groups[groupName]) >= 10 { // Clients -> Groups[groupName]
		clientsMu.Unlock()
		_, err := conn.Write([]byte("Server full. Try again later.\n"))
		if err != nil {
			fmt.Println("erorr writing to the client", err)
		}
		return
	}

	Groups[groupName][conn] = name // Clients -> Groups[groupName]

	clientsMu.Unlock()

	//////////////////////////

	sendHistory(conn)

	joinMsg := fmt.Sprintf("✅ %s has joined our chat...", name)
	broadcast(groupName, joinMsg, conn, isSystemMessage) // groupName added
	isSystemMessage = true
	logs(joinMsg + "\n")
	addToHistory(joinMsg)

	// commented
	// reader := bufio.NewReader(conn)
	////////

	flag := true
	for {
		if flag {
			prompt(groupName)
			isSystemMessage = false
		}
		message, err := reader.ReadString('\n')
		if err != nil {
			clientsMu.Lock()
			delete(Groups[groupName], conn) // Clients -> Groups[groupName]
			clientsMu.Unlock()

			leaveMsg := fmt.Sprintf("🔴 %s has left our chat...", name)
			broadcast(groupName, leaveMsg, conn, isSystemMessage) // groupName added
			prompt(groupName)
			isSystemMessage = true
			logs(leaveMsg + "\n")
			addToHistory(leaveMsg)
			flag = true
			return
		}
		message = strings.TrimSpace(message)

		if message == "" || !Isvalidmessage(message) {
			flag = false
			clientsMu.Lock()
			clientName, ok := Groups[groupName][conn] // Clients -> Groups[groupName]
			clientsMu.Unlock()
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
		broadcast(groupName, formatted, conn, isSystemMessage) // groupName added
		isSystemMessage = false
		logs(formatted + "\n")
		flag = true

	}
}
