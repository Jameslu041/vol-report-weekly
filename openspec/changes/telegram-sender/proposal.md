# Proposal: telegram-sender

## Motivation

实现 Telegram 发送功能，将生成的报告发送到指定群组。

## Scope

- Telegram Bot API 调用
- 长消息分段（4096 字符限制）
- 本地文件保存

## Non-goals

- 不实现消息编辑/删除
- 不实现媒体发送
