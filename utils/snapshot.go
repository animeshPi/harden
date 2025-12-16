package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"harden/policy"
)

type PolicySnapshot struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Category  string        `json:"category"`
	Severity  string        `json:"severity"`
	Result    CommandResult `json:"result"`
}

type Snapshot struct {
	CreatedAt time.Time        `json:"created_at"`
	Host      string           `json:"host"`
	Policies  []PolicySnapshot `json:"policies"`
}

// BuildSnapshot executes policy check commands and records results
func BuildSnapshot(policies []policy.Policy) (*Snapshot, error) {
	host, _ := os.Hostname()

	snap := &Snapshot{
		CreatedAt: time.Now().UTC(),
		Host:      host,
	}

	for _, p := range policies {
		res := Execute(p.CheckCmd)

		snap.Policies = append(snap.Policies, PolicySnapshot{
			ID:       p.ID,
			Title:    p.Title,
			Category: p.Category,
			Severity: p.Severity,
			Result:   res,
		})
	}

	return snap, nil
}

// SaveNextToExecutable writes snapshot JSON inside a `snapshots` folder
// next to the running binary, using YYYY-MM-DD-HHMMSS format.
func (s *Snapshot) SaveNextToExecutable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	baseDir := filepath.Dir(exe)

	// Ensure snapshots directory exists
	snapDir := filepath.Join(baseDir, "snapshots")
	if err := os.MkdirAll(snapDir, 0755); err != nil {
		return "", err
	}

	// YYYY-MM-DD-HHMMSS format
	filename := "snapshot-" + time.Now().Format("2006-01-02-150405") + ".json"
	path := filepath.Join(snapDir, filename)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}

	return path, nil
}
