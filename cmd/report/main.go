package main

import (
	"fmt"
	"log/slog"
	"os"

	"vol-report-weekly/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting vol-report")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("config loaded",
		"api_base_url", cfg.APIBaseURL,
		"data_dir", cfg.DataDir,
		"report_dir", cfg.ReportDir,
		"cron_schedule", cfg.CronSchedule,
	)

	if len(os.Args) < 2 {
		fmt.Println("Usage: vol-report <command>")
		fmt.Println("Commands:")
		fmt.Println("  run        Run report generation once")
		fmt.Println("  cron       Run in scheduled mode")
		os.Exit(0)
	}

	command := os.Args[1]
	switch command {
	case "run":
		slog.Info("running report generation")
		// TODO: implement report generation
	case "cron":
		slog.Info("starting cron mode", "schedule", cfg.CronSchedule)
		// TODO: implement cron scheduling
	default:
		slog.Error("unknown command", "command", command)
		os.Exit(1)
	}
}
