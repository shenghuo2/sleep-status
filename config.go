package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	CurrentConfigVersion = 2 // Increment this when adding new versions
)

var (
	configFilePath      = "./config.json"
	sleepRecordFilePath = "./sleep_record.json"
	ConfigData          Config
	mutex               = &sync.Mutex{}
)

// BaseConfig contains common fields across all versions
type BaseConfig struct {
	Version int `json:"version"`
}

// Config represents the current version of configuration
type Config struct {
	BaseConfig
	Sleep            bool   `json:"sleep"`
	Key              string `json:"key"`
	HeartbeatEnabled bool   `json:"heartbeat_enabled"`
	HeartbeatTimeout int    `json:"heartbeat_timeout"`
}

// ConfigV1 represents version 1 of configuration
type ConfigV1 struct {
	Sleep bool   `json:"sleep"`
	Key   string `json:"key"`
}

// migrationFunc defines the signature for version migration functions
type migrationFunc func([]byte) (Config, error)

// migrationMap stores all version upgrade paths
var migrationMap = map[int]migrationFunc{
	0: migrateFromV0ToLatest, // For configs with no version field
	1: migrateFromV1ToLatest,
}

// migrateFromV0ToLatest handles migration from the original version (no version field)
func migrateFromV0ToLatest(data []byte) (Config, error) {
	var oldConfig ConfigV1
	if err := json.Unmarshal(data, &oldConfig); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal v0 config: %v", err)
	}

	return Config{
		BaseConfig: BaseConfig{
			Version: CurrentConfigVersion,
		},
		Sleep:            oldConfig.Sleep,
		Key:              oldConfig.Key,
		HeartbeatEnabled: false,
		HeartbeatTimeout: 60,
	}, nil
}

// migrateFromV1ToLatest handles migration from version 1
func migrateFromV1ToLatest(data []byte) (Config, error) {
	var v1Config ConfigV1
	if err := json.Unmarshal(data, &v1Config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal v1 config: %v", err)
	}

	return Config{
		BaseConfig: BaseConfig{
			Version: CurrentConfigVersion,
		},
		Sleep:            v1Config.Sleep,
		Key:              v1Config.Key,
		HeartbeatEnabled: false,
		HeartbeatTimeout: 60,
	}, nil
}

// LoadConfig loads the configuration from the config file
func LoadConfig() error {
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Create new config with default values
		randomKey, err := generateRandomKey(16)
		if err != nil {
			return err
		}

		ConfigData = Config{
			BaseConfig: BaseConfig{
				Version: CurrentConfigVersion,
			},
			Sleep:            false,
			Key:              randomKey,
			HeartbeatEnabled: false,
			HeartbeatTimeout: 60,
		}

		if err := SaveConfig(); err != nil {
			return err
		}
		fmt.Println("Config file created with default values.")
		fmt.Println("配置文件已创建默认值。")
		return nil
	}

	// Read existing config file
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	// Try to determine version
	var baseConfig BaseConfig
	if err := json.Unmarshal(content, &baseConfig); err != nil {
		baseConfig.Version = 0 // Assume version 0 if no version field
	}

	// Check if migration is needed
	if baseConfig.Version < CurrentConfigVersion {
		fmt.Printf("Upgrading config from version %d to %d...\n", baseConfig.Version, CurrentConfigVersion)
		fmt.Printf("正在将配置从版本 %d 升级到 %d...\n", baseConfig.Version, CurrentConfigVersion)

		// Get migration function
		migrationFunc, exists := migrationMap[baseConfig.Version]
		if !exists {
			return fmt.Errorf("no migration path from version %d", baseConfig.Version)
		}

		// Perform migration
		newConfig, err := migrationFunc(content)
		if err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}

		ConfigData = newConfig
		if err := SaveConfig(); err != nil {
			return fmt.Errorf("failed to save migrated config: %v", err)
		}

		fmt.Println("Config upgraded successfully!")
		fmt.Println("配置升级成功！")
		return nil
	}

	// Current version, just load it
	decoder := json.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(&ConfigData); err != nil {
		return fmt.Errorf("failed to decode config: %v", err)
	}

	return nil
}

// SaveConfig saves the current configuration to the config file
func SaveConfig() error {
	file, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(&ConfigData)
}

// SaveSleepRecord saves the sleep record to the sleep_record.json file
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
func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
