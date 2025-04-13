"use client"

import { useState } from "react"
import { motion, AnimatePresence } from "framer-motion"
import { ProfileCardSwipeable } from "@/components/profile-card-swipeable"
import { useToast } from "@/hooks/use-toast"

// Mock data for developer profiles
const mockProfiles = [
  {
    id: 1,
    name: "Alex Johnson",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Full-stack developer with 5 years of experience in React and Node.js",
    interests: ["React", "TypeScript", "GraphQL", "Node.js"],
    github: "alexjohnson",
    linkedin: "alex-johnson",
  },
  {
    id: 2,
    name: "Sarah Chen",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Frontend developer specializing in UI/UX and accessibility",
    interests: ["Vue.js", "CSS", "Accessibility", "Design Systems"],
    github: "sarahchen",
    linkedin: "sarah-chen",
  },
  {
    id: 3,
    name: "Miguel Rodriguez",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Backend engineer with expertise in distributed systems",
    interests: ["Go", "Kubernetes", "Microservices", "System Design"],
    github: "miguelrodriguez",
    linkedin: "miguel-rodriguez",
  },
  {
    id: 4,
    name: "Priya Patel",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Machine learning engineer focused on NLP applications",
    interests: ["Python", "TensorFlow", "NLP", "Data Science"],
    github: "priyapatel",
    linkedin: "priya-patel",
  },
  {
    id: 5,
    name: "David Kim",
    image: "/placeholder.svg?height=300&width=300",
    summary: "DevOps engineer with a passion for automation",
    interests: ["AWS", "Terraform", "CI/CD", "Docker"],
    github: "davidkim",
    linkedin: "david-kim",
  },
]

export function ProfileCarouselSwipeable() {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [direction, setDirection] = useState(0)
  const [favorites, setFavorites] = useState<number[]>([])
  const { toast } = useToast()

  const currentProfile = mockProfiles[currentIndex]

  const handleSwipeLeft = () => {
    if (currentIndex < mockProfiles.length - 1) {
      setDirection(1)
      setTimeout(() => {
        setCurrentIndex(currentIndex + 1)
      }, 300)
    } else {
      toast({
        title: "No more profiles",
        description: "You've seen all available profiles",
      })
    }
  }

  const handleSwipeRight = () => {
    if (currentIndex < mockProfiles.length - 1) {
      setDirection(1)
      setTimeout(() => {
        setCurrentIndex(currentIndex + 1)
      }, 300)
    } else {
      toast({
        title: "No more profiles",
        description: "You've seen all available profiles",
      })
    }
  }

  const handleSwipeUp = () => {
    if (currentIndex < mockProfiles.length - 1) {
      setDirection(1)
      setTimeout(() => {
        setCurrentIndex(currentIndex + 1)
      }, 300)
    } else {
      toast({
        title: "No more profiles",
        description: "You've seen all available profiles",
      })
    }
  }

  const toggleFavorite = (id: number) => {
    if (favorites.includes(id)) {
      setFavorites(favorites.filter((favId) => favId !== id))
      toast({
        title: "Removed from Favorites",
        description: `${currentProfile.name} has been removed from your favorites`,
      })
    } else {
      setFavorites([...favorites, id])
      toast({
        title: "Added to Favorites",
        description: `${currentProfile.name} has been added to your favorites`,
      })
    }
  }

  const variants = {
    enter: (direction: number) => ({
      x: direction > 0 ? 1000 : -1000,
      opacity: 0,
      scale: 0.8,
    }),
    center: {
      x: 0,
      opacity: 1,
      scale: 1,
      transition: {
        duration: 0.4,
      },
    },
    exit: (direction: number) => ({
      x: direction < 0 ? 1000 : -1000,
      opacity: 0,
      scale: 0.8,
      transition: {
        duration: 0.4,
      },
    }),
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-[70vh]">
      <div className="relative w-full max-w-md h-[500px] flex items-center justify-center">
        <AnimatePresence initial={false} custom={direction} mode="wait">
          <motion.div
            key={currentIndex}
            custom={direction}
            variants={variants}
            initial="enter"
            animate="center"
            exit="exit"
            className="absolute w-full"
          >
            <ProfileCardSwipeable
              profile={currentProfile}
              isFavorite={favorites.includes(currentProfile.id)}
              onToggleFavorite={() => toggleFavorite(currentProfile.id)}
              onSwipeLeft={handleSwipeLeft}
              onSwipeRight={handleSwipeRight}
              onSwipeUp={handleSwipeUp}
            />
          </motion.div>
        </AnimatePresence>
      </div>

      <div className="flex justify-center mt-6">
        {mockProfiles.map((_, index) => (
          <div
            key={index}
            className={`h-2 w-2 rounded-full mx-1 ${index === currentIndex ? "bg-pink-600" : "bg-pink-200"}`}
          />
        ))}
      </div>

      <div className="text-center mt-8 max-w-md">
        <p className="text-sm text-gray-500 mb-2">Swipe gestures:</p>
        <div className="flex justify-center space-x-6 text-xs text-gray-500">
          <div className="flex flex-col items-center">
            <span className="mb-1">←</span>
            <span>Skip</span>
          </div>
          <div className="flex flex-col items-center">
            <span className="mb-1">→</span>
            <span>Like</span>
          </div>
          <div className="flex flex-col items-center">
            <span className="mb-1">↑</span>
            <span>Connect</span>
          </div>
        </div>
      </div>
    </div>
  )
}
