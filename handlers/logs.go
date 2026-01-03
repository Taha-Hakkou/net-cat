package zone

import (
	"fmt"
	"log"
	"os"
)

// a function to store logs
func logs(groupName, text string) {
	_, err := os.Stat("logs/")
	if err != nil {
		os.MkdirAll("logs/", 0o755)
	}
	logFile := fmt.Sprintf("logs/%s.txt", groupName)
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		// log in LOGS
	}
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		log.Fatal(err) //
	}
}
