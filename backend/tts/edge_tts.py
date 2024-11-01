import asyncio
import os
import shutil
from concurrent.futures import ThreadPoolExecutor

import edge_tts

from backend.util.constant import audio_dir, novel_fragments_dir
from backend.util.file import read_lines_from_directory


# async def by_edge_tts():
#     if os.path.exists(audio_dir):
#         shutil.rmtree(audio_dir)
#     os.makedirs(audio_dir, exist_ok=True)
#     lines, err = read_lines_from_directory(novel_fragments_dir)
#     if err:
#         raise "Failed to read fragments"
#     if lines is None:
#         print(f"Failed to read novel fragments from {novel_fragments_dir}")
#         return
#     for i, line in enumerate(lines):
#         try:
#             # zh-CN-YunxiNeural  YunjianNeural  rate='25%'
#             communicate = edge_tts.Communicate(text=line, voice="zh-CN-YunxiNeural",rate='+35%')
#             full_path = os.path.join(audio_dir, f"{i}.mp3")
#             await communicate.save(full_path)
#         except Exception as e:
#             print(f"TTS conversion failed for line {i}, error: {e}")

async def by_edge_tts():
    if os.path.exists(audio_dir):
        shutil.rmtree(audio_dir)
    os.makedirs(audio_dir, exist_ok=True)
    lines, err = read_lines_from_directory(novel_fragments_dir)
    if err:
        raise "Failed to read fragments"
    if lines is None:
        print(f"Failed to read novel fragments from {novel_fragments_dir}")
        return

    with ThreadPoolExecutor(max_workers=20) as executor:
        for i, line in enumerate(lines):
            executor.submit(convert_text_to_speech, line, audio_dir, i)

def convert_text_to_speech(line, audio_dir, i):
    try:
        # zh-CN-YunxiNeural  YunjianNeural  rate='25%' YunyangNeural
        communicate = edge_tts.Communicate(
            text=line,
            voice="zh-CN-XiaoxiaoNeural",
            rate="+35%"
        )

        full_path = os.path.join(audio_dir, f"{i}.mp3")
        asyncio.run(communicate.save(full_path))

    except Exception as e:
        # Handle any other unexpected errors
        print(f"An unexpected error occurred for line {i}, error: {e}")