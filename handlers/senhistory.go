package zone
import "net"
// sending history for new joigned clients
func sendHistory(conn net.Conn, groupName string) {
	logMu.Lock()
	defer logMu.Unlock()
	
	if history, ok := messageLog[groupName]; ok {
		for _, msg := range history {
			conn.Write([]byte(msg + "\n"))
		}
	}
}