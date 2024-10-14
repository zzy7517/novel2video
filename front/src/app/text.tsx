'use client'

import { useState, useEffect } from 'react'

export default function TextEditor() {
    const [novelContent, setNovelContent] = useState('')
    const [promptContent, setPromptContent] = useState('')
    const [novelMessage, setNovelMessage] = useState('')
    const [promptMessage, setPromptMessage] = useState('')

    useEffect(() => {
        loadContent('novel')
        loadContent('prompt')
    }, [])

    const loadContent = async (type: 'novel' | 'prompt') => {
        try {
            const response = await fetch(`http://localhost:1198/api/${type}/load`)
            if (response.ok) {
                const data = await response.json()
                if (data.content) {
                    type === 'novel' ? setNovelContent(data.content) : setPromptContent(data.content)
                    setMessage(type, '内容已加载')
                } else {
                    type === 'novel' ? setNovelContent('') : setPromptContent('')
                    setMessage(type, '没有找到保存的内容')
                }
            } else {
                throw new Error('加载失败')
            }
        } catch (error) {
            setMessage(type, '加载失败，请稍后重试。')
        }
    }

    const handleSave = async (type: 'novel' | 'prompt') => {
        try {
            const content = type === 'novel' ? novelContent : promptContent
            const response = await fetch(`http://localhost:1198/api/${type}/save`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ content }),
            })

            if (response.ok) {
                setMessage(type, '保存成功！')
            } else {
                throw new Error('保存失败')
            }
        } catch (error) {
            setMessage(type, '保存失败，请稍后重试。')
        }
    }

    const setMessage = (type: 'novel' | 'prompt', message: string) => {
        type === 'novel' ? setNovelMessage(message) : setPromptMessage(message)
    }

    const containerStyle: React.CSSProperties = {
        maxWidth: '800px',
        margin: '40px auto',
        padding: '24px',
        backgroundColor: 'white',
        boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
        borderRadius: '8px',
    }

    const sectionStyle: React.CSSProperties = {
        marginBottom: '32px',
    }

    const headingStyle: React.CSSProperties = {
        fontSize: '24px',
        fontWeight: 'bold',
        textAlign: 'center',
        marginBottom: '16px',
    }

    const textareaStyle: React.CSSProperties = {
        width: '100%',
        minHeight: '200px',
        padding: '12px',
        marginBottom: '16px',
        border: '1px solid #ccc',
        borderRadius: '4px',
        resize: 'vertical',
    }

    const buttonStyle: React.CSSProperties = {
        width: '100%',
        padding: '12px',
        backgroundColor: '#007bff',
        color: 'white',
        border: 'none',
        borderRadius: '4px',
        cursor: 'pointer',
        fontSize: '16px',
    }

    const messageStyle = (message: string): React.CSSProperties => ({
        textAlign: 'center',
        color: message.includes('成功') ? 'green' : 'red',
        marginTop: '16px',
    })

    return (
        <div style={containerStyle}>
            <div style={sectionStyle}>
                <h2 style={headingStyle}>小说文本</h2>
                <textarea
                    value={novelContent}
                    onChange={(e) => setNovelContent(e.target.value)}
                    placeholder="在这里输入您的小说文本..."
                    style={textareaStyle}
                />
                <button
                    onClick={() => handleSave('novel')}
                    style={buttonStyle}
                >
                    保存小说
                </button>
                {novelMessage && (
                    <p style={messageStyle(novelMessage)}>
                        {novelMessage}
                    </p>
                )}
            </div>

            <div style={sectionStyle}>
                <h2 style={headingStyle}>提示文本</h2>
                <textarea
                    value={promptContent}
                    onChange={(e) => setPromptContent(e.target.value)}
                    placeholder="在这里输入您的提示文本..."
                    style={textareaStyle}
                />
                <button
                    onClick={() => handleSave('prompt')}
                    style={buttonStyle}
                >
                    保存提示
                </button>
                {promptMessage && (
                    <p style={messageStyle(promptMessage)}>
                        {promptMessage}
                    </p>
                )}
            </div>
        </div>
    )
}