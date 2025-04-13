"use client"

import { useState, useEffect } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Github, Linkedin } from "lucide-react"

interface UserProfile {
  firstName: string
  lastName: string
  nickname?: string
  age: string
  gender: string
  summary: string
  github?: string
  linkedin?: string
  interests: string[]
  image: string
}

export function UserProfileCard() {
  const [userProfile, setUserProfile] = useState<UserProfile | null>(null)

  useEffect(() => {
    // Load user profile from localStorage
    const savedUserData = localStorage.getItem("userProfile")
    if (savedUserData) {
      setUserProfile(JSON.parse(savedUserData))
    }
  }, [])

  if (!userProfile) {
    return (
      <Card className="w-full overflow-hidden shadow-xl backdrop-blur-sm bg-white/90 border-none">
        <CardContent className="p-6 text-center">
          <p>Please complete your profile details first.</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="w-full overflow-hidden shadow-xl backdrop-blur-sm bg-white/90 border-none">
      <CardContent className="p-0">
        <div className="flex flex-col items-center p-6 pb-4">
          <div className="relative mb-6">
            <div className="h-36 w-36 rounded-full overflow-hidden border-4 border-purple-200 shadow-lg">
              <img
                src={userProfile.image || "/placeholder.svg?height=300&width=300"}
                alt={`${userProfile.firstName} ${userProfile.lastName}`}
                className="h-full w-full object-cover"
              />
            </div>
          </div>

          <h2 className="text-xl font-semibold text-center mb-3 bg-clip-text text-transparent bg-gradient-to-r from-purple-700 to-indigo-700">
            {userProfile.firstName} {userProfile.lastName}
            {userProfile.nickname && <span className="text-gray-500 text-sm ml-1">({userProfile.nickname})</span>}
          </h2>

          <div className="flex flex-wrap gap-2 mb-2 text-xs text-gray-600">
            <span>Age: {userProfile.age}</span>
            <span>•</span>
            <span>Gender: {userProfile.gender}</span>
          </div>

          <p className="text-gray-600 text-center mb-4 line-clamp-3">{userProfile.summary}</p>

          <div className="flex flex-wrap justify-center gap-2 mb-4">
            {userProfile.interests.map((interest, index) => (
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
            {userProfile.github && (
              <a
                href={`https://github.com/${userProfile.github}`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-700 hover:text-purple-600 transition-colors"
              >
                <Github className="h-5 w-5" />
              </a>
            )}
            {userProfile.linkedin && (
              <a
                href={`https://linkedin.com/in/${userProfile.linkedin}`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-700 hover:text-purple-600 transition-colors"
              >
                <Linkedin className="h-5 w-5" />
              </a>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
