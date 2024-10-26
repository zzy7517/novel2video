import json
from flask import Flask, request, jsonify
import os
import shutil
import re
import concurrent.futures
import logging
from backend.llm.llm import llm_translate, query_llm
from backend.util.constant import character_dir, novel_fragments_dir, prompts_dir, prompts_en_dir, prompt_path
from backend.util.file import read_lines_from_directory, save_list_to_files, read_file

fragmentsLen = 80

# Function to generate input prompts
def generate_input_prompts(lines, step):
    prompts = []
    for i in range(0, len(lines), step):
        end = min(i + step, len(lines))
        prompt = "\n".join(f"{j}. {lines[j]}" for j in range(i, end))
        logging.info(f"prompt is {prompt}")
        prompts.append(prompt)
    return prompts

# Function to translate prompts
def translate_prompts(lines):
    def translate_line(line):
        res = llm_translate(line)
        return res

    with concurrent.futures.ThreadPoolExecutor() as executor:
        translated_lines = list(executor.map(translate_line, lines))
    return translated_lines

def extract_scene_from_texts():
    try:
        lines, err = read_lines_from_directory(novel_fragments_dir)
        if err:
            return jsonify({"error": "Failed to read fragments"}), 500
    except Exception as e:
        return jsonify({"error": "Failed to read fragments"}), 500

    try:
        if os.path.exists(prompts_dir):
            shutil.rmtree(prompts_dir)
        os.makedirs(prompts_dir)
    except Exception as e:
        return jsonify({"error": "Failed to manage directory"}), 500

    prompts_mid = generate_input_prompts(lines, fragmentsLen)
    offset = 0
    sys = read_file(prompt_path)
    for p in prompts_mid:
        res = query_llm(p, sys, "x", 1, 8192)
        logging.info(res)
        lines = res.split("\n")
        re_pattern = re.compile(r'^\d+\.\s*')
        t2i_prompts = [re_pattern.sub('', line) for line in lines if line.strip()]
        offset += len(t2i_prompts)
        try:
            save_list_to_files(t2i_prompts, prompts_dir, offset - len(t2i_prompts))
        except Exception as e:
            return jsonify({"error": "save list to file failed"}), 500

    
    lines, err = read_lines_from_directory(prompts_dir)
    if err:
        return jsonify({"error": "Failed to read fragments"}), 500
    
       

    logging.info("extract prompts from novel fragments finished")
    return jsonify(lines), 200

def get_prompts_en():
    try:
        if os.path.exists(prompts_en_dir):
            shutil.rmtree(prompts_en_dir)
        os.makedirs(prompts_en_dir)
    except Exception as e:
        return jsonify({"error": "Failed to manage directory"}), 500

    
    lines, err = read_lines_from_directory(prompts_dir)
    if err:
        return jsonify({"error": "Failed to read fragments"}), 500

    try:
        character_map = {}
        p = os.path.join(character_dir, 'characters.txt')
        if os.path.exists(p):
            with open(p, 'r', encoding='utf8') as file:
                    character_map = json.load(file)

        for i, line in enumerate(lines):
            for key, value in character_map.items():
                if key in line:
                    lines[i] = lines[i].replace(key, value)
    except Exception as e:
        logging.error(f"translate prompts failed, err {e}")
        return jsonify({"error": "translate failed"}), 500

    try:
        lines = translate_prompts(lines)
    except Exception as e:
        logging.error(f"translate prompts failed, err {e}")
        return jsonify({"error": "translate failed"}), 500

    try:
        save_list_to_files(lines, prompts_en_dir, 0)
    except Exception as e:
        return jsonify({"error": "Failed to save promptsEn"}), 500

    logging.info("translate prompts to English finished")
    return jsonify(lines), 200

def save_prompt_en():
    req = request.get_json()
    if not req or 'index' not in req or 'content' not in req:
        return jsonify({"error": "parse request body failed"}), 400

    file_path = os.path.join(prompts_en_dir, f"{req['index']}.txt")
    try:
        os.makedirs(prompts_en_dir, exist_ok=True)
        with open(file_path, 'w') as file:
            file.write(req['content'])
    except Exception as e:
        return jsonify({"error": "Failed to write file"}), 500

    return jsonify({"message": "Attachment saved successfully"}), 200

def save_prompt_zh():
    req = request.get_json()
    if not req or 'index' not in req or 'content' not in req:
        return jsonify({"error": "parse request body failed"}), 400

    file_path = os.path.join(prompts_dir, f"{req['index']}.txt")
    try:
        os.makedirs(prompts_dir, exist_ok=True)
        with open(file_path, 'w') as file:
            file.write(req['content'])
    except Exception as e:
        return jsonify({"error": "Failed to write file"}), 500

    return jsonify({"message": "Attachment saved successfully"}), 200