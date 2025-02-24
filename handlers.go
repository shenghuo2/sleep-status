package main

import (
	"encoding/json"
	"net/http"
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

// StatusHandler handles the /status route and returns the current sleep status
// StatusHandler 处理 /status 路由并返回当前的 sleep 状态
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	LogAccess(clientIP, "/status")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"sleep": ConfigData.Sleep})
}

// ChangeHandler handles the /change route and updates the sleep status if the key is correct
// ChangeHandler 处理 /change 路由并在 key 正确时更新 sleep 状态
func ChangeHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	LogAccess(clientIP, "/change")

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
	LogAccess(clientIP, "/heartbeat")

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
	
	// If currently sleeping, wake up
	if ConfigData.Sleep {
		ConfigData.Sleep = false
		record := SleepRecord{
			Action: "wake",
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
	}

	json.NewEncoder(w).Encode(Response{Success: true, Result: "Heartbeat received"})
}
