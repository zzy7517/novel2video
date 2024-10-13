import json
from pathlib import Path

from flask import Flask, request, jsonify
import os
import shutil
import re
import concurrent.futures
import logging

from backend.llm.llm import llm_translate, query_llm
from backend.util.constant import CharacterDir, NovelFragmentsDir, PromptsDir, PromptsEnDir
from backend.util.file import read_lines_from_directory, save_list_to_files

sys = """
#Task: #
从输入中提取画面信息

#Rules:#
1. 提取画面感强烈的元素，输出的内容是一组词语
2. 优先输出画面感强烈的名词，比如“谁知道呢，或许做了什么亏心事，惹得神灵降怒了吧…”这句输入中，神灵的画面感最强，则输出神灵
3. 输出画面感强的动词，比如开车，不输出画面感不强的动词，比如嘲讽
4. 不要输出形容词
5. 关于文本中描述的场景，可以适当结合上下文发散
6. 不要输出心理描写，不要输出情绪
7. 避免使用模糊或不明确的描述
8. 每个输入一定要有一个输出
9. 如果出现了人，需要具体到名字
10. 如果没有具体的名字，先结合上下文推断，如果还是没有，则不输出名字
11. 如果没有出现人，可以只描述下场景
12. 如果句子中不是很具体的内容，可以挑选文本中某个词语发散一下
13. 如果某一行输入无法提取出内容，输出“无”，但不要不输出

#example#
input
1.炎炎八月。
2.滴滴滴——！
3.刺耳的蝉鸣混杂着此起彼伏的鸣笛声，回荡在人流湍急的街道上，灼热的阳光炙烤着灰褐色的沥青路面，热量涌动，整个街道仿佛都扭曲了起来。
4.路边为数不多的几团树荫下，几个小年轻正簇在一起，叼着烟等待着红绿灯。
5.突然，一个正在吞云吐雾的小年轻似乎是发现了什么，轻咦了一声，目光落在了街角某处。
6.“阿诺，你在看什么？”他身旁的同伴问道。
7.那个名为阿诺的年轻人呆呆的望着街角，半晌才开口，“你说……盲人怎么过马路？”
8.同伴一愣，迟疑了片刻之后，缓缓开口：“一般来说，盲人出门都有人照看，或者导盲犬引导，要是在现代点的城市的话，马路边上也有红绿灯的语音播报，实在不行的话，或许能靠着声音和导盲杖一点点挪过去？”
9.阿诺摇了摇头，“那如果即没人照看，又没导盲犬，也没有语音播报，甚至连导盲杖都用来拎花生油了呢？”
10. 中年男子话刚刚脱口，便是不出意外的在人头汹涌的广场上带起了一阵嘲讽的骚动。
11. 当初的少年，自信而且潜力无可估量，不知让得多少少女对其春心荡漾，当然，这也包括以前的萧媚。
12. “唉…”莫名的轻叹了一口气，萧媚脑中忽然浮现出三年前那意气风发的少年，四岁练气，十岁拥有九段斗之气，十一岁突破十段斗之气，成功凝聚斗之气旋，一跃成为家族百年之内最年轻的斗者！

output
1. 夏天,炎热的街道
2. 街道,很多汽车
3. 街道,很多汽车
4. 红绿灯,马路,几个年轻人
5. 阿诺,向远处看
6. 阿诺,街角
6. 阿诺,盲人过马路
7. 导盲犬,盲人
8. 导盲犬,红绿灯,盲人
9. 阿诺
10. 广场,人群
11. 少女,喜欢
13. 萧媚,叹气

#Input Format:#
0.输入0
1.输入1
2.输入2
每个数字开头的行代表一个输入，每行输入必须对应一行输出

#Output Format:#
每行输入需要对应一行输出，每个输出用空行隔开
输入和输出的序号需要对应
不要输出多余内容
这是一个输出示例
0. 输出0
1. 输出1
2. 输出2

# 检查 #
输出的行数是否和输入一致，如果不一致，则重新生成输出内容
"""

fragmentsLen = 30

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
        lines, err = read_lines_from_directory(NovelFragmentsDir)
        if err:
            return jsonify({"error": "Failed to read fragments"}), 500
    except Exception as e:
        return jsonify({"error": "Failed to read fragments"}), 500

    try:
        if os.path.exists(PromptsDir):
            shutil.rmtree(PromptsDir)
        os.makedirs(PromptsDir)
    except Exception as e:
        return jsonify({"error": "Failed to manage directory"}), 500

    prompts_mid = generate_input_prompts(lines, fragmentsLen)
    offset = 0
    for p in prompts_mid:
        res = query_llm(p, sys, "x", 0.01, 8192)
        lines = res.split("\n")
        re_pattern = re.compile(r'^\d+\.\s*')
        t2i_prompts = [re_pattern.sub('', line) for line in lines if line.strip()]
        offset += len(t2i_prompts)
        try:
            save_list_to_files(t2i_prompts, PromptsDir, offset - len(t2i_prompts))
        except Exception as e:
            return jsonify({"error": "save list to file failed"}), 500

    
    lines, err = read_lines_from_directory(PromptsDir)
    if err:
        return jsonify({"error": "Failed to read fragments"}), 500
    
       

    logging.info("extract prompts from novel fragments finished")
    return jsonify(lines), 200

def get_prompts_en():
    try:
        if os.path.exists(PromptsEnDir):
            shutil.rmtree(PromptsEnDir)
        os.makedirs(PromptsEnDir)
    except Exception as e:
        return jsonify({"error": "Failed to manage directory"}), 500

    
    lines, err = read_lines_from_directory(PromptsDir)
    if err:
        return jsonify({"error": "Failed to read fragments"}), 500

    try:
        character_map = {}
        p = os.path.join(CharacterDir, 'characters.txt')
        if os.path.exists(p):
            with open(p, 'r', encoding='utf8') as file:
                    character_map = json.load(file)

        for i, line in enumerate(lines):
            for key, value in character_map.items():
                if key in line:
                    lines[i] = line.replace(key, value)
    except Exception as e:
        logging.error(f"translate prompts failed, err {e}")
        return jsonify({"error": "translate failed"}), 500

    try:
        lines = translate_prompts(lines)
    except Exception as e:
        logging.error(f"translate prompts failed, err {e}")
        return jsonify({"error": "translate failed"}), 500

    try:
        save_list_to_files(lines, PromptsEnDir, 0)
    except Exception as e:
        return jsonify({"error": "Failed to save promptsEn"}), 500

    logging.info("translate prompts to English finished")
    return jsonify(lines), 200

def save_prompt_en():
    req = request.get_json()
    if not req or 'index' not in req or 'content' not in req:
        return jsonify({"error": "parse request body failed"}), 400

    file_path = os.path.join(PromptsEnDir, f"{req['index']}.txt")
    try:
        os.makedirs(PromptsEnDir, exist_ok=True)
        with open(file_path, 'w') as file:
            file.write(req['content'])
    except Exception as e:
        return jsonify({"error": "Failed to write file"}), 500

    return jsonify({"message": "Attachment saved successfully"}), 200

def save_prompt_zh():
    req = request.get_json()
    if not req or 'index' not in req or 'content' not in req:
        return jsonify({"error": "parse request body failed"}), 400

    file_path = os.path.join(PromptsDir, f"{req['index']}.txt")
    try:
        os.makedirs(PromptsDir, exist_ok=True)
        with open(file_path, 'w') as file:
            file.write(req['content'])
    except Exception as e:
        return jsonify({"error": "Failed to write file"}), 500

    return jsonify({"message": "Attachment saved successfully"}), 200