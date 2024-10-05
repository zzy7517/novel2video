"use client";

import { useState } from "react";
import AIImageGenerator from "./image";
import TextEditor from "@/app/text";
import CharacterExtractor from "@/app/character";

export default function Home() {
    const [activeTab, setActiveTab] = useState("text");

    return (
        <div style={styles.container}>
            <div style={styles.sidebar}>
                <div style={styles.item} onClick={() => setActiveTab("text")}>
                    保存文本
                </div>
                <div style={styles.item} onClick={() => setActiveTab("character")}>
                    提取角色
                </div>
                <div style={styles.item} onClick={() => setActiveTab("image")}>
                    提取图像
                </div>
            </div>
            <div style={styles.content}>
                {activeTab === "text" && <TextEditor/>}
                {activeTab === "character" && <CharacterExtractor/>}
                {activeTab === "image" && <AIImageGenerator/>}
            </div>
        </div>
    );
}

const styles = {
    container: {
        display: "flex",
        height: "100vh",
    },
    sidebar: {
        width: "200px",
        padding: "10px",
        boxShadow: "2px 0 5px rgba(0,0,0,0.1)",
    },
    item: {
        padding: "10px",
        margin: "5px 0",
        cursor: "pointer",
        borderRadius: "5px",
        transition: "background-color 0.3s",
    },
    content: {
        flex: 1,
        padding: "20px",
    },
};
