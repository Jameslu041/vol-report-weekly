# Implementation Plan: BTC 每周波动率市场报告

## Background

构建一个 Go 语言命令行工具 `vol-report`，每周自动从 Deribit DataLab API 获取 BTC 期权市场波动率数据，生成中文 Markdown 格式的周报，并通过 Telegram Bot 发送到指定群组。

## Goals

1. 实现完整的 Deribit DataLab API 客户端（5 个接口）
2. 实现数据快照存储与历史对比
3. 实现 Markdown 报告生成器（含模板引擎）
4. 实现 Telegram 发送功能（支持长消息分段）
5. 实现 CLI 入口（run/cron/dry-run 模式）

## Non-goals

- 不实现 Web UI
- 不支持其他交易所（仅 Deribit）
- 不支持其他币种（仅 BTC）
- 不实现实时监控/告警

## Blueprint Alignment

本计划严格遵循 `docs/blueprint.md` 定义的：
- API 接口规范（5 个 POST 接口）
- 报告结构（6 个章节）
- 项目结构（cmd/internal 布局）
- 运行模式（run/cron/dry-run）
- 配置方式（环境变量）

## Repo Reality

- **当前状态**: 全新项目，仅有 `go.mod` 声明模块
- **技术栈**: Go 1.25.1，已声明 `github.com/robfig/cron/v3` 依赖
- **目录结构**: 需要从零创建 `cmd/` 和 `internal/` 目录
- **OpenSpec**: 已初始化 `openspec/` 目录，可用于 Spec-Driven 开发

## Implementation Strategy

分 6 个阶段交付，每阶段可独立验证：

### Phase 1: 项目骨架与配置
- 建立 `cmd/report/main.go` 入口
- 实现配置加载（环境变量）
- 实现结构化日志

### Phase 2: API 客户端
- 实现 Deribit DataLab HTTP 客户端
- 实现 5 个接口的调用与响应解析
- 实现重试机制（3次，2秒间隔）

### Phase 3: 数据存储
- 实现快照存储（JSON 文件）
- 实现历史数据加载与对比

### Phase 4: 报告生成
- 实现 Markdown 模板
- 实现 6 个章节的数据分析逻辑
- 实现变化趋势判断（IV/VRP/Skew/Flow）

### Phase 5: Telegram 发送
- 实现 Bot API 调用
- 实现长消息分段（4096 字符限制）
- 实现本地文件保存

### Phase 6: CLI 完善
- 实现 cron 定时模式
- 实现 --dry-run 参数
- 实现 --date 参数

## Demo 目标与流程

### 价值目标
证明工具能够：完整获取 Deribit 波动率数据 → 生成结构化中文周报 → 成功发送到 Telegram。

### Demo 前置条件
- 环境变量已配置：`TG_BOT_TOKEN`、`TG_CHAT_ID`
- 网络可访问：`dev-api.yazhan.vip`、`api.telegram.org`
- Go 运行环境可用

### 关键环境变量清单
| 变量名 | 必填 | 说明 |
|--------|------|------|
| TG_BOT_TOKEN | 是 | Telegram Bot Token |
| TG_CHAT_ID | 是 | 目标群组 Chat ID |
| API_BASE_URL | 否 | API 地址，默认 https://dev-api.yazhan.vip |
| DATA_DIR | 否 | 数据目录，默认 ./data |
| REPORT_DIR | 否 | 报告目录，默认 ./reports |

### 如何运行
```bash
# 安装依赖
go mod tidy

# 构建
go build -o vol-report ./cmd/report

# 配置环境变量
export TG_BOT_TOKEN="your-bot-token"
export TG_CHAT_ID="your-chat-id"

# 运行（生成并发送）
./vol-report run

# Dry-run 模式（仅生成不发送）
./vol-report run --dry-run
```

### Demo 流程（Step-by-step）
1. **API 连通性验证**: 运行 `./vol-report run --dry-run`，检查日志是否成功获取 5 个接口数据
2. **报告生成验证**: 检查 `reports/YYYY-MM-DD.md` 文件内容，确认 6 个章节完整
3. **Telegram 发送验证**: 运行 `./vol-report run`，在 Telegram 群组中确认收到消息
4. **定时模式验证**: 运行 `./vol-report cron`，确认进入等待状态并输出下次执行时间

### 失败兜底
若 Telegram 发送失败，验收口径降级为：
- 报告文件已成功生成在 `reports/` 目录
- 日志输出报告内容摘要

## Work Breakdown（模块级拆分）

| 模块 | 文件 | 依赖 |
|------|------|------|
| CLI 入口 | cmd/report/main.go | config |
| 配置 | internal/config/config.go | - |
| API 客户端 | internal/api/client.go, types.go | config |
| 存储 | internal/storage/snapshot.go | - |
| 报告生成 | internal/report/generator.go, template.go, analyzer.go | api, storage |
| Telegram | internal/telegram/sender.go | config |

## Backlog Draft

以下为原子化的 backlog items，可直接用于生成 `BACKLOG.md`：

### project-scaffold
**描述**: 建立项目骨架：cmd/report/main.go 入口 + internal/config 配置加载 + 结构化日志
**验收**: `go build ./cmd/report` 成功，运行后输出配置加载日志
**依赖**: 无

### api-client
**描述**: 实现 Deribit DataLab API 客户端：HTTP 客户端 + 5 个接口 + 重试机制
**验收**: 单元测试覆盖 5 个接口，mock server 验证重试逻辑
**依赖**: project-scaffold

### data-storage
**描述**: 实现数据快照存储：JSON 文件读写 + 历史数据加载
**验收**: 单元测试覆盖存储/加载/对比场景
**依赖**: api-client

### report-generator
**描述**: 实现报告生成器：Markdown 模板 + 6 章节分析逻辑 + 变化判断
**验收**: 给定 mock 数据，生成的报告符合 blueprint 定义的格式
**依赖**: data-storage

### telegram-sender
**描述**: 实现 Telegram 发送：Bot API 调用 + 长消息分段 + 本地文件保存
**验收**: 集成测试验证消息发送成功（或 dry-run 模式下跳过发送）
**依赖**: report-generator

### cli-complete
**描述**: 完善 CLI：cron 定时模式 + --dry-run + --date 参数
**验收**: `vol-report cron` 进入定时等待，`--dry-run` 仅生成不发送
**依赖**: telegram-sender

## Risks & TBD

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| API 接口变更 | 数据解析失败 | 响应结构做容错处理，日志记录原始响应 |
| Telegram 发送限频 | 消息发送失败 | 分段发送间隔 1 秒 |
| 无历史数据 | 首次运行无法对比 | 首次运行标注"无历史数据对比" |

## Test Plan

1. **单元测试**: 覆盖 API 响应解析、存储读写、分析判断逻辑
2. **集成测试**: Mock API server，端到端验证报告生成
3. **手动验收**: 真实环境运行，验证 Telegram 消息内容
