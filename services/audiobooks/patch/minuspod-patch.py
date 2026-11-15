from pathlib import Path

config = Path("/app/src/config.py")
source = config.read_text()

old = "LOCAL_MODEL_IDLE_UNLOAD_SECONDS = 300"
new = "LOCAL_MODEL_IDLE_UNLOAD_SECONDS = 0"

if old in source:
    config.write_text(source.replace(old, new, 1))
elif new not in source:
    raise RuntimeError("Whisper timeout setting not found")

pipeline = Path("/app/src/main_app/pipeline/transcribe.py")
source = pipeline.read_text()

old_return = "    return audio_path, segments"
new_return = (
    "    from transcriber.local import unload_whisper_if_idle\n"
    "    unload_whisper_if_idle()\n"
    "    return audio_path, segments"
)

if new_return not in source:
    if old_return not in source:
        raise RuntimeError("Transcription return point not found")
    source = source.replace(old_return, new_return, 1)
    pipeline.write_text(source)
