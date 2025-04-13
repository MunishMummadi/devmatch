import { FavoritesList } from "@/components/favorites-list"
import { DashboardHeader } from "@/components/dashboard-header"
import { AnimatedBackground } from "@/components/animated-background"

export default function FavoritesPage() {
  return (
    <main className="min-h-screen relative">
      <AnimatedBackground />
      <DashboardHeader />
      <div className="container mx-auto px-4 py-8 relative z-10">
        <h1 className="text-2xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-pink-600 to-purple-600">
          Your Favorites
        </h1>
        <FavoritesList />
      </div>
    </main>
  )
}
