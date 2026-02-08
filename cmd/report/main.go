package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"

	"vol-report-weekly/internal/api"
	"vol-report-weekly/internal/config"
	"vol-report-weekly/internal/report"
	"vol-report-weekly/internal/storage"
	"vol-report-weekly/internal/telegram"
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
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]
	switch command {
	case "run":
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		dryRun := runCmd.Bool("dry-run", false, "Generate report without sending")
		dateStr := runCmd.String("date", time.Now().Format("2006-01-02"), "Report date (YYYY-MM-DD)")
		runCmd.Parse(os.Args[2:])

		if err := runReport(cfg, *dateStr, *dryRun); err != nil {
			slog.Error("report generation failed", "error", err)
			os.Exit(1)
		}
	case "cron":
		runCron(cfg)
	default:
		slog.Error("unknown command", "command", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: vol-report <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run [--dry-run] [--date YYYY-MM-DD]")
	fmt.Println("      Run report generation once")
	fmt.Println()
	fmt.Println("  cron")
	fmt.Println("      Run in scheduled mode (uses CRON_SCHEDULE env)")
}

func runReport(cfg *config.Config, dateStr string, dryRun bool) error {
	slog.Info("running report generation", "date", dateStr, "dry_run", dryRun)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Initialize components
	apiClient := api.NewClient(cfg.APIBaseURL)
	store := storage.NewStorage(cfg.DataDir)
	gen := report.NewGenerator()
	sender := telegram.NewSender(cfg.TGBotToken, cfg.TGChatID, cfg.ReportDir)

	// Fetch data from API
	slog.Info("fetching data from API")
	snapshot, err := fetchAllData(ctx, apiClient, dateStr)
	if err != nil {
		return fmt.Errorf("fetch data: %w", err)
	}

	// Save snapshot
	slog.Info("saving snapshot")
	if err := store.Save(snapshot); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}

	// Load previous snapshot
	slog.Info("loading previous snapshot")
	previous, err := store.LoadPrevious(dateStr)
	if err != nil {
		slog.Warn("failed to load previous snapshot", "error", err)
	}
	if previous == nil {
		slog.Info("no previous snapshot found, first run")
	}

	// Generate report
	slog.Info("generating report")
	content, err := gen.Generate(snapshot, previous)
	if err != nil {
		return fmt.Errorf("generate report: %w", err)
	}

	// Send report
	slog.Info("sending report", "dry_run", dryRun)
	if err := sender.Send(content, dateStr, dryRun); err != nil {
		return fmt.Errorf("send report: %w", err)
	}

	slog.Info("report generation completed")
	return nil
}

func fetchAllData(ctx context.Context, client *api.Client, dateStr string) (*storage.Snapshot, error) {
	snapshot := &storage.Snapshot{Date: dateStr}

	var err error

	slog.Info("fetching IV history")
	snapshot.IVHistory, err = client.GetIVHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("get IV history: %w", err)
	}

	slog.Info("fetching IV/RV")
	snapshot.IVRV, err = client.GetIVRV(ctx)
	if err != nil {
		return nil, fmt.Errorf("get IV/RV: %w", err)
	}

	slog.Info("fetching skew chart")
	snapshot.SkewChart, err = client.GetSkewChart(ctx)
	if err != nil {
		return nil, fmt.Errorf("get skew chart: %w", err)
	}

	slog.Info("fetching FIV matrix")
	snapshot.FIVMatrix, err = client.GetFIVMatrix(ctx)
	if err != nil {
		return nil, fmt.Errorf("get FIV matrix: %w", err)
	}

	slog.Info("fetching option flows")
	snapshot.OptionFlows, err = client.GetOptionFlows(ctx)
	if err != nil {
		return nil, fmt.Errorf("get option flows: %w", err)
	}

	return snapshot, nil
}

func runCron(cfg *config.Config) {
	slog.Info("starting cron mode", "schedule", cfg.CronSchedule)

	c := cron.New()
	_, err := c.AddFunc(cfg.CronSchedule, func() {
		dateStr := time.Now().Format("2006-01-02")
		slog.Info("cron triggered, running report", "date", dateStr)
		if err := runReport(cfg, dateStr, false); err != nil {
			slog.Error("cron report failed", "error", err)
		}
	})
	if err != nil {
		slog.Error("invalid cron schedule", "schedule", cfg.CronSchedule, "error", err)
		os.Exit(1)
	}

	c.Start()
	slog.Info("cron scheduler started, waiting for schedule")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("received shutdown signal, stopping cron")
	c.Stop()
}
