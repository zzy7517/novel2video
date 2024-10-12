import asyncio
from backend_py.tts.edge_tts import by_edge_tts


def generate_audio_files():
    asyncio.run(by_edge_tts())