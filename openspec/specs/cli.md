# Spec: CLI

## Commands

```
vol-report run [--dry-run] [--date YYYY-MM-DD]
vol-report cron
```

## Flags

- --dry-run: 仅生成报告，不发送 Telegram
- --date: 指定报告日期（默认今天）

## Run Flow

1. 加载配置
2. 调用 API 获取数据
3. 保存快照
4. 加载上周快照
5. 生成报告
6. 发送 Telegram（除非 dry-run）

## Cron Mode

- 使用 robfig/cron 库
- 按 CRON_SCHEDULE 定时执行 run
