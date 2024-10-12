import os
import time
from flask import jsonify

from backend.util.file import make_dir, remove_all
from backend.util.movie import create_video_with_audio_images


def get_video():
    """
    Endpoint to fetch the initial video.
    """
    try:
        now = int(time.time()) 
        video_data = {
            "videoUrl": os.path.join("/video", "video.mp4") + f"?v={now}"
        }
        return jsonify(video_data), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500

def generate_video():
    """
    Endpoint to generate a new video.
    """
    try:
        remove_all("temp/video")
        make_dir("temp/video")
        create_video_with_audio_images()
        now = int(time.time()) 
        new_video_data = {
            "videoUrl": os.path.join("/videos", "video.mp4") + f"?v={now}"
        }
        return jsonify(new_video_data), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500