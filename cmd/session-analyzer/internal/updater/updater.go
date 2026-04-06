package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/meganerd/pai-universal/cmd/session-analyzer/internal/llm"
)

type HarnessTarget struct {
	Name       string
	MemoryPath string
}

var DefaultHarnessTargets = []HarnessTarget{
	{Name: "pai-universal", MemoryPath: ""}, // Uses baseDir
	{Name: "claude", MemoryPath: ""},        // Will be set dynamically
	{Name: "opencode", MemoryPath: ""},
	{Name: "codex", MemoryPath: ""},
	{Name: "cursor", MemoryPath: ""},
}

type Updater struct {
	baseDir string
	dryRun  bool
	targets []string // which harnesses to update, empty = all
}

func NewUpdater(baseDir string, dryRun bool) *Updater {
	return &Updater{
		baseDir: baseDir,
		dryRun:  dryRun,
		targets: []string{},
	}
}

func NewUpdaterWithTargets(baseDir string, dryRun bool, targets []string) *Updater {
	u := &Updater{
		baseDir: baseDir,
		dryRun:  dryRun,
		targets: targets,
	}
	return u
}

func (u *Updater) SetTargets(targets []string) {
	u.targets = targets
}

func (u *Updater) shouldUpdate(name string) bool {
	if len(u.targets) == 0 {
		return true // no targets specified = update all
	}
	for _, t := range u.targets {
		if t == name || t == "all" {
			return true
		}
	}
	return false
}

func (u *Updater) Apply(insights []llm.Insight) error {
	if u.dryRun {
		fmt.Println("DRY RUN - No changes will be made")
	}

	// Categorize insights
	var memoryInsights []llm.Insight
	var goalInsights []llm.Insight
	var beliefInsights []llm.Insight

	for _, insight := range insights {
		switch insight.Type {
		case "learning", "preference", "pattern":
			memoryInsights = append(memoryInsights, insight)
		case "goal":
			goalInsights = append(goalInsights, insight)
		case "belief":
			beliefInsights = append(beliefInsights, insight)
		}
	}

	// Update pai-universal memory (always, unless targets explicitly exclude it)
	if u.shouldUpdate("pai-universal") || u.shouldUpdate("all") {
		// Apply memory updates
		if len(memoryInsights) > 0 {
			if err := u.updateMemory(memoryInsights); err != nil {
				return fmt.Errorf("memory update failed: %w", err)
			}
		}

		// Apply goal updates
		if len(goalInsights) > 0 {
			if err := u.updateGoals(goalInsights); err != nil {
				return fmt.Errorf("goal update failed: %w", err)
			}
		}

		// Apply belief updates
		if len(beliefInsights) > 0 {
			if err := u.updateBeliefs(beliefInsights); err != nil {
				return fmt.Errorf("belief update failed: %w", err)
			}
		}
	}

	// Update Claude memory
	if u.shouldUpdate("claude") {
		if err := u.updateClaudeMemory(insights); err != nil {
			return fmt.Errorf("claude memory update failed: %w", err)
		}
	}

	// Update opencode memory
	if u.shouldUpdate("opencode") {
		if err := u.updateOpencodeMemory(insights); err != nil {
			return fmt.Errorf("opencode memory update failed: %w", err)
		}
	}

	// Update Codex memory (writes to same format)
	if u.shouldUpdate("codex") {
		if err := u.updateClaudeMemory(insights); err != nil {
			return fmt.Errorf("codex memory update failed: %w", err)
		}
	}

	// Update Cursor memory (writes to same format)
	if u.shouldUpdate("cursor") {
		if err := u.updateClaudeMemory(insights); err != nil {
			return fmt.Errorf("cursor memory update failed: %w", err)
		}
	}

	return nil
}

