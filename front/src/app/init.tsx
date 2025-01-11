'use client'

import React, { useState, useEffect } from 'react'
import { showToast } from "@/app/toast"
import { ToastContainer } from "react-toastify"

export default function Component() {
  const [url, setAddress0] = useState('')
  const [apikey, setAddress1] = useState('')
  const [model, setModel] = useState('')
  const [address2, setAddress2] = useState('')
  const [address3, setAddress3] = useState('')
  const [address3Type, setAddress3Type] = useState('stable_diffusion_web_ui')
  const [comfyuiNodeApi, setComfyuiNodeApi] = useState('')
  const [savingStates, setSavingStates] = useState({
    url:false,
    apikey: false,
    model:false,
    address2: false,
    address3: false,
    comfyuiNodeApi: false
  })

  useEffect(() => {
    fetchSavedAddresses()
  }, [])

  const fetchSavedAddresses = async () => {
    try {
      showToast(`读取本地配置`)
      const response = await fetch('http://localhost:1198/api/model/config')
      if (response.ok) {
        const data = await response.json()
        setAddress0(data.url || '')
        setAddress1(data.apikey || '')
        setModel(data.model || '')
        setAddress2(data.address2 || '')
        setAddress3(data.address3 || '')
        setAddress3Type(data.address3Type || 'stable_diffusion_web_ui')
        setComfyuiNodeApi(JSON.stringify(data.comfyuiNodeApi) || '')
      } else {
        showToast(`读取本地配置出错`)
        console.error('Failed to fetch addresses')
      }
    } catch (error) {
      showToast(`读取本地配置出错 ${error}`)
      console.error('Error fetching addresses:', error)
    }
  }

  const saveAddress = async (key: 'url' | 'apikey' | 'model' | 'address2' | 'address3' | 'comfyuiNodeApi', value: string) => {
    setSavingStates(prev => ({ ...prev, [key]: true }))
    try {
      const response = await fetch('http://localhost:1198/api/model/config', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ key, value }),
      })

      if (response.ok) {
        showToast(`${key} 已成功保存`)
      } else {
        showToast(`保存 ${key} 时出错`)
      }
    } catch (error) {
      console.error(`Error saving ${key}:`, error)
      showToast(`保存 ${key} 时出错 ${error}`)
    } finally {
      setSavingStates(prev => ({ ...prev, [key]: false }))
    }
  }

  const saveAddress3Type = async (value: string) => {
    try {
      const response = await fetch('http://localhost:1198/api/model/config', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ key: 'address3Type', value }),
      })

      if (response.ok) {
        showToast(`保存成功`)
        setAddress3Type(value)
      } else {
        showToast(`保存出错`)
      }
    } catch (error) {
      console.error(`Error saving address type:`, error)
      showToast(`保存出错 ${error}`)
    }
  }

  return (
    <div className="w-full max-w-3xl mx-auto p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-6 text-gray-900">模型配置</h2>
      <div className="space-y-6">
        <div className="space-y-2">
          <label htmlFor="url" className="block text-sm font-medium text-gray-800">
            url
          </label>
          <div className="flex space-x-2">
            <input
                id="url"
                type="text"
                value={url}
                onChange={(e) => setAddress0(e.target.value)}
                placeholder="url"
                className="flex-grow px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 placeholder-gray-500"
            />
            <button
                onClick={() => saveAddress('url', url)}
                disabled={savingStates.url}
                className="px-4 py-2 bg-black text-white rounded-md shadow-sm hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {savingStates.url ? '保存中...' : '保存'}
            </button>
          </div>
          <p className="text-sm text-gray-600">默认使用sambanova提供的免费llama3.1-405b，也支持openai标准格式</p>
        </div>
        <div className="space-y-2">
          <label htmlFor="apikey" className="block text-sm font-medium text-gray-800">
            apikey
          </label>
          <div className="flex space-x-2">
            <input
                id="apikey"
                type="text"
                value={apikey}
                onChange={(e) => setAddress1(e.target.value)}
                placeholder="apikey"
                className="flex-grow px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 placeholder-gray-500"
            />
            <button
                onClick={() => saveAddress('apikey', apikey)}
                disabled={savingStates.apikey}
                className="px-4 py-2 bg-black text-white rounded-md shadow-sm hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {savingStates.apikey ? '保存中...' : '保存'}
            </button>
          </div>
          <p className="text-sm text-gray-600">默认的key，也可以去这个网站申请:
            https://cloud.sambanova.ai/apis，或者使用其他公司提供的服务</p>
        </div>
        <div className="space-y-2">
          <label htmlFor="model" className="block text-sm font-medium text-gray-800">
            model
          </label>
          <div className="flex space-x-2">
            <input
                id="model"
                type="text"
                value={model}
                onChange={(e) => setModel(e.target.value)}
                placeholder="url"
                className="flex-grow px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 placeholder-gray-500"
            />
            <button
                onClick={() => saveAddress('model', model)}
                disabled={savingStates.model}
                className="px-4 py-2 bg-black text-white rounded-md shadow-sm hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {savingStates.apikey ? '保存中...' : '保存'}
            </button>
          </div>
          <p className="text-sm text-gray-600">默认的可以用 Meta-Llama-3.1-405B-Instruct Qwen2.5-72B-Instruct 等等</p>
        </div>
        <div className="space-y-2">
          <label htmlFor="address2" className="block text-sm font-medium text-gray-800">
            硅基流动
          </label>
          <div className="flex space-x-2">
            <input
                id="address2"
                type="text"
                value={address2}
                onChange={(e) => setAddress2(e.target.value)}
                placeholder="apikey"
                className="flex-grow px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 placeholder-gray-500"
            />
            <button
                onClick={() => saveAddress('address2', address2)}
                disabled={savingStates.address2}
                className="px-4 py-2 bg-black text-white rounded-md shadow-sm hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {savingStates.address2 ? '保存中...' : '保存'}
            </button>
          </div>
          <p className="text-sm text-gray-600">可以使用默认的key，也可以去这个网站申请你自己的apikey:
            https://cloud.siliconflow.cn/account/ak，可以使用免费的小模型进行翻译 (需要实名认证）</p>
        </div>
        <div className="space-y-2">
          <label htmlFor="address3Type" className="block text-sm font-medium text-gray-800">
            文生图工具
          </label>
          <div className="flex space-x-2 mb-2">
            <select
                id="address3Type"
                value={address3Type}
                onChange={(e) => saveAddress3Type(e.target.value)}
                className="flex-grow px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900"
            >
              <option value="stable_diffusion_web_ui">Stable Diffusion Web Ui</option>
              <option value="comfyui">ComfyUI</option>
            </select>
          </div>
          <label htmlFor="address3" className="block text-sm font-medium text-gray-800">
            地址
          </label>
          <div className="flex space-x-2">
            <input
                id="address3"
                type="text"
                value={address3}
                onChange={(e) => setAddress3(e.target.value)}
                placeholder={address3Type === 'stable_diffusion_web_ui' ? 'Stable Diffusion Web Ui 地址' : 'ComfyUI 地址'}
                className="flex-grow px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 placeholder-gray-500"
            />
            <button
                onClick={() => saveAddress('address3', address3)}
                disabled={savingStates.address3}
                className="px-4 py-2 bg-black text-white rounded-md shadow-sm hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {savingStates.address3 ? '保存中...' : '保存'}
            </button>
          </div>
          {address3Type === 'comfyui' && (
              <div className="mt-2">
                <label htmlFor="comfyuiNodeApi" className="block text-sm font-medium text-gray-800">
                  ComfyUI API
                </label>
                <div className="flex flex-col space-y-2">
                <textarea
                    id="comfyuiNodeApi"
                    value={comfyuiNodeApi}
                    onChange={(e) => setComfyuiNodeApi(e.target.value)}
                    placeholder="在此粘贴 ComfyUI API JSON..."
                    className="w-full h-40 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 placeholder-gray-500 font-mono whitespace-pre resize-none"
                />
                  <button
                      onClick={() => saveAddress('comfyuiNodeApi', comfyuiNodeApi)}
                      disabled={savingStates.comfyuiNodeApi}
                      className="px-4 py-2 bg-black text-white rounded-md shadow-sm hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {savingStates.comfyuiNodeApi ? '保存中...' : '保存'}
                  </button>
                </div>
              </div>
          )}
          <p className="text-sm text-gray-600">地址可以是本地的，也可以是云端的，如果使用你自己的comfyuiapi，需要在节点里填prompt的地方加上占位符$prompt$，</p>
        </div>
      </div>
      <ToastContainer/>
    </div>
  )
}