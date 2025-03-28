package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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

// SleepPeriod 结构体表示一个完整的睡眠周期（入睡和醒来）
type SleepPeriod struct {
	SleepTime    int64  `json:"sleep_time"`               // Unix 时间戳（秒）
	WakeTime     int64  `json:"wake_time"`                // Unix 时间戳（秒）
	SleepTimeStr string `json:"sleep_time_str,omitempty"` // 原始时间字符串，可选
	WakeTimeStr  string `json:"wake_time_str,omitempty"`  // 原始时间字符串，可选
	Duration     int    `json:"duration_minutes"`         // 睡眠时长（分钟）
	IsShort      bool   `json:"is_short,omitempty"`       // 标记是否为短时间睡眠（小于3小时）
}

// SleepStats 结构体保存睡眠统计信息
type SleepStats struct {
	AvgSleepTime string        `json:"avg_sleep_time"`
	AvgWakeTime  string        `json:"avg_wake_time"`
	AvgDuration  int           `json:"avg_duration_minutes"`
	Periods      []SleepPeriod `json:"periods"`
	ActualDays   int           `json:"-"` // 实际包含的天数，不输出到JSON
}

// StatsResponse 结构体保存睡眠统计响应
type StatsResponse struct {
	Success        bool       `json:"success"`
	Stats          SleepStats `json:"stats"`
	Days           int        `json:"days"`                   // 实际统计的天数
	RequestDays    int        `json:"request_days,omitempty"` // 请求的天数，如果与实际天数不同才显示
	Sleep          *bool      `json:"sleep,omitempty"`       // 当前睡眠状态，可选
	CurrentSleepAt int64      `json:"current_sleep_at,omitempty"` // 当前睡眠开始时间（Unix 时间戳），仅当 sleep 为 true 时显示
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

// StatsHandler 处理 /sleep-stats 路由并返回睡眠统计数据
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	// 异步记录访问日志，避免阻塞主处理流程
	go LogAccess(clientIP, "/sleep-stats")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// 获取天数参数，默认为7天
	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	// 获取是否显示原始时间字符串的参数，默认不显示
	showTimeStr := false
	showTimeStrParam := r.URL.Query().Get("show_time_str")
	if showTimeStrParam == "1" || showTimeStrParam == "true" {
		showTimeStr = true
	}
	
	// 获取是否显示当前睡眠状态的参数，默认不显示
	showSleep := false
	showSleepParam := r.URL.Query().Get("show_sleep")
	if showSleepParam == "1" || showSleepParam == "true" {
		showSleep = true
	}

	// 读取足够多的睡眠记录以覆盖所需天数
	// 由于我们不知道确切需要多少记录才能覆盖指定的天数，先获取较大数量的记录
	records, err := LoadSleepRecords(1000) // 获取足够多的记录
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Success: false, Result: "Failed to load sleep records"})
		return
	}

	// 计算睡眠统计数据
	stats, err := calculateSleepStats(records, days, showTimeStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Success: false, Result: fmt.Sprintf("Failed to calculate sleep stats: %v", err)})
		return
	}

	// 返回统计数据
	w.WriteHeader(http.StatusOK)

	// 构建响应
	response := StatsResponse{
		Success: true,
		Stats:   stats,
		Days:    stats.ActualDays,
	}

	// 如果实际天数与请求天数不同，则显示请求天数
	if stats.ActualDays != days {
		response.RequestDays = days
	}
	
	// 如果需要显示当前睡眠状态
	if showSleep {
		// 直接使用配置中的状态
		sleepStatus := ConfigData.Sleep
		response.Sleep = &sleepStatus
	}
	
	// current_sleep_at 字段与 show_sleep 参数无关，只由 ConfigData.Sleep 控制
	if ConfigData.Sleep {
		// 获取最新的一条入睡记录
		latestSleepRecord, err := findLatestSleepRecord()
		if err == nil && latestSleepRecord.Action == "sleep" {
			// 将时间字符串转换为时间对象
			sleepTimeObj, err := convertToUTC8(latestSleepRecord.Time)
			if err == nil {
				// 转换为 Unix 时间戳
				response.CurrentSleepAt = sleepTimeObj.Unix()
			}
		}
	}

	json.NewEncoder(w).Encode(response)
}

