package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var (
	flagDryRun  bool
	flagSource  string
	flagTarget  string
	flagVerbose bool
)

type SyncItem struct {
	Path     string
	Selected bool
	Type     string
}

type HarnessTarget struct {
	Name       string
	ConfigPath string
	MemoryPath string
	SkillsPath string
}

var harnessTargets = map[string]HarnessTarget{
	"claude": {
		Name:       "Claude Code",
		ConfigPath: "~/.claude/",
		MemoryPath: "~/.claude/projects/",
		SkillsPath: "~/.claude/skills/",
	},
	"opencode": {
		Name:       "OpenCode",
		ConfigPath: "~/.config/opencode/",
		MemoryPath: "~/.local/share/opencode/",
		SkillsPath: "~/.config/opencode/skills/",
	},
	"codex": {
		Name:       "Codex",
		ConfigPath: "~/.codex/",
		MemoryPath: "~/.codex/memories/",
		SkillsPath: "~/.codex/skills/",
	},
	"cursor": {
		Name:       "Cursor",
		ConfigPath: "~/.cursor/",
		MemoryPath: "~/.cursor/memories/",
		SkillsPath: "~/.cursor/rules/",
	},
	"pigo": {
		Name:       "pi-go",
		ConfigPath: "~/src/Code/pi-go/",
		MemoryPath: "~/src/Code/pi-go/memory/",
		SkillsPath: "~/src/Code/pi-go/skills/",
	},
	"gemini": {
		Name:       "Gemini CLI",
		ConfigPath: "~/.gemini/",
		MemoryPath: "~/.gemini/memories/",
		SkillsPath: "~/.gemini/skills/",
	},
	"goose": {
		Name:       "Goose",
		ConfigPath: "~/.config/goose/",
		MemoryPath: "~/.config/goose/memory/",
		SkillsPath: "~/.config/goose/skills/",
	},
}

