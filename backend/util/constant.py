import os
from pathlib import Path

base_dir = Path(os.getcwd()+ "\\" + "temp")

image_dir = os.path.join(base_dir, "image")
character_dir = os.path.join(base_dir, "character")
novel_fragments_dir = os.path.join(base_dir, "fragments")
prompts_dir = os.path.join(base_dir, "prompts")
prompts_en_dir = os.path.join(base_dir, "promptsEn")
audio_dir = os.path.join(base_dir, "audio")
video_dir = os.path.join(base_dir, "video")
NovelPath = "novel.txt"