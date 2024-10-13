import asyncio
import logging

from flask import jsonify

from backend.tts.edge_tts import by_edge_tts


def generate_audio_files():
    try:
        asyncio.run(by_edge_tts())
        return jsonify("suc"), 200
    except Exception as e:
        logging.error(f"generate_audio_files error: {e}")
        return jsonify(f"failed {e}"), 500