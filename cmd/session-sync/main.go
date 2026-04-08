package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	flagDryRun  bool
	flagVerbose bool
	flagSource  string
	flagTarget  string
)

func main() {
	flag.BoolVar(&flagDryRun, "dry-run", false, "Show what would be synced without making changes")
	flag.BoolVar(&flagVerbose, "v", false, "Verbose output")
	flag.StringVar(&flagSource, "source", "", "Source harness (claude, opencode, codex, cursor, pigo)")
	flag.StringVar(&flagTarget, "target", "", "Target harness (comma-separated: claude,opencode,etc)")
	flag.Parse()

	baseDir := os.Getenv("PAI_BASE_DIR")
	if baseDir == "" {
		home, _ := os.UserHomeDir()
		baseDir = filepath.Join(home, "src", "Code", "pai-universal")
	}

	if flagSource == "" || flagTarget == "" {
		log.Fatal("Both --source and --target are required")
	}

	sessions, err := loadSourceSessions(flagSource, baseDir)
	if err != nil {
		log.Fatalf("Failed to load sessions from %s: %v", flagSource, err)
	}

	if flagVerbose {
		fmt.Printf("Loaded %d sessions from %s\n", len(sessions), flagSource)
	}

	targets := parseTargets(flagTarget)
	for _, target := range targets {
		if err := syncToTarget(target, sessions, baseDir, flagDryRun); err != nil {
			log.Printf("Warning: Failed to sync to %s: %v", target, err)
		}
	}

	fmt.Println("Session sync complete")
}

type Session struct {
	ID        string
	Timestamp time.Time
	Project   string
	Content   string
	Source    string
}

func loadSourceSessions(source, baseDir string) ([]Session, error) {
	switch source {
	case "claude":
		return loadClaudeSessions()
	case "opencode":
		return loadOpencodeSessions()
	case "codex", "cursor":
		return loadCodexSessions(source)
	case "pigo":
		return loadPigoSessions()
	default:
		return nil, fmt.Errorf("unknown source: %s", source)
	}
}

func loadClaudeSessions() ([]Session, error) {
	home, _ := os.UserHomeDir()
	historyPath := filepath.Join(home, ".claude", "history.jsonl")

	if _, err := os.Stat(historyPath); err != nil {
		return nil, err
	}

	// TODO: Parse history.jsonl
	return []Session{}, nil
}

func loadOpencodeSessions() ([]Session, error) {
	home, _ := os.UserHomeDir()
	dbPath := filepath.Join(home, ".local", "share", "opencode", "opencode.db")

	if _, err := os.Stat(dbPath); err != nil {
		return nil, err
	}

	// TODO: Query opencode.db
	return []Session{}, nil
}

func loadCodexSessions(source string) ([]Session, error) {
	home, _ := os.UserHomeDir()
	var dbPath string

	if source == "codex" {
		dbPath = filepath.Join(home, ".codex", "logs_1.sqlite")
	} else {
		dbPath = filepath.Join(home, ".cursor", "logs_1.sqlite")
	}

	if _, err := os.Stat(dbPath); err != nil {
		return nil, err
	}

	// TODO: Query Codex/Cursor SQLite
	return []Session{}, nil
}

func loadPigoSessions() ([]Session, error) {
	home, _ := os.UserHomeDir()
	sessionsDir := filepath.Join(home, ".local", "share", "pi-go", "sessions")

	if _, err := os.Stat(sessionsDir); err != nil {
		return nil, err
	}

	// TODO: Parse pi-go JSONL sessions
	return []Session{}, nil
}

func parseTargets(targets string) []string {
	var result []string
	for _, t := range splitCommas(targets) {
		t = trimSpace(t)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}

func splitCommas(s string) []string {
	var result []string
	var current []rune

	for _, r := range s {
		if r == ',' {
			result = append(result, string(current))
			current = nil
		} else {
			current = append(current, r)
		}
	}
	if len(current) > 0 {
		result = append(result, string(current))
	}

	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for ; start < end && (s[start] == ' ' || s[start] == '\t'); start++ {
	}
	for ; end > start && (s[end-1] == ' ' || s[end-1] == '\t'); end-- {
	}

	return s[start:end]
}

func syncToTarget(target string, sessions []Session, baseDir string, dryRun bool) error {
	if flagVerbose {
		fmt.Printf("Syncing %d sessions to %s\n", len(sessions), target)
	}

	memoryPath := getTargetMemoryPath(target, baseDir)

	if err := os.MkdirAll(memoryPath, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("sync-%s-%s.md", target, time.Now().Format("20060102-150405"))

	content := "# Cross-harness Session Sync\n\n"
	content += fmt.Sprintf("Source harness: session-sync\n")
	content += fmt.Sprintf("Timestamp: %s\n\n", time.Now().Format("2006-01-02 15:04"))
	content += fmt.Sprintf("Synced %d sessions\n\n", len(sessions))

	dest := filepath.Join(memoryPath, filename)

	if dryRun {
		fmt.Printf("Would create: %s\n", dest)
		return nil
	}

	return os.WriteFile(dest, []byte(content), 0644)
}

func getTargetMemoryPath(target, baseDir string) string {
	home, _ := os.UserHomeDir()

	switch target {
	case "claude":
		return filepath.Join(home, ".claude", "MEMORY", "WORK")
	case "opencode":
		return filepath.Join(home, ".local", "share", "opencode", "storage", "memory")
	case "pigo":
		return filepath.Join(home, ".local", "share", "pi-go", "memory")
	case "pai-universal":
		return filepath.Join(baseDir, "MEMORY", "warm")
	default:
		return filepath.Join(baseDir, "MEMORY", "warm")
	}
}
