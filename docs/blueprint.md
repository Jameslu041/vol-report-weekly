# BTC 每周波动率市场报告 — Blueprint

## 1. 项目目标

构建一个 Go 语言命令行工具，每周自动从 Deribit DataLab API 获取 BTC 期权市场波动率数据，生成中文 Markdown 格式的周报，并通过 Telegram Bot 发送到指定群组。

## 2. 数据源

### 2.1 API 基础信息

- **Base URL**: `https://dev-api.yazhan.vip`
- **鉴权**: Header `User-Agent: greeks_news_bot`
- **通用响应格式**: `{"code": 200, "message": "", "data": <业务数据>}`

### 2.2 需要调用的接口

#### 2.2.1 IV History — 隐含波动率历史

- **POST** `/api/v1/deribit/datalab/iv_history`
- 请求: `{"currency": "BTC", "gap": "7D", "exchange": "deribit"}`
- 用途: 获取各期限（1D/1W/1M/2M/3M/6M/1Y）的 IV 数据及 HV 对照
- 报告用途: 展示当前 IV 期限结构、与上周对比变化

#### 2.2.2 IV/RV — 隐含波动率与实现波动率

- **POST** `/api/v1/deribit/datalab/iv_rv`
- 请求: `{"currency": "BTC", "gap": "1M", "exchange": "deribit"}`
- 用途: 获取各周期的 IV、RV、VRP（波动率风险溢价）
- 报告用途: 分析波动率溢价水平，判断期权定价偏贵/偏便宜

#### 2.2.3 Skew Chart — IV 偏斜曲线

- **POST** `/api/v1/deribit/datalab/skew_chart`
- 请求: `{"currency": "BTC", "gap": "1M", "exchange": "deribit"}`
- 用途: 获取各期限偏斜值（put/call IV 之比）
- 报告用途: 分析市场情绪偏向（>1 看跌偏斜，<1 看涨偏斜）

#### 2.2.4 FIV Matrix — 远期隐含波动率矩阵

- **POST** `/api/v1/deribit/datalab/fiv_matrix`
- 请求: `{"currency": "BTC", "page_num": 0, "exchange": "deribit"}`
- 用途: 获取远期 IV 期限结构矩阵
- 报告用途: 展示市场对未来不同时段波动率的预期

#### 2.2.5 Option Flows — 期权流向

- **POST** `/api/v1/deribit/datalab/option_flows`
- 请求: `{"currency": "BTC", "exchange": "deribit"}`
- 用途: 获取 24h 内 call/put 买卖量、大宗交易量
- 报告用途: 分析市场资金流向和机构动向

## 3. 报告输出

### 3.1 格式

Markdown 格式，适合 Telegram 消息发送（使用 MarkdownV2 或 HTML parse mode）。

### 3.2 报告结构

报告包含以下章节:
- 一、IV 期限结构概览: 各期限 IV 当前值，与上周变化
- 二、IV vs RV 分析: 各周期 IV/RV/VRP，VRP 正负及含义解读
- 三、波动率偏斜（Skew）: 各期限 skew 值，市场情绪判断
- 四、远期波动率矩阵（Forward IV）: 关键期限组合的 Forward IV，期限结构形态
- 五、期权流向: 24h call/put 成交量对比，大宗交易方向，资金流向解读
- 六、总结与交易观点: 综合以上数据的市场观点

### 3.3 数据持久化

使用本地 JSON 文件存储每次运行的数据快照（data/snapshots/YYYY-MM-DD.json），用于与上周数据对比生成变化趋势和历史数据回溯。

## 4. Telegram 发送

- 使用 Telegram Bot API 发送消息
- Bot Token 和 Chat ID 通过环境变量配置: TG_BOT_TOKEN, TG_CHAT_ID
- 超过 4096 字符自动分段发送
- 同时保存为本地文件: reports/YYYY-MM-DD.md

## 5. 技术设计

### 5.1 技术栈

- 语言: Go
- HTTP 客户端: net/http 标准库
- JSON 处理: encoding/json 标准库
- Telegram: 直接调用 Bot API
- 模板引擎: text/template 标准库

### 5.2 项目结构

- cmd/report/main.go — 入口
- internal/api/deribit.go — Deribit DataLab API 客户端
- internal/report/generator.go — 报告生成器
- internal/report/template.go — Markdown 模板
- internal/telegram/sender.go — Telegram 发送
- internal/storage/snapshot.go — 数据快照存储
- data/snapshots/ — 历史数据快照
- reports/ — 生成的报告文件

### 5.3 运行模式

- vol-report run — 立即运行一次
- vol-report cron — 定时模式（每周日 10:00 UTC+8）
- vol-report run --dry-run — 仅生成不发送
- vol-report run --date 2026-02-01 — 指定日期

### 5.4 配置（环境变量）

- TG_BOT_TOKEN: Telegram Bot Token（必填）
- TG_CHAT_ID: Telegram 群组 Chat ID（必填）
- API_BASE_URL: API 基础地址（默认 https://dev-api.yazhan.vip）
- CRON_SCHEDULE: Cron 表达式（默认每周日10:00）
- DATA_DIR: 数据目录（默认 ./data）
- REPORT_DIR: 报告输出目录（默认 ./reports）

## 6. 报告生成规则

### 6.1 IV 变化判断
- 变化 > 5%: 大幅上升/下降
- 变化 2-5%: 小幅上升/下降
- 变化 < 2%: 基本持平

### 6.2 VRP 判断
- VRP > 5: 期权溢价偏高，卖方有利
- VRP 0-5: 期权定价合理
- VRP < 0: 期权折价，买方有利

### 6.3 Skew 判断
- skew > 1.05: 明显看跌偏斜，市场偏谨慎
- skew 0.95-1.05: 偏斜中性
- skew < 0.95: 看涨偏斜，市场偏乐观

### 6.4 流向判断
- call_buys / put_buys > 1.5: 看涨情绪浓厚
- put_buys / call_buys > 1.5: 看跌情绪浓厚
- 否则: 多空博弈均衡

## 7. 非功能需求

- API 调用失败时重试 3 次，间隔 2 秒
- 所有 API 调用超时 30 秒
- 日志输出到 stdout，结构化日志
- 报告生成时间 < 30 秒
