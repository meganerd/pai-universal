# yt-transcript

Get YouTube video transcript or transcribe via Whisper.

## Installation

```bash
go install ./tools/yt-transcript
```

Or build locally:
```bash
go build -o yt-transcript ./tools/yt-transcript
```

## Usage

```bash
# Get transcript (prefers YouTube captions, falls back to Whisper)
./yt-transcript --url "https://youtube.com/watch?v=..."

# Output to file (use -o, not --output)
./yt-transcript --url "URL" -o transcript.txt

# Use OpenAI provider (requires OPENAI_API_KEY)
./yt-transcript --url "URL" --provider openai

# Use local Whisper (default - requires whisper CLI installed)
./yt-transcript --url "URL" --provider local --model base

# Verbose
./yt-transcript --url "URL" -v
```

## Providers

- `openai` (default): Uses OpenAI Whisper API, requires `OPENAI_API_KEY` env var
- `local`: Uses `whisper` CLI (must be installed separately)

## Requirements

- `yt-dlp` - for downloading audio
- For local provider: `whisper` CLI (https://github.com/openai/whisper)
- For OpenAI provider: `OPENAI_API_KEY` env var

## Priority

1. Try YouTube auto-generated captions
2. Try YouTube manually created captions  
3. Fall back to Whisper transcription

## Flags

- `--url` - YouTube URL (required)
- `-o` - Output file (default: stdout)
- `-v` - Verbose output
- `--model` - Local Whisper model (default: base)
- `--provider` - Provider: local (default), openai