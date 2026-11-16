"""
My GPU isn't very big (8 GB), so this patch makes it so that:
1. The parakeet model is unloaded as soon as it's done with transcription
2. The Ollama-provided ad detection does not start until the GPU is free

Fixes this for both passes
"""

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


# Release Whisper before second-pass LLM detection.
verification_file = Path("/app/src/verification_pass.py")
source = verification_file.read_text()

target = (
    "        verification_result = "
    "self.ad_detector.run_verification_detection(\n"
)

replacement = (
    "        # Release local transcription VRAM before LLM verification.\n"
    "        from transcriber.local import unload_whisper_if_idle\n"
    "        unload_whisper_if_idle()\n"
    "        verification_result = "
    "self.ad_detector.run_verification_detection(\n"
)

if replacement not in source:
    if source.count(target) != 1:
        raise RuntimeError(
            "Could not uniquely locate second-pass LLM detection"
        )

    verification_file.write_text(
        source.replace(target, replacement, 1)
    )

print("MinusPod second-pass Whisper unload patch applied")
