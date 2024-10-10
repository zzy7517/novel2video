import logging
import requests
import json
from typing import Tuple, Optional

from backend_py.llm.keys import SAMBA_NOVA_API_KEY

LLAMA_405B = "Meta-Llama-3.1-405B-Instruct"

def query_samba_nova(input_text: str, sys: str, model_name: str, temperature: float) -> str:
    url = "https://api.sambanova.ai/v1/chat/completions"
    messages = []
    if sys:
        messages.append({"role": "system", "content": sys})
    messages.append({"role": "user", "content": input_text})

    request_body = {
        "temperature": temperature,
        "messages": messages,
        "model": LLAMA_405B,
    }

    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {SAMBA_NOVA_API_KEY}",
    }

    try:
        response = requests.post(url, headers=headers, json=request_body)
        response.raise_for_status()  # Raises an HTTPError for bad responses

        response_data = response.json()
        if 'choices' in response_data and len(response_data['choices']) > 0:
            return response_data['choices'][0]['message']['content']
        else:
            logging.error(f"No choices found in response: {response_data}")
            return ""
    except requests.exceptions.RequestException as e:
        logging.error(f"Error querying SambaNova: {e}")
        return ""

# Example usage:
# result, error = query_samba_nova("Hello, how are you?", "", LLAMA_405B, 0.5)
# if error:
#     print(f"Error: {error}")
# else:
#     print(result)
