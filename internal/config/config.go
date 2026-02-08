package config

import (
	"errors"
	"os"
)

type Config struct {
	TGBotToken   string
	TGChatID     string
	APIBaseURL   string
	DataDir      string
	ReportDir    string
	CronSchedule string
}

func Load() (*Config, error) {
	cfg := &Config{
		TGBotToken:   os.Getenv("TG_BOT_TOKEN"),
		TGChatID:     os.Getenv("TG_CHAT_ID"),
		APIBaseURL:   getEnvOrDefault("API_BASE_URL", "https://dev-api.yazhan.vip"),
		DataDir:      getEnvOrDefault("DATA_DIR", "./data"),
		ReportDir:    getEnvOrDefault("REPORT_DIR", "./reports"),
		CronSchedule: getEnvOrDefault("CRON_SCHEDULE", "0 10 * * 0"),
	}

	if cfg.TGBotToken == "" {
		return nil, errors.New("TG_BOT_TOKEN is required")
	}
	if cfg.TGChatID == "" {
		return nil, errors.New("TG_CHAT_ID is required")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
