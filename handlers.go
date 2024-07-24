package main

import (
	"encoding/json"
	"net/http"
)

// Response struct to standardize API responses
// Response 结构体标准化 API 响应
type Response struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
}

// StatusHandler handles the /status route and returns the current sleep status
// StatusHandler 处理 /status 路由并返回当前的 sleep 状态
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	LogAccess(clientIP)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"sleep": ConfigData.Sleep})
}

// ChangeHandler handles the /change route and updates the sleep status if the key is correct
// ChangeHandler 处理 /change 路由并在 key 正确时更新 sleep 状态
func ChangeHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	LogAccess(clientIP)

	w.Header().Set("Content-Type", "application/json")
	key := r.URL.Query().Get("key")
	status := r.URL.Query().Get("status")

	if key != ConfigData.Key {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: false, Result: "Key Wrong"})
		return
	}

	var desiredSleepState bool
	if status == "1" {
		desiredSleepState = true
	} else if status == "0" {
		desiredSleepState = false
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
		if err := SaveConfig(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Success: false, Result: "Failed to Save Config"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Success: true, Result: "status changed! now sleep is " + status})
	}
}
