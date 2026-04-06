package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("PAI Universal - Personal AI Infrastructure")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("CLI wrapper for PAI with universal AI tool support.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  pai-universal              # Launch with default AI tool")
	fmt.Println("  pai-universal opencode     # Launch with opencode")
	fmt.Println("  pai-universal install      # Run installation wizard")
	fmt.Println("  pai-universal init         # Initialize PAI in current directory")
	fmt.Println("  pai-universal backup       # Backup PAI data")
	fmt.Println("  pai-universal restore      # Restore from backup")
	fmt.Println("  pai-universal tool [name]  # Switch default AI tool")
	fmt.Println()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			fmt.Println("Running installation wizard...")
			os.Exit(0)
		case "init":
			fmt.Println("Initializing PAI in current directory...")
			os.Exit(0)
		case "backup":
			fmt.Println("Backing up PAI data...")
			os.Exit(0)
		case "restore":
			fmt.Println("Restoring from backup...")
			os.Exit(0)
		case "tool":
			if len(os.Args) > 2 {
				fmt.Printf("Setting default tool to: %s\n", os.Args[2])
			} else {
				fmt.Println("Usage: pai-universal tool [opencode|cursor|codex|gemini]")
			}
			os.Exit(0)
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			os.Exit(1)
		}
	}
}
