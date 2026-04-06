package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ClaudeHistoryPath string `json:"claude_history_path"`
	ClaudeMemoryPath  string `json:"claude_memory_path"`
	OpencodePath      string `json:"opencode_path"`
	CodexPath         string `json:"codex_path"`
	CursorPath        string `json:"cursor_path"`
	GeminiPath        string `json:"gemini_path"`
	PigoPath          string `json:"pigo_path"`
	LLMProvider       string `json:"llm_provider"`
	LLMModel          string `json:"llm_model"`
	LLMAPIKey         string `json:"llm_api_key"`
	LLMEndpoint       string `json:"llm_endpoint"`
	AnthropicKey      string `json:"anthropic_key"`
}

func Default() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		ClaudeHistoryPath: filepath.Join(home, ".claude", "history.jsonl"),
		ClaudeMemoryPath:  filepath.Join(home, ".claude", "MEMORY", "WORK"),
		OpencodePath:      filepath.Join(home, ".local", "share", "opencode", "opencode.db"),
		CodexPath:         filepath.Join(home, ".codex", "logs_1.sqlite"),
		CursorPath:        "",
		GeminiPath:        "",
		PigoPath:          filepath.Join(home, "src", "Code", "pi-go"),
		LLMProvider:       "openrouter",
		LLMModel:          "google/gemini-2.0-flash-001",
		AnthropicKey:      os.Getenv("ANTHROPIC_API_KEY"),
	}
}

func Load(path string) (*Config, error) {
	cfg := Default()

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
