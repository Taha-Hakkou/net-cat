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
		answer, _ := reader.ReadString('\n') // error handling like in get client
		var g, n string
		var e error
		if strings.TrimSpace(answer) == "2" {
			// create new group
			g, n, e = CreateGroup(conn)
		} else {
			// join existing group
			g, n, e = JoinGroup(conn)
		}
		if g == "" && n == "" {
			if e == nil {
				continue
			} else {
				return
			}
		}
		groupName = g
		name = n
		break
	}

	sendHistory(conn)

	joinMsg := fmt.Sprintf("✅ %s has joined our chat...", name)
	broadcast(groupName, joinMsg, conn, isSystemMessage) // groupName added
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
			clientsMu.Lock()
			delete(Groups[groupName], conn) // Clients -> Groups[groupName]
			clientsMu.Unlock()

			leaveMsg := fmt.Sprintf("🔴 %s has left our chat...", name)
			broadcast(groupName, leaveMsg, conn, isSystemMessage) // groupName added
			prompt(groupName)
			isSystemMessage = true
			logs(groupName, leaveMsg+"\n")
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
		logs(groupName, formatted+"\n")
		flag = true

	}
}

// Prompts client to create a new group chat
func CreateGroup(conn net.Conn) (string, string, error) {
	groupsMu.Lock()
	if len(Groups) >= MAX_GROUPS {
		groupsMu.Unlock()
		/*_, err := */ conn.Write([]byte("Server is full. Choose again!\n"))
		// if err != nil {
		// 	fmt.Println("error writing to the client", err) // logs should be in terminal ui
		// }
		// write with error handling could be wrapped in separate function !
		// continue
		return "", "", nil
	}

	groupName := fmt.Sprintf("room%d", len(Groups)+1) // handle group names when all users exit
	Groups[groupName] = make(map[net.Conn]string)

	name, err := getClientName(conn, groupName)
	if err != nil {
		fmt.Println("Invalid name. Disconnecting client.")
		// return
		return "", "", err
	}
	Groups[groupName][conn] = name
	groupsMu.Unlock()

	return groupName, name, nil
}

// Prompts client to join an existing group chat
func JoinGroup(conn net.Conn) (string, string, error) {
	reader := bufio.NewReader(conn)
	// groupsMu.Lock()
	if len(Groups) == 0 {
		// groupsMu.Unlock() // must be related to group creation logic !!!!!!!!!!!!!
		conn.Write([]byte("No groups found. Choose again!\n"))
		// continue
		return "", "", nil
	}
	for {
		for grp := range Groups {
			s := fmt.Sprintf("* %s [%d users]\n", grp, len(Groups[grp]))
			conn.Write([]byte(s))
		}
		conn.Write([]byte("[ENTER GROUP NAME]:"))
		groupName, _ := reader.ReadString('\n')
		groupName = strings.ToLower(strings.TrimSpace(groupName))
		groupsMu.Lock() // in case it got deleted, until client joins group
		_, ok := Groups[groupName]
		if !ok {
			conn.Write([]byte("Group not found. Choose again!\n")) // locks + group delete logic
			continue                                               // choose another group
		}

		name, err := getClientName(conn, groupName)
		if err != nil {
			fmt.Println("Invalid name. Disconnecting client.")
			// return
			return "", "", err
		}

		clientsMu.Lock()                           // must be specific to that group
		if len(Groups[groupName]) >= MAX_CLIENTS { // Clients -> Groups[groupName]
			clientsMu.Unlock()
			_, err := conn.Write([]byte("Server full. Choose again!\n"))
			if err != nil {
				fmt.Println("error writing to the client", err)
			}
			continue
		}
		clientsMu.Unlock()

		Groups[groupName][conn] = name
		groupsMu.Unlock()

		return groupName, name, nil
	}
}
