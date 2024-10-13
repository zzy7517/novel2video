import os
import time
from flask import jsonify

from backend.util.constant import video_dir
from backend.util.file import make_dir, remove_all
from backend.util.movie import create_video_with_audio_images


def get_video():
    """
    Endpoint to fetch the initial video.
    """
    try:
        video_data = {
            "videoUrl": os.path.join("/videos", "video.mp4")
        }
        return jsonify(video_data), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500

def generate_video():
    """
    Endpoint to generate a new video.
    """
    try:
        remove_all(video_dir)
        make_dir(video_dir)
        create_video_with_audio_images()
        now = int(time.time()) 
        new_video_data = {
            "videoUrl": os.path.join("/videos", "video.mp4") + f"?v={now}"
        }
        return jsonify(new_video_data), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500