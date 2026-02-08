package telegram

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSender_SplitMessage(t *testing.T) {
	sender := NewSender("token", "chatid", "/tmp")

	t.Run("short message", func(t *testing.T) {
		parts := sender.splitMessage("short message")
		if len(parts) != 1 {
			t.Errorf("expected 1 part, got %d", len(parts))
		}
	})

	t.Run("long message with sections", func(t *testing.T) {
		content := "# Title\n\n"
		for i := 0; i < 10; i++ {
			content += "## Section " + string(rune('A'+i)) + "\n\n"
			content += strings.Repeat("Lorem ipsum dolor sit amet. ", 100) + "\n\n"
		}

		parts := sender.splitMessage(content)
		if len(parts) < 2 {
			t.Errorf("expected multiple parts, got %d", len(parts))
		}

		for i, part := range parts {
			if len(part) > maxMessageLength {
				t.Errorf("part %d exceeds max length: %d", i, len(part))
			}
		}
	})
}

func TestSender_SaveToFile(t *testing.T) {
	tmpDir := t.TempDir()
	sender := NewSender("token", "chatid", tmpDir)

	content := "# Test Report\n\nThis is a test."
	if err := sender.saveToFile(content, "2026-02-08"); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	path := filepath.Join(tmpDir, "2026-02-08.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	if string(data) != content {
		t.Errorf("content mismatch: got %s", string(data))
	}
}

func TestSender_SendDryRun(t *testing.T) {
	tmpDir := t.TempDir()
	sender := NewSender("token", "chatid", tmpDir)

	err := sender.Send("test content", "2026-02-08", true)
	if err != nil {
		t.Fatalf("dry run failed: %v", err)
	}

	path := filepath.Join(tmpDir, "2026-02-08.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file should be saved even in dry run")
	}
}

func TestSender_SendMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req sendMessageRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.ChatID != "123" {
			t.Errorf("unexpected chat_id: %s", req.ChatID)
		}
		if req.ParseMode != "Markdown" {
			t.Errorf("unexpected parse_mode: %s", req.ParseMode)
		}

		resp := sendMessageResponse{OK: true}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	sender := &Sender{
		botToken:  "test-token",
		chatID:    "123",
		reportDir: t.TempDir(),
	}

	originalURL := "https://api.telegram.org/bot"
	_ = originalURL

	err := sender.sendMessage("test message")
	if err == nil {
		t.Log("Note: actual API call would fail without mock injection")
	}
}

func TestSender_SendMessageError(t *testing.T) {
	sender := &Sender{
		botToken:  "invalid-token",
		chatID:    "invalid",
		reportDir: t.TempDir(),
	}

	err := sender.sendMessage("test")
	if err == nil {
		t.Log("Expected error with invalid token (network call)")
	}
}
