'use client'

import React, { useState } from 'react'

interface Role {
    name: string
    description: string
}

export default function CharacterExtractor() {
    const [roles, setRoles] = useState<Record<string, Role>>({})
    const [editedDescriptions, setEditedDescriptions] = useState<Record<string, string>>({})
    const [isLoading, setIsLoading] = useState(false)

    const extractRoles = async () => {
        setIsLoading(true)
        try {
            // 这里应该是实际的API调用
            const response = await fetch('http://localhost:1198/api/novel/characters')
            const data = await response.json()
            setRoles(data)
            setEditedDescriptions({})
        } catch (error) {
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

    const saveChanges = async () => {
        setIsLoading(true)
        try {
            // 这里应该是实际的API调用来保存修改
            await fetch('/api/save-roles', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(editedDescriptions),
            })
            // 更新本地状态
            setRoles(prev => {
                const newRoles = { ...prev }
                Object.entries(editedDescriptions).forEach(([name, description]) => {
                    if (newRoles[name]) {
                        newRoles[name] = { ...newRoles[name], description }
                    }
                })
                return newRoles
            })
            setEditedDescriptions({})
        } catch (error) {
            console.error('Failed to save changes:', error)
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div style={{ fontFamily: 'Arial, sans-serif', maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
            <h1 style={{ textAlign: 'center', marginBottom: '20px' }}>角色提取器</h1>
            <button
                onClick={extractRoles}
                disabled={isLoading}
                style={{
                    padding: '10px 20px',
                    fontSize: '16px',
                    backgroundColor: '#0051a2',
                    color: 'white',
                    border: 'none',
                    borderRadius: '5px',
                    cursor: 'pointer',
                    marginBottom: '20px'
                }}
            >
                {isLoading ? '加载中...' : '提取角色'}
            </button>
            {Object.entries(roles).map(([name, role]) => (
                <div key={name} style={{ marginBottom: '20px', border: '1px solid #ddd', padding: '15px', borderRadius: '5px' }}>
                    <h3 style={{ marginTop: '0' }}>{name}</h3>
                    <textarea
                        value={editedDescriptions[name] ?? role.description}
                        onChange={(e) => handleDescriptionChange(name, e.target.value)}
                        style={{
                            width: '100%',
                            minHeight: '100px',
                            padding: '10px',
                            marginBottom: '10px',
                            borderRadius: '5px',
                            border: '1px solid #ddd'
                        }}
                    />
                </div>
            ))}
            {Object.keys(roles).length > 0 && (
                <button
                    onClick={saveChanges}
                    disabled={isLoading || Object.keys(editedDescriptions).length === 0}
                    style={{
                        padding: '10px 20px',
                        fontSize: '16px',
                        backgroundColor: '#008CBA',
                        color: 'white',
                        border: 'none',
                        borderRadius: '5px',
                        cursor: 'pointer'
                    }}
                >
                    {isLoading ? '保存中...' : '保存修改'}
                </button>
            )}
        </div>
    )
}