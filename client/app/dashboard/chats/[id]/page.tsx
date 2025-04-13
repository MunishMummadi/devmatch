import { ChatInterface } from "@/components/chat-interface"
import { DashboardHeader } from "@/components/dashboard-header"
import { AnimatedBackground } from "@/components/animated-background"

export default function ChatPage({ params }: { params: { id: string } }) {
  return (
    <main className="min-h-screen flex flex-col relative">
      <AnimatedBackground />
      <DashboardHeader />
      <div className="flex-1 container mx-auto px-4 py-4 flex flex-col relative z-10">
        <ChatInterface chatId={params.id} />
      </div>
    </main>
  )
}
