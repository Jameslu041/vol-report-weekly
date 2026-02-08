# Spec: Config

## Overview

配置模块负责从环境变量加载运行时配置。

## Config Structure

```go
type Config struct {
    // Telegram 配置
    TGBotToken string // TG_BOT_TOKEN (required)
    TGChatID   string // TG_CHAT_ID (required)

    // API 配置
    APIBaseURL string // API_BASE_URL (default: https://dev-api.yazhan.vip)

    // 存储配置
    DataDir   string // DATA_DIR (default: ./data)
    ReportDir string // REPORT_DIR (default: ./reports)

    // Cron 配置
    CronSchedule string // CRON_SCHEDULE (default: 0 10 * * 0)
}
```

## Behavior

1. `Load()` 函数读取环境变量并返回 `*Config`
2. 必填字段（TGBotToken, TGChatID）为空时返回错误
3. 可选字段使用默认值

## Validation

- TGBotToken: 非空
- TGChatID: 非空
- APIBaseURL: 非空（有默认值）
