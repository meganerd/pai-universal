#!/usr/bin/env python3
"""
yt-transcript: Get YouTube video transcript or transcribe via Whisper

Usage:
    yt-transcript "https://youtube.com/watch?v=..."
    yt-transcript "https://youtube.com/watch?v=..." -o transcript.txt
    yt-transcript "https://youtube.com/watch?v=..." --whisper large
    yt-transcript "https://youtube.com/watch?v=..." --provider openai --model whisper-1

Priority:
    1. Try YouTube captions (auto-generated preferred, then manually created)
    2. Fall back to Whisper transcription

Providers:
    - local: faster-whisper (local, free)
    - openai: OpenAI Whisper API (paid)
    - google: Google Cloud Speech-to-Text (paid)
    - assemblyai: AssemblyAI (paid)
"""

import argparse
import os
import re
import subprocess
import sys
import tempfile
from pathlib import Path

USER_SITE_PACKAGES = os.path.expanduser("~/.local/lib/python3.12/site-packages")
if USER_SITE_PACKAGES not in sys.path:
    sys.path.insert(0, USER_SITE_PACKAGES)

try:
    import faster_whisper
except ImportError:
    faster_whisper = None
    print("[yt-transcript] Warning: faster-whisper not available, install with: pip install --user faster-whisper", file=sys.stderr)

def extract_video_id(url: str) -> str:
    patterns = [
        r'(?:v=|\.be/)([a-zA-Z0-9_-]{11})',
        r'(?:embed/)([a-zA-Z0-9_-]{11})',
        r'(?:watch\?v=)([a-zA-Z0-9_-]{11})',
    ]
    for pattern in patterns:
        match = re.search(pattern, url)
        if match:
            return match.group(1)
    return None

def get_youtube_subtitles(url: str) -> str | None:
    """Try to get YouTube subtitles via yt-dlp"""
    with tempfile.TemporaryDirectory() as tmpdir:
        cmd = [
            "yt-dlp",
            "--write-subs", "--write-auto-subs",
            "--sub-lang", "en",
            "--skip-download",
            "--output", f"{tmpdir}/%(id)s.%(ext)s",
            url
        ]
        
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        except (subprocess.TimeoutExpired, FileNotFoundError):
            return None
        
        if result.returncode != 0:
            return None
        
        for ext in ["srt", "vtt", "txt"]:
            subtitle_file = Path(tmpdir)
            for f in subtitle_file.glob(f"*.{ext}"):
                try:
                    content = f.read_text(errors="ignore")
                    if content.strip():
                        return content
                except Exception:
                    continue
        
        return None

def transcribe_local(audio_path: str, model: str = "base") -> str:
    """Transcribe audio using faster-whisper (local)"""
    import sys
    sys.path.insert(0, os.path.expanduser("~/.local/lib/python3.12/site-packages"))
    
    try:
        import faster_whisper
    except ImportError:
        raise ImportError("faster-whisper not installed. Install with: pip install --user --breakage-packages faster-whisper")
    
    model_size = model if model in ["tiny", "base", "small", "medium", "large-v3", "large-v3-turbo"] else "base"
    
    print(f"[yt-transcript] Loading Whisper model: {model_size}", file=sys.stderr)
    whisper_model = faster_whisper.CTranslate2WhisperModel(model_size, download_root=os.path.expanduser("~/.cache/whisper"))
    
    segments, info = whisper_model.transcribe(audio_path, language="en", beam_size=5)
    
    print(f"[yt-transcript] Detected language: {info.language} (probability: {info.language_probability:.2f})", file=sys.stderr)
    
    text_parts = []
    for segment in segments:
        text_parts.append(segment.text.strip())
    
    return " ".join(text_parts)

def transcribe_openai(audio_path: str, model: str = "whisper-1") -> str:
    """Transcribe audio using OpenAI Whisper API"""
    try:
        from openai import OpenAI
    except ImportError:
        raise ImportError("openai not installed: pip install openai")
    
    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        raise ValueError("OPENAI_API_KEY not set")
    
    client = OpenAI(api_key=api_key)
    
    with open(audio_path, "rb") as audio_file:
        response = client.audio.transcriptions.create(
            model=model,
            file=audio_file,
            response_format="text"
        )
    
    return response.text.strip()

def transcribe_google(audio_path: str) -> str:
    """Transcribe audio using Google Cloud Speech-to-Text"""
    try:
        from google.cloud import speech
    except ImportError:
        raise ImportError("google-cloud-speech not installed: pip install google-cloud-speech")
    
    if not os.environ.get("GOOGLE_APPLICATION_CREDENTIALS"):
        raise ValueError("GOOGLE_APPLICATION_CREDENTIALS not set")
    
    client = speech.SpeechClient()
    
    with open(audio_path, "rb") as audio_file:
        content = audio_file.read()
    
    audio = speech.RecognitionAudio(content=content)
    config = speech.RecognitionConfig(
        encoding=speech.RecognitionConfig.AudioEncoding.MP3,
        language_code="en-US",
    )
    
    response = client.recognize(config=config, audio=audio)
    
    if not response.results:
        return ""
    
    return response.results[0].alternatives[0].transcript

