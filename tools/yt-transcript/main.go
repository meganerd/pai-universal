// yt-transcript gets YouTube video transcript or transcribes via Whisper
//
// Usage:
//
//	yt-transcript "https://youtube.com/watch?v=..."
//	yt-transcript -o transcript.txt "URL"
//
// Priority:
//  1. Try YouTube captions (auto-generated preferred)
//  2. Fall back to Whisper transcription
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	flagURL      string
	flagOutput   string
	flagVerbose  bool
	flagModel    string
	flagProvider string
)

func main() {
	flag.StringVar(&flagURL, "url", "", "YouTube URL (required)")
	flag.StringVar(&flagOutput, "o", "", "Output file (default: stdout)")
	flag.BoolVar(&flagVerbose, "v", false, "Verbose output")
	flag.StringVar(&flagModel, "model", "base", "Whisper model (tiny, base, small, medium, large)")
	flag.StringVar(&flagProvider, "provider", "openai", "Provider: local, openai")
	flag.Parse()

	if flagURL == "" {
		log.Fatal("Error: --url is required")
	}

	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	videoID := extractVideoID(flagURL)
	if videoID == "" {
		return fmt.Errorf("invalid YouTube URL: %s", flagURL)
	}

	if flagVerbose {
		fmt.Printf("[yt-transcript] Video ID: %s\n", videoID)
	}

	subtitle, err := getYouTubeSubtitles(flagURL)
	if err == nil && subtitle != "" {
		if flagVerbose {
			fmt.Printf("[yt-transcript] Found YouTube subtitles\n")
		}
		return writeOutput(subtitle)
	}

	if flagVerbose {
		fmt.Printf("[yt-transcript] No subtitles found, downloading audio...\n")
	}

	audioPath, err := downloadAudio(flagURL)
	if err != nil {
		return fmt.Errorf("download audio: %w", err)
	}
	defer os.Remove(audioPath)

	var transcript string
	switch flagProvider {
	case "local":
		transcript, err = transcribeLocal(audioPath)
	case "openai":
		transcript, err = transcribeOpenAI(audioPath)
	default:
		return fmt.Errorf("unknown provider: %s", flagProvider)
	}

	if err != nil {
		return fmt.Errorf("transcribe: %w", err)
	}

	return writeOutput(transcript)
}

func extractVideoID(url string) string {
	patterns := []string{
		`(?:v=|\.be/)([a-zA-Z0-9_-]{11})`,
		`(?:embed/)([a-zA-Z0-9_-]{11})`,
	}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		m := re.FindStringSubmatch(url)
		if len(m) > 1 {
			return m[1]
		}
	}
	return ""
}

func getYouTubeSubtitles(url string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "yt-subs-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("yt-dlp",
		"--write-subs", "--write-auto-subs",
		"--sub-lang", "en",
		"--skip-download",
		"--output", filepath.Join(tmpDir, "subs"),
		url,
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	_ = cmd.Run()

	for _, ext := range []string{".en.vtt", ".vtt", ".en.srt", ".srt"} {
		files, _ := filepath.Glob(filepath.Join(tmpDir, "subs*"+ext))
		for _, f := range files {
			if data, err := os.ReadFile(f); err == nil && len(data) > 50 {
				return string(data), nil
			}
		}
	}
	return "", nil
}

func downloadAudio(url string) (string, error) {
	videoID := extractVideoID(url)
	if videoID == "" {
		return "", fmt.Errorf("invalid URL")
	}

	tempTemplate := fmt.Sprintf("/tmp/yt-audio-%s.%%(ext)s", videoID)

	cmd := exec.Command("yt-dlp",
		"-f", "bestaudio",
		"--output", tempTemplate,
		"--force-overwrite",
		"--print", "filename",
		url,
	)
	output, _ := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	audioPath := strings.TrimSpace(lines[len(lines)-1])

	if audioPath == "" || !strings.HasPrefix(audioPath, "/tmp") {
		fmt.Printf("[yt-transcript] yt-dlp output: %s\n", string(output))
	}
	return "", fmt.Errorf("no audio file created for video: %s", videoID)
}

func transcribeLocal(audioPath string) (string, error) {
	if flagVerbose {
		fmt.Printf("[yt-transcript] Transcribing with local Whisper (%s)...\n", flagModel)
	}

	cmd := exec.Command("whisper",
		"--model", flagModel,
		"--language", "en",
		audioPath,
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	parts := strings.Split(audioPath, ".")
	parts[len(parts)-1] = "txt"
	txtPath := strings.Join(parts, ".")

	data, err := os.ReadFile(txtPath)
	if err != nil {
		return "", err
	}

	os.Remove(txtPath)
	return string(data), nil
}

func transcribeOpenAI(audioPath string) (string, error) {
	if flagVerbose {
		fmt.Printf("[yt-transcript] Transcribing with OpenAI Whisper...\n")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	cmd := exec.Command("openai", "audio", "transcriptions", "create",
		"--file", audioPath,
		"--model", "whisper-1",
		"--response-format", "text",
		"--api-key", apiKey,
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	data, err := os.ReadFile(audioPath[:len(audioPath)-len(filepath.Ext(audioPath))] + ".txt")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func writeOutput(content string) error {
	if flagOutput == "" {
		fmt.Println(content)
		return nil
	}
	return os.WriteFile(flagOutput, []byte(content), 0644)
}
