"use client"

import type React from "react"

import { useState, useRef, useEffect } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Github, Linkedin, Heart, X } from "lucide-react"
import { motion, AnimatePresence } from "framer-motion"

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
}

export function ProfileCard({ profile, isFavorite, onToggleFavorite }: ProfileCardProps) {
  const [showDetails, setShowDetails] = useState(false)
  const detailsRef = useRef<HTMLDivElement>(null)

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

  return (
    <Card className="w-full overflow-hidden shadow-xl backdrop-blur-sm bg-white/90 border-none">
      <CardContent className="p-0">
        <div className="flex flex-col items-center p-6 pb-4">
          <div className="relative mb-6">
            <div className="h-36 w-36 rounded-full overflow-hidden border-4 border-purple-200 shadow-lg">
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

          <h2 className="text-xl font-semibold text-center mb-3 bg-clip-text text-transparent bg-gradient-to-r from-purple-700 to-indigo-700">
            {profile.name}
          </h2>

          <p className="text-gray-600 text-center mb-4 line-clamp-3">{profile.summary}</p>

          <div className="flex flex-wrap justify-center gap-2 mb-4">
            {profile.interests.map((interest, index) => (
              <Badge
                key={index}
                variant="secondary"
                className="bg-gradient-to-r from-purple-100 to-indigo-100 text-purple-800 hover:from-purple-200 hover:to-indigo-200"
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
              className="text-gray-700 hover:text-purple-600 transition-colors"
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

          <Button variant="ghost" size="sm" className="text-purple-600 mt-2" onClick={toggleDetails}>
            {showDetails ? "Hide Details" : "View Details"}
          </Button>
        </div>

        <AnimatePresence>
          {showDetails && (
            <motion.div
              ref={detailsRef}
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              transition={{ duration: 0.3 }}
              className="relative"
            >
              <div className="p-6 pt-0 bg-gradient-to-b from-white/50 to-purple-50/50">
                <Button
                  variant="ghost"
                  size="icon"
                  className="absolute top-0 right-2 text-gray-400 hover:text-gray-600"
                  onClick={toggleDetails}
                >
                  <X className="h-4 w-4" />
                </Button>

                <h3 className="font-medium mb-2 text-purple-700">About {profile.name}</h3>
                <p className="text-gray-600 mb-4">
                  {profile.summary} Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, nisl eget
                  ultricies tincidunt, nisl nisl aliquam nisl, eget ultricies nisl nisl eget nisl.
                </p>

                <h3 className="font-medium mb-2 text-purple-700">Experience</h3>
                <ul className="list-disc list-inside text-gray-600 mb-4">
                  <li>Senior Developer at TechCorp (2020-Present)</li>
                  <li>Frontend Developer at WebSolutions (2018-2020)</li>
                  <li>Junior Developer at StartupXYZ (2016-2018)</li>
                </ul>

                <Button className="w-full bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-700 hover:to-indigo-700">
                  Send Connection Request
                </Button>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </CardContent>
    </Card>
  )
}
