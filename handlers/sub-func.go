package zone

func Isnameexist(name string) bool {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for _, k := range Clients {
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

func Isvalidmessage(msg string) bool {
	for _, i := range msg {
		if i < 32 || i > 126 {
			return false
		}
	}
	return true
}