func main() {
	flag.BoolVar(&flagDryRun, "dry-run", false, "Show what would be synced without making changes")
	flag.BoolVar(&flagVerbose, "v", false, "Verbose output")
	flag.StringVar(&flagSource, "source", "", "Source category (all, beliefs, goals, projects, learnings, skills)")
	flag.StringVar(&flagTarget, "target", "", "Target harness (claude, opencode, codex, cursor, pigo, gemini, goose)")
	flag.Parse()

	baseDir := getBaseDir()
	items := scanSyncableItems(baseDir)

	// Handle CLI mode with target
	if flagTarget != "" {
		showItemsList(items, baseDir)
		syncToHarnessByName(items, flagTarget, baseDir)
		return
	}

	// Handle non-interactive mode
	if !isInteractive() {
		showItemsList(items, baseDir)
		fmt.Println("\nRun with -target <harness> to sync to specific harness")
		fmt.Println("Available harnesses: claude, opencode, codex, cursor, pigo, gemini, goose")
		return
	}

	if flagVerbose {
		log.Printf("Found %d syncable items\n", len(items))
		log.Printf("Available harnesses: ")
		for name := range harnessTargets {
			log.Printf("%s ", name)
		}
		log.Println()
	}

	if err := runTUI(items, baseDir); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func getBaseDir() string {
	if dir := os.Getenv("PAI_BASE_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "src", "Code", "pai-universal")
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

func isInteractive() bool {
	// Check if stdout is a terminal
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

func scanSyncableItems(baseDir string) []SyncItem {
	var items []SyncItem

	telosDir := filepath.Join(baseDir, "USER", "TELOS")
	if dirs, err := os.ReadDir(telosDir); err == nil {
		for _, d := range dirs {
			if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
				continue
			}
			path := filepath.Join(telosDir, d.Name())
			var itemType string
			switch {
			case strings.Contains(d.Name(), "belief"):
				itemType = "belief"
			case strings.Contains(d.Name(), "goal"):
				itemType = "goal"
			case strings.Contains(d.Name(), "project"):
				itemType = "project"
			case strings.Contains(d.Name(), "mission"):
				itemType = "mission"
			default:
				itemType = "other"
			}
			items = append(items, SyncItem{Path: path, Type: itemType})
		}
	}

	coldDir := filepath.Join(baseDir, "MEMORY", "cold")
	if dirs, err := os.ReadDir(coldDir); err == nil {
		for _, d := range dirs {
			if d.IsDir() {
				items = append(items, SyncItem{
					Path:     filepath.Join(coldDir, d.Name()),
					Type:     "learning",
					Selected: true,
				})
			}
		}
	}

	skillsDir := filepath.Join(baseDir, "skills")
	if dirs, err := os.ReadDir(skillsDir); err == nil {
		for _, d := range dirs {
			if d.IsDir() {
				items = append(items, SyncItem{
					Path:     filepath.Join(skillsDir, d.Name()),
					Type:     "skill",
					Selected: false,
				})
			}
		}
	}

	return items
}

func filterItems(items []SyncItem, source string) []SyncItem {
	if source == "all" || source == "" {
		return items
	}
	// Handle plural forms (skills -> skill, beliefs -> belief, etc.)
	sourceSingular := strings.TrimSuffix(source, "s")
	var filtered []SyncItem
	for _, item := range items {
		if item.Type == source || item.Type == sourceSingular {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func showItemsList(items []SyncItem, baseDir string) {
	filtered := filterItems(items, flagSource)
	fmt.Println("Available syncable items:")
	if len(filtered) == 0 {
		fmt.Println("  (none)")
		return
	}
	for _, item := range filtered {
		relPath, _ := filepath.Rel(baseDir, item.Path)
		checkbox := "[ ]"
		if item.Selected {
			checkbox = "[x]"
		}
		fmt.Printf("  %s %s (%s)\n", checkbox, relPath, item.Type)
	}
}

func syncToHarnessByName(items []SyncItem, targetName, baseDir string) {
	target, ok := harnessTargets[targetName]
	if !ok {
		fmt.Printf("Unknown target: %s\n", targetName)
		fmt.Printf("Available: claude, opencode, codex, cursor\n")
		return
	}

	var selected []SyncItem
	for _, item := range items {
		if item.Selected {
			selected = append(selected, item)
		}
	}

	if len(selected) == 0 {
		// If nothing selected, sync all filtered items
		selected = filterItems(items, flagSource)
		if len(selected) == 0 {
			fmt.Println("No items to sync")
			return
		}
	}

	fmt.Printf("\nSyncing %d items to %s (%s):\n", len(selected), target.Name, target.ConfigPath)

	targetPath := expandPath(target.ConfigPath)

	for _, item := range selected {
		relPath, _ := filepath.Rel(baseDir, item.Path)
		destPath := filepath.Join(targetPath, filepath.Base(relPath))

		if flagDryRun {
			fmt.Printf("  [dry-run] Would copy: %s → %s\n", relPath, destPath)
		} else {
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				fmt.Printf("  ✗ Failed to create directory: %v\n", err)
				continue
			}

			data, err := os.ReadFile(item.Path)
			if err != nil {
				fmt.Printf("  ✗ Failed to read: %v\n", err)
				continue
			}

			if err := os.WriteFile(destPath, data, 0644); err != nil {
				fmt.Printf("  ✗ Failed to write: %v\n", err)
				continue
			}

			fmt.Printf("  ✓ Synced: %s → %s\n", relPath, destPath)
		}
	}

	if flagDryRun {
		fmt.Println("\n(dry-run - no changes made)")
	} else {
		fmt.Println("\nSync complete!")
	}
}

func showSyncPlan(items []SyncItem, source, target, baseDir string) {
	targetInfo, ok := harnessTargets[target]
	if !ok {
		log.Printf("Unknown target harness: %s", target)
		log.Printf("Available: claude, opencode, codex, cursor")
		return
	}

	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  SYNC PLAN")
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("  Source:     %s\n", source)
	fmt.Printf("  Target:     %s (%s)\n", targetInfo.Name, targetInfo.ConfigPath)
	fmt.Println("═══════════════════════════════════════")

	filtered := filterItems(items, source)
	fmt.Printf("\n%d items would be synced:\n\n", len(filtered))

	for _, item := range filtered {
		relPath, _ := filepath.Rel(baseDir, item.Path)
		checkbox := "[ ]"
		if item.Selected {
			checkbox = "[x]"
		}
		fmt.Printf("  %s %s (%s)\n", checkbox, relPath, item.Type)
	}

	if flagDryRun {
		fmt.Println("\n(dry-run - no changes made)")
	}
}

func runTUI(items []SyncItem, baseDir string) error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		return err
	}

	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	titleStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite).Bold(true)
	selectedStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)

	screen.SetStyle(style)

	selectedIndex := 0
	quit := false

	for !quit {
		screen.Clear()

		title := " PAI Sync - Select items (space=toggle, arrows=navigate, t=target, enter=sync, q=quit) "
		for x, r := range title {
			if x < 80 {
				screen.SetCell(x, 0, titleStyle, r)
			}
		}

		for i, item := range items {
			relPath, _ := filepath.Rel(baseDir, item.Path)
			checkbox := "[ ]"
			if item.Selected {
				checkbox = "[x]"
			}
			displayLine := fmt.Sprintf("  %s %-40s (%s)", checkbox, filepath.Base(relPath), item.Type)

			rowStyle := style
			if i == selectedIndex {
				rowStyle = selectedStyle
			}

			for x, r := range displayLine {
				if x < 80 {
					screen.SetCell(x, i+2, rowStyle, r)
				}
			}
		}

		help := " Space: toggle | T: select target | Enter: sync | Q: quit "
		_, helpY := screen.Size()
		helpY = helpY - 1
		for x, r := range help {
			if x < 80 {
				screen.SetCell(x, helpY, titleStyle, r)
			}
		}

		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				if selectedIndex > 0 {
					selectedIndex--
				}
			case tcell.KeyDown:
				if selectedIndex < len(items)-1 {
					selectedIndex++
				}
			case tcell.KeyEnter:
				syncSelected(items, baseDir)
				quit = true
			case tcell.KeyEsc:
				quit = true
			default:
				r := ev.Rune()
				if r == 'q' || r == 'Q' {
					quit = true
				} else if r == ' ' {
					items[selectedIndex].Selected = !items[selectedIndex].Selected
				} else if r == 't' || r == 'T' {
					// Show target selection
					quit = true
					runTargetSelector(baseDir, items)
				}
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}

	return nil
}

func runTargetSelector(baseDir string, items []SyncItem) error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		return err
	}

	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	titleStyle := tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorWhite).Bold(true)
	selectedStyle := tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorWhite)

	screen.SetStyle(style)

	// Build list of targets
	var targetNames []string
	for name := range harnessTargets {
		targetNames = append(targetNames, name)
	}

	selectedIndex := 0
	quit := false

	for !quit {
		screen.Clear()

		title := " Select Target Harness (arrows=navigate, enter=select, q=cancel) "
		for x, r := range title {
			if x < 60 {
				screen.SetCell(x, 0, titleStyle, r)
			}
		}

		for i, name := range targetNames {
			target := harnessTargets[name]
			displayLine := fmt.Sprintf("  %-10s %-20s", name, target.Name)

			rowStyle := style
			if i == selectedIndex {
				rowStyle = selectedStyle
			}

			for x, r := range displayLine {
				if x < 40 {
					screen.SetCell(x, i+2, rowStyle, r)
				}
			}
		}

		help := " Enter: sync to selected | Q: cancel "
		_, helpY := screen.Size()
		helpY = helpY - 1
		for x, r := range help {
			if x < 60 {
				screen.SetCell(x, helpY, titleStyle, r)
			}
		}

		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				if selectedIndex > 0 {
					selectedIndex--
				}
			case tcell.KeyDown:
				if selectedIndex < len(targetNames)-1 {
					selectedIndex++
				}
			case tcell.KeyEnter:
				targetName := targetNames[selectedIndex]
				syncToHarness(items, targetName, baseDir)
				quit = true
			case tcell.KeyEsc:
				quit = true
			default:
				if ev.Rune() == 'q' || ev.Rune() == 'Q' {
					quit = true
				}
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}

	return nil
}

