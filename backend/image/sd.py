import json
import logging
import requests
import base64

from backend.util.constant import image_dir

async def generate_image(prompt: str, seed: int, width: int, height: int, order):
    url = "http://10.193.239.248:7860"
    payload = {
        "prompt": "anime" +  prompt + " <lora:超级玄幻:0.7> ",
        "negative_prompt": "(painting by bad-artist-anime:0.9), (painting by bad-artist:0.9), watermark, text, error, blurry, jpeg artifacts, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, artist name,deformed,distorted,disfigured,doll,poorly drawn,bad anatomy,wrong anatomy,bad hand,bad fingers,NSFW",
        "cfg_scale": 7,
        "steps": 35,
        "width": width,
        "height": height,
        "override_settings": {
            "sd_vae": "Automatic",
        },
        # "enable_hr": True,
        # "denoising_strength": 0.7,
        # "hr_upscaler": "Latent",
        # "hr_resize_x": 1024,
        # "hr_resize_y": 1024,
        # "hr_sampler_name": "Euler",
        # "hr_second_pass_steps": 10,
    }

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

    output_filename = f"{image_dir}/{order}.png"

    try:
        with open(output_filename, "wb") as image_file:
            image_file.write(image_data)
    except IOError as e:
        logging.error(f"Failed to write image file: {e}")
        return

    logging.info(f"Image saved to {output_filename}")