"use client"

import { useState, useRef, useEffect } from "react"
import { motion, type PanInfo, useAnimation } from "framer-motion"
import { ProfileCard } from "@/components/profile-card"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight } from "lucide-react"

// Mock data for developer profiles
const mockProfiles = [
  {
    id: 1,
    name: "Alex Johnson",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Full-stack developer with 5 years of experience in React and Node.js",
    interests: ["React", "TypeScript", "GraphQL"],
    github: "alexjohnson",
    linkedin: "alex-johnson",
  },
  {
    id: 2,
    name: "Sarah Chen",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Frontend developer specializing in UI/UX and accessibility",
    interests: ["Vue.js", "CSS", "Accessibility"],
    github: "sarahchen",
    linkedin: "sarah-chen",
  },
  {
    id: 3,
    name: "Miguel Rodriguez",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Backend engineer with expertise in distributed systems",
    interests: ["Go", "Kubernetes", "Microservices"],
    github: "miguelrodriguez",
    linkedin: "miguel-rodriguez",
  },
  {
    id: 4,
    name: "Priya Patel",
    image: "/placeholder.svg?height=300&width=300",
    summary: "Machine learning engineer focused on NLP applications",
    interests: ["Python", "TensorFlow", "NLP"],
    github: "priyapatel",
    linkedin: "priya-patel",
  },
  {
    id: 5,
    name: "David Kim",
    image: "/placeholder.svg?height=300&width=300",
    summary: "DevOps engineer with a passion for automation",
    interests: ["AWS", "Terraform", "CI/CD"],
    github: "davidkim",
    linkedin: "david-kim",
  },
]

export function ProfileDeck() {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [direction, setDirection] = useState<string | null>(null)
  const controls = useAnimation()
  const constraintsRef = useRef(null)

  const currentProfile = mockProfiles[currentIndex]

  const handleDragEnd = (event: MouseEvent | TouchEvent | PointerEvent, info: PanInfo) => {
    const threshold = 100

    if (info.offset.x > threshold) {
      // Swiped right - previous
      handlePrevious()
    } else if (info.offset.x < -threshold) {
      // Swiped left - next
      handleNext()
    } else if (info.offset.y < -threshold) {
      // Swiped up - send request
      handleSendRequest()
    } else {
      // Reset position
      controls.start({ x: 0, y: 0, transition: { type: "spring", stiffness: 300, damping: 20 } })
    }
  }

  const handleNext = () => {
    if (currentIndex < mockProfiles.length - 1) {
      setDirection("right")
      controls
        .start({
          x: -300,
          opacity: 0,
          transition: { duration: 0.3 },
        })
        .then(() => {
          setCurrentIndex(currentIndex + 1)
          controls.start({
            x: 0,
            opacity: 1,
            transition: { duration: 0.3 },
          })
        })
    }
  }

  const handlePrevious = () => {
    if (currentIndex > 0) {
      setDirection("left")
      controls
        .start({
          x: 300,
          opacity: 0,
          transition: { duration: 0.3 },
        })
        .then(() => {
          setCurrentIndex(currentIndex - 1)
          controls.start({
            x: 0,
            opacity: 1,
            transition: { duration: 0.3 },
          })
        })
    }
  }

  const handleSendRequest = () => {
    controls
      .start({
        y: -300,
        opacity: 0,
        transition: { duration: 0.3 },
      })
      .then(() => {
        // In a real app, this would send a connection request
        alert(`Connection request sent to ${currentProfile.name}!`)

        if (currentIndex < mockProfiles.length - 1) {
          setCurrentIndex(currentIndex + 1)
        } else {
          setCurrentIndex(0) // Loop back to the first profile
        }

        controls.start({
          y: 0,
          opacity: 1,
          transition: { duration: 0.3 },
        })
      })
  }

  useEffect(() => {
    // Reset direction after animation
    if (direction) {
      const timer = setTimeout(() => {
        setDirection(null)
      }, 300)
      return () => clearTimeout(timer)
    }
  }, [direction])

  return (
    <div className="flex flex-col items-center justify-center min-h-[70vh]" ref={constraintsRef}>
      <div className="relative w-full max-w-md h-[500px] flex items-center justify-center">
        <motion.div
          drag
          dragConstraints={constraintsRef}
          onDragEnd={handleDragEnd}
          animate={controls}
          initial={{ opacity: 1, x: 0, y: 0 }}
          className="w-full"
        >
          <ProfileCard profile={currentProfile} />
        </motion.div>
      </div>

      <div className="flex items-center justify-center mt-8 space-x-8">
        <Button
          variant="outline"
          size="icon"
          className="rounded-full h-12 w-12 bg-white"
          onClick={handlePrevious}
          disabled={currentIndex === 0}
        >
          <ChevronLeft className="h-6 w-6" />
        </Button>

        <Button
          variant="default"
          size="lg"
          className="rounded-full h-14 w-14 bg-purple-600 hover:bg-purple-700"
          onClick={handleSendRequest}
        >
          <span className="sr-only">Send Connection Request</span>
          <span className="text-xl">↑</span>
        </Button>

        <Button
          variant="outline"
          size="icon"
          className="rounded-full h-12 w-12 bg-white"
          onClick={handleNext}
          disabled={currentIndex === mockProfiles.length - 1}
        >
          <ChevronRight className="h-6 w-6" />
        </Button>
      </div>

      <p className="text-center text-sm text-gray-500 mt-4">
        Swipe up to send a connection request, or left/right to browse profiles
      </p>
    </div>
  )
}
