"use client"

import type React from "react"

import { useState, useRef, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { Mic, MicOff, Send } from "lucide-react"
import { motion, AnimatePresence } from "framer-motion"

// Mock data for a specific chat
const mockChatData = {
  1: {
    id: 1,
    name: "Sarah Chen",
    image: "/placeholder.svg?height=50&width=50",
    messages: [
      {
        id: 1,
        sender: "them",
        text: "Hi there! I saw your profile and I'm impressed with your React work.",
        time: "10:15 AM",
      },
      {
        id: 2,
        sender: "you",
        text: "Thanks! I've been working on some interesting projects lately.",
        time: "10:20 AM",
      },
      {
        id: 3,
        sender: "them",
        text: "I'm working on a new project that might interest you. It's a developer networking platform.",
        time: "10:25 AM",
      },
      { id: 4, sender: "you", text: "That sounds like a great project idea!", time: "10:30 AM" },
    ],
  },
  2: {
    id: 2,
    name: "Miguel Rodriguez",
    image: "/placeholder.svg?height=50&width=50",
    messages: [
      { id: 1, sender: "them", text: "Hey, have you worked with Kubernetes before?", time: "Yesterday" },
      { id: 2, sender: "you", text: "Yes, I've set up a few clusters for my projects.", time: "Yesterday" },
      { id: 3, sender: "them", text: "Great! I'm trying to optimize our deployment pipeline.", time: "Yesterday" },
      {
        id: 4,
        sender: "you",
        text: "I can help with that. Let me know what you're working on specifically.",
        time: "Yesterday",
      },
      { id: 5, sender: "them", text: "Can you share that article about microservices?", time: "Yesterday" },
    ],
  },
}

interface ChatInterfaceProps {
  chatId: string
}

export function ChatInterface({ chatId }: ChatInterfaceProps) {
  const [message, setMessage] = useState("")
  const [messages, setMessages] = useState<any[]>([])
  const [chatData, setChatData] = useState<any>(null)
  const [isRecording, setIsRecording] = useState(false)
  const [isRecognizing, setIsRecognizing] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // In a real app, this would fetch chat data from an API
    const chat = mockChatData[Number(chatId) as keyof typeof mockChatData]
    if (chat) {
      setChatData(chat)
      setMessages(chat.messages)
    }
  }, [chatId])

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }

  const handleSendMessage = () => {
    if (message.trim()) {
      const newMessage = {
        id: messages.length + 1,
        sender: "you",
        text: message,
        time: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
      }

      setMessages([...messages, newMessage])
      setMessage("")
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  const toggleRecording = () => {
    if (!isRecording) {
      // In a real app, this would use the Web Speech API
      setIsRecording(true)
      setIsRecognizing(true)

      // Simulate speech recognition after 3 seconds
      setTimeout(() => {
        setMessage((prev) => prev + "This is a simulated speech-to-text message.")
        setIsRecording(false)
        setIsRecognizing(false)
      }, 3000)
    } else {
      setIsRecording(false)
      setIsRecognizing(false)
    }
  }

  if (!chatData) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-purple-600 animate-pulse">Loading chat...</div>
      </div>
    )
  }

  return (
    <Card className="flex-1 flex flex-col h-[calc(100vh-150px)] backdrop-blur-sm bg-white/90 border-none shadow-xl">
      <CardHeader className="border-b bg-white/50">
        <div className="flex items-center space-x-4">
          <Avatar>
            <AvatarImage src={chatData.image || "/placeholder.svg"} alt={chatData.name} />
            <AvatarFallback>{chatData.name.substring(0, 2)}</AvatarFallback>
          </Avatar>
          <div>
            <h2 className="text-lg font-semibold text-purple-700">{chatData.name}</h2>
            <p className="text-sm text-gray-500">Online</p>
          </div>
        </div>
      </CardHeader>

      <CardContent className="flex-1 overflow-y-auto p-4 space-y-4">
        <AnimatePresence initial={false}>
          {messages.map((msg) => (
            <motion.div
              key={msg.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3 }}
              className={`flex ${msg.sender === "you" ? "justify-end" : "justify-start"}`}
            >
              <div
                className={`max-w-[70%] rounded-lg p-3 ${
                  msg.sender === "you"
                    ? "bg-gradient-to-r from-purple-600 to-indigo-600 text-white"
                    : "bg-white shadow-sm text-gray-800"
                }`}
              >
                <p>{msg.text}</p>
                <p className={`text-xs mt-1 ${msg.sender === "you" ? "text-purple-200" : "text-gray-500"}`}>
                  {msg.time}
                </p>
              </div>
            </motion.div>
          ))}
        </AnimatePresence>
        <div ref={messagesEndRef} />
      </CardContent>

      <div className="p-4 border-t bg-white/50">
        <div className="flex space-x-2">
          <Button
            variant="outline"
            size="icon"
            className={`rounded-full ${isRecording ? "bg-red-100 text-red-500 animate-pulse" : ""}`}
            onClick={toggleRecording}
          >
            {isRecording ? <MicOff className="h-5 w-5" /> : <Mic className="h-5 w-5" />}
          </Button>

          <Input
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={handleKeyPress}
            placeholder="Type a message..."
            className="flex-1 bg-white/50"
          />

          <Button
            variant="default"
            size="icon"
            className="rounded-full bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-700 hover:to-indigo-700"
            onClick={handleSendMessage}
            disabled={!message.trim()}
          >
            <Send className="h-5 w-5" />
          </Button>
        </div>

        {isRecognizing && (
          <div className="mt-2 text-sm text-center text-red-500 animate-pulse">Listening... Speak now</div>
        )}
      </div>
    </Card>
  )
}
