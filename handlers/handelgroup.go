package zone

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

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

	// Initialize history for new group
	logMu.Lock()
	messageLog[groupName] = []string{}
	logMu.Unlock()

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
