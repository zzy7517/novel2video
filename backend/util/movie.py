import logging
import os
import re
from moviepy.editor import ImageClip, concatenate_videoclips, AudioFileClip

from backend.util.constant import AudioDir, ImageDir

def extract_number(filename):
    match = re.search(r'\d+', filename)
    return int(match.group()) if match else float('inf')

def create_video_with_audio_images():
    """
    根据提供的图片集长度，生成一个视频。
    图片和音频列表将在函数内部生成。

    Parameters:
    - length: 图片集的长度，也是音频文件的数量。
    - txt_name: 视频名称。
    - lang: 音频的语言（'en' 或 'cn'）。
    """
    try:
        images = [os.path.join(ImageDir, file) for file in os.listdir(ImageDir) if file.endswith('.png')]
        images.sort(key=lambda x: extract_number(os.path.basename(x)))  

        audios = [os.path.join(AudioDir, file) for file in os.listdir(AudioDir) if file.endswith('.mp3')]
        audios.sort(key=lambda x: extract_number(os.path.basename(x))) 
        clips = []

        for img_path, audio_path in zip(images, audios):
            audio_clip = AudioFileClip(audio_path)
            image_clip = ImageClip(img_path).set_duration(audio_clip.duration).set_audio(audio_clip)
            clips.append(image_clip)
        final_clip = concatenate_videoclips(clips, method="compose")
        final_clip.write_videofile("temp/video/video.mp4", fps=30)
    except Exception as e :
        logging.err("gen video failed{e}")
        raise

