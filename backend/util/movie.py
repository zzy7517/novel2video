import logging
import os
import re
from moviepy.editor import ImageClip, concatenate_videoclips, AudioFileClip

from backend.util.constant import audio_dir, image_dir, video_dir


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
        images = [os.path.join(image_dir, file) for file in os.listdir(image_dir) if file.endswith('.png')]
        images.sort(key=lambda x: extract_number(os.path.basename(x)))  

        audios = [os.path.join(audio_dir, file) for file in os.listdir(audio_dir) if file.endswith('.mp3')]
        audios.sort(key=lambda x: extract_number(os.path.basename(x))) 
        clips = []

        for img_path, audio_path in zip(images, audios):
            audio_clip = AudioFileClip(audio_path)
            image_clip = ImageClip(img_path).set_duration(audio_clip.duration).set_audio(audio_clip)
            clips.append(image_clip)
        final_clip = concatenate_videoclips(clips, method="compose")
        video = os.path.join(video_dir, 'video.mp4')
        final_clip.write_videofile(video, fps=24)
    except Exception as e :
        logging.error(f"gen video failed{e}")
        raise

# todo 字幕 网站https://news.miracleplus.com/share_link/45459