# Proposal: data-storage

## Motivation

实现数据快照存储，用于保存每次 API 调用结果并支持与历史数据对比。

## Scope

- JSON 文件存储（data/snapshots/YYYY-MM-DD.json）
- 加载历史快照
- 提供上周数据对比

## Non-goals

- 不实现数据库存储
- 不实现数据压缩

## Risks

- 文件系统权限问题
