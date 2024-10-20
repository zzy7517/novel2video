import json
import logging
import os

import requests
import base64

from backend.util.constant import image_dir
from backend.util.file import get_config


async def generate_image(prompt: str, seed: int, width: int, height: int, order):
    try:
        url = get_config()['address3']
        payload = {
            "prompt": "anime" +  prompt + " ",
            "negative_prompt": "(painting by bad-artist-anime:0.9), (painting by bad-artist:0.9), watermark, text, error, blurry, jpeg artifacts, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, artist name,deformed,distorted,disfigured,doll,poorly drawn,bad anatomy,wrong anatomy,bad hand,bad fingers,NSFW",
            "cfg_scale": 7,
            "steps": 25,
            "width": width,
            "height": height,
            "override_settings": {
                # "sd_vae": "Automatic",
            },
            "scheduler": "Simple",
            "forge_additional_modules": [
                "E:\\sd\\FORGE-V2-\\forge\\models\\VAE\\ae.safetensors",
                "E:\\sd\\FORGE-V2-\\forge\\models\\VAE\\clip_l.safetensors",
                "E:\\sd\\FORGE-V2-\\forge\\models\\VAE\\t5xxl_fp16.safetensors",
            ],
            "enable_hr": True,
            "denoising_strength": 0.7,
            "hr_upscaler": "Latent",
            # "hr_resize_x": 1024,
            # "hr_resize_y": 1024,
            "hr_sampler_name": "Euler",
            "hr_second_pass_steps": 10,
        }
    except Exception as e:
        logging.error(e)
        return

    try:
        response = requests.post(f"{url}/sdapi/v1/txt2img", json=payload)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        logging.error(f"Failed to make request: {e}")
        return

    try:
        response_data = response.json()
    except json.JSONDecodeError as e:
        logging.error(f"Failed to decode JSON response: {e}")
        return

    images = response_data.get("images")
    if not images:
        logging.error("No images found in response")
        return

    try:
        image_data = base64.b64decode(images[0])
    except (IndexError, base64.binascii.Error) as e:
        logging.error(f"Failed to decode image: {e}")
        return

    if not os.path.exists(image_dir):
        os.makedirs(image_dir)
    output_filename = os.path.join(image_dir, f"{order}.png")

    try:
        with open(output_filename, "wb") as image_file:
            image_file.write(image_data)
    except IOError as e:
        logging.error(f"Failed to write image file: {e}")
        return

    logging.info(f"Image saved to {output_filename}")