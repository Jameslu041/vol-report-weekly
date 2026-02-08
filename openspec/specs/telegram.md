# Spec: Telegram Sender

## Overview

Telegram 消息发送模块。

## Interface

```go
type Sender struct {
    botToken  string
    chatID    string
    reportDir string
}

func NewSender(botToken, chatID, reportDir string) *Sender
func (s *Sender) Send(content, date string, dryRun bool) error
```

## Message Splitting

- Telegram 消息限制：4096 字符
- 按段落分割，保持 Markdown 结构
- 分段发送间隔 1 秒

## File Saving

- 路径：{reportDir}/{date}.md
