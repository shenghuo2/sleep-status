package main

import (
	"fmt"
	"os"
	"time"
)

var logFilePath = "./access.log"

// LogAccess logs the access information to access.log
// LogAccess 将访问信息记录到 access.log
func LogAccess(ip string) error {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	logEntry := fmt.Sprintf("%s - %s\n", time.Now().Format(time.RFC3339), ip)
	if _, err := file.WriteString(logEntry); err != nil {
		return err
	}
	fmt.Println(logEntry) // Print log entry to console
	return nil
}
