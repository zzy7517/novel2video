import requests
from typing import Tuple, Optional

from backend.llm.sambanova import query_samba_nova
from backend.llm.siliconflow import query_silicon_flow

def query_llm(input_text: str, system_content: str, model_name: str, temperature: float, max_output_tokens: int) -> str:
    # Example of querying a specific model
    # return query_gemini(input_text, system_content, "gemini-1.5-pro-002", temperature, max_output_tokens)
    return query_samba_nova(input_text, system_content, model_name, temperature)
    # if model_name == "doubao":
    #     return query_vol_engine(input_text, system_content, 0.98)
    # elif model_name == "gemini-1.5-pro-002":
    #     return query_gemini(input_text, system_content, model_name, temperature, max_output_tokens)
    # else:
    #     return "", Exception("wrong llm")


def llm_translate(input_text: str) -> str:
    translate_sys = "把输入完全翻译成英文，不要输出翻译文本以外的内容，只需要输出翻译后的文本。如果包含翻译之外的内容，则重新输出"
    return query_silicon_flow(input_text, translate_sys, 0.01)


# Example usage:
# result, error = llm_translate("你好")
# if error:
#     print(f"Error: {error}")
# else:
#     print(result)