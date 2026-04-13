package main

import (
	"flag"
	"fmt"
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

func main() {
	flag.BoolVar(&flagDryRun, "dry-run", false, "Show what would be synced without making changes")
	flag.BoolVar(&flagVerbose, "v", false, "Verbose output")
	flag.StringVar(&flagSource, "source", "", "Source category (all, beliefs, goals, projects, learnings, skills)")
	flag.StringVar(&flagTarget, "target", "", "Target harness (claude, opencode, codex, cursor)")
	flag.Parse()

	baseDir := getBaseDir()
	items := scanSyncableItems(baseDir)

	if flagDryRun && flagSource != "" && flagTarget != "" {
		showSyncPlan(items, flagSource, flagTarget, baseDir)
		return
	}

	if err := runTUI(items, baseDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getBaseDir() string {
	if dir := os.Getenv("PAI_BASE_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "src", "Code", "pai-universal")
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

func showSyncPlan(items []SyncItem, source, target, baseDir string) {
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  SYNC PLAN")
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("  Source: %s\n", source)
	fmt.Printf("  Target: %s\n", target)
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

func filterItems(items []SyncItem, source string) []SyncItem {
	if source == "all" {
		return items
	}
	var filtered []SyncItem
	for _, item := range items {
		if item.Type == source {
			filtered = append(filtered, item)
		}
	}
	return filtered
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

		// Title
		title := " PAI Sync - Select items to sync (space=toggle, arrows=navigate, enter=sync, q=quit) "
		for x, r := range title {
			if x < 80 {
				screen.SetCell(x, 0, titleStyle, r)
			}
		}

		// Items
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

		// Help
		help := " Space: toggle | Enter: sync | Q: quit "
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
				if ev.Rune() == 'q' || ev.Rune() == 'Q' {
					quit = true
				} else if ev.Rune() == ' ' {
					items[selectedIndex].Selected = !items[selectedIndex].Selected
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
