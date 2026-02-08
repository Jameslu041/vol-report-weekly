package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxMessageLength = 4096
	sendInterval     = time.Second
)

type Sender struct {
	botToken  string
	chatID    string
	reportDir string
}

func NewSender(botToken, chatID, reportDir string) *Sender {
	return &Sender{
		botToken:  botToken,
		chatID:    chatID,
		reportDir: reportDir,
	}
}

func (s *Sender) Send(content, date string, dryRun bool) error {
	if err := s.saveToFile(content, date); err != nil {
		slog.Warn("failed to save report file", "error", err)
	}

	if dryRun {
		slog.Info("dry run mode, skipping telegram send")
		return nil
	}

	parts := s.splitMessage(content)
	for i, part := range parts {
		if i > 0 {
			time.Sleep(sendInterval)
		}

		if err := s.sendMessage(part); err != nil {
			return fmt.Errorf("send message part %d: %w", i+1, err)
		}
		slog.Info("sent message part", "part", i+1, "total", len(parts))
	}

	return nil
}

func (s *Sender) saveToFile(content, date string) error {
	if err := os.MkdirAll(s.reportDir, 0755); err != nil {
		return fmt.Errorf("create report directory: %w", err)
	}

	path := filepath.Join(s.reportDir, date+".md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write report file: %w", err)
	}

	slog.Info("saved report to file", "path", path)
	return nil
}

func (s *Sender) splitMessage(content string) []string {
	if len(content) <= maxMessageLength {
		return []string{content}
	}

	var parts []string
	sections := strings.Split(content, "\n## ")

	var current strings.Builder
	for i, section := range sections {
		if i > 0 {
			section = "## " + section
		}

		if current.Len()+len(section)+1 > maxMessageLength {
			if current.Len() > 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			}

			if len(section) > maxMessageLength {
				for len(section) > maxMessageLength {
					parts = append(parts, section[:maxMessageLength])
					section = section[maxMessageLength:]
				}
			}
		}

		if current.Len() > 0 {
			current.WriteString("\n")
		}
		current.WriteString(section)
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type sendMessageResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}

func (s *Sender) sendMessage(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	req := sendMessageRequest{
		ChatID:    s.chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	var result sendMessageResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if !result.OK {
		return fmt.Errorf("telegram api error: %s", result.Description)
	}

	return nil
}
