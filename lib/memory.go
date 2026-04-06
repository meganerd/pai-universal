package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

// MemoryTier represents the three-tier memory system
type MemoryTier string

const (
	MemoryTierHot  MemoryTier = "hot"
	MemoryTierWarm MemoryTier = "warm"
	MemoryTierCold MemoryTier = "cold"
)

// MemoryStore provides 3-tier memory storage
type MemoryStore struct {
	BasePath string
}

// NewMemoryStore creates a new memory store
func NewMemoryStore(basePath string) *MemoryStore {
	return &MemoryStore{BasePath: basePath}
}

// GetPath returns the path for a given memory tier
func (m *MemoryStore) GetPath(tier MemoryTier) string {
	return filepath.Join(m.BasePath, string(tier))
}

// Ensure creates the memory directories if they don't exist
func (m *MemoryStore) Ensure() error {
	for _, tier := range []MemoryTier{MemoryTierHot, MemoryTierWarm, MemoryTierCold} {
		path := m.GetPath(tier)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("create memory directory %s: %w", path, err)
		}
	}
	return nil
}

// Add adds a memory entry to the specified tier
func (m *MemoryStore) Add(tier MemoryTier, content string) error {
	path := m.GetPath(tier)

	// Create unique filename
	filename := fmt.Sprintf("memory-%d.md", os.Getpid())
	filepath := filepath.Join(path, filename)

	return os.WriteFile(filepath, []byte(content), 0644)
}

// List returns all memory entries in a tier
func (m *MemoryStore) List(tier MemoryTier) ([]string, error) {
	path := m.GetPath(tier)

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, entry := range entries {
		if !entry.IsDir() {
			results = append(results, entry.Name())
		}
	}
	return results, nil
}
