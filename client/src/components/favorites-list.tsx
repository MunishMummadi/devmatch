"use client"

import { useState } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Github, Linkedin, Heart, MessageSquare } from "lucide-react"
import { useRouter } from "next/navigation"
import { motion } from "framer-motion"

// Mock data for favorite profiles
const mockFavorites = [
  {
    id: 1,
    name: "Sarah Chen",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Frontend developer specializing in UI/UX and accessibility",
    interests: ["Vue.js", "CSS", "Accessibility"],
    github: "sarahchen",
    linkedin: "sarah-chen",
    connected: true,
  },
  {
    id: 3,
    name: "Miguel Rodriguez",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Backend engineer with expertise in distributed systems",
    interests: ["Go", "Kubernetes", "Microservices"],
    github: "miguelrodriguez",
    linkedin: "miguel-rodriguez",
    connected: true,
  },
  {
    id: 5,
    name: "David Kim",
    image: "/placeholder.svg?height=300&width=300",
    summary: "DevOps engineer with a passion for automation",
    interests: ["AWS", "Terraform", "CI/CD"],
    github: "davidkim",
    linkedin: "david-kim",
    connected: false,
  },
]

export function FavoritesList() {
  const router = useRouter()
  const [favorites, setFavorites] = useState(mockFavorites)

  const removeFavorite = (id: number) => {
    setFavorites(favorites.filter((fav) => fav.id !== id))
  }

  const startChat = (id: number) => {
    router.push(`/dashboard/chats/${id}`)
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {favorites.map((profile, index) => (
        <motion.div
          key={profile.id}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: index * 0.1 }}
        >
          <Card className="overflow-hidden shadow-md backdrop-blur-sm bg-white/90 border-none">
            <CardContent className="p-6">
              <div className="flex flex-col items-center">
                <div className="relative mb-4">
                  <div className="h-24 w-24 rounded-full overflow-hidden border-4 border-purple-200 shadow-lg">
                    <img
                      src={profile.image || "/placeholder.svg"}
                      alt={profile.name}
                      className="h-full w-full object-cover"
                    />
                  </div>
                  <Button
                    variant="outline"
                    size="icon"
                    className="absolute -right-2 bottom-0 rounded-full bg-pink-100 text-pink-500 border-pink-200 shadow-sm"
                    onClick={() => removeFavorite(profile.id)}
                  >
                    <Heart className="h-4 w-4 fill-pink-500" />
                  </Button>
                </div>

                <h2 className="text-lg font-semibold text-center mb-2 bg-clip-text text-transparent bg-gradient-to-r from-purple-700 to-indigo-700">
                  {profile.name}
                </h2>

                <p className="text-gray-600 text-center mb-4 text-sm line-clamp-2">{profile.summary}</p>

                <div className="flex flex-wrap justify-center gap-1 mb-4">
                  {profile.interests.slice(0, 2).map((interest, index) => (
                    <Badge
                      key={index}
                      variant="secondary"
                      className="bg-gradient-to-r from-purple-100 to-indigo-100 text-purple-800 text-xs"
                    >
                      {interest}
                    </Badge>
                  ))}
                  {profile.interests.length > 2 && (
                    <Badge
                      variant="secondary"
                      className="bg-gradient-to-r from-purple-100 to-indigo-100 text-purple-800 text-xs"
                    >
                      +{profile.interests.length - 2}
                    </Badge>
                  )}
                </div>

                <div className="flex space-x-4 mb-4">
                  <a
                    href={`https://github.com/${profile.github}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-gray-700 hover:text-purple-600 transition-colors"
                  >
                    <Github className="h-5 w-5" />
                  </a>
                  <a
                    href={`https://linkedin.com/in/${profile.linkedin}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-gray-700 hover:text-purple-600 transition-colors"
                  >
                    <Linkedin className="h-5 w-5" />
                  </a>
                </div>

                <div className="w-full space-y-2">
                  {profile.connected ? (
                    <Button
                      variant="default"
                      className="w-full bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-700 hover:to-indigo-700"
                      onClick={() => startChat(profile.id)}
                    >
                      <MessageSquare className="h-4 w-4 mr-2" />
                      Message
                    </Button>
                  ) : (
                    <Button variant="outline" className="w-full">
                      Send Connection Request
                    </Button>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        </motion.div>
      ))}

      {favorites.length === 0 && (
        <div className="col-span-full text-center py-12">
          <h3 className="text-lg font-medium text-purple-700 mb-2">No favorites yet</h3>
          <p className="text-gray-500 mb-4">
            When you find developers you'd like to connect with, add them to your favorites.
          </p>
          <Button
            onClick={() => router.push("/dashboard")}
            className="bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-700 hover:to-indigo-700"
          >
            Discover Developers
          </Button>
        </div>
      )}
    </div>
  )
}
