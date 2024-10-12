'use client'

import { useState, useEffect } from 'react'

export default function TextEditor() {
    const [content, setContent] = useState('')
    const [message, setMessage] = useState('')

    useEffect(() => {
        loadContent()
    }, [])

    const loadContent = async () => {
        try {
            const response = await fetch(`http://localhost:1198/api/novel/load`)
            if (response.ok) {
                const data = await response.json()
                if (data.content) {
                    setContent(data.content)
                    setMessage('内容已加载')
                } else {
                    setContent('')
                    setMessage('没有找到保存的内容')
                }
            } else {
                throw new Error('加载失败')
            }
        } catch (error) {
            setMessage('加载失败，请稍后重试。')
        }
    }

    const handleSave = async () => {
        try {
            const response = await fetch('http://localhost:1198/api/novel/save', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ content }),
            })

            if (response.ok) {
                setMessage('保存成功！')
            } else {
                throw new Error('保存失败')
            }
        } catch (error) {
            setMessage('保存失败，请稍后重试。')
        }
    }

    return (
        <div style={{
            fontFamily: 'Arial, sans-serif',
            maxWidth: '600px',
            margin: '40px auto',
            padding: '20px',
            boxShadow: '0 0 10px rgba(0,0,0,0.1)',
            borderRadius: '8px',
        }}>
            <h1 style={{ textAlign: 'center', marginBottom: '20px' }}>小说文本</h1>
            <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                style={{
                    width: '100%',
                    height: '300px',
                    padding: '10px',
                    marginBottom: '20px',
                    border: '1px solid #ccc',
                    borderRadius: '4px',
                    resize: 'vertical',
                }}
                placeholder="在这里输入您的文本..."
            />
            <button
                onClick={handleSave}
                style={{
                    display: 'block',
                    width: '100%',
                    padding: '10px',
                    backgroundColor: '#007bff',
                    color: 'white',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: 'pointer',
                }}
            >
                保存
            </button>
            {message && (
                <p style={{
                    marginTop: '20px',
                    textAlign: 'center',
                    color: message.includes('成功') ? 'green' : 'red',
                }}>
                    {message}
                </p>
            )}
        </div>
    )
}