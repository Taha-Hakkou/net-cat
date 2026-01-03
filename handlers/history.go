package zone

// Updated: Add message to specific group's history
func addToHistory(groupName, message string) {
	logMu.Lock()
	defer logMu.Unlock()
	messageLog[groupName] = append(messageLog[groupName], message)
}
