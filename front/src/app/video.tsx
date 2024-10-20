'use client'

import { useState, useEffect } from 'react'

export default function VideoGenerator() {
  const [videoUrl, setVideoUrl] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchInitialVideo()
  }, [])

  const fetchInitialVideo = async () => {
    try {
      const response = await fetch('http://localhost:1198/api/novel/video', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error('Failed to fetch initial video')
      }

      const data = await response.json()
      if (data.videoUrl) {
        setVideoUrl(`http://localhost:1198${data.videoUrl}`)
      }
    } catch (err) {
      console.error('Error fetching initial video:', err)
    }
  }

  const generateVideo = async () => {
    setIsLoading(true)
    setError(null)

    try {
      const response = await fetch('http://localhost:1198/api/novel/video', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        // Add any necessary body parameters here
        body: JSON.stringify({}),
      })

      if (!response.ok) {
        throw new Error('Failed to generate video')
      }

      const data = await response.json()
      setVideoUrl(data.videoUrl)
    } catch (err) {
      setError('An error occurred while generating the video')
      console.error(err)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex flex-col items-center space-y-4 p-4 bg-gray-100">
      <button
        onClick={generateVideo}
        disabled={isLoading}
        className="px-4 py-2 bg-gray-900 text-white rounded hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isLoading ? 'Generating Video...' : 'Generate New Video'}
      </button>

      {error && (
        <p className="text-gray-900">{error}</p>
      )}

      {videoUrl && (
        <div className="mt-4 w-full max-w-md">
          <video controls className="w-full h-auto">
            <source src={videoUrl} type="video/mp4" />
            Your browser does not support the video tag.
          </video>
        </div>
      )}
    </div>
  )
}