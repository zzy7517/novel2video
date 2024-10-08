"use client"

import React, { useState, useEffect } from 'react';
import Image from "next/image";

export default function AIImageGenerator() {
    const [images, setImages] = useState<string[]>([]);
    const [fragments, setFragments] = useState<string[]>([]);
    const [prompts, setPrompts] = useState<string[]>([]);
    const [loaded, setLoaded] = useState<boolean>(false);
    const [promptsEn, setPromptsEn] = useState<string[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(false);

    useEffect(() => {
        initialize();
    }, []);

    const initialize = () => {
        fetch('http://localhost:1198/api/novel/initial')
            .then(response => response.json())
            .then(data => {
                setFragments(data.fragments || []);
                const updatedImages = (data.images || []).map((imageUrl:string) => `http://localhost:1198${imageUrl}`);
                setImages(updatedImages);
                setPrompts(data.prompts || []);
                setPromptsEn(data.promptsEn)
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
        fetch('http://localhost:1198/api/get/novel/prompts')
            .then(response => response.json())
            .then(data => {
                setPrompts(data);
            })
            .catch(error => console.error('Error fetching prompts:', error));
    };

    const generateImage = (index: number) => {
        console.log(`Generate image for fragment ${index + 1}`);
    };

    const generateAllImages = () => {
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
            .catch(error => console.error('Error generating all images:', error));
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

        if (direction === 'up') {
            newFragments[index - 1] += ' ' + newFragments[index];
            newFragments.splice(index, 1);
            newImages.splice(index, 1);
        } else if (direction === 'down') {
            newFragments[index] += ' ' + newFragments[index + 1];
            newFragments.splice(index + 1, 1);
            newImages.splice(index + 1, 1);
        }

        setFragments(newFragments);
        setImages(newImages);
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
        } catch (error) {
            console.error('Error saving attachment:', error);
        }
    };

    const generatePromptsEn = () => {
        fetch('http://localhost:1198/api/novel/prompts/en')
            .then(response => response.json())
            .then(data => {
                setPromptsEn(data);
            })
            .catch(error => console.error('Error fetching prompts:', error));
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
                            <div className="description-section">
                                <textarea
                                    value={prompts[index] || ''}
                                    placeholder="prompt"
                                    rows={4}
                                    className="scrollable"
                                    readOnly
                                ></textarea>
                                <button onClick={() => generateImage(index)}>Generate Image</button>
                            </div>
                            <div className="attachment-section">
                                <textarea
                                    value={promptsEn[index]}
                                    onChange={(e) => {
                                        const newAttachments = [...promptsEn];
                                        newAttachments[index] = e.target.value;
                                        setPromptsEn(newAttachments);
                                    }}
                                    placeholder="Attachment"
                                    rows={4}
                                    className="scrollable"
                                ></textarea>
                                <button onClick={() => savePromptEn(index)}>Save Attachment</button>
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
                .input-section, .description-section, .attachment-section, .image-section {
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
                .generate-all, .refresh-images, .generate-promptsEn {
                    padding: 10px 20px;
                    font-size: 16px;
                }
            `}</style>
        </div>
    );
}