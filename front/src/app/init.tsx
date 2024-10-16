'use client'

import React, { useState, useEffect } from 'react'
import {showToast} from "@/app/toast";
import {ToastContainer} from "react-toastify";

export default function Component() {
  const [address1, setAddress1] = useState('')
  const [address2, setAddress2] = useState('')
  const [address3, setAddress3] = useState('')
  const [savingStates, setSavingStates] = useState({
    address1: false,
    address2: false,
    address3: false
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
        setAddress1(data.address1 || '')
        setAddress2(data.address2 || '')
        setAddress3(data.address3 || '')
      } else {
        showToast(`读取本地配置出错`)
        console.error('Failed to fetch addresses')
      }
    } catch (error) {
      showToast(`读取本地配置出错 ${error}`)
      console.error('Error fetching addresses:', error)
    }
  }

  const saveAddress = async (key: 'address1' | 'address2' | 'address3', value: string) => {
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

  return (
    <div className="w-full max-w-3xl mx-auto p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-6 text-gray-900">模型配置</h2>
      <div className="space-y-6">
        <div className="space-y-2">
          <label htmlFor="address1" className="block text-sm font-medium text-gray-800">
            SambaNova
          </label>
          <div className="flex space-x-2">
            <input
              id="address1"
              type="text"
              value={address1}
              onChange={(e) => setAddress1(e.target.value)}
              placeholder="apikey"
              className="flex-grow px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 placeholder-gray-500"
            />
            <button
              onClick={() => saveAddress('address1', address1)}
              disabled={savingStates.address1}
              className="px-4 py-2 bg-black text-white rounded-md shadow-sm hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {savingStates.address1 ? '保存中...' : '保存'}
            </button>
          </div>
          <p className="text-sm text-gray-600">可以使用默认的key，也可以去这个网站申请你自己的apikey:
            https://cloud.sambanova.ai/apis，可以使用免费的llama3.1-405b</p>
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
          <label htmlFor="address3" className="block text-sm font-medium text-gray-800">
            stable diffusion 地址
          </label>
          <div className="flex space-x-2">
            <input
              id="address3"
              type="text"
              value={address3}
              onChange={(e) => setAddress3(e.target.value)}
              placeholder="sd address"
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
          <p className="text-sm text-gray-600">可以是本地的，也可以是云端的</p>
        </div>
      </div>
      <ToastContainer />
    </div>
  )
}