package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/meganerd/pai-universal/cmd/session-analyzer/internal/config"
	"github.com/meganerd/pai-universal/cmd/session-analyzer/internal/parser"
)

type Insight struct {
	Type        string
	Category    string
	Description string
	Evidence    []string
	Confidence  float64
}

type Analyzer struct {
	cfg *config.Config
}

func NewAnalyzer(cfg *config.Config) *Analyzer {
	return &Analyzer{cfg: cfg}
}

func (a *Analyzer) Analyze(sessions []parser.Session) ([]Insight, error) {
	text := parser.SessionToText(sessions)

	prompt := buildAnalysisPrompt(text)

	var response string
	var err error

	switch strings.ToLower(a.cfg.LLMProvider) {
	case "openrouter":
		response, err = a.callOpenRouter(prompt)
	case "anthropic":
		response, err = a.callAnthropic(prompt)
	default:
		response, err = a.callAnthropic(prompt)
	}

	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	return parseInsights(response)
}

func buildAnalysisPrompt(sessionsText string) string {
	return `You are an AI assistant analyzing a user's session history to extract insights about their work, preferences, and patterns.

## User Profile Context
The user is a developer who:
- Prefers Go (Golang) for projects
- Uses "beads" (bd) as their issue tracker
- Works on homelab infrastructure
- Builds AI infrastructure tools

## Task
Analyze these session logs and extract:
1. **Technology Preferences** - Tools, languages, frameworks the user prefers
2. **Dev Patterns** - How they work (e.g., commits often, uses testing, etc.)
3. **Active Projects** - Projects they're working on
4. **Infrastructure** - Homelab, servers, services they manage
5. **Learnings** - New things they've learned or tried
6. **Goals** - Things they want to do (from what they ask for)

Session logs:
` + sessionsText + `

## Output Format
Return a JSON array of insights, each with:
- type: "preference", "pattern", "project", "infrastructure", "learning", "goal"
- category: specific category within that type
- description: what was found
- evidence: 1-2 example sessions that show this
- confidence: 0.0-1.0 how certain you are

Only include insights with confidence > 0.5. Max 20 insights.

JSON Output:`
}

func (a *Analyzer) callAnthropic(prompt string) (string, error) {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" && a.cfg.AnthropicKey != "" {
		key = a.cfg.AnthropicKey
	}
	if key == "" {
		key = os.Getenv("CLAUDE_API_KEY")
	}

	if key == "" {
		return "", fmt.Errorf("no API key found (ANTHROPIC_API_KEY or CLAUDE_API_KEY)")
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":      a.cfg.LLMModel,
		"max_tokens": 4096,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", key)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	content := result["content"].([]interface{})
	if len(content) == 0 {
		return "", fmt.Errorf("empty response")
	}

	text := content[0].(map[string]interface{})["text"].(string)
	return text, nil
}

func (a *Analyzer) callOpenRouter(prompt string) (string, error) {
	key := os.Getenv("OPENROUTER_API_KEY")
	if key == "" && a.cfg.LLMAPIKey != "" {
		key = a.cfg.LLMAPIKey
	}

	if key == "" {
		return "", fmt.Errorf("no API key found (OPENROUTER_API_KEY)")
	}

	// Default to a good openrouter model if not set
	model := a.cfg.LLMModel
	if model == "" {
		model = "anthropic/claude-sonnet-4-20250514"
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("HTTP-Referer", "https://github.com/meganerd/pai-universal")
	req.Header.Set("X-Title", "PAI Session Analyzer")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices := result["choices"].([]interface{})
	if len(choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	text := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return text, nil
}

func parseInsights(jsonStr string) ([]Insight, error) {
	// Find JSON array in response
	start := strings.Index(jsonStr, "[")
	end := strings.LastIndex(jsonStr, "]")
	if start == -1 || end == -1 {
		return nil, fmt.Errorf("no JSON array found in response")
	}

	jsonStr = jsonStr[start : end+1]

	var insights []Insight
	if err := json.Unmarshal([]byte(jsonStr), &insights); err != nil {
		// Try to fix common issues
		// Remove markdown code blocks
		jsonStr = strings.ReplaceAll(jsonStr, "```json", "")
		jsonStr = strings.ReplaceAll(jsonStr, "```", "")
		// Remove control characters that can break JSON
		jsonStr = strings.Map(func(r rune) rune {
			if r < 32 && r != '\n' && r != '\t' && r != '\r' {
				return -1
			}
			return r
		}, jsonStr)
		if err := json.Unmarshal([]byte(jsonStr), &insights); err != nil {
			return nil, fmt.Errorf("failed to parse insights: %w", err)
		}
	}

	return insights, nil
}
