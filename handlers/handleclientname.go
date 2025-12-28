package zone

import (
	"bufio"
	"net"
	"strings"
	"sync"
)

var (
	clients    = make(map[net.Conn]string)
	clientsMu  sync.Mutex
	messageLog []string
	logMu      sync.Mutex
)

// the main function to handle client name validity
func getClientName(conn net.Conn) (string, error) {


	reader := bufio.NewReader(conn)

	for {
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		name = strings.TrimSpace(name)
		booln := Isnameexist(name)
		bool2 := Validname(name)

		if !booln {
			conn.Write([]byte("this name is exist .\n"))
			conn.Write([]byte("[ENTER YOUR NAME]:"))
			continue
		}
		if !bool2 {
			conn.Write([]byte("this is not valid name.\n"))
			conn.Write([]byte("[ENTER YOUR NAME]:"))
			continue
		}

		if name == "" {
			conn.Write([]byte("Name cannot be empty.\n"))
			conn.Write([]byte("[ENTER YOUR NAME]:"))
			continue
		}

		return name, nil
	}
}

func Isnameexist(name string) bool {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for _, k := range clients {
		if name == k {
			return false
		}
	}
	return true
}

func Validname(name string) bool {
	if len(name) > 20 {
		return false
	}
	for _, i := range name {
		if i < 32 || i > 126 {
			return false
		}
	}

	return true
}
