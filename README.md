## 项目介绍

Novel2Video 是一个工具，旨在将小说内容批量转换为图片和音频，最终生成小说推文。

- 使用免费的llama-3.1-405b提取小说场景
- 支持 Stable Diffusion web UI 和 ComfyUI
- 支持人物锁定，确保角色一致性
- 支持批量出图以及单张重绘
- 使用 EdgeTTS 进行文本到语音转换

## 安装

```bash
# 克隆仓库
git clone https://github.com/zzy7517/novel2video.git
cd novel2video

# 后端 python版本3.10以上就行
pip install -r requirements.txt
python main.py

# 前端
## 到下面这个地址下载nvm-setup.exe
https://github.com/coreybutler/nvm-windows/releases

## 按照安装向导的指示进行安装
使用命令 nvm install <version> 来安装特定版本的 Node.js。例如，安装最新的 LTS 版本可以使用 nvm install lts

## 安装依赖
npm install next --registry=https://registry.npmmirror.com
npm install toastify-js --registry=https://registry.npmmirror.com

## 运行
cd front
npm run dev
```

## 使用说明
以comfyui为例 <br>
1. 如图所示，先保存一下你的配置 <br>
<img width="564" alt="Snipaste_2024-10-21_22-35-46" src="https://github.com/user-attachments/assets/00feb1d9-6213-425d-8747-90bd64566cd9"> <br>
2. 然后在保存文本页面 保存你的小说文本和提示词, 提示词用来提取小说的场景 <br>
![image](https://github.com/user-attachments/assets/d5dc1a80-5db4-4e00-b722-1c959fcb32a9)
3. 先点击分割章节，分割完成后，点击提取中文的prompts，刚刚保存的提示词就作用在这里，理论上每段分割的文本会有一个中文提示词 <br>
![image](https://github.com/user-attachments/assets/15b49f2f-4924-4115-8051-a3cc3b2dc1b9)
4. 为了保证人物的一致性，需要写死角色，这一步一定要在第3步之后，如果之前没有生成过角色，点击 '提取角色' 按钮 <br>
![image](https://github.com/user-attachments/assets/d0ecd807-eba1-406f-9f47-ec9ae103ee94)
5. 配置好角色之后，点击 '翻译成英文' 按钮, 这个时候可以点击 '一键生成'  或者 '重新生成'，生成全部或者单张图片，'一键生成' 的过程中，可以点击 '刷新' 按钮加载本地的图片  <br>
![image](https://github.com/user-attachments/assets/f5496226-0876-4d4b-8c3d-ca2d55089947)
6. 与此同时，可以点击 '生成音频' 生成声音 <br>
7. 生成的文本文件/图像/音频都在temp目录下
8. 等文本和语音生成完成后，可以一键生成视频，后续的字幕，BGM等视频处理可以使用剪映

- [ ] TBD
    - [ ] 一键反推
    - [ ] midjourney 支持
    - [ ] 更丰富的语音合成功能
