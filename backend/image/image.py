from backend.image import sd, comfyui
from backend.util.file import get_config

width = 512
height = 512

async def generate_image(lines):
    try:
        type = get_config()['address3Type']
        if type == 'stable_diffusion_web_ui':
            for i, p in enumerate(lines):
                await sd.generate_image(p, 114514191981, width, height, i)

        elif type == 'comfyui':
            await comfyui.generate_image(lines)
        else:
            raise ValueError(f"Unknown tool type: {type}")
    except Exception as e:
        print(f"An error occurred: {e}")
        raise

async def generate_images_single(content, i):
    try:
        type = get_config()['address3Type']
        if type == 'stable_diffusion_web_ui':
            await sd.generate_image(content, 114514191981, width, height, i)
        elif type == 'comfyui':
            await comfyui.generate_single_image(content, i)
        else:
            raise ValueError(f"Unknown tool type: {type}")
    except Exception as e:
        print(f"An error occurred: {e}")
        raise
