"use client"

import { useState, useRef } from "react"
import { motion, AnimatePresence } from "framer-motion"
import { ProfileCard } from "@/components/profile-card"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight, Heart } from "lucide-react"
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

export function ProfileCarousel() {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [direction, setDirection] = useState(0)
  const [favorites, setFavorites] = useState<number[]>([])
  const constraintsRef = useRef(null)
  const { toast } = useToast()

  const currentProfile = mockProfiles[currentIndex]

  const handleNext = () => {
    if (currentIndex < mockProfiles.length - 1) {
      setDirection(1)
      setCurrentIndex(currentIndex + 1)
    }
  }

  const handlePrevious = () => {
    if (currentIndex > 0) {
      setDirection(-1)
      setCurrentIndex(currentIndex - 1)
    }
  }

  const handleSendRequest = () => {
    toast({
      title: "Connection Request Sent",
      description: `You've sent a connection request to ${currentProfile.name}`,
    })
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
    <div className="flex flex-col items-center justify-center min-h-[70vh]" ref={constraintsRef}>
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
            <ProfileCard
              profile={currentProfile}
              isFavorite={favorites.includes(currentProfile.id)}
              onToggleFavorite={() => toggleFavorite(currentProfile.id)}
            />
          </motion.div>
        </AnimatePresence>
      </div>

      <div className="flex items-center justify-center mt-8 space-x-8">
        <Button
          variant="outline"
          size="icon"
          className="rounded-full h-12 w-12 backdrop-blur-sm bg-white/80 shadow-md"
          onClick={handlePrevious}
          disabled={currentIndex === 0}
        >
          <ChevronLeft className="h-6 w-6" />
        </Button>

        <Button
          variant="default"
          size="lg"
          className="rounded-full h-14 w-14 bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-700 hover:to-indigo-700 shadow-md"
          onClick={handleSendRequest}
        >
          <span className="sr-only">Send Connection Request</span>
          <Heart className="h-6 w-6" />
        </Button>

        <Button
          variant="outline"
          size="icon"
          className="rounded-full h-12 w-12 backdrop-blur-sm bg-white/80 shadow-md"
          onClick={handleNext}
          disabled={currentIndex === mockProfiles.length - 1}
        >
          <ChevronRight className="h-6 w-6" />
        </Button>
      </div>

      <div className="flex justify-center mt-6">
        {mockProfiles.map((_, index) => (
          <button
            key={index}
            className={`h-2 w-2 rounded-full mx-1 ${index === currentIndex ? "bg-purple-600" : "bg-purple-200"}`}
            onClick={() => {
              setDirection(index > currentIndex ? 1 : -1)
              setCurrentIndex(index)
            }}
          />
        ))}
      </div>

      <p className="text-center text-sm text-gray-500 mt-4">
        Use the arrows to browse profiles or the heart button to send a connection request
      </p>
    </div>
  )
}
