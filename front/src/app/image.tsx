"use client"

import React, { useState, useEffect } from 'react';
import Image from "next/image";
import {showToast} from "@/app/toast";
import {ToastContainer} from "react-toastify";

export default function AIImageGenerator() {
    const [images, setImages] = useState<string[]>([]);
    const [fragments, setFragments] = useState<string[]>([]);
    const [prompts, setPrompts] = useState<string[]>([]);
    const [loaded, setLoaded] = useState<boolean>(false);
    const [promptsEn, setPromptsEn] = useState<string[]>([]);
    const [isLoading] = useState<boolean>(false);

    useEffect(() => {
        initialize();
    }, []);

    const addCacheBuster = (url: string) => {
        const cacheBuster = `?v=${Date.now()}`
        return url.includes('?') ? `${url}&${cacheBuster}` : `${url}${cacheBuster}`
  }

    const initialize = () => {
        fetch('http://localhost:1198/api/novel/initial')
            .then(response => response.json())
            .then(data => {
                setFragments(data.fragments || []);
                const updatedImages = (data.images || []).map((imageUrl: string) =>
                  addCacheBuster(`http://localhost:1198${imageUrl}`)
                )
                setImages(updatedImages);
                setPrompts(data.prompts || []);
                setPromptsEn(data.promptsEn || [])
                setLoaded(true);
            })
            .catch(error => {
                console.error('Error initializing data:', error);
                setLoaded(false);
            });
    };

    const extractChapterFragments = () => {
        fetch('http://localhost:1198/api/get/novel/fragments')
            .then(response => response.json())
            .then(data => {
                setFragments(data);
                setImages(data.map(() => "http://localhost:1198/images/placeholder.png"));
                setLoaded(true);
            })
            .catch(error => {
                console.error('Error fetching file content:', error);
                setLoaded(false);
            });
    };

    const extractPrompts = () => {
        showToast('开始生成');
        fetch('http://localhost:1198/api/get/novel/prompts')
            .then(response => response.json())
            .then(data => {
            setPrompts(data || []);
            console.log('Prompts fetched successfully');
        })
        .catch(error => {
            console.error('Error fetching prompts:', error);
            showToast('失败');
        });
    };

    const generateAllImages = () => {
        showToast('开始生成，请等待');
        fetch('http://localhost:1198/api/novel/images', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
        })
            .then(response => response.json())
            .then(() => {
                console.log('Images generation initiated');
                refreshImages();
            })
            .catch(error => {
                showToast('失败，请检查日志');
                console.error('Error generating all images:', error)
            });
    };

    type ImageMap = Record<string, string>;

    const refreshImages = () => {
        fetch('http://localhost:1198/api/novel/images')
            .then(response => response.json() as Promise<ImageMap>)
            .then((imageMap: ImageMap) => {
                const updatedImages = [...images];
                for (const [index, imageUrl] of Object.entries(imageMap)) {
                    const numericIndex = Number(index);
                    if (!isNaN(numericIndex)) {
                        updatedImages[numericIndex] = `http://localhost:1198${imageUrl}`;
                    }
                }
                setImages(updatedImages);
            })
            .catch(error => console.error('Error fetching image:', error));
    };

    const mergeFragments = (index: number, direction: 'up' | 'down') => {
        if ((direction === 'up' && index === 0) || (direction === 'down' && index === fragments.length - 1)) {
            return;
        }

        const newFragments = [...fragments];
        const newImages = [...images];
        const newPrompts = [...prompts]
        const newPromptsEn = [...promptsEn]

        if (direction === 'up') {
            newFragments[index - 1] += ' ' + newFragments[index];
            newFragments.splice(index, 1);
            newImages.splice(index, 1);
            newPrompts.splice(index, 1)
            newPromptsEn.splice(index, 1)
        } else if (direction === 'down') {
            newFragments[index] += ' ' + newFragments[index + 1];
            newFragments.splice(index + 1, 1);
            newImages.splice(index + 1, 1);
            newPrompts.splice(index+1, 1)
            newPromptsEn.splice(index+1, 1)
        }

        setFragments(newFragments);
        setImages(newImages);
        setPromptsEn(prompts)
        setPrompts(prompts)
        // todo 是不是最好都重新保存一下
        saveFragments(newFragments);
    };

    const saveFragments = async (fragments: string[]) => {
        try {
            const response = await fetch('http://localhost:1198/api/save/novel/fragments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(fragments),
            });

            if (!response.ok) {
                throw new Error('Failed to save fragments');
            }
            console.log('Fragments saved successfully');
        } catch (error) {
            console.error('Error saving fragments:', error);
        }
    };

    const savePromptEn = async (index: number) => {
        try {
            const response = await fetch('http://localhost:1198/api/novel/prompt/en', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    index: index,
                    content: promptsEn[index]
                }),
            });

            if (!response.ok) {
                throw new Error('Failed to save attachment');
            }
            console.log(`Attachment for fragment ${index + 1} saved successfully.`);
            showToast('保存成功');
        } catch (error) {
            console.error('Error saving attachment:', error);
            showToast('保存失败');
        }
    };

    const savePromptZh = async  (index: number) => {
        try {
            const response = await fetch('http://localhost:1198/api/novel/prompt/zh', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    index: index,
                    content: prompts[index]
                }),
            });

            if (!response.ok) {
                throw new Error('Failed to save attachment');
            }
            console.log(`Attachment for fragment ${index + 1} saved successfully.`);
            showToast('保存成功');
        } catch (error) {
            console.error('Error saving attachment:', error);
            showToast('保存失败');
        }
    };

    const generatePromptsEn = () => {
        showToast('开始生成，请等待');
        fetch('http://localhost:1198/api/novel/prompts/en')
        .then(response => response.json())
        .then(data => {
            setPromptsEn(data || []);
        })
        .catch(error => {
            showToast('失败');
            console.error('Error fetching prompts:', error)
        });
    };

    const generateAudio = () => {
        showToast('开始生成，请等待');
        fetch('http://localhost:1198/api/novel/audio', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ fragments })
        })
            .then(response => response.json())
            .then(() => {
                console.log('Audio generation initiated');
            })
            .catch(error => {
                console.error('Error generating audio:', error);
                showToast('失败');
            });
    };

    return (
        <div className="container">
            <div className="header">
                <h1>AI Image Generator</h1>
            </div>
            <div className="button-container">
                <button onClick={extractChapterFragments}>分割章节</button>
                {loaded && (
                    <>
                        <button onClick={extractPrompts} className="extract-prompts-button">提取文生图prompts</button>
                        <button onClick={generatePromptsEn} className="generate-promptsEn" disabled={isLoading}>
                            {isLoading ? 'Generating...' : 'Translate Prompts'}
                        </button>
                        <button onClick={generateAllImages} className="generate-all">一键生成</button>
                        <button onClick={initialize} className="refresh-images">刷新</button>
                        <button onClick={generateAudio} className="generate-audio">生成音频</button>
                    </>
                )}
            </div>
            {loaded && (
                <>
                    {fragments.map((line, index) => (
                        <div key={index} className="card">
                            <div className="input-section">
                                <textarea value={line} readOnly rows={4} className="scrollable"></textarea>
                                <div className="button-group">
                                    {index !== 0 && (
                                        <button className="merge-button" onClick={() => mergeFragments(index, 'up')}>Merge Up</button>
                                    )}
                                    {index !== fragments.length - 1 && (
                                        <button className="merge-button" onClick={() => mergeFragments(index, 'down')}>Merge Down</button>
                                    )}
                                </div>
                            </div>
                            <div className="prompt-section">
                                <textarea
                                    value={prompts[index] || ''}
                                    placeholder="prompt"
                                    onChange={(e) => {
                                        const newAttachments = [...prompts];
                                        newAttachments[index] = e.target.value;
                                        setPrompts(newAttachments);
                                    }}
                                    rows={4}
                                    className="scrollable"
                                ></textarea>
                                <button onClick={() => savePromptZh(index)}>保存</button>
                            </div>
                            <div className="promptEn-section">
                                <textarea
                                    value={promptsEn[index] || ''}
                                    onChange={(e) => {
                                        const newAttachments = [...promptsEn];
                                        newAttachments[index] = e.target.value;
                                        setPromptsEn(newAttachments);
                                    }}
                                    placeholder="Attachment"
                                    rows={4}
                                    className="scrollable"
                                ></textarea>
                                <button onClick={() => savePromptEn(index)}>保存</button>
                            </div>
                            <div className="image-section">
                                <Image
                                    src={images[index]}
                                    key={images[index]}
                                    alt={`Generated image ${index + 1}`}
                                    width={300}
                                    height={200}
                                />
                            </div>
                            <ToastContainer />
                        </div>
                    ))}
                </>
            )}
            <style jsx>{`
                .container {
                    max-width: 1200px;
                    margin: 0 auto;
                    padding: 20px;
                    font-family: Arial, sans-serif;
                }
                .header {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    margin-bottom: 20px;
                }
                .button-container {
                    display: flex;
                    gap: 20px;
                    margin-bottom: 20px;
                }
                .card {
                    display: flex;
                    justify-content: space-between;
                    border: 1px solid #ddd;
                    border-radius: 8px;
                    padding: 20px;
                    margin-bottom: 20px;
                }
                .input-section, .prompt-section, .promptEn-section, .image-section {
                    width: 23%;
                }
                textarea {
                    width: 100%;
                    padding: 10px;
                    margin-bottom: 10px;
                    border: 1px solid #ddd;
                    border-radius: 4px;
                    resize: vertical;
                    overflow-y: auto;
                    color: #333;
                    background-color: #fff;
                }
                .button-group {
                    display: flex;
                    flex-direction: column;
                }
                .button-group .merge-button {
                    margin-bottom: 5px;
                    padding: 5px 10px;
                    font-size: 14px;
                }
                button {
                    background-color: #0070f3;
                    color: white;
                    border: none;
                    padding: 10px 20px;
                    border-radius: 4px;
                    cursor: pointer;
                    font-size: 16px;
                }
                button:hover {
                    background-color: #0051a2;
                }
                button:disabled {
                    background-color: #ccc;
                    cursor: not-allowed;
                }
                .generate-all, .refresh-images, .generate-promptsEn, .generate-audio {
                    padding: 10px 20px;
                    font-size: 16px;
                }
            `}</style>
        </div>
    );
}