func (u *Updater) updateMemory(insights []llm.Insight) error {
	memoryPath := filepath.Join(u.baseDir, "MEMORY", "cold")

	if err := os.MkdirAll(memoryPath, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("insights-%s.md", time.Now().Format("20060102"))
	content := "# Session Insights\n\n"
	content += fmt.Sprintf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04"))

	for _, insight := range insights {
		content += fmt.Sprintf("## %s: %s\n", insight.Type, insight.Category)
		content += fmt.Sprintf("**Confidence**: %.0f%%\n\n", insight.Confidence*100)
		content += fmt.Sprintf("%s\n\n", insight.Description)
		content += "**Evidence:**\n"
		for _, e := range insight.Evidence {
			content += fmt.Sprintf("- %s\n", e)
		}
		content += "\n---\n\n"
	}

	dest := filepath.Join(memoryPath, filename)
	if u.dryRun {
		fmt.Printf("Would create: %s\n", dest)
		fmt.Printf("Content preview:\n%s\n", content[:min(len(content), 500)])
	} else {
		return os.WriteFile(dest, []byte(content), 0644)
	}

	return nil
}

func (u *Updater) updateGoals(insights []llm.Insight) error {
	goalsPath := filepath.Join(u.baseDir, "USER", "TELOS", "GOALS.md")

	if _, err := os.Stat(goalsPath); err != nil {
		return err
	}

	// Read existing goals
	existing, _ := os.ReadFile(goalsPath)

	// Add inferred goals section
	newContent := string(existing)
	if !strings.Contains(newContent, "## Inferred from Sessions") {
		newContent += "\n\n## Inferred from Sessions\n"
		newContent += fmt.Sprintf("*Last updated: %s*\n\n", time.Now().Format("2006-01-02"))

		for _, insight := range insights {
			newContent += fmt.Sprintf("- [%s] %s\n", insight.Category, insight.Description)
		}

		if u.dryRun {
			fmt.Printf("Would update: %s\n", goalsPath)
			fmt.Printf("New content:\n%s\n", newContent)
		} else {
			return os.WriteFile(goalsPath, []byte(newContent), 0644)
		}
	}

	return nil
}

func (u *Updater) updateBeliefs(insights []llm.Insight) error {
	beliefsPath := filepath.Join(u.baseDir, "USER", "TELOS", "BELIEFS.md")

	if _, err := os.Stat(beliefsPath); err != nil {
		return err
	}

	existing, _ := os.ReadFile(beliefsPath)
	newContent := string(existing)

	hasTechPrefs := strings.Contains(newContent, "## Technology Preferences")

	if !hasTechPrefs {
		newContent += "\n\n## Technology Preferences\n"
		newContent += fmt.Sprintf("*Inferred from session analysis - %s*\n\n", time.Now().Format("2006-01-02"))

		for _, insight := range insights {
			if insight.Type == "preference" {
				newContent += fmt.Sprintf("- **%s**: %s\n", insight.Category, insight.Description)
			}
		}

		if u.dryRun {
			fmt.Printf("Would update: %s\n", beliefsPath)
		} else {
			return os.WriteFile(beliefsPath, []byte(newContent), 0644)
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// updateClaudeMemory writes insights to Claude's MEMORY/WORK format
func (u *Updater) updateClaudeMemory(insights []llm.Insight) error {
	home, _ := os.UserHomeDir()
	claudeMemoryPath := filepath.Join(home, ".claude", "MEMORY", "WORK")

	if err := os.MkdirAll(claudeMemoryPath, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("insights-%s.md", time.Now().Format("20060102-150405"))
	content := "# Session Insights\n\n"
	content += fmt.Sprintf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04"))
	content += "Source: PAI Universal Session Analyzer\n\n"

	for _, insight := range insights {
		content += fmt.Sprintf("## %s: %s\n", insight.Type, insight.Category)
		content += fmt.Sprintf("**Confidence**: %.0f%%\n\n", insight.Description)
		content += "**Evidence:**\n"
		for _, e := range insight.Evidence {
			content += fmt.Sprintf("- %s\n", e)
		}
		content += "\n---\n\n"
	}

	dest := filepath.Join(claudeMemoryPath, filename)
	if u.dryRun {
		fmt.Printf("Would update Claude memory: %s\n", dest)
	} else {
		return os.WriteFile(dest, []byte(content), 0644)
	}
	return nil
}

// updateOpencodeMemory writes to opencode's memory format
func (u *Updater) updateOpencodeMemory(insights []llm.Insight) error {
	home, _ := os.UserHomeDir()
	opencodeMemoryPath := filepath.Join(home, ".local", "share", "opencode", "storage", "memory")

	if err := os.MkdirAll(opencodeMemoryPath, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("insights-%s.md", time.Now().Format("20060102"))
	content := "# Session Insights\n\n"
	content += fmt.Sprintf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04"))
	content += "Source: PAI Universal Session Analyzer\n\n"

	for _, insight := range insights {
		content += fmt.Sprintf("## %s: %s\n", insight.Type, insight.Category)
		content += fmt.Sprintf("**Confidence**: %.0f%%\n\n", insight.Description)
		content += "**Evidence:**\n"
		for _, e := range insight.Evidence {
			content += fmt.Sprintf("- %s\n", e)
		}
		content += "\n---\n\n"
	}

	dest := filepath.Join(opencodeMemoryPath, filename)
	if u.dryRun {
		fmt.Printf("Would update opencode memory: %s\n", dest)
	} else {
		return os.WriteFile(dest, []byte(content), 0644)
	}
	return nil
}
