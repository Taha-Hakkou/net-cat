package zone

import (
	"bufio"
	"net"
	"strings"
)

// the main function to handle client name validity
func getClientName(conn net.Conn, groupName string) (string, error) {
	reader := bufio.NewReader(conn)

	for {
		conn.Write([]byte("[ENTER YOUR NAME]:"))
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		name = strings.TrimSpace(name)
		booln := Isnameexist(name, groupName)
		bool2 := Validname(name)

		if !booln {
			conn.Write([]byte("this name is exist .\n"))
			continue
		}
		if !bool2 {
			conn.Write([]byte("this is not valid name.\n"))
			continue
		}

		if name == "" {
			conn.Write([]byte("Name cannot be empty.\n"))
			continue
		}

		return name, nil
	}
}
