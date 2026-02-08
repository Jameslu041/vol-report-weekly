# Proposal: api-client

## Motivation

实现 Deribit DataLab API 客户端，用于获取 BTC 期权波动率市场数据。

## Scope

- 创建 HTTP 客户端，支持 POST 请求
- 实现 5 个 API 接口调用
- 实现重试机制（3 次，间隔 2 秒）
- 定义响应数据结构

## Non-goals

- 不实现数据持久化（由 data-storage 负责）
- 不实现报告生成逻辑
- 不实现其他交易所支持

## Risks

- API 响应格式变化可能导致解析失败
- 网络超时需要合理处理
