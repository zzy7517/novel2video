from flask import Flask, request, jsonify
import os
import shutil
import time

from backend_py.util.constant import ImageDir, NovelFragmentsDir, PromptsDir, PromptsEnDir
from backend_py.util.file import read_lines_from_directory, save_list_to_files

def handle_error(status_code, message, error):
    response = jsonify({'error': message, 'details': str(error)})
    response.status_code = status_code
    return response

def save_lines_to_files(file_name):
    try:
        with open(file_name, 'r') as file:
            lines = file.readlines()
            for i, line in enumerate(lines):
                line = line.strip()
                if line:
                    file_path = os.path.join(NovelFragmentsDir, f"{i + 1}.txt")
                    with open(file_path, 'w') as f:
                        f.write(line)
    except Exception as e:
        return e
    return None


def save_combined_fragments():
    fragments = request.json
    if not isinstance(fragments, list):
        return handle_error(400, "Invalid request", "Expected a list of strings")

    try:
        shutil.rmtree(NovelFragmentsDir, ignore_errors=True)
        os.makedirs(NovelFragmentsDir, exist_ok=True)
        error = save_list_to_files(fragments, NovelFragmentsDir)
        if error:
            return handle_error(500, "Failed to save", error)
    except Exception as e:
        return handle_error(500, "Failed to process request", e)

    return jsonify({"message": "Fragments saved successfully"}), 200

def get_novel_fragments():
    try:
        shutil.rmtree(NovelFragmentsDir, ignore_errors=True)
        os.makedirs(NovelFragmentsDir, exist_ok=True)
        error = save_lines_to_files('a.txt')
        if error:
            return handle_error(500, "Failed to process file", error)

        lines, error = read_lines_from_directory(NovelFragmentsDir)
        if error:
            return handle_error(500, "Failed to read fragments", error)
    except Exception as e:
        return handle_error(500, "Failed to process request", e)

    return jsonify(lines), 200

def get_initial():
    try:
        novels, error = read_lines_from_directory(NovelFragmentsDir)
        if error:
            return handle_error(500, "Failed to read fragments", error)

        prompts, error = read_lines_from_directory(PromptsDir)
        if error:
            return handle_error(500, "Failed to read prompts", error)

        prompts_en, error = read_lines_from_directory(PromptsEnDir)
        if error:
            return handle_error(500, "Failed to read prompts", error)

        files = os.listdir(ImageDir)

        images = []
        now = int(time.time())  # Get the current Unix timestamp

        for file in files:
            if not os.path.isdir(file):  # Check if the file is not a directory
                image_path = os.path.join("/static", file) + f"?v={now}"
                images.append(image_path)

        data = {
            "fragments": novels,
            "images": images,
            "prompts": prompts,
            "promptsEn": prompts_en
        }
    except Exception as e:
        return handle_error(500, "Failed to process request", e)

    return jsonify(data), 200

