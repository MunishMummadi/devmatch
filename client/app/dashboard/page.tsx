import { ProfileCarousel } from "@/components/profile-carousel"
import { DashboardHeader } from "@/components/dashboard-header"
import { FloatingShapes } from "@/components/floating-shapes"

export default function DashboardPage() {
  return (
    <main className="min-h-screen relative">
      <FloatingShapes />
      <DashboardHeader />
      <div className="container mx-auto px-4 py-8 relative z-10">
        <ProfileCarousel />
      </div>
    </main>
  )
}
