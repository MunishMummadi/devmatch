import { ProfileCarouselSwipeable } from "@/components/profile-carousel-swipeable"
import { DashboardHeader } from "@/components/dashboard-header"
import { AnimatedBackground } from "@/components/animated-background"

export default function DashboardPage() {
  return (
    <main className="min-h-screen relative">
      <AnimatedBackground />
      <DashboardHeader />
      <div className="container mx-auto px-4 py-8 relative z-10">
        <ProfileCarouselSwipeable />
      </div>
    </main>
  )
}
