import asyncio
from flask import Flask, jsonify, request
import os
import re
import time

from backend.image.sd import generate_image
from backend.util.constant import image_dir, prompts_en_dir
from backend.util.file import make_dir, read_lines_from_directory, remove_all

def handle_error(message, err):
    return jsonify({"error": message}), 500

async def async_generate_images(lines):
    for i, p in enumerate(lines):
        try:
            await generate_image(p, 114514191981, 540, 960, i)
        except Exception as e:
            print("Error:", e)

async def async_generate_image_single(content, index):
    try:
        await generate_image(content, 114514191981, 540, 960, index)
    except Exception as e:
        print("Error:", e)
                
def generate_images():
    try:
        remove_all(image_dir)
        make_dir(image_dir)
    except Exception as e:
        return handle_error("Failed to manage directory", e)
    try:
        lines, err = read_lines_from_directory(prompts_en_dir)
        if err:
            return handle_error("Failed to read fragments", err)
        asyncio.run(async_generate_images(lines))
    except Exception as e:
        return handle_error("Failed to read fragments", e)
    return jsonify({"status": "Image generation started"}), 200

def get_local_images():
    try:
        files = os.listdir(image_dir)
    except Exception as e:
        return jsonify({"error": "Failed to read image directory"}), 500

    image_map = {}
    for file in files:
        if not os.path.isdir(file):
            matches = re.match(r'(\d+)\.png', file)
            if matches:
                key = matches.group(1)
                abs_path = os.path.join("/images", file)
                image_map[key] = abs_path

    return jsonify(image_map), 200

def generate_single_image():
    try:
        req = request.get_json()
        if not req or 'index' not in req or 'content' not in req:
            return jsonify({"error": "parse request body failed"}), 400
        file = os.path.join("/images", str(req['index'])+'.png')
        asyncio.run(async_generate_image_single(req['content'], req['index']))
    except Exception as e:
        return handle_error("Failed to read fragments", e)
    return jsonify({"status": "Image generation started", "url":file}), 200
