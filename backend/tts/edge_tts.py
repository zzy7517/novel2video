import os
import shutil

import edge_tts

from backend_py.util.constant import AudioDir, NovelFragmentsDir
from backend_py.util.file import read_lines_from_directory


async def by_edge_tts():
    # Remove and recreate the audio directory
    if os.path.exists(AudioDir):
        shutil.rmtree(AudioDir)
    os.makedirs(AudioDir, exist_ok=True)

    # Read lines from the novel fragments directory
    lines, err = read_lines_from_directory(NovelFragmentsDir)
    if err:
        raise "Failed to read fragments"
    if lines is None:
        print(f"Failed to read novel fragments from {NovelFragmentsDir}")
        return

    # Convert each line to speech and save as an mp3 file
    for i, line in enumerate(lines):
        try:
            communicate = edge_tts.Communicate(text=line, voice="en-US-EmmaMultilingualNeural", rate='+25%')
            full_path = os.path.join(AudioDir, f"{i}.mp3")
            await communicate.save(full_path)
        except Exception as e:
            print(f"TTS conversion failed for line {i}, error: {e}")