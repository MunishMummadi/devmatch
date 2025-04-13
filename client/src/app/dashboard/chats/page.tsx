import { ChatList } from "@/components/chat-list"
import { DashboardHeader } from "@/components/dashboard-header"
import { FloatingShapes } from "@/components/floating-shapes"

export default function ChatsPage() {
  return (
    <main className="min-h-screen relative">
      <FloatingShapes />
      <DashboardHeader />
      <div className="container mx-auto px-4 py-8 relative z-10">
        <h1 className="text-2xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-purple-700 to-indigo-700">
          Your Chats
        </h1>
        <ChatList />
      </div>
    </main>
  )
}
