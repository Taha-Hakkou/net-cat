package zone

import (
	"fmt"
	"net"
	"time"
)

// a function to broadcast(send) messages to other users
func broadcastToOthers(groupName string, message string, senderConn net.Conn) {
	groupsMu.RLock()
	group, ok := Groups[groupName]
	if !ok {
		groupsMu.RUnlock()
		return
	}

	// Create a copy of connections to avoid holding lock during writes
	connections := make([]struct {
		conn net.Conn
		name string
	}, 0, len(group))

	for conn, name := range group {
		if conn != senderConn {
			connections = append(connections, struct {
				conn net.Conn
				name string
			}{conn, name})
		}
	}
	groupsMu.RUnlock()

	// Send message to all other users
	for _, c := range connections {
		// Send the message
		c.conn.Write([]byte("\n" + message + "\n"))
		// Send their prompt back
		sendPrompt(c.conn, c.name)
	}
}

// Send prompt to a SPECIFIC user only (not broadcast)
func sendPrompt(conn net.Conn, name string) {
	formatted := fmt.Sprintf("[%s][%s]:",
		time.Now().Format("2006-01-02 15:04:05"),
		name)
	conn.Write([]byte(formatted))
}
