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

// Removes log file after group deletion
func removeLogFile(groupName string) {
	err := os.Remove(fmt.Sprintf("logs/%s.txt", groupName))
	if err != nil {
		msg := fmt.Sprintf("Error removing file: %s", err)
		WriteToMainLog(msg)
		return
	}
	WriteToMainLog("File removed successfully")
}

func WriteToMainLog(msg string) {
	file, _ := os.OpenFile("nc.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	defer file.Close()
	file.WriteString(msg + "\n")
}
