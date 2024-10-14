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

  return (
    <div className="max-w-full w-[90vw] mx-auto my-10 px-4 py-6 bg-gray-100 rounded-lg shadow-lg">
      <div className="mb-8 bg-blue-50 p-4 rounded-lg">
        <h2 className="text-3xl font-bold text-center mb-4 text-blue-800 border-b-2 border-blue-300 pb-2">小说</h2>
        <textarea
          value={novelContent}
          onChange={(e) => setNovelContent(e.target.value)}
          placeholder="在这里输入小说文本..."
          className="w-full min-h-[400px] p-4 mb-4 border border-blue-300 rounded-md text-gray-800 bg-white resize-vertical focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg"
        />
        <button
          onClick={() => handleSave('novel')}
          className="w-full py-3 bg-blue-600 text-white border-none rounded-md cursor-pointer text-lg hover:bg-blue-700 transition-colors"
        >
          保存小说
        </button>
        {novelMessage && (
          <p className={`text-center mt-4 ${novelMessage.includes('成功') ? 'text-green-600' : 'text-red-600'}`}>
            {novelMessage}
          </p>
        )}
      </div>

      <div className="bg-green-50 p-4 rounded-lg">
        <h2 className="text-3xl font-bold text-center mb-4 text-green-800 border-b-2 border-green-300 pb-2">提示词</h2>
        <textarea
          value={promptContent}
          onChange={(e) => setPromptContent(e.target.value)}
          placeholder="在这里输入提示词..."
          className="w-full min-h-[400px] p-4 mb-4 border border-green-300 rounded-md text-gray-800 bg-white resize-vertical focus:outline-none focus:ring-2 focus:ring-green-500 text-lg"
        />
        <button
          onClick={() => handleSave('prompt')}
          className="w-full py-3 bg-green-600 text-white border-none rounded-md cursor-pointer text-lg hover:bg-green-700 transition-colors"
        >
          保存提示词
        </button>
        {promptMessage && (
          <p className={`text-center mt-4 ${promptMessage.includes('成功') ? 'text-green-600' : 'text-red-600'}`}>
            {promptMessage}
          </p>
        )}
      </div>
    </div>
  )
}