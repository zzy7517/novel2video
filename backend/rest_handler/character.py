import os
import json
import logging
import shutil

from flask import Flask, request, jsonify

from backend.llm.llm import query_llm
from backend.util.constant import character_dir, prompts_dir, novel_fragments_dir

extract_character_sys = """
	#Task: #
	Extract characters from the novel fragment
	
	#Rule#
	1. 提取出所有出现过的角色
	2. 所有的人名，别名，称呼，包括对话中引用到的名字都需要提取
    3. 所有出现过的和人有关的称呼都需要提取
    4. 如果文本以第一人称或者第二人称叙述，我/你这种人称代词也需要提取
	5. 如果一个人出现过不止一种称呼，则都提取
	
	#Output Format:#
	名字1/名字2/名字3/...
"""

def get_new_characters():
    try:
        # Remove and recreate the character directory
        if os.path.exists(character_dir):
            shutil.rmtree(character_dir)
        os.makedirs(character_dir)

        # Read lines from the prompts directory
        lines = []
        for file_name in os.listdir(novel_fragments_dir):
            with open(os.path.join(novel_fragments_dir, file_name), 'r') as file:
                lines.extend(file.readlines())

        # Process lines in chunks
        character_map = {}
        for i in range(0, len(lines), 500):
            end = min(i + 500, len(lines))
            prompt = ''.join(lines[i:end])
            response = query_llm(prompt, extract_character_sys, "doubao", 0.01, 8192)
            for character in response.split('/'):
                character_map[character.strip()] = character.strip()

        # Save characters to a file
        with open(os.path.join(character_dir, 'characters.txt'), 'w') as file:
            json.dump(character_map, file)

        return jsonify(character_map), 200

    except Exception as e:
        logging.error(f"Failed to get new characters: {e}")
        return jsonify({"error": "Failed to get new characters"}), 500


def get_local_characters():
    try:
        if not os.path.exists(character_dir):
            return jsonify({"error":"no local characters"}), 40401
        with open(os.path.join(character_dir, 'characters.txt'), 'r') as file:
            character_map = json.load(file)
        return jsonify(character_map), 200
    except Exception as e:
        logging.error(f"Failed to get local characters: {e}")
        return jsonify({"error": "Failed to get local characters"}), 500

def put_characters():
    try:
        descriptions = request.json
        if not descriptions:
            return jsonify({"error": "Invalid JSON"}), 400

        with open(os.path.join(character_dir, 'characters.txt'), 'r') as file:
            character_map = json.load(file)

        character_map.update(descriptions)
        # Save descriptions to a file
        with open(os.path.join(character_dir, 'characters.txt'), 'w') as file:
            json.dump(character_map, file)

        return jsonify({"message": "Descriptions updated successfully"}), 200

    except Exception as e:
        logging.error(f"Failed to put characters: {e}")
        return jsonify({"error": "Failed to put characters"}), 500

appearance_prompt = """
    随机生成动漫角色的外形描述，输出简练，以一组描述词的形式输出，每个描述用逗号隔开
    数量：一个
    包含：衣着，眼睛，发色，发型等等。
    生成男性和女性的概率都为50%
    根据生成的年龄和性别, 输出时在最前方标明1girl/1man/1boy/1lady等等
    示例1: 1boy, white school uniform, brown eyes, short messy black hair.
    示例2: 1girl, short skirt, skinny, blue eyes, blonde hair, twin tails, knee-high socks
    示例3: 1man, Navy suit, Green eyes, short, slicked-back blonde hair
    示例4: 1elderly man, Grey cardigan, Grey eyes, balding with white hair
    使用英文输出，不要输出额外内容
"""

def get_random_appearance():
    try:
        prompt = appearance_prompt
        appearance = query_llm(prompt, "", "doubao", 1, 100)
        return jsonify(appearance), 200
    except Exception as e:
        logging.error(f"Failed to get random appearance: {e}")
        return jsonify({"error": "Failed to get random appearance"}), 400