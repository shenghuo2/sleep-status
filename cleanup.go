package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// runCleanup 执行清理操作，删除无效的睡眠记录
func runCleanup() {
	// 加载配置
	if err := LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	// 确定 cleanup 使用的最小间隔
	cleanupDuration := ConfigData.CleanupMinDuration
	if cleanupDuration == 0 {
		cleanupDuration = ConfigData.MinSleepDuration
	}

	fmt.Printf("Starting cleanup with CleanupMinDuration = %d minutes\n", cleanupDuration)
	fmt.Printf("开始清理，清理最小间隔 = %d 分钟\n", cleanupDuration)

	// 读取所有睡眠记录
	records, err := loadAllSleepRecords()
	if err != nil {
		fmt.Printf("Error loading sleep records: %v\n", err)
		fmt.Printf("加载睡眠记录失败: %v\n", err)
		return
	}

	if len(records) == 0 {
		fmt.Println("No sleep records found.")
		fmt.Println("没有找到睡眠记录。")
		return
	}

	fmt.Printf("Found %d records before cleanup\n", len(records))
	fmt.Printf("清理前共有 %d 条记录\n", len(records))

	// 执行清理
	cleanedRecords, removedCount := cleanupRecords(records, cleanupDuration)

	if removedCount == 0 {
		fmt.Println("No invalid records found. Nothing to clean.")
		fmt.Println("没有发现无效记录，无需清理。")
		return
	}

	fmt.Printf("Removed %d invalid records\n", removedCount)
	fmt.Printf("删除了 %d 条无效记录\n", removedCount)

	// 保存清理后的记录
	if err := saveAllSleepRecords(cleanedRecords); err != nil {
		fmt.Printf("Error saving cleaned records: %v\n", err)
		fmt.Printf("保存清理后的记录失败: %v\n", err)
		return
	}

	fmt.Printf("Cleanup completed. %d records remaining.\n", len(cleanedRecords))
	fmt.Printf("清理完成，剩余 %d 条记录。\n", len(cleanedRecords))
}

// loadAllSleepRecords 读取所有睡眠记录（不限制数量）
func loadAllSleepRecords() ([]SleepRecord, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// 检查文件是否存在
	if _, err := os.Stat(sleepRecordFilePath); os.IsNotExist(err) {
		return []SleepRecord{}, nil
	}

	// 读取文件内容
	data, err := os.ReadFile(sleepRecordFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sleep record file: %v", err)
	}

	var records []SleepRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return []SleepRecord{}, nil
	}

	return records, nil
}

// saveAllSleepRecords 保存所有睡眠记录
func saveAllSleepRecords(records []SleepRecord) error {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sleep records: %v", err)
	}

	if err := os.WriteFile(sleepRecordFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write sleep records: %v", err)
	}

	return nil
}

// cleanupRecords 清理无效的睡眠记录
// 智能合并逻辑：
// 1. 孤立的有效睡眠段（如午休）：睡眠时长 >= minDuration，且前后间隔都 >= minDuration，保留
// 2. 连续的短周期（网络波动）：合并为从第一次 sleep 到最后一次 wake
func cleanupRecords(records []SleepRecord, minDuration int) ([]SleepRecord, int) {
	if len(records) == 0 {
		return records, 0
	}

	// 按时间排序（从旧到新）
	sort.Slice(records, func(i, j int) bool {
		return records[i].Time < records[j].Time
	})

	// 解析所有记录的时间
	type parsedRecord struct {
		record SleepRecord
		time   time.Time
	}
	var parsed []parsedRecord
	for _, r := range records {
		t, err := parseRecordTime(r.Time)
		if err != nil {
			continue
		}
		parsed = append(parsed, parsedRecord{record: r, time: t})
	}

	if len(parsed) == 0 {
		return records, 0
	}

	// 找出需要合并的连续短周期区间
	// 区间定义：从一个 sleep 开始，到下一个间隔 >= minDuration 的 wake 结束
	var result []SleepRecord
	originalCount := len(parsed)

	i := 0
	for i < len(parsed) {
		current := parsed[i]

		// 如果是 wake 记录，检查是否需要保留
		if current.record.Action == "wake" {
			// 检查与下一条记录的间隔
			if i+1 < len(parsed) {
				gap := int(parsed[i+1].time.Sub(current.time).Minutes())
				if gap >= minDuration {
					// 间隔足够大，保留这个 wake
					result = append(result, current.record)
				} else {
					fmt.Printf("  Removing isolated wake: %s (gap to next: %d min)\n",
						current.record.Time, gap)
					fmt.Printf("  删除孤立的 wake: %s (到下一条间隔: %d 分钟)\n",
						current.record.Time, gap)
				}
			} else {
				// 最后一条记录，保留
				result = append(result, current.record)
			}
			i++
			continue
		}

		// 当前是 sleep 记录，开始寻找这个睡眠周期的结束
		sleepStart := current
		sleepStartIdx := i

		// 找到这个睡眠周期的结束（最后一个 wake，之后有大间隔或没有更多记录）
		j := i + 1
		lastWakeIdx := -1
		lastWake := parsedRecord{}

		for j < len(parsed) {
			rec := parsed[j]

			if rec.record.Action == "wake" {
				lastWakeIdx = j
				lastWake = rec

				// 检查这个 wake 之后的间隔
				if j+1 < len(parsed) {
					gapAfterWake := int(parsed[j+1].time.Sub(rec.time).Minutes())
					if gapAfterWake >= minDuration {
						// 找到了大间隔，这是睡眠周期的结束
						break
					}
				} else {
					// 没有更多记录，这是最后的 wake
					break
				}
			}
			j++
		}

		// 如果没有找到对应的 wake，只保留 sleep
		if lastWakeIdx == -1 {
			result = append(result, sleepStart.record)
			i++
			continue
		}

		// 计算整个睡眠周期的时长
		totalDuration := int(lastWake.time.Sub(sleepStart.time).Minutes())

		// 检查是否有中间记录被跳过（需要合并）
		recordsInBetween := lastWakeIdx - sleepStartIdx - 1

		if recordsInBetween > 0 {
			// 有中间记录，需要合并
			fmt.Printf("  Merging sleep cycle: %s -> %s (%d min, %d records merged)\n",
				sleepStart.record.Time, lastWake.record.Time, totalDuration, recordsInBetween)
			fmt.Printf("  合并睡眠周期: %s -> %s (%d 分钟, 合并了 %d 条记录)\n",
				sleepStart.record.Time, lastWake.record.Time, totalDuration, recordsInBetween)
		}

		// 只保留第一个 sleep 和最后一个 wake
		result = append(result, sleepStart.record)
		result = append(result, lastWake.record)

		// 跳到 wake 之后继续处理
		i = lastWakeIdx + 1
	}

	removedCount := originalCount - len(result)
	return result, removedCount
}

// parseRecordTime 解析记录中的时间字符串
func parseRecordTime(timeStr string) (time.Time, error) {
	// 尝试 RFC3339 格式
	t, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return t, nil
	}

	// 尝试其他常见格式
	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}
