# Backlog

## vol-report-weekly

> BTC 每周波动率市场报告 CLI 工具

**Epic branch**: `epic/vol-report-weekly`
**Epic branch URL**: https://github.com/Jameslu041/vol-report-weekly/tree/epic/vol-report-weekly

### project-scaffold

> Issue: #1

建立项目骨架：cmd/report/main.go 入口 + internal/config 配置加载 + 结构化日志

- **验收**: `go build ./cmd/report` 成功，运行后输出配置加载日志
- **依赖**: 无

### api-client

> Issue: #2

实现 Deribit DataLab API 客户端：HTTP 客户端 + 5 个接口 + 重试机制

- **验收**: 单元测试覆盖 5 个接口，mock server 验证重试逻辑
- **依赖**: project-scaffold

### data-storage

> Issue: #3

实现数据快照存储：JSON 文件读写 + 历史数据加载

- **验收**: 单元测试覆盖存储/加载/对比场景
- **依赖**: api-client

### report-generator

> Issue: #4

实现报告生成器：Markdown 模板 + 6 章节分析逻辑 + 变化判断

- **验收**: 给定 mock 数据，生成的报告符合 blueprint 定义的格式
- **依赖**: data-storage

### telegram-sender

> Issue: #5

实现 Telegram 发送：Bot API 调用 + 长消息分段 + 本地文件保存

- **验收**: 集成测试验证消息发送成功（或 dry-run 模式下跳过发送）
- **依赖**: report-generator

### cli-complete

> Issue: #6

完善 CLI：cron 定时模式 + --dry-run + --date 参数

- **验收**: `vol-report cron` 进入定时等待，`--dry-run` 仅生成不发送
- **依赖**: telegram-sender
