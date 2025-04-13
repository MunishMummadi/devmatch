"use client"

import type React from "react"

import { useState, useEffect, useRef } from "react"
import { useRouter } from "next/navigation"
import { AnimatedBackground } from "@/components/animated-background"
import { DashboardHeader } from "@/components/dashboard-header"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Github, Linkedin, Edit, ArrowRight, Camera } from "lucide-react"
import { motion } from "framer-motion"
import { useToast } from "@/hooks/use-toast"

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

export default function ProfilePage() {
  const router = useRouter()
  const { toast } = useToast()
  const [userProfile, setUserProfile] = useState<UserProfile | null>(null)
  const [loading, setLoading] = useState(true)
  const fileInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    // Load user profile from localStorage
    const savedUserData = localStorage.getItem("userProfile")
    if (savedUserData) {
      setUserProfile(JSON.parse(savedUserData))
    }
    setLoading(false)
  }, [])

  const handleEditProfile = () => {
    router.push("/details")
  }

  const handleContinue = () => {
    router.push("/dashboard")
  }

  const handleImageClick = () => {
    fileInputRef.current?.click()
  }

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (event) => {
        if (event.target?.result && userProfile) {
          const updatedProfile = {
            ...userProfile,
            image: event.target.result as string,
          }
          setUserProfile(updatedProfile)
          localStorage.setItem("userProfile", JSON.stringify(updatedProfile))

          toast({
            title: "Profile Picture Updated",
            description: "Your profile picture has been successfully updated.",
          })
        }
      }
      reader.readAsDataURL(file)
    }
  }

  if (loading) {
    return (
      <main className="min-h-screen relative">
        <AnimatedBackground />
        <DashboardHeader />
        <div className="container mx-auto px-4 py-8 flex items-center justify-center">
          <div className="text-purple-600 animate-pulse">Loading profile...</div>
        </div>
      </main>
    )
  }

  if (!userProfile) {
    return (
      <main className="min-h-screen relative">
        <AnimatedBackground />
        <DashboardHeader />
        <div className="container mx-auto px-4 py-8 flex flex-col items-center justify-center">
          <h1 className="text-2xl font-bold mb-4">Profile Not Found</h1>
          <p className="mb-6">Please complete your profile details first.</p>
          <Button onClick={() => router.push("/details")}>Create Profile</Button>
        </div>
      </main>
    )
  }

  return (
    <main className="min-h-screen relative">
      <AnimatedBackground />
      <DashboardHeader />
      <div className="container mx-auto px-4 py-8 relative z-10">
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5 }}>
          <h1 className="text-2xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-pink-600 to-purple-600">
            Your Profile
          </h1>

          <Card className="max-w-3xl mx-auto backdrop-blur-sm bg-white/90 border-none shadow-xl">
            <CardContent className="p-8">
              <div className="flex flex-col md:flex-row gap-8">
                <div className="flex flex-col items-center">
                  <div
                    className="relative h-48 w-48 rounded-full overflow-hidden border-4 border-purple-200 shadow-lg mb-4 cursor-pointer group"
                    onClick={handleImageClick}
                  >
                    <img
                      src={userProfile.image || "/placeholder.svg?height=300&width=300"}
                      alt={`${userProfile.firstName} ${userProfile.lastName}`}
                      className="h-full w-full object-cover"
                    />
                    <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                      <Camera className="h-8 w-8 text-white" />
                    </div>
                  </div>
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={handleImageUpload}
                  />
                  <Button variant="outline" size="sm" className="mb-4 text-purple-600" onClick={handleEditProfile}>
                    <Edit className="h-4 w-4 mr-2" /> Edit Profile
                  </Button>

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

                <div className="flex-1">
                  <h2 className="text-2xl font-semibold mb-2 bg-clip-text text-transparent bg-gradient-to-r from-purple-700 to-indigo-700">
                    {userProfile.firstName} {userProfile.lastName}
                    {userProfile.nickname && (
                      <span className="text-gray-500 text-lg ml-2">({userProfile.nickname})</span>
                    )}
                  </h2>

                  <div className="flex flex-wrap gap-4 mb-4 text-sm text-gray-600">
                    <div>Age: {userProfile.age}</div>
                    <div>Gender: {userProfile.gender}</div>
                  </div>

                  <div className="mb-6">
                    <h3 className="text-lg font-medium mb-2 text-purple-700">About Me</h3>
                    <p className="text-gray-700">{userProfile.summary}</p>
                  </div>

                  <div className="mb-6">
                    <h3 className="text-lg font-medium mb-2 text-purple-700">Interests</h3>
                    <div className="flex flex-wrap gap-2">
                      {userProfile.interests.map((interest, index) => (
                        <Badge
                          key={index}
                          variant="secondary"
                          className="bg-gradient-to-r from-purple-100 to-indigo-100 text-purple-800"
                        >
                          {interest}
                        </Badge>
                      ))}
                      {userProfile.interests.length === 0 && (
                        <p className="text-sm text-gray-500">No interests added yet</p>
                      )}
                    </div>
                  </div>

                  <Button
                    onClick={handleContinue}
                    className="w-full bg-gradient-to-r from-pink-500 to-purple-600 hover:from-pink-600 hover:to-purple-700"
                  >
                    Continue to Discover <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </motion.div>
      </div>
    </main>
  )
}
