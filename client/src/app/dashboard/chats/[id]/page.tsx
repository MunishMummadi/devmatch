import { ChatInterface } from "@/components/chat-interface"
import { DashboardHeader } from "@/components/dashboard-header"
import { FloatingShapes } from "@/components/floating-shapes"

export default function ChatPage({ params }: { params: { id: string } }) {
  return (
    <main className="min-h-screen flex flex-col relative">
      <FloatingShapes />
      <DashboardHeader />
      <div className="flex-1 container mx-auto px-4 py-4 flex flex-col relative z-10">
        <ChatInterface chatId={params.id} />
      </div>
    </main>
  )
}
