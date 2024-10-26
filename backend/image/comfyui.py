import asyncio
import base64
import copy
import logging
import json
import os
import random

import requests

from backend.util.constant import image_dir
from backend.util.file import get_config

async def generate_image(lines):
    url = get_config()['address3']
    api = get_config()['comfyuiNodeApi']
    task_id_to_file = {}
    task_status = {}
    for i, p in enumerate(lines):
        try:
            prompt_id = post_prompt(p, api, url)
            task_id_to_file[prompt_id] = i
            task_status[prompt_id] = False
        except Exception as e:
            logging.error(e)
            raise
    await monitor_tasks(task_id_to_file, task_status, url)


async def generate_single_image(content, i):
    url = get_config()['address3']
    api = get_config()['comfyuiNodeApi']
    task_id_to_file = {}
    task_status = {}

    prompt_id = post_prompt(content, api, url)
    task_id_to_file[prompt_id] = i
    task_status[prompt_id] = False
    await monitor_tasks(task_id_to_file, task_status, url)

async def monitor_tasks(task_id_to_file, task_status, url):
    while len(task_id_to_file) > 0:
        keys_to_delete = []
        await asyncio.sleep(2)
        for key, value in task_id_to_file.items():
            status = prompt_history(key, url, value)
            if status:
                keys_to_delete.append(key)
        for key in keys_to_delete:
            del task_id_to_file[key]
            del task_status[key]


def prompt_history(prompt_id, url, order):
    try:
        response = requests.get(f"{url}/history/{prompt_id}")
        response.raise_for_status()
        response_data = response.json()
        if len(response_data) == 0:
            logging.info("waiting...")
            return False
        if "status" not in response_data[prompt_id]:
            return False
        status = response_data[prompt_id]["status"]
        if status["status_str"] != "success" or status["completed"] != True:
            logging.error(f"status {status} ")
            return False
        filename = None
        outputs = response_data[prompt_id]["outputs"]
        for key, value in outputs.items():
            if "images" in value and len(value["images"]) > 0:
                filename = value["images"][0]["filename"]
                break
        if filename:
            response = requests.get(f"{url}/view?filename={filename}&type=output")
            if response.status_code == 200:
                if not os.path.exists(image_dir):
                    os.makedirs(image_dir)
                output_filename = os.path.join(image_dir, f"{order}.png")
                try:
                    with open(output_filename, "wb") as image_file:
                        image_file.write(response.content)
                except IOError as e:
                    logging.error(f"Failed to write image file: {e}")
                    return False
            else:
                return False
            logging.info(f"Image saved to {output_filename}")
            return True
        else:
            logging.error(f"No images found for prompt {prompt_id}")
            return True
    except Exception as e:
        logging.error(e)
        raise

def generate_random_seed(copyed_data, p):
    random_number = random.randint(1, 10 ** 15 - 1)
    if isinstance(copyed_data, dict):
        for key, value in copyed_data.items():
            if key == "noise_seed":
                copyed_data[key] = random_number
            elif isinstance(value, (dict, list)):
                generate_random_seed(value, p)
    elif isinstance(copyed_data, list):
        for index, item in enumerate(copyed_data):
            if isinstance(item, (dict, list)):
                generate_random_seed(item, p)
    return copyed_data

def replace_prompt_in_map(copyed_data, p):
    target_value = "$prompt$"
    if isinstance(copyed_data, dict):
        for key, value in copyed_data.items():
            if isinstance(value, (dict, list)):
                replace_prompt_in_map(value, p)
            elif isinstance(value, str) and target_value in value:
                copyed_data[key] = value.replace(target_value, p)
    elif isinstance(copyed_data, list):
        for index, item in enumerate(copyed_data):
            if isinstance(item, (dict, list)):
                replace_prompt_in_map(item, p)
            elif isinstance(item, str) and target_value in item:
                copyed_data[index] = item.replace(target_value, p)
    return copyed_data

def post_prompt(prompt, apiRaw, url):
    try:
        api1 = replace_prompt_in_map(copy.deepcopy(apiRaw), prompt)
        api = generate_random_seed(copy.deepcopy(api1), prompt)
        payload = {
            "prompt": api,
        }
    except Exception as e:
        logging.error(e)
        raise
    try:
        response = requests.post(f"{url}/prompt", json=payload)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        logging.error(f"Failed to make request: {e}")
        raise

    try:
        response_data = response.json()
        prompt_id = response_data.get("prompt_id")
        if not prompt_id:
            logging.error("No prompt_id found in response")
            raise Exception("No images found")
    except json.JSONDecodeError as e:
        logging.error(f"Failed to decode JSON response: {e}")
        raise
    return prompt_id



