package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var configFilePath = "./config.json"
var sleepRecordFilePath = "./sleep_record.json"
var ConfigData Config
var mutex = &sync.Mutex{}

// Config struct to hold the config file data
// Config 结构体保存配置文件数据
type Config struct {
	Sleep bool   `json:"sleep"`
	Key   string `json:"key"`
}

// LoadConfig loads the configuration from the config file
// LoadConfig 从配置文件加载配置
func LoadConfig() error {
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the config file exists
	// 检查配置文件是否存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// If the config file does not exist, create it with default values
		// 如果配置文件不存在，则创建它并写入默认值
		randomKey, err := generateRandomKey(16)
		if err != nil {
			return err
		}

		ConfigData = Config{
			Sleep: false,
			Key:   randomKey,
		}
		// Save the default config
		// 保存默认配置
		if err := SaveConfig(); err != nil {
			return err
		}
		fmt.Println("Config file created with default values. You can change the key in the config file.")
		fmt.Println("配置文件已创建默认值。可以通过修改config文件修改密码。")
		return nil
	}

	// If the config file exists, load it
	// 如果配置文件存在，则加载它
	file, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ConfigData)
	if err != nil {
		return err
	}

	return nil
}

// SaveConfig saves the current configuration to the config file
// SaveConfig 将当前配置保存到配置文件
func SaveConfig() error {
	file, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&ConfigData)
	if err != nil {
		return err
	}

	return nil
}

// SaveSleepRecord saves the sleep record to the sleep_record.json file
// SaveSleepRecord 将睡眠记录保存到 sleep_record.json 文件
func SaveSleepRecord(record SleepRecord) error {
	file, err := os.OpenFile(sleepRecordFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(record); err != nil {
		return err
	}

	return nil
}

// generateRandomKey generates a random key of the specified length
// generateRandomKey 生成指定长度的随机密钥
func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
