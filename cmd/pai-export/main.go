package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	flagDryRun  bool
	flagVerbose bool
	flagOutput  string
	flagInclude string
)

func main() {
	flag.BoolVar(&flagDryRun, "dry-run", false, "Show what would be exported without creating archive")
	flag.BoolVar(&flagVerbose, "v", false, "Verbose output")
	flag.StringVar(&flagOutput, "o", "", "Output file (default: PAI-export-YYYYMMDD.tar.gz)")
	flag.StringVar(&flagInclude, "include", "full", "What to include: full, infra, skills, learnings, telos")
	flag.Parse()

	outputFile := flagOutput
	if outputFile == "" {
		outputFile = fmt.Sprintf("PAI-export-%s.tar.gz", time.Now().Format("20060102"))
	}
	if !strings.HasSuffix(outputFile, ".tar.gz") {
		outputFile += ".tar.gz"
	}

	items := collectExportItems()
	filtered := filterItems(items, flagInclude)

	if flagDryRun {
		fmt.Println("═══════════════════════════════════════")
		fmt.Println("  EXPORT PLAN (dry-run)")
		fmt.Println("═══════════════════════════════════════")
		fmt.Printf("  Include: %s\n", flagInclude)
		fmt.Printf("  Output:  %s\n", outputFile)
		fmt.Println("═══════════════════════════════════════")
		fmt.Printf("\n%d items would be exported:\n\n", len(filtered))
		for _, item := range filtered {
			fmt.Printf("  %s (%s)\n", item.Path, item.Type)
		}
		return
	}

	if err := createExportArchive(filtered, outputFile); err != nil {
		log.Fatalf("Export failed: %v", err)
	}

	fmt.Printf("Export complete: %s\n", outputFile)
	fmt.Printf("Items exported: %d\n", len(filtered))
}

type ExportItem struct {
	Path   string
	Type   string
	Source string
}

func collectExportItems() []ExportItem {
	var items []ExportItem
	home, _ := os.UserHomeDir()
	claudeDir := filepath.Join(home, ".claude")

	// PAI directory (infrastructure)
	paiDir := filepath.Join(claudeDir, "PAI")
	if dirs, err := os.ReadDir(paiDir); err == nil {
		for _, d := range dirs {
			path := filepath.Join(paiDir, d.Name())
			var itemType string
			if d.IsDir() {
				itemType = "pai-dir"
			} else {
				itemType = "pai-file"
			}
			items = append(items, ExportItem{Path: path, Type: itemType, Source: "PAI"})
		}
	}

	// PAI/USER/TELOS
	telosDir := filepath.Join(paiDir, "USER", "TELOS")
	if dirs, err := os.ReadDir(telosDir); err == nil {
		for _, d := range dirs {
			if d.IsDir() || strings.HasSuffix(d.Name(), ".md") {
				path := filepath.Join(telosDir, d.Name())
				items = append(items, ExportItem{Path: path, Type: "telos", Source: "PAI/USER/TELOS"})
			}
		}
	}

	// PAI/USER/* (other USER subdirs)
	userDir := filepath.Join(paiDir, "USER")
	if dirs, err := os.ReadDir(userDir); err == nil {
		for _, d := range dirs {
			if d.Name() == "TELOS" {
				continue
			}
			path := filepath.Join(userDir, d.Name())
			if d.IsDir() {
				items = append(items, ExportItem{Path: path, Type: "user-dir", Source: "PAI/USER"})
			} else if strings.HasSuffix(d.Name(), ".md") {
				items = append(items, ExportItem{Path: path, Type: "user-file", Source: "PAI/USER"})
			}
		}
	}

	// Skills
	skillsDir := filepath.Join(claudeDir, "skills")
	if dirs, err := os.ReadDir(skillsDir); err == nil {
		for _, d := range dirs {
			path := filepath.Join(skillsDir, d.Name())
			if d.IsDir() {
				items = append(items, ExportItem{Path: path, Type: "skill", Source: "skills"})
			}
		}
	}

	// MEMORY/LEARNING
	memLearnDir := filepath.Join(claudeDir, "MEMORY", "LEARNING")
	if dirs, err := os.ReadDir(memLearnDir); err == nil {
		for _, d := range dirs {
			path := filepath.Join(memLearnDir, d.Name())
			if d.IsDir() {
				items = append(items, ExportItem{Path: path, Type: "learning", Source: "MEMORY/LEARNING"})
			}
		}
	}

	// Hooks
	hooksDir := filepath.Join(claudeDir, "hooks")
	if dirs, err := os.ReadDir(hooksDir); err == nil {
		for _, d := range dirs {
			path := filepath.Join(hooksDir, d.Name())
			if d.IsDir() {
				items = append(items, ExportItem{Path: path, Type: "hook", Source: "hooks"})
			} else if strings.HasSuffix(d.Name(), ".ts") || strings.HasSuffix(d.Name(), ".sh") {
				items = append(items, ExportItem{Path: path, Type: "hook-file", Source: "hooks"})
			}
		}
	}

	return items
}

func filterItems(items []ExportItem, include string) []ExportItem {
	switch include {
	case "full":
		return items
	case "infra":
		var filtered []ExportItem
		for _, item := range items {
			if item.Source == "PAI" || item.Source == "hooks" {
				filtered = append(filtered, item)
			}
		}
		return filtered
	case "skills":
		var filtered []ExportItem
		for _, item := range items {
			if item.Source == "skills" {
				filtered = append(filtered, item)
			}
		}
		return filtered
	case "learnings":
		var filtered []ExportItem
		for _, item := range items {
			if item.Source == "MEMORY/LEARNING" {
				filtered = append(filtered, item)
			}
		}
		return filtered
	case "telos":
		var filtered []ExportItem
		for _, item := range items {
			if item.Type == "telos" {
				filtered = append(filtered, item)
			}
		}
		return filtered
	default:
		return items
	}
}

func createExportArchive(items []ExportItem, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	for _, item := range items {
		if err := addToArchive(tw, item.Path); err != nil {
			if flagVerbose {
				log.Printf("Warning: could not add %s: %v", item.Path, err)
			}
		}
	}

	return nil
}

func addToArchive(tw *tar.Writer, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return addDirectoryToArchive(tw, path)
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = filepath.Base(path)

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(tw, f)
	return err
}

func addDirectoryToArchive(tw *tar.Writer, dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(filepath.Dir(dirPath), path)
		if err != nil {
			return err
		}

		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})
}
