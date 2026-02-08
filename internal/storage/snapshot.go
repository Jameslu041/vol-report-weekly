package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"vol-report-weekly/internal/api"
)

type Snapshot struct {
	Date        string                   `json:"date"`
	IVHistory   *api.IVHistoryResponse   `json:"iv_history"`
	IVRV        *api.IVRVResponse        `json:"iv_rv"`
	SkewChart   *api.SkewChartResponse   `json:"skew_chart"`
	FIVMatrix   *api.FIVMatrixResponse   `json:"fiv_matrix"`
	OptionFlows *api.OptionFlowsResponse `json:"option_flows"`
}

type Storage struct {
	dataDir string
}

func NewStorage(dataDir string) *Storage {
	return &Storage{dataDir: dataDir}
}

func (s *Storage) snapshotPath(date string) string {
	return filepath.Join(s.dataDir, "snapshots", date+".json")
}

func (s *Storage) Save(snapshot *Snapshot) error {
	dir := filepath.Join(s.dataDir, "snapshots")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create snapshots directory: %w", err)
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	path := s.snapshotPath(snapshot.Date)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write snapshot file: %w", err)
	}

	return nil
}

func (s *Storage) Load(date string) (*Snapshot, error) {
	path := s.snapshotPath(date)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read snapshot file: %w", err)
	}

	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}

	return &snapshot, nil
}

func (s *Storage) LoadPrevious(currentDate string) (*Snapshot, error) {
	current, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		return nil, fmt.Errorf("parse current date: %w", err)
	}

	previousDate := current.AddDate(0, 0, -7).Format("2006-01-02")
	return s.Load(previousDate)
}
