# Proposal: project-scaffold

## Motivation

建立项目基础骨架，为后续功能开发提供稳定的代码结构、配置加载和日志输出能力。

## Scope

- 创建 `cmd/report/main.go` 作为 CLI 入口
- 创建 `internal/config/config.go` 实现环境变量配置加载
- 实现结构化日志输出（使用标准库 log/slog）
- 创建基本目录结构：`cmd/`, `internal/`, `data/`, `reports/`

## Non-goals

- 不实现具体的 API 调用逻辑
- 不实现报告生成逻辑
- 不实现 Telegram 发送逻辑

## Risks

- 无重大风险，这是标准的 Go 项目初始化
