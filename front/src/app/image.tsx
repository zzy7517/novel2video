"use client";
import React, { useState } from 'react';
import Image from "next/image";

export default function AIImageGenerator() {
    const [images, setImages] = useState<string[]>(["https://via.placeholder.com/400x300"]);
    const [fragments, setFragments] = useState<string[]>([""]);
    const [loaded, setLoaded] = useState<boolean>(false);

    const extractChapterFragments = () => {
        fetch('http://localhost:1198/api/get/novel/fragments')
            .then(response => response.json())
            .then(data => {
                setFragments(data);
                setImages(data.map(() => "https://via.placeholder.com/400x300"));
                setLoaded(true); // Set loaded to true after data is fetched
            })
            .catch(error => console.error('Error fetching file content:', error));
    };

    const generateImage = (index: number) => {
        console.log(`Generate image for fragment ${index + 1}`);
    };

    const mergeFragments = (index: number, direction: 'up' | 'down') => {
        if ((direction === 'up' && index === 0) || (direction === 'down' && index === fragments.length - 1)) {
            return; // Cannot merge if at the boundary
        }

        const newFragments = [...fragments];
        const newImages = [...images];

        if (direction === 'up') {
            newFragments[index - 1] += ' ' + newFragments[index];
            newFragments.splice(index, 1);
            newImages.splice(index, 1); // Remove the image for the merged fragment
        } else if (direction === 'down') {
            newFragments[index] += ' ' + newFragments[index + 1];
            newFragments.splice(index + 1, 1);
            newImages.splice(index + 1, 1); // Remove the image for the merged fragment
        }

        setFragments(newFragments);
        setImages(newImages);
    };

    return (
        <div className="container">
            <div className="header">
                <h1>AI Image Generator</h1>
            </div>
            <button onClick={extractChapterFragments}>分割章节</button>
            {loaded && (
                <>
                    {fragments.map((line, index) => (
                        <div key={index} className="card">
                            <div className="input-section">
                                <textarea value={line} readOnly rows={4} className="scrollable"></textarea>
                                <div className="checkbox-group">
                                    {index !== 0 && (
                                        <label>
                                            <input
                                                type="checkbox"
                                                onChange={() => mergeFragments(index, 'up')}
                                            /> Merge Up
                                        </label>
                                    )}
                                    {index !== fragments.length - 1 && (
                                        <label>
                                            <input
                                                type="checkbox"
                                                onChange={() => mergeFragments(index, 'down')}
                                            /> Merge Down
                                        </label>
                                    )}
                                </div>
                            </div>
                            <div className="description-section">
                                <textarea placeholder="Image prompt" rows={4} className="scrollable"></textarea>
                                <button onClick={() => generateImage(index)}>Generate Image</button>
                            </div>
                            <div className="image-section">
                                <Image
                                    src={images[0]} // Use the single placeholder image
                                    alt={`Generated image ${index + 1}`}
                                    width={300}
                                    height={200}
                                />
                            </div>
                        </div>
                    ))}
                    <button className="generate-all">Generate All Images</button>
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
                .card {
                    display: flex;
                    justify-content: space-between;
                    border: 1px solid #ddd;
                    border-radius: 8px;
                    padding: 20px;
                    margin-bottom: 20px;
                }
                .input-section, .description-section, .image-section {
                    width: 30%;
                }
                textarea {
                    width: 100%;
                    padding: 10px;
                    margin-bottom: 10px;
                    border: 1px solid #ddd;
                    border-radius: 4px;
                    resize: vertical;
                    overflow-y: auto;
                }
                .checkbox-group {
                    display: flex;
                    flex-direction: column;
                }
                .checkbox-group label {
                    margin-bottom: 5px;
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
                .generate-all {
                    display: block;
                    width: 100%;
                    margin-top: 20px;
                }
            `}</style>
        </div>
    );
}
