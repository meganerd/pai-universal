package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/meganerd/pai-universal/cmd/session-analyzer/internal/config"
	"github.com/meganerd/pai-universal/cmd/session-analyzer/internal/llm"
	"github.com/meganerd/pai-universal/cmd/session-analyzer/internal/parser"
	"github.com/meganerd/pai-universal/cmd/session-analyzer/internal/updater"

	siftrank "github.com/meganerd/siftrank/pkg/siftrank"
)

var (
	flagConfig         string
	flagDryRun         bool
	flagVerbose        bool
	flagProvider       string
	flagModel          string
	flagUseSiftrank    bool
	flagUpdateAll      bool
	flagUpdateClaude   bool
	flagUpdateOpencode bool
	flagUpdateCodex    bool
	flagUpdateCursor   bool
	flagUpdateGemini   bool
	flagUpdatePigo     bool
)

func main() {
	flag.StringVar(&flagConfig, "c", "", "Config file path")
	flag.BoolVar(&flagDryRun, "dry-run", false, "Show what would be done without making changes")
	flag.BoolVar(&flagVerbose, "v", false, "Verbose output")
	flag.StringVar(&flagProvider, "provider", "openrouter", "LLM provider (anthropic, openrouter)")
	flag.StringVar(&flagModel, "model", "google/gemini-2.0-flash-001", "LLM model (defaults to provider-specific)")
	flag.BoolVar(&flagUseSiftrank, "siftrank", false, "Use siftrank to auto-select optimal model")
	flag.BoolVar(&flagUpdateAll, "all", true, "Update all harnesses (default: true)")
	flag.BoolVar(&flagUpdateClaude, "claude", false, "Update Claude Code memory")
	flag.BoolVar(&flagUpdateOpencode, "opencode", false, "Update opencode memory")
	flag.BoolVar(&flagUpdateCodex, "codex", false, "Update Codex memory")
	flag.BoolVar(&flagUpdateCursor, "cursor", false, "Update Cursor memory")
	flag.BoolVar(&flagUpdateGemini, "gemini", false, "Update Gemini CLI memory")
	flag.BoolVar(&flagUpdatePigo, "pigo", false, "Update pi-go memory")

	cfg, err := config.Load(flagConfig)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override from flags
	if flagProvider != "" {
		cfg.LLMProvider = flagProvider
	}
	if flagModel != "" {
		cfg.LLMModel = flagModel
	}

	flag.Parse()

	// Use siftrank to select optimal model if requested
	if flagUseSiftrank {
		selectedModel, err := selectOptimalModel(cfg.LLMModel)
		if err != nil {
			log.Printf("Siftrank selection failed, using default: %v", err)
		} else {
			if flagVerbose {
				fmt.Printf("Siftrank selected: %s (was: %s)\n", selectedModel, cfg.LLMModel)
			}
			cfg.LLMModel = selectedModel
		}
	}

	if flagVerbose {
		fmt.Println("Session Analyzer v0.1.0")
		fmt.Printf("Config: %+v\n", cfg)
	}

	baseDir := os.Getenv("PAI_BASE_DIR")
	if baseDir == "" {
		home, _ := os.UserHomeDir()
		baseDir = filepath.Join(home, "src", "Code", "pai-universal")
	}

	sessions, err := parser.ParseAll(cfg)
	if err != nil {
		log.Fatalf("Failed to parse sessions: %v", err)
	}

	if flagVerbose {
		fmt.Printf("Parsed %d sessions\n", len(sessions))
	}

	analyzer := llm.NewAnalyzer(cfg)
	insights, err := analyzer.Analyze(sessions)
	if err != nil {
		log.Fatalf("Failed to analyze sessions: %v", err)
	}

	if flagVerbose {
		fmt.Printf("Found %d insights\n", len(insights))
		for _, i := range insights {
			fmt.Printf("  - %s: %s\n", i.Type, i.Description)
		}
	}

	// Determine which harnesses to update
	var targets []string
	if flagUpdateAll {
		// Default: update pai-universal + all available harnesses
		targets = []string{"all"}
	} else {
		// Only update specified harnesses
		if flagUpdateClaude {
			targets = append(targets, "claude")
		}
		if flagUpdateOpencode {
			targets = append(targets, "opencode")
		}
		if flagUpdateCodex {
			targets = append(targets, "codex")
		}
		if flagUpdateCursor {
			targets = append(targets, "cursor")
		}
		if flagUpdateGemini {
			targets = append(targets, "gemini")
		}
		if flagUpdatePigo {
			targets = append(targets, "pigo")
		}
	}

	updater := updater.NewUpdater(baseDir, flagDryRun)
	updater.SetTargets(targets)

	if flagVerbose {
		fmt.Printf("Updating harnesses: %v\n", targets)
	}

	if err := updater.Apply(insights); err != nil {
		log.Fatalf("Failed to apply insights: %v", err)
	}

	fmt.Println("Session analysis complete")
}

// selectOptimalModel uses siftrank to recommend the best model for session analysis
func selectOptimalModel(baselineModel string) (string, error) {
	catalog := siftrank.NewModelCatalog()

	// Add models we know about (in practice, siftrank has a built-in catalog)
	models := []siftrank.ModelInfo{
		{ID: "google/gemini-2.0-flash-001", Provider: siftrank.ProviderTypeOpenRouter, ContextWindow: 1000000, Pricing: siftrank.ModelPricing{InputPricePerToken: 0.000000, OutputPricePerToken: 0.000000}},
		{ID: "google/gemini-2.5-pro-preview-05-20", Provider: siftrank.ProviderTypeOpenRouter, ContextWindow: 1000000, Pricing: siftrank.ModelPricing{InputPricePerToken: 0.000000, OutputPricePerToken: 0.000000}},
		{ID: "anthropic/claude-sonnet-4-20250514", Provider: siftrank.ProviderTypeOpenRouter, ContextWindow: 200000, Pricing: siftrank.ModelPricing{InputPricePerToken: 0.000003, OutputPricePerToken: 0.000015}},
		{ID: "anthropic/claude-3-5-sonnet-20241022", Provider: siftrank.ProviderTypeOpenRouter, ContextWindow: 200000, Pricing: siftrank.ModelPricing{InputPricePerToken: 0.000003, OutputPricePerToken: 0.000015}},
		{ID: "openai/gpt-4o-mini", Provider: siftrank.ProviderTypeOpenRouter, ContextWindow: 128000, Pricing: siftrank.ModelPricing{InputPricePerToken: 0.00000015, OutputPricePerToken: 0.0000006}},
		{ID: "openai/gpt-4o", Provider: siftrank.ProviderTypeOpenRouter, ContextWindow: 128000, Pricing: siftrank.ModelPricing{InputPricePerToken: 0.0000025, OutputPricePerToken: 0.00001}},
	}

	for _, m := range models {
		catalog.Add(m)
	}

	// Session analysis is typically ~5000 tokens in, ~4000 tokens out
	// Use a reasonable baseline for comparison
	baseline := baselineModel
	if baseline == "" {
		baseline = "anthropic/claude-3-5-sonnet-20241022"
	}

	recs := siftrank.Recommend(5000, 4000, baseline, catalog, nil)

	if len(recs) == 0 {
		return baselineModel, fmt.Errorf("no recommendations available")
	}

	// Return the cheapest model (first in sorted list)
	return recs[0].ModelID, nil
}
