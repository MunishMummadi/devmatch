"use client"

import { useEffect, useRef } from "react"

export function AnimatedBackground() {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext("2d")
    if (!ctx) return

    // Set canvas to full screen
    const resizeCanvas = () => {
      canvas.width = window.innerWidth
      canvas.height = window.innerHeight
    }

    resizeCanvas()
    window.addEventListener("resize", resizeCanvas)

    // Wave parameters
    const waves = [
      {
        color: "rgba(236, 72, 153, 0.2)", // pink-500
        amplitude: 50,
        frequency: 0.01,
        speed: 0.02,
        phase: 0,
      },
      {
        color: "rgba(168, 85, 247, 0.15)", // purple-500
        amplitude: 70,
        frequency: 0.008,
        speed: 0.015,
        phase: 2,
      },
      {
        color: "rgba(139, 92, 246, 0.1)", // violet-500
        amplitude: 90,
        frequency: 0.006,
        speed: 0.01,
        phase: 4,
      },
    ]

    // Animation loop
    let animationFrameId: number
    const animate = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height)

      // Update wave phases
      waves.forEach((wave) => {
        wave.phase += wave.speed
      })

      // Draw waves
      waves.forEach((wave) => {
        ctx.fillStyle = wave.color
        ctx.beginPath()
        ctx.moveTo(0, canvas.height)

        for (let x = 0; x <= canvas.width; x += 10) {
          const y = canvas.height - 200 - wave.amplitude * Math.sin(wave.frequency * x + wave.phase)
          ctx.lineTo(x, y)
        }

        ctx.lineTo(canvas.width, canvas.height)
        ctx.closePath()
        ctx.fill()
      })

      animationFrameId = requestAnimationFrame(animate)
    }

    animate()

    return () => {
      window.removeEventListener("resize", resizeCanvas)
      cancelAnimationFrame(animationFrameId)
    }
  }, [])

  return (
    <canvas
      ref={canvasRef}
      className="fixed inset-0 w-full h-full -z-10 bg-gradient-to-br from-pink-50 to-purple-100"
    />
  )
}
