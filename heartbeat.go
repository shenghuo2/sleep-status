package main

import (
	"sync"
	"time"
)

var (
	lastHeartbeat     time.Time
	heartbeatMutex    sync.RWMutex
	heartbeatChecker  *time.Timer
	isCheckerRunning  bool
	checkerMutex      sync.Mutex
)

// updateLastHeartbeat updates the timestamp of the last heartbeat
func updateLastHeartbeat() {
	heartbeatMutex.Lock()
	lastHeartbeat = time.Now()
	heartbeatMutex.Unlock()
}

// getLastHeartbeat returns the timestamp of the last heartbeat
func getLastHeartbeat() time.Time {
	heartbeatMutex.RLock()
	defer heartbeatMutex.RUnlock()
	return lastHeartbeat
}

// startHeartbeatChecker starts the heartbeat checking routine
func startHeartbeatChecker() {
	checkerMutex.Lock()
	defer checkerMutex.Unlock()

	if !ConfigData.HeartbeatEnabled || isCheckerRunning {
		return
	}

	isCheckerRunning = true
	updateLastHeartbeat() // Initialize last heartbeat time

	heartbeatChecker = time.NewTimer(time.Duration(ConfigData.HeartbeatTimeout) * time.Second)

	go func() {
		for {
			<-heartbeatChecker.C
			
			if !ConfigData.HeartbeatEnabled {
				stopHeartbeatChecker()
				return
			}

			timeSinceLastHeartbeat := time.Since(getLastHeartbeat())
			if timeSinceLastHeartbeat > time.Duration(ConfigData.HeartbeatTimeout)*time.Second {
				// Mark as sleeping if heartbeat timeout exceeded
				mutex.Lock()
				if !ConfigData.Sleep {
					ConfigData.Sleep = true
					record := SleepRecord{
						Action: "sleep",
						Time:   time.Now().Format(time.RFC3339),
					}
					SaveConfig()
					SaveSleepRecord(record)
				}
				mutex.Unlock()
			}

			// Reset timer for next check
			heartbeatChecker.Reset(time.Duration(ConfigData.HeartbeatTimeout) * time.Second)
		}
	}()
}

// stopHeartbeatChecker stops the heartbeat checking routine
func stopHeartbeatChecker() {
	checkerMutex.Lock()
	defer checkerMutex.Unlock()

	if isCheckerRunning {
		heartbeatChecker.Stop()
		isCheckerRunning = false
	}
}