def transcribe_assemblyai(audio_path: str) -> str:
    """Transcribe audio using AssemblyAI"""
    try:
        import requests
    except ImportError:
        raise ImportError("requests not installed: pip install requests")
    
    api_key = os.environ.get("ASSEMBLYAI_API_KEY")
    if not api_key:
        raise ValueError("ASSEMBLYAI_API_KEY not set")
    
    upload_url = "https://api.assemblyai.com/v3/upload"
    
    with open(audio_path, "rb") as f:
        response = requests.post(upload_url, files={"file": f}, headers={"Authorization": api_key})
        response.raise_for_status()
        audio_url = response.json()["upload_url"]
    
    transcript_response = requests.post(
        "https://api.assemblyai.com/v3/transcript",
        json={"audio_url": audio_url, "language_code": "en"},
        headers={"Authorization": api_key}
    )
    transcript_response.raise_for_status()
    transcript_id = transcript_response.json()["id"]
    
    while True:
        status_response = requests.get(
            f"https://api.assemblyai.com/v3/transcript/{transcript_id}",
            headers={"Authorization": api_key}
        )
        status = status_response.json()["status"]
        
        if status == "completed":
            break
        elif status == "error":
            raise RuntimeError("Transcription failed")
        
        import time
        time.sleep(3)
    
    final_response = requests.get(
        f"https://api.assemblyai.com/v3/transcript/{transcript_id}",
        headers={"Authorization": api_key}
    )
    
    return final_response.json()["text"]

def download_audio(url: str) -> str:
    """Download audio from YouTube"""
    video_id = extract_video_id(url)
    if not video_id:
        raise ValueError(f"Invalid YouTube URL: {url}")
    
    tempdir = tempfile.mkdtemp()
    output_path = os.path.join(tempdir, "audio.opus")
    
    cmd = [
        "yt-dlp",
        "-x", "--audio-format", "opus",
        "--output", output_path,
        "--keep-video",
        url
    ]
    
    result = subprocess.run(cmd, capture_output=True, text=True, cwd=tempdir)
    if result.returncode != 0:
        raise RuntimeError(f"Failed to download audio: {result.stderr}")
    
    files = list(Path(tempdir).glob("audio.opus"))
    if files:
        return str(files[0])
    
    webm_files = list(Path(tempdir).glob("audio.webm"))
    if webm_files:
        return str(webm_files[0])
    
    raise RuntimeError(f"Audio file not created in {tempdir}")

def main():
    parser = argparse.ArgumentParser(
        description="Get YouTube transcript or transcribe via Whisper",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    
    parser.add_argument("url", help="YouTube URL")
    parser.add_argument("-o", "--output", help="Output file (default: stdout)")
    parser.add_argument(
        "--provider",
        choices=["local", "openai", "google", "assemblyai"],
        default="local",
        help="Transcription provider (default: local)"
    )
    parser.add_argument(
        "--whisper-model",
        dest="whisper_model",
        default="base",
        help="Local Whisper model size or API model (default: base)"
    )
    parser.add_argument(
        "-v", "--verbose",
        action="store_true",
        help="Verbose output"
    )
    
    args = parser.parse_args()
    
    video_id = extract_video_id(args.url)
    if not video_id:
        print(f"[yt-transcript] Error: Invalid YouTube URL: {args.url}", file=sys.stderr)
        sys.exit(1)
    
    if args.verbose:
        print(f"[yt-transcript] Video ID: {video_id}", file=sys.stderr)
    
    subtitle = get_youtube_subtitles(args.url)
    
    if subtitle:
        if args.verbose:
            print("[yt-transcript] Found YouTube subtitles", file=sys.stderr)
        
        output = subtitle
    else:
        if args.verbose:
            print("[yt-transcript] No subtitles, downloading audio...", file=sys.stderr)
        
        audio_path = download_audio(args.url)
        
        if args.verbose:
            print(f"[yt-transcript] Transcribing with provider: {args.provider}", file=sys.stderr)
        
        if args.provider == "local":
            output = transcribe_local(audio_path, args.whisper_model)
        elif args.provider == "openai":
            output = transcribe_openai(audio_path, args.whisper_model)
        elif args.provider == "google":
            output = transcribe_google(audio_path)
        elif args.provider == "assemblyai":
            output = transcribe_assemblyai(audio_path)
        
        os.unlink(audio_path)
    
    if args.output:
        Path(args.output).write_text(output)
        if args.verbose:
            print(f"[yt-transcript] Wrote transcript to: {args.output}", file=sys.stderr)
    else:
        print(output)

if __name__ == "__main__":
    main()