from flask import Flask, request, jsonify
import os
import shutil
import time

from backend.util.constant import image_dir, novel_fragments_dir, NovelPath, prompts_dir, prompts_en_dir
from backend.util.file import read_files_from_directory, read_lines_from_directory, save_list_to_files

def handle_error(status_code, message, error):
    response = jsonify({'error': message, 'details': str(error)})
    response.status_code = status_code
    return response

def save_lines_to_files(file_name):
    try:
        with open(file_name, 'r', encoding='utf-8') as file:
            lines = file.readlines()
            linesWithContent = []
            for line in lines:
                line = line.strip()
                if line:
                    linesWithContent.append(line)
            for i, line in enumerate(linesWithContent):
                file_path = os.path.join(novel_fragments_dir, f"{i}.txt")
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
        if os.path.exists(novel_fragments_dir):
            shutil.rmtree(novel_fragments_dir, ignore_errors=True)
        os.makedirs(novel_fragments_dir, exist_ok=True)
        error = save_list_to_files(fragments, novel_fragments_dir, 0)
        if error:
            return handle_error(500, "Failed to save", error)
    except Exception as e:
        return handle_error(500, "Failed to process request", e)

    return jsonify({"message": "Fragments saved successfully"}), 200

def get_novel_fragments():
    try:
        if os.path.exists(novel_fragments_dir):
            shutil.rmtree(novel_fragments_dir, ignore_errors=True)
        os.makedirs(novel_fragments_dir, exist_ok=True)
        error = save_lines_to_files(NovelPath)
        if error:
            return handle_error(500, "Failed to process file", error)

        lines, error = read_lines_from_directory(novel_fragments_dir)
        if error:
            return handle_error(500, "Failed to read fragments", error)
    except Exception as e:
        return handle_error(500, "Failed to process request", e)

    return jsonify(lines), 200

def get_initial():
    try:
        novels, error = read_lines_from_directory(novel_fragments_dir)
        if error:
            return handle_error(500, "Failed to read fragments", error)

        prompts, error = read_lines_from_directory(prompts_dir)
        if error:
            return handle_error(500, "Failed to read prompts", error)

        prompts_en, error = read_lines_from_directory(prompts_en_dir)
        if error:
            return handle_error(500, "Failed to read prompts", error)

        files = read_files_from_directory(image_dir)
        images = []

        for file in files:
            if not os.path.isdir(file):  
                image_path = os.path.join("/images", file)
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


def load_novel():
    try:
        with open(NovelPath, 'r', encoding='utf-8') as file:
            content = file.read()
        return jsonify({'content': content}), 200
    except FileNotFoundError:
        return jsonify({'content': ''}), 200  
    except Exception as e:
        return jsonify({'error': str(e)}), 500

def save_novel():
    try:
        data = request.get_json()
        content = data.get('content', '')

        with open(NovelPath, 'w', encoding='utf-8') as file:
            file.write(content)

        return jsonify({'message': '保存成功！'}), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500