package zone

import (
	"fmt"
	"net"
	"time"
)

// a function to broadcast(send) messages to other users
func broadcast(groupName string, message string, excludeConn net.Conn, flg bool) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range Groups[groupName] { // Clients -> Groups[groupName]
		if conn != excludeConn {
			if flg {
				conn.Write([]byte(message))
			}

			if !flg {
				conn.Write([]byte("\n" + message + "\n"))
			}
		}
	}
}

// a function to send the propmt to users
func prompt(groupName string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range Groups[groupName] { // Clients -> Groups[groupName]
		clientName, ok := Groups[groupName][conn] // Clients -> Groups[groupName]
		if !ok {
			continue
		}

		formatted1 := fmt.Sprintf("[%s][%s]:",
			time.Now().Format("2006-01-02 15:04:05"),
			clientName)
		_, err := conn.Write([]byte(formatted1))
		if err != nil {
			fmt.Println("error print the propmt", err)
		}

	}
}
