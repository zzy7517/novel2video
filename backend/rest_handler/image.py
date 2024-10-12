import asyncio
from flask import Flask, jsonify, request
import os
import re
import time
import shutil
from threading import Thread

from backend.image.sd import generate_image
from backend.util.constant import ImageDir, PromptsEnDir
from backend.util.file import read_lines_from_directory


def remove_all(directory):
    shutil.rmtree(directory, ignore_errors=True)

def make_dir(directory):
    os.makedirs(directory, exist_ok=True)

def handle_error(message, err):
    return jsonify({"error": message}), 500

async def async_generate_images(lines):
    for i, p in enumerate(lines):
        try:
            await generate_image(p, 114514191981, 540, 960, i)
        except Exception as e:
            print("Error:", e)
                
def generate_images():
    try:
        remove_all(ImageDir)
        make_dir(ImageDir)
    except Exception as e:
        return handle_error("Failed to manage directory", e)
    try:
        lines, err = read_lines_from_directory(PromptsEnDir)
        if err:
            return handle_error("Failed to read fragments", err)
        asyncio.run(async_generate_images(lines))
    except Exception as e:
        return handle_error("Failed to read fragments", e)
    return jsonify({"status": "Image generation started"}), 200

def get_local_images():
    try:
        files = os.listdir(ImageDir)
    except Exception as e:
        return jsonify({"error": "Failed to read image directory"}), 500

    image_map = {}
    now = int(time.time())
    for file in files:
        if not os.path.isdir(file):
            matches = re.match(r'(\d+)\.png', file)
            if matches:
                key = matches.group(1)
                abs_path = os.path.join("/images", file)
                image_map[key] = f"{abs_path}?v={now}"

    return jsonify(image_map), 200
