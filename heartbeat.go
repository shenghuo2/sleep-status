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
				setStatusToSleep()
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

// setStatusToSleep 安全地将状态设置为睡眠，避免死锁
func setStatusToSleep() {
	// 先检查当前状态，避免不必要的锁定和文件操作
	if ConfigData.Sleep {
		return
	}
	
	// 使用单独的锁，避免与其他函数产生死锁
	mutex.Lock()
	defer mutex.Unlock()
	
	// 再次检查状态，避免在获取锁期间状态被改变
	if ConfigData.Sleep {
		return
	}
	
	ConfigData.Sleep = true
	record := SleepRecord{
		Action: "sleep",
		Time:   time.Now().Format(time.RFC3339),
	}
	
	// 保存配置但不要在这里处理错误，避免阻塞心跳检测器
	_ = SaveConfig()
	
	// 使用无锁版本保存记录，避免死锁
	go func(r SleepRecord) {
		_ = SaveSleepRecordNoLock(r)
	}(record)
}
