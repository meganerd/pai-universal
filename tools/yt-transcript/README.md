# yt-transcript

Get YouTube video transcript or transcribe via Whisper.

## Installation

```bash
pip install --user faster-whisper openai requests
```

## Usage

```bash
# Get transcript (prefers YouTube captions, falls back to Whisper)
python yt_transcript.py "https://youtube.com/watch?v=..."

# Specify output file
python yt_transcript.py "URL" -o transcript.txt

# Use different provider
python yt_transcript.py "URL" --provider openai    # OpenAI Whisper API
python yt_transcript.py "URL" --provider local     # Local faster-whisper (default)
python yt_transcript.py "URL" --provider google     # Google Cloud Speech-to-Text
python yt_transcript.py "URL" --provider assemblyai  # AssemblyAI
```

## Providers

- `local` (default): faster-whisper - runs locally, free once model downloaded
- `openai`: OpenAI Whisper API - paid per minute, requires OPENAI_API_KEY
- `google`: Google Cloud Speech-to-Text - paid, requires GOOGLE_APPLICATION_CREDENTIALS
- `assemblyai`: AssemblyAI - paid, requires ASSEMBLYAI_API_KEY

## Environment Variables

- `OPENAI_API_KEY` - for OpenAI provider
- `GOOGLE_APPLICATION_CREDENTIALS` - for Google provider  
- `ASSEMBLYAI_API_KEY` - for AssemblyAI provider

## Priority

1. Try YouTube auto-generated captions
2. Try YouTube manually created captions
3. Fall back to Whisper transcription (configurable provider)