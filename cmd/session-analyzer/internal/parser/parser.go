package parser

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/meganerd/pai-universal/cmd/session-analyzer/internal/config"
)

type Session struct {
	ID        string
	Timestamp time.Time
	Project   string
	Prompt    string
	Type      string
}

type parsedSession struct {
	Display   string `json:"display"`
	Project   string `json:"project"`
	SessionID string `json:"sessionId"`
	Timestamp int64  `json:"timestamp"`
}

func ParseHistory(path string) ([]Session, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sessions []Session
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var entry parsedSession
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if entry.Display == "" {
			continue
		}

		sessions = append(sessions, Session{
			ID:        entry.SessionID,
			Timestamp: time.UnixMilli(entry.Timestamp),
			Project:   entry.Project,
			Prompt:    entry.Display,
			Type:      "history",
		})
	}

	return sessions, scanner.Err()
}

func ParseMemoryWork(path string) ([]Session, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var sessions []Session

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		parts := parseMemoryDirName(entry.Name())
		if parts.date == "" {
			continue
		}

		timestamp, err := time.Parse("20060102-150405", parts.date)
		if err != nil {
			continue
		}

		sessions = append(sessions, Session{
			ID:        entry.Name(),
			Timestamp: timestamp,
			Project:   parts.project,
			Prompt:    parts.activity,
			Type:      "memory_work",
		})
	}

	return sessions, nil
}

type memoryDirParts struct {
	date     string
	activity string
	project  string
}

func parseMemoryDirName(name string) memoryDirParts {
	// Format: 20260228-120028_ssh-key-auth-phoenix
	var parts memoryDirParts

	// Find the first dash followed by 6 digits (date pattern)
	// Format is: YYYYMMDD-HHMMSS_description
	if len(name) < 15 {
		return parts
	}

	parts.date = name[:15] // "20260228-120028"
	if len(name) > 16 {
		parts.activity = name[16:] // everything after underscore
	}

	// Extract project from activity (last segment after underscore)
	parts.project = parts.activity
	return parts
}

func ParseAll(cfg *config.Config) ([]Session, error) {
	var allSessions []Session

	// Parse Claude history
	if cfg.ClaudeHistoryPath != "" {
		if _, err := os.Stat(cfg.ClaudeHistoryPath); err == nil {
			sessions, err := ParseHistory(cfg.ClaudeHistoryPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to parse history: %v\n", err)
			} else {
				allSessions = append(allSessions, sessions...)
			}
		}
	}

	// Parse MEMORY/WORK
	if cfg.ClaudeMemoryPath != "" {
		if _, err := os.Stat(cfg.ClaudeMemoryPath); err == nil {
			sessions, err := ParseMemoryWork(cfg.ClaudeMemoryPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to parse memory work: %v\n", err)
			} else {
				allSessions = append(allSessions, sessions...)
			}
		}
	}

	// Parse opencode
	if cfg.OpencodePath != "" {
		if _, err := os.Stat(cfg.OpencodePath); err == nil {
			sessions, err := ParseOpencode(cfg.OpencodePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to parse opencode: %v\n", err)
			} else {
				allSessions = append(allSessions, sessions...)
			}
		}
	}

	// Parse Codex
	if cfg.CodexPath != "" {
		if _, err := os.Stat(cfg.CodexPath); err == nil {
			sessions, err := ParseCodex(cfg.CodexPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to parse codex: %v\n", err)
			} else {
				allSessions = append(allSessions, sessions...)
			}
		}
	}

	// Parse Cursor (same format as Codex)
	if cfg.CursorPath != "" {
		if _, err := os.Stat(cfg.CursorPath); err == nil {
			sessions, err := ParseCodex(cfg.CursorPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to parse cursor: %v\n", err)
			} else {
				allSessions = append(allSessions, sessions...)
			}
		}
	}

	// Parse pi-go sessions
	if cfg.PigoPath != "" {
		if sessions, err := ParsePigo(cfg.PigoPath); err == nil && len(sessions) > 0 {
			allSessions = append(allSessions, sessions...)
		}
	}

	return allSessions, nil
}

func SessionToText(sessions []Session) string {
	var text string
	for _, s := range sessions {
		text += fmt.Sprintf("[%s] %s: %s\n", s.Timestamp.Format("2006-01-02"), s.Project, s.Prompt)
	}
	return text
}

func ParseOpencode(dbPath string) ([]Session, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT s.id, s.title, s.directory, s.time_created
		FROM session s
		ORDER BY s.time_created DESC
		LIMIT 500
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var id, title, directory string
		var timeCreated int64

		if err := rows.Scan(&id, &title, &directory, &timeCreated); err != nil {
			continue
		}

		sessions = append(sessions, Session{
			ID:        id,
			Timestamp: time.UnixMilli(timeCreated),
			Project:   directory,
			Prompt:    title,
			Type:      "opencode",
		})
	}

	return sessions, nil
}

func ParseCodex(dbPath string) ([]Session, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Codex/Cursor logs table - extract user feedback logs as sessions
	query := `
		SELECT id, ts, target, feedback_log_body
		FROM logs
		WHERE feedback_log_body IS NOT NULL
		ORDER BY ts DESC
		LIMIT 500
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var id int64
		var ts int64
		var target, body string

		if err := rows.Scan(&id, &ts, &target, &body); err != nil {
			continue
		}

		// Use target as project (often the module/path)
		project := target
		if len(project) > 100 {
			project = project[:100]
		}

		sessions = append(sessions, Session{
			ID:        fmt.Sprintf("codex-%d", id),
			Timestamp: time.Unix(ts, 0),
			Project:   project,
			Prompt:    body,
			Type:      "codex",
		})
	}

	return sessions, nil
}

func ParsePigo(pigoDir string) ([]Session, error) {
	// Parse pi-go sessions from ~/.local/share/pi-go/sessions
	home, _ := os.UserHomeDir()
	sessionsDir := filepath.Join(home, ".local", "share", "pi-go", "sessions")

	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return nil, err
	}

	var sessions []Session
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		path := filepath.Join(sessionsDir, entry.Name())
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		// Extract timestamp from filename (format: 1773081361658587.jsonl)
		ts, err := strconv.ParseInt(strings.TrimSuffix(entry.Name(), ".jsonl"), 10, 64)
		if err != nil {
			ts = time.Now().Unix()
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var msg map[string]interface{}
			if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
				continue
			}

			// Extract prompt from user messages
			role, _ := msg["role"].(string)
			content, _ := msg["content"].(string)

			if role == "user" && content != "" {
				sessions = append(sessions, Session{
					ID:        entry.Name(),
					Timestamp: time.UnixMilli(ts / 1000),
					Project:   pigoDir,
					Prompt:    content,
					Type:      "pi-go",
				})
				break // Only take first user message per session
			}
		}
		file.Close()
	}

	return sessions, nil
}
