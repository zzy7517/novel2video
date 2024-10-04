"use client";
import React, { useState } from 'react';
import Image from "next/image";

export default function AIImageGenerator() {
  const [images, setImages] = useState(["/placeholder.svg?height=300&width=400"]);

  const generateImage = (index: number) => {
    const newImages = [...images];
    newImages[index] = `/placeholder.svg?height=300&width=400&text=Generated${index + 1}`;
    setImages(newImages);
  };

  const addImageBlock = () => {
    setImages([...images, "/placeholder.svg?height=300&width=400"]);
  };

  const removeImageBlock = () => {
    if (images.length > 1) {
      setImages(images.slice(0, -1));
    }
  };

  return (
    <div className="container">
      <div className="header">
        <h1>AI Image Generator Interface</h1>
        <div className="controls">
          <button onClick={removeImageBlock} disabled={images.length === 1}>-</button>
          <button onClick={addImageBlock}>+</button>
        </div>
      </div>
      {images.map((image, index) => (
        <div key={index} className="card">
          <div className="input-section">
            <input type="text" placeholder="Enter prompt" />
            <div className="checkbox-group">
              <label>
                <input type="checkbox" /> Option 1
              </label>
              <label>
                <input type="checkbox" /> Option 2
              </label>
            </div>
          </div>
          <div className="description-section">
            <textarea placeholder="Image description" rows={4}></textarea>
            <button onClick={() => generateImage(index)}>Generate Image</button>
          </div>
          <div className="image-section">
            <Image
              src={image}
              alt={`Generated image ${index + 1}`}
              width={300}
              height={200}
            />
          </div>
        </div>
      ))}
      <button className="generate-all">Generate All Images</button>
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
        .controls button {
          margin-left: 10px;
          padding: 5px 10px;
          font-size: 16px;
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
        input[type="text"], textarea {
          width: 100%;
          padding: 10px;
          margin-bottom: 10px;
          border: 1px solid #ddd;
          border-radius: 4px;
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