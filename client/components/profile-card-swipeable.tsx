"use client"

import type React from "react"

import { useState, useRef, useEffect } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Github, Linkedin, Heart, X } from "lucide-react"
import { motion, type PanInfo, useAnimation } from "framer-motion"
import { useToast } from "@/hooks/use-toast"

interface Profile {
  id: number
  name: string
  image: string
  summary: string
  interests: string[]
  github: string
  linkedin: string
}

interface ProfileCardProps {
  profile: Profile
  isFavorite: boolean
  onToggleFavorite: () => void
  onSwipeLeft: () => void
  onSwipeRight: () => void
  onSwipeUp: () => void
}

export function ProfileCardSwipeable({
  profile,
  isFavorite,
  onToggleFavorite,
  onSwipeLeft,
  onSwipeRight,
  onSwipeUp,
}: ProfileCardProps) {
  const [showDetails, setShowDetails] = useState(false)
  const detailsRef = useRef<HTMLDivElement>(null)
  const controls = useAnimation()
  const { toast } = useToast()
  const [isDragging, setIsDragging] = useState(false)

  const toggleDetails = (e: React.MouseEvent) => {
    e.stopPropagation()
    setShowDetails(!showDetails)
  }

  const handleFavorite = (e: React.MouseEvent) => {
    e.stopPropagation()
    onToggleFavorite()
  }

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (showDetails && detailsRef.current && !detailsRef.current.contains(event.target as Node)) {
        setShowDetails(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    return () => {
      document.removeEventListener("mousedown", handleClickOutside)
    }
  }, [showDetails])

  const handleDragEnd = (event: MouseEvent | TouchEvent | PointerEvent, info: PanInfo) => {
    const threshold = 100
    setIsDragging(false)

    if (info.offset.x > threshold) {
      // Swiped right
      controls.start({
        x: 500,
        opacity: 0,
        transition: { duration: 0.3 },
      })
      onSwipeRight()
      toast({
        title: "Liked Profile",
        description: `You liked ${profile.name}'s profile`,
      })
    } else if (info.offset.x < -threshold) {
      // Swiped left
      controls.start({
        x: -500,
        opacity: 0,
        transition: { duration: 0.3 },
      })
      onSwipeLeft()
      toast({
        title: "Skipped Profile",
        description: `You skipped ${profile.name}'s profile`,
      })
    } else if (info.offset.y < -threshold) {
      // Swiped up
      controls.start({
        y: -500,
        opacity: 0,
        transition: { duration: 0.3 },
      })
      onSwipeUp()
      toast({
        title: "Connection Request Sent",
        description: `You've sent a connection request to ${profile.name}`,
      })
    } else {
      // Reset position
      controls.start({
        x: 0,
        y: 0,
        opacity: 1,
        transition: { type: "spring", stiffness: 300, damping: 20 },
      })
    }
  }

  return (
    <motion.div
      drag
      dragConstraints={{ left: 0, right: 0, top: 0, bottom: 0 }}
      dragElastic={0.9}
      onDragStart={() => setIsDragging(true)}
      onDragEnd={handleDragEnd}
      animate={controls}
      initial={{ opacity: 1, x: 0, y: 0 }}
      className="w-full touch-none"
      whileDrag={{ scale: 1.05 }}
    >
      <Card className="w-full overflow-hidden shadow-xl backdrop-blur-sm bg-white/90 border-none">
        <CardContent className="p-0">
          <div className="flex flex-col items-center p-6 pb-4">
            <div className="relative mb-6">
              <div className="h-36 w-36 rounded-full overflow-hidden border-4 border-pink-200 shadow-lg">
                <img
                  src={profile.image || "/placeholder.svg"}
                  alt={profile.name}
                  className="h-full w-full object-cover"
                />
              </div>
              <Button
                variant="outline"
                size="icon"
                className={`absolute -right-2 bottom-0 rounded-full shadow-md ${
                  isFavorite ? "bg-pink-100 text-pink-500 border-pink-200" : "bg-white"
                }`}
                onClick={handleFavorite}
              >
                <Heart className={`h-4 w-4 ${isFavorite ? "fill-pink-500" : ""}`} />
              </Button>
            </div>

            <h2 className="text-xl font-semibold text-center mb-3 bg-clip-text text-transparent bg-gradient-to-r from-pink-600 to-purple-600">
              {profile.name}
            </h2>

            <p className="text-gray-600 text-center mb-4 line-clamp-3">{profile.summary}</p>

            <div className="flex flex-wrap justify-center gap-2 mb-4">
              {profile.interests.map((interest, index) => (
                <Badge
                  key={index}
                  variant="secondary"
                  className="bg-gradient-to-r from-pink-100 to-purple-100 text-purple-800 hover:from-pink-200 hover:to-purple-200"
                >
                  {interest}
                </Badge>
              ))}
            </div>

            <div className="flex space-x-4 mb-4">
              <a
                href={`https://github.com/${profile.github}`}
                target="_blank"
                rel="noopener noreferrer"
                onClick={(e) => e.stopPropagation()}
                className="text-gray-700 hover:text-pink-600 transition-colors"
              >
                <Github className="h-5 w-5" />
              </a>
              <a
                href={`https://linkedin.com/in/${profile.linkedin}`}
                target="_blank"
                rel="noopener noreferrer"
                onClick={(e) => e.stopPropagation()}
                className="text-gray-700 hover:text-purple-600 transition-colors"
              >
                <Linkedin className="h-5 w-5" />
              </a>
            </div>

            {!isDragging && (
              <Button variant="ghost" size="sm" className="text-purple-600 mt-2" onClick={toggleDetails}>
                {showDetails ? "Hide Details" : "View Details"}
              </Button>
            )}
          </div>

          {showDetails && (
            <motion.div
              ref={detailsRef}
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              transition={{ duration: 0.3 }}
              className="relative"
            >
              <div className="p-6 pt-0 bg-gradient-to-b from-white/50 to-pink-50/50">
                <Button
                  variant="ghost"
                  size="icon"
                  className="absolute top-0 right-2 text-gray-400 hover:text-gray-600"
                  onClick={toggleDetails}
                >
                  <X className="h-4 w-4" />
                </Button>

                <h3 className="font-medium mb-2 text-pink-700">About {profile.name}</h3>
                <p className="text-gray-600 mb-4">
                  {profile.summary} Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, nisl eget
                  ultricies tincidunt, nisl nisl aliquam nisl, eget ultricies nisl nisl eget nisl.
                </p>

                <h3 className="font-medium mb-2 text-pink-700">Experience</h3>
                <ul className="list-disc list-inside text-gray-600 mb-4">
                  <li>Senior Developer at TechCorp (2020-Present)</li>
                  <li>Frontend Developer at WebSolutions (2018-2020)</li>
                  <li>Junior Developer at StartupXYZ (2016-2018)</li>
                </ul>

                <Button className="w-full bg-gradient-to-r from-pink-500 to-purple-600 hover:from-pink-600 hover:to-purple-700">
                  Send Connection Request
                </Button>
              </div>
            </motion.div>
          )}
        </CardContent>
      </Card>

      {/* Swipe indicators */}
      {isDragging && (
        <div className="absolute inset-0 pointer-events-none flex items-center justify-center">
          <motion.div
            className="absolute top-1/2 left-8 transform -translate-y-1/2 bg-red-500 text-white rounded-full p-3"
            animate={{ opacity: controls.x?.get() < -50 ? 1 : 0 }}
          >
            <X className="h-8 w-8" />
          </motion.div>
          <motion.div
            className="absolute top-1/2 right-8 transform -translate-y-1/2 bg-green-500 text-white rounded-full p-3"
            animate={{ opacity: controls.x?.get() > 50 ? 1 : 0 }}
          >
            <Heart className="h-8 w-8" />
          </motion.div>
          <motion.div
            className="absolute top-8 left-1/2 transform -translate-x-1/2 bg-purple-500 text-white rounded-full p-3"
            animate={{ opacity: controls.y?.get() < -50 ? 1 : 0 }}
          >
            <Github className="h-8 w-8" />
          </motion.div>
        </div>
      )}
    </motion.div>
  )
}
