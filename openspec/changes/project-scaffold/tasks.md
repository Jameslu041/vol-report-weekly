# Tasks: project-scaffold

## Implementation Tasks

- [x] 创建目录结构：`cmd/report/`, `internal/config/`, `data/snapshots/`, `reports/`
- [x] 实现 `internal/config/config.go`：定义 Config 结构体，从环境变量加载配置
- [x] 实现 `cmd/report/main.go`：CLI 入口，加载配置并输出日志
- [x] 添加 `.gitkeep` 到空目录
- [x] 验证 `go build ./cmd/report` 成功
