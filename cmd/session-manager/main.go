package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	flagVerbose   bool
	flagMessage   string
	flagShowScore bool
)

func main() {
	flag.BoolVar(&flagVerbose, "v", false, "Verbose output")
	flag.StringVar(&flagMessage, "m", "", "Message to analyze")
	flag.BoolVar(&flagShowScore, "score", false, "Show complexity score and exit")
	flag.Parse()

	baseDir := getBaseDir()

	if flagShowScore && flagMessage != "" {
		score := calculateComplexity(flagMessage)
		fmt.Printf("Complexity score: %d\n", score)
		os.Exit(0)
	}

	// Default: show current task complexity if in a session
	showCurrentTask(baseDir)
}

func getBaseDir() string {
	if dir := os.Getenv("PAI_BASE_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "src", "Code", "pai-universal")
}

type ComplexityResult struct {
	Score      int
	Level      string
	Factors    []string
	Suggestion string
}

func calculateComplexity(prompt string) int {
	score := 0
	words := strings.Fields(prompt)

	// File mentions
	for _, w := range words {
		lower := strings.ToLower(w)
		if strings.HasSuffix(lower, ".go") || strings.HasSuffix(lower, ".ts") ||
			strings.HasSuffix(lower, ".js") || strings.HasSuffix(lower, ".py") ||
			strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".json") {
			score++
		}
	}

	// Directory mentions
	dirs := []string{"cmd/", "lib/", "internal/", "pkg/", "api/", "ui/", "hooks/"}
	for _, dir := range dirs {
		if strings.Contains(prompt, dir) {
			score += 2
		}
	}

	// Keywords
	highComplex := map[string]int{
		"rewrite": 5, "migrate": 5, "refactor": 4,
		"create new": 4, "design": 3, "architecture": 4,
		"database": 3, "api": 2, "service": 2,
	}

	for kw, pts := range highComplex {
		if strings.Contains(strings.ToLower(prompt), kw) {
			score += pts
		}
	}

	// Testing/CI
	if strings.Contains(strings.ToLower(prompt), "test") {
		score += 2
	}
	if strings.Contains(strings.ToLower(prompt), "ci") || strings.Contains(strings.ToLower(prompt), "deploy") {
		score += 3
	}

	// Prompt length (significant requests)
	if len(words) > 200 {
		score += 3
	} else if len(words) > 100 {
		score += 2
	}

	// Multiple projects
	if strings.Count(prompt, "../") > 1 || strings.Count(prompt, "/home/") > 1 {
		score += 4
	}

	return score
}

func getLevel(score int) string {
	switch {
	case score <= 3:
		return "Standard"
	case score <= 8:
		return "Extended"
	case score <= 16:
		return "Advanced"
	default:
		return "Deep"
	}
}

func getSuggestion(score int) string {
	switch {
	case score <= 3:
		return "Normal mode - proceed directly"
	case score <= 8:
		return "Consider breaking into smaller tasks"
	case score <= 16:
		return "Recommend using Algorithm - create PRD"
	default:
		return "Use Algorithm - this is a complex task requiring ISC breakdown"
	}
}

func showCurrentTask(baseDir string) {
	hotPath := filepath.Join(baseDir, "MEMORY", "hot")
	currentTaskFile := filepath.Join(hotPath, "current-task.md")

	if _, err := os.Stat(currentTaskFile); err != nil {
		fmt.Println("No active task found.")
		fmt.Println("Run with -m 'your prompt' -score to analyze complexity")
		return
	}

	data, err := os.ReadFile(currentTaskFile)
	if err != nil {
		fmt.Printf("Error reading task: %v\n", err)
		return
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	var prompt string
	for i, line := range lines {
		if strings.HasPrefix(line, "## Prompt") {
			prompt = strings.Join(lines[i+1:], "\n")
			break
		}
	}

	if prompt != "" {
		score := calculateComplexity(prompt)
		level := getLevel(score)
		suggestion := getSuggestion(score)

		fmt.Printf("Current Task Complexity\n")
		fmt.Printf("========================\n")
		fmt.Printf("Level:    %s\n", level)
		fmt.Printf("Score:    %d\n", score)
		fmt.Printf("Suggest:  %s\n", suggestion)
	} else {
		fmt.Println("Could not parse prompt from current task")
	}
}

// LogCurrentTask saves the current task prompt for complexity tracking
func LogCurrentTask(prompt, baseDir string) error {
	hotPath := filepath.Join(baseDir, "MEMORY", "hot")
	if err := os.MkdirAll(hotPath, 0755); err != nil {
		return err
	}

	score := calculateComplexity(prompt)
	level := getLevel(score)

	content := fmt.Sprintf(`# Current Task

## Prompt
%s

## Metadata
- Score: %d
- Level: %s
- Started: %s

## Progress
- [ ] 

`, prompt, score, level, time.Now().Format("2006-01-02 15:04"))

	return os.WriteFile(filepath.Join(hotPath, "current-task.md"), []byte(content), 0644)
}
