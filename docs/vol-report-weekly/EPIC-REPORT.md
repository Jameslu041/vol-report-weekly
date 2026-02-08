# Epic Report: vol-report-weekly

## 概述

构建一个 Go 语言命令行工具 `vol-report`，每周自动从 Deribit DataLab API 获取 BTC 期权市场波动率数据，生成中文 Markdown 格式的周报，并通过 Telegram Bot 发送到指定群组。

## 交付清单

| Item | Issue | PR | Status |
|------|-------|-----|--------|
| project-scaffold | #1 | #7 | ✅ Done |
| api-client | #2 | #9 | ✅ Done |
| data-storage | #3 | #11 | ✅ Done |
| report-generator | #4 | #13 | ✅ Done |
| telegram-sender | #5 | #15 | ✅ Done |
| cli-complete | #6 | #17 | ✅ Done |

## 技术架构

```
cmd/report/main.go          # CLI 入口
internal/
├── api/                    # Deribit DataLab API 客户端
│   ├── client.go          # HTTP 客户端 + 重试机制
│   └── types.go           # 响应数据结构
├── config/                 # 配置加载（环境变量）
├── report/                 # 报告生成器
│   ├── analyzer.go        # 数据分析逻辑
│   ├── generator.go       # Markdown 生成
│   └── template.go        # 报告模板
├── storage/               # 数据快照存储
│   └── snapshot.go        # JSON 文件读写
└── telegram/              # Telegram 发送
    └── sender.go          # Bot API + 消息分段
```

## 关键决策

1. **API 响应解析**: 原始 API 返回复杂嵌套结构，需要自定义解析逻辑转换为简化格式
2. **消息分段**: Telegram 4096 字符限制，按 Markdown 章节分段
3. **Cron 调度**: 使用 robfig/cron 库实现定时执行

## 测试覆盖

- API 客户端：7 个单元测试（含重试机制验证）
- 数据存储：4 个单元测试
- 报告生成：6 个单元测试
- Telegram 发送：5 个单元测试

## 运行指南

```bash
# 环境变量
export TG_BOT_TOKEN="your-bot-token"
export TG_CHAT_ID="your-chat-id"

# 构建
go build -o vol-report ./cmd/report

# 立即运行
./vol-report run

# Dry-run（仅生成不发送）
./vol-report run --dry-run

# 定时模式
./vol-report cron
```

## 链接

- Epic Branch: https://github.com/Jameslu041/vol-report-weekly/tree/epic/vol-report-weekly
- Blueprint: docs/blueprint.md
- Implementation Plan: docs/vol-report-weekly/Implementation Plan.md
- OpenSpec Specs: openspec/specs/

## Demo 验证

已通过 `--dry-run` 模式验证：
- ✅ API 连通性（5 个接口全部成功）
- ✅ 报告生成（6 个章节完整）
- ✅ 文件保存（reports/YYYY-MM-DD.md）
