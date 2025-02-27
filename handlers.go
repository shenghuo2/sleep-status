package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Response struct to standardize API responses
// Response 结构体标准化 API 响应
type Response struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
}

// SleepRecord struct to hold the sleep records
// SleepRecord 结构体保存入睡和醒来时间的记录
type SleepRecord struct {
	Action string `json:"action"`
	Time   string `json:"time"`
}

// RecordsResponse struct to hold the sleep records response
// RecordsResponse 结构体保存睡眠记录的响应
type RecordsResponse struct {
	Success bool          `json:"success"`
	Records []SleepRecord `json:"records"`
}

// StatusHandler handles the /status route and returns the current sleep status
// StatusHandler 处理 /status 路由并返回当前的 sleep 状态
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	// 异步记录访问日志，避免阻塞主处理流程
	go LogAccess(clientIP, "/status")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"sleep": ConfigData.Sleep})
}

// ChangeHandler handles the /change route and updates the sleep status if the key is correct
// ChangeHandler 处理 /change 路由并在 key 正确时更新 sleep 状态
func ChangeHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	// 异步记录访问日志，避免阻塞主处理流程
	go LogAccess(clientIP, "/change")

	key := r.URL.Query().Get("key")
	status := r.URL.Query().Get("status")

	if key != ConfigData.Key {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: false, Result: "Key Wrong"})
		return
	}

	var desiredSleepState bool
	var action string
	if status == "1" {
		desiredSleepState = true
		action = "sleep"
	} else if status == "0" {
		desiredSleepState = false
		action = "wake"
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: false, Result: "Invalid Status"})
		return
	}

	if ConfigData.Sleep == desiredSleepState {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: false, Result: "sleep already " + status})
	} else {
		ConfigData.Sleep = desiredSleepState
		record := SleepRecord{
			Action: action,
			Time:   time.Now().Format(time.RFC3339),
		}
		if err := SaveConfig(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Success: false, Result: "Failed to Save Config"})
			return
		}
		if err := SaveSleepRecord(record); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Success: false, Result: "Failed to Save Sleep Record"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Success: true, Result: "status changed! now sleep is " + status})
	}
}

// HeartbeatHandler handles the /heartbeat route for receiving heartbeat signals
// HeartbeatHandler 处理 /heartbeat 路由接收心跳信号
func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	// 异步记录访问日志，避免阻塞主处理流程
	go LogAccess(clientIP, "/heartbeat")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if !ConfigData.HeartbeatEnabled {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(Response{Success: false, Result: "Heartbeat monitoring is disabled"})
		return
	}

	key := r.URL.Query().Get("key")
	if key != ConfigData.Key {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: false, Result: "Key Wrong"})
		return
	}

	updateLastHeartbeat()
	
	// 检查当前状态，如果是睡眠状态则唤醒
	// 使用原子操作检查，避免不必要的锁争用
	needWakeUp := ConfigData.Sleep
	
	if needWakeUp {
		// 使用单独的函数处理唤醒操作，避免在处理请求时持有锁太久
		go wakeUpFromSleep()
	}

	json.NewEncoder(w).Encode(Response{Success: true, Result: "Heartbeat received"})
}

// wakeUpFromSleep 安全地将状态设置为醒来
func wakeUpFromSleep() {
	mutex.Lock()
	defer mutex.Unlock()
	
	// 再次检查状态，避免在获取锁期间状态被改变
	if !ConfigData.Sleep {
		return
	}
	
	ConfigData.Sleep = false
	record := SleepRecord{
		Action: "wake",
		Time:   time.Now().Format(time.RFC3339),
	}
	
	// 保存配置
	_ = SaveConfig()
	
	// 保存记录
	_ = SaveSleepRecordNoLock(record)
}

// RecordsHandler handles the /records route and returns the latest sleep records
// RecordsHandler 处理 /records 路由并返回最新的睡眠记录
func RecordsHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	// 异步记录访问日志，避免阻塞主处理流程
	go LogAccess(clientIP, "/records")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// 获取限制数量参数，默认为30
	limitStr := r.URL.Query().Get("limit")
	limit := 30
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// 读取睡眠记录
	records, err := LoadSleepRecords(limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Success: false, Result: "Failed to load sleep records"})
		return
	}

	// 返回记录
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RecordsResponse{Success: true, Records: records})
}
