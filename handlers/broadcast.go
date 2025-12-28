package zone

import (
	"net"
)

// a function to broadcast(send) messages to other users
func broadcast(message string, excludeConn net.Conn, flg bool) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
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
