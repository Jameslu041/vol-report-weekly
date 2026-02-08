package storage

import (
	"os"
	"path/filepath"
	"testing"

	"vol-report-weekly/internal/api"
)

func TestStorage_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	snapshot := &Snapshot{
		Date: "2026-02-08",
		IVHistory: &api.IVHistoryResponse{
			Items: []api.IVHistoryData{
				{Tenor: "1D", IV: 45.5, HV: 42.0},
			},
		},
		IVRV: &api.IVRVResponse{
			Items: []api.IVRVData{
				{Period: "7D", IV: 50.0, RV: 45.0, VRP: 5.0},
			},
		},
	}

	if err := storage.Save(snapshot); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	path := filepath.Join(tmpDir, "snapshots", "2026-02-08.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("snapshot file not created")
	}

	loaded, err := storage.Load("2026-02-08")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if loaded.Date != snapshot.Date {
		t.Errorf("expected date %s, got %s", snapshot.Date, loaded.Date)
	}
	if len(loaded.IVHistory.Items) != 1 {
		t.Errorf("expected 1 IV history item, got %d", len(loaded.IVHistory.Items))
	}
	if loaded.IVHistory.Items[0].IV != 45.5 {
		t.Errorf("expected IV 45.5, got %f", loaded.IVHistory.Items[0].IV)
	}
}

func TestStorage_LoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	snapshot, err := storage.Load("2020-01-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snapshot != nil {
		t.Error("expected nil snapshot for non-existent file")
	}
}

func TestStorage_LoadPrevious(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	previousSnapshot := &Snapshot{
		Date: "2026-02-01",
		IVHistory: &api.IVHistoryResponse{
			Items: []api.IVHistoryData{
				{Tenor: "1D", IV: 40.0},
			},
		},
	}
	if err := storage.Save(previousSnapshot); err != nil {
		t.Fatalf("failed to save previous: %v", err)
	}

	loaded, err := storage.LoadPrevious("2026-02-08")
	if err != nil {
		t.Fatalf("failed to load previous: %v", err)
	}

	if loaded == nil {
		t.Fatal("expected previous snapshot")
	}
	if loaded.Date != "2026-02-01" {
		t.Errorf("expected date 2026-02-01, got %s", loaded.Date)
	}
}

func TestStorage_LoadPreviousNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	loaded, err := storage.LoadPrevious("2026-02-08")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded != nil {
		t.Error("expected nil for non-existent previous snapshot")
	}
}
