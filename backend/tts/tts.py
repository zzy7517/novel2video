import asyncio
from backend.tts.edge_tts import by_edge_tts


def generate_audio_files():
    asyncio.run(by_edge_tts())