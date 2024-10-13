import os
from pathlib import Path

base_dir = Path(os.getcwd()+ "\\" + "temp")

ImageDir = os.path.join(base_dir, "image")
CharacterDir = os.path.join(base_dir, "character")
NovelFragmentsDir = os.path.join(base_dir, "fragments")
PromptsDir = os.path.join(base_dir, "prompts")
PromptsEnDir = os.path.join(base_dir, "promptsEn")
AudioDir = os.path.join(base_dir, "audio")
NovelPath = "novel.txt"