func syncSelected(items []SyncItem, baseDir string) {
	var selected []SyncItem
	for _, item := range items {
		if item.Selected {
			selected = append(selected, item)
		}
	}

	if len(selected) == 0 {
		fmt.Println("\nNo items selected")
		return
	}

	fmt.Println("\nSyncing selected items:")
	for _, item := range selected {
		relPath, _ := filepath.Rel(baseDir, item.Path)
		fmt.Printf("  → %s (%s)\n", relPath, item.Type)
	}
	fmt.Println("\nSync complete!")
}

func syncToHarness(items []SyncItem, targetName, baseDir string) {
	target, ok := harnessTargets[targetName]
	if !ok {
		fmt.Printf("Unknown target: %s\n", targetName)
		return
	}

	var selected []SyncItem
	for _, item := range items {
		if item.Selected {
			selected = append(selected, item)
		}
	}

	if len(selected) == 0 {
		fmt.Println("\nNo items selected")
		return
	}

	fmt.Printf("\nSyncing %d items to %s (%s):\n", len(selected), target.Name, target.ConfigPath)

	targetPath := expandPath(target.ConfigPath)

	for _, item := range selected {
		relPath, _ := filepath.Rel(baseDir, item.Path)
		destPath := filepath.Join(targetPath, filepath.Base(relPath))

		if flagDryRun {
			fmt.Printf("  [dry-run] Would copy: %s → %s\n", relPath, destPath)
		} else {
			// Ensure directory exists
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				fmt.Printf("  ✗ Failed to create directory: %v\n", err)
				continue
			}

			// Copy file
			data, err := os.ReadFile(item.Path)
			if err != nil {
				fmt.Printf("  ✗ Failed to read: %v\n", err)
				continue
			}

			if err := os.WriteFile(destPath, data, 0644); err != nil {
				fmt.Printf("  ✗ Failed to write: %v\n", err)
				continue
			}

			fmt.Printf("  ✓ Synced: %s → %s\n", relPath, destPath)
		}
	}

	if flagDryRun {
		fmt.Println("\n(dry-run - no changes made)")
	} else {
		fmt.Println("\nSync complete!")
	}
}
