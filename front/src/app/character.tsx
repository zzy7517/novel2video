'use client'

import React, {useEffect, useState} from 'react'
import {showToast} from "@/app/toast";
import {ToastContainer} from "react-toastify";

export default function CharacterExtractor() {
    const [roles, setRoles] = useState<Record<string, string>>({})
    const [editedDescriptions, setEditedDescriptions] = useState<Record<string, string>>({})
    const [isLoading, setIsLoading] = useState(false)

    useEffect(() => {
        showToast("提取本地角色");
        extractRoles(true);
    }, []);

    const extractRoles = async (isLocal: boolean) => {
        setIsLoading(true)
        try {
            const endpoint = isLocal
                ? 'http://localhost:1198/api/novel/characters/local'
                : 'http://localhost:1198/api/novel/characters'
            const response = await fetch(endpoint)
            const data = await response.json()
            if (response.status == 40401) {
                showToast("本地没有角色");
                return
            }
            setRoles(data)
            setEditedDescriptions({})
        } catch (error) {
            showToast("失败");
            console.error('Failed to extract roles:', error)
        } finally {
            setIsLoading(false)
        }
    }

    const handleDescriptionChange = (roleName: string, newDescription: string) => {
        setEditedDescriptions(prev => ({
            ...prev,
            [roleName]: newDescription
        }))
    }

    const generateRandomDescription = async (roleName: string) => {
        try {
            const response = await fetch('http://localhost:1198/api/novel/characters/random')
            const randomDescription = await response.json()
            setEditedDescriptions(prev => ({
                ...prev,
                [roleName]: randomDescription
            }))
            showToast("成功");
        } catch (error) {
            showToast("失败");
            console.error('Failed to generate random description:', error)
        }
    }

    const saveChanges = async () => {
        setIsLoading(true)
        try {
            await fetch('http://localhost:1198/api/novel/characters', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(editedDescriptions),
            })
            setRoles(prev => {
                const newRoles = { ...prev }
                Object.entries(editedDescriptions).forEach(([name, description]) => {
                    if (newRoles[name]) {
                        newRoles[name] = description
                    }
                })
                return newRoles
            })
            showToast("成功");
        } catch (error) {
            showToast("失败");
            console.error('Failed to save changes:', error)
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div style={{
            fontFamily: 'Arial, sans-serif',
            maxWidth: '800px',
            margin: '0 auto',
            padding: '20px',
            backgroundColor: '#f7f7f7',
            borderRadius: '8px',
            boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)'
        }}>
            <h1 style={{
                textAlign: 'center',
                marginBottom: '20px',
                color: '#2c3e50'
            }}>生成中文prompts之后提取里面的角色，用于文生图时锁定人物</h1>
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                gap: '10px',
                marginBottom: '20px'
            }}>
                <button
                    onClick={() => extractRoles(true)}
                    disabled={isLoading}
                    style={{
                        padding: '10px 20px',
                        fontSize: '16px',
                        backgroundColor: '#27ae60',
                        color: 'white',
                        border: 'none',
                        borderRadius: '5px',
                        cursor: 'pointer',
                        transition: 'background-color 0.3s'
                    }}
                >
                    {isLoading ? '加载中...' : '提取本地描述'}
                </button>
                <button
                    onClick={() => extractRoles(false)}
                    disabled={isLoading}
                    style={{
                        padding: '10px 20px',
                        fontSize: '16px',
                        backgroundColor: '#3498db',
                        color: 'white',
                        border: 'none',
                        borderRadius: '5px',
                        cursor: 'pointer',
                        transition: 'background-color 0.3s'
                    }}
                >
                    {isLoading ? '加载中...' : '提取角色'}
                </button>
            </div>
            {Object.entries(roles).map(([name, description]) => (
                <div key={name} style={{
                    marginBottom: '20px',
                    border: '1px solid #bdc3c7',
                    padding: '15px',
                    borderRadius: '5px',
                    backgroundColor: 'white'
                }}>
                    <h3 style={{
                        marginTop: '0',
                        color: '#34495e'
                    }}>{name}</h3>
                    <textarea
                        value={editedDescriptions[name] ?? description}
                        onChange={(e) => handleDescriptionChange(name, e.target.value)}
                        style={{
                            width: '100%',
                            minHeight: '100px',
                            padding: '10px',
                            marginBottom: '10px',
                            borderRadius: '5px',
                            border: '1px solid #bdc3c7',
                            backgroundColor: '#ecf0f1',
                            color: '#2c3e50',
                            fontSize: '14px',
                            resize: 'vertical'
                        }}
                    />
                    <button
                        onClick={() => generateRandomDescription(name)}
                        style={{
                            padding: '5px 10px',
                            fontSize: '14px',
                            backgroundColor: '#e67e22',
                            color: 'white',
                            border: 'none',
                            borderRadius: '5px',
                            cursor: 'pointer',
                            transition: 'background-color 0.3s'
                        }}
                    >
                        生成随机描述
                    </button>
                </div>
            ))}
            {Object.keys(roles).length > 0 && (
                <button
                    onClick={saveChanges}
                    disabled={isLoading || Object.keys(editedDescriptions).length === 0}
                    style={{
                        padding: '10px 20px',
                        fontSize: '16px',
                        backgroundColor: '#9b59b6',
                        color: 'white',
                        border: 'none',
                        borderRadius: '5px',
                        cursor: 'pointer',
                        transition: 'background-color 0.3s',
                        display: 'block',
                        margin: '0 auto'
                    }}
                >
                    {isLoading ? '保存中...' : '保存修改'}
                </button>
            )}
            <ToastContainer />
        </div>
    )
}