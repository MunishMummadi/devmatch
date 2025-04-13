"use client"
import { useRouter } from "next/navigation"
import { Card, CardContent } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { motion } from "framer-motion"

// Mock data for chats
const mockChats = [
  {
    id: 1,
    name: "Sarah Chen",
    image: "/placeholder.svg?height=50&width=50",
    lastMessage: "That sounds like a great project idea!",
    time: "10:30 AM",
    unread: 2,
  },
  {
    id: 2,
    name: "Miguel Rodriguez",
    image: "/placeholder.svg?height=50&width=50",
    lastMessage: "Can you share that article about microservices?",
    time: "Yesterday",
    unread: 0,
  },
  {
    id: 3,
    name: "Priya Patel",
    image: "/placeholder.svg?height=50&width=50",
    lastMessage: "I'd love to collaborate on your ML project",
    time: "Yesterday",
    unread: 0,
  },
  {
    id: 4,
    name: "David Kim",
    image: "/placeholder.svg?height=50&width=50",
    lastMessage: "Let's catch up at the conference next week",
    time: "Monday",
    unread: 0,
  },
]

export function ChatList() {
  const router = useRouter()

  const handleChatClick = (chatId: number) => {
    router.push(`/dashboard/chats/${chatId}`)
  }

  return (
    <div className="space-y-4">
      {mockChats.map((chat, index) => (
        <motion.div
          key={chat.id}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: index * 0.1 }}
        >
          <Card
            className="cursor-pointer hover:shadow-md transition-shadow backdrop-blur-sm bg-white/90 border-none"
            onClick={() => handleChatClick(chat.id)}
          >
            <CardContent className="p-4">
              <div className="flex items-center space-x-4">
                <Avatar>
                  <AvatarImage src={chat.image || "/placeholder.svg"} alt={chat.name} />
                  <AvatarFallback>{chat.name.substring(0, 2)}</AvatarFallback>
                </Avatar>

                <div className="flex-1 min-w-0">
                  <div className="flex justify-between items-baseline">
                    <h3 className="text-sm font-medium truncate text-purple-700">{chat.name}</h3>
                    <span className="text-xs text-gray-500">{chat.time}</span>
                  </div>
                  <p className="text-sm text-gray-500 truncate">{chat.lastMessage}</p>
                </div>

                {chat.unread > 0 && (
                  <div className="flex-shrink-0">
                    <span className="inline-flex items-center justify-center h-5 w-5 rounded-full bg-gradient-to-r from-purple-600 to-indigo-600 text-xs font-medium text-white">
                      {chat.unread}
                    </span>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </motion.div>
      ))}
    </div>
  )
}