// calculateSleepStats 计算睡眠统计数据
func calculateSleepStats(records []SleepRecord, requestedDays int, showTimeStr bool) (SleepStats, error) {
	// 将记录按时间排序（从旧到新）
	sort.Slice(records, func(i, j int) bool {
		return records[i].Time < records[j].Time
	})

	// 找出成对的睡眠-醒来记录
	var periods []SleepPeriod
	var sleepTime string

	// 当前时间，用于计算天数限制
	now := time.Now()
	cutoffTime := now.AddDate(0, 0, -requestedDays)

	for i := 0; i < len(records); i++ {
		record := records[i]

		// 将时间字符串转换为时间对象
		recordTime, err := convertToUTC8(record.Time)
		if err != nil {
			continue // 跳过无效的时间记录
		}

		if record.Action == "sleep" {
			sleepTime = record.Time
		} else if record.Action == "wake" && sleepTime != "" {
			// 找到一对睡眠-醒来记录
			sleepTimeObj, err := convertToUTC8(sleepTime)
			if err != nil {
				sleepTime = "" // 重置睡眠时间
				continue
			}

			// 只包含指定天数内的记录
			if sleepTimeObj.Before(cutoffTime) {
				sleepTime = "" // 重置睡眠时间
				continue
			}

			// 计算睡眠时长（分钟）
			duration := int(recordTime.Sub(sleepTimeObj).Minutes())

			// 将所有睡眠时长大于10分钟且小于24小时的睡眠段都记录下来
			if duration >= 10 && duration <= 1440 {
				// 标记短时间睡眠（小于3小时）
				isShort := duration < 180

				// 创建睡眠周期对象
				period := SleepPeriod{
					SleepTime: sleepTimeObj.Unix(),
					WakeTime:  recordTime.Unix(),
					Duration:  duration,
					IsShort:   isShort,
				}

				// 如果需要显示原始时间字符串，则添加相应字段
				if showTimeStr {
					period.SleepTimeStr = sleepTime
					period.WakeTimeStr = record.Time
				}

				periods = append(periods, period)
			}

			sleepTime = "" // 重置睡眠时间，准备下一对
		}
	}

	// 如果没有有效的睡眠周期，返回空结果
	if len(periods) == 0 {
		return SleepStats{
			AvgSleepTime: "",
			AvgWakeTime:  "",
			AvgDuration:  0,
			Periods:      []SleepPeriod{},
		}, nil
	}

	// 计算平均睡眠时间、醒来时间和睡眠时长
	var totalSleepMinutes, totalWakeMinutes, totalDuration int
	var normalPeriodCount int // 正常睡眠周期计数（大于3小时）

	for _, period := range periods {
		// 使用 Unix 时间戳创建时间对象
		sleepTimeObj := time.Unix(period.SleepTime, 0)
		wakeTimeObj := time.Unix(period.WakeTime, 0)

		// 将时间转换为 UTC+8 时区
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			// 如果无法加载时区，手动创建 UTC+8 时区
			loc = time.FixedZone("UTC+8", 8*60*60)
		}

		sleepTimeObj = sleepTimeObj.In(loc)
		wakeTimeObj = wakeTimeObj.In(loc)

		// 计算睡眠时间的分钟数（从当天0点开始）
		sleepMinutes := sleepTimeObj.Hour()*60 + sleepTimeObj.Minute()

		// 处理跨天的情况
		if sleepTimeObj.Hour() < 12 {
			sleepMinutes += 24 * 60 // 将凌晨时间视为前一天的延续
		}

		// 计算醒来时间的分钟数
		wakeMinutes := wakeTimeObj.Hour()*60 + wakeTimeObj.Minute()
		if wakeTimeObj.Hour() < 12 {
			wakeMinutes += 24 * 60 // 将凌晨时间视为前一天的延续
		}

		// 对于所有睡眠周期，都计入总睡眠时长
		totalDuration += period.Duration

		// 只有非短时间睡眠（大于3小时）才计入平均入睡和醒来时间
		if !period.IsShort {
			totalSleepMinutes += sleepMinutes
			totalWakeMinutes += wakeMinutes
			normalPeriodCount++
		}
	}

	// 计算平均值
	var avgSleepMinutes, avgWakeMinutes int

	// 如果有正常睡眠周期，计算平均入睡和醒来时间
	if normalPeriodCount > 0 {
		avgSleepMinutes = totalSleepMinutes / normalPeriodCount
		avgWakeMinutes = totalWakeMinutes / normalPeriodCount
	}

	// 对于所有睡眠周期（包括短时间睡眠），计算平均睡眠时长
	avgDuration := totalDuration / len(periods)

	// 将分钟数转换回小时:分钟格式
	avgSleepHour := (avgSleepMinutes % (24 * 60)) / 60
	avgSleepMinute := avgSleepMinutes % 60

	avgWakeHour := (avgWakeMinutes % (24 * 60)) / 60
	avgWakeMinute := avgWakeMinutes % 60

	// 格式化时间
	avgSleepTimeStr := fmt.Sprintf("%02d:%02d", avgSleepHour, avgSleepMinute)
	avgWakeTimeStr := fmt.Sprintf("%02d:%02d", avgWakeHour, avgWakeMinute)

	// 计算实际的天数范围
	var actualDays int
	if len(periods) > 0 {
		// 找出最早和最晚的睡眠记录时间
		earliest := time.Unix(periods[0].SleepTime, 0)
		latest := time.Unix(periods[len(periods)-1].WakeTime, 0)

		// 计算实际跨越的天数
		duration := latest.Sub(earliest)
		actualDays = int(duration.Hours()/24) + 1 // 加 1 是因为包含当天

		// 如果实际天数为 0（同一天内的记录），设置为 1
		if actualDays < 1 {
			actualDays = 1
		}
	} else {
		actualDays = 0
	}

	return SleepStats{
		AvgSleepTime: avgSleepTimeStr,
		AvgWakeTime:  avgWakeTimeStr,
		AvgDuration:  avgDuration,
		Periods:      periods,
		ActualDays:   actualDays,
	}, nil
}

// convertToUTC8 将时间字符串转换为UTC+8时区的时间对象
func convertToUTC8(timeStr string) (time.Time, error) {
	// 解析时间字符串
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, err
	}

	// 将时间转换为UTC+8时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 如果无法加载时区，手动创建UTC+8时区
		loc = time.FixedZone("UTC+8", 8*60*60)
	}

	return t.In(loc), nil
}



// findLatestSleepRecord 查找最近的一条入睡记录
func findLatestSleepRecord() (SleepRecord, error) {
	// 读取所有睡眠记录
	records, err := LoadSleepRecords(100) // 只获取最近的100条记录就足够了
	if err != nil {
		return SleepRecord{}, err
	}
	
	// 如果没有记录，返回错误
	if len(records) == 0 {
		return SleepRecord{}, fmt.Errorf("no sleep records found")
	}
	
	// 将记录按时间排序（从新到旧）
	sort.Slice(records, func(i, j int) bool {
		return records[i].Time > records[j].Time
	})
	
	// 找出最新的一条入睡记录
	for _, record := range records {
		if record.Action == "sleep" {
			return record, nil
		}
	}
	
	// 如果没有找到入睡记录，返回错误
	return SleepRecord{}, fmt.Errorf("no sleep record found")
}
