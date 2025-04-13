"use client"

import type React from "react"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { AnimatedBackground } from "@/components/animated-background"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { X, Github, Linkedin } from "lucide-react"
import { motion } from "framer-motion"

export default function DetailsPage() {
  const router = useRouter()
  const [firstName, setFirstName] = useState("")
  const [lastName, setLastName] = useState("")
  const [nickname, setNickname] = useState("")
  const [age, setAge] = useState("")
  const [gender, setGender] = useState("")
  const [summary, setSummary] = useState("")
  const [github, setGithub] = useState("")
  const [linkedin, setLinkedin] = useState("")
  const [interest, setInterest] = useState("")
  const [interests, setInterests] = useState<string[]>([])
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [profileImage, setProfileImage] = useState<string>("/placeholder.svg?height=300&width=300")

  // Load existing data if available
  useEffect(() => {
    const savedUserData = localStorage.getItem("userProfile")
    if (savedUserData) {
      const userData = JSON.parse(savedUserData)
      setFirstName(userData.firstName || "")
      setLastName(userData.lastName || "")
      setNickname(userData.nickname || "")
      setAge(userData.age || "")
      setGender(userData.gender || "")
      setSummary(userData.summary || "")
      setGithub(userData.github || "")
      setLinkedin(userData.linkedin || "")
      setInterests(userData.interests || [])
      setProfileImage(userData.image || "/placeholder.svg?height=300&width=300")
    }
  }, [])

  const handleAddInterest = () => {
    if (interest.trim() && !interests.includes(interest.trim())) {
      setInterests([...interests, interest.trim()])
      setInterest("")
    }
  }

  const handleRemoveInterest = (interestToRemove: string) => {
    setInterests(interests.filter((i) => i !== interestToRemove))
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      e.preventDefault()
      handleAddInterest()
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)

    // Save user data to localStorage
    const userData = {
      firstName,
      lastName,
      nickname,
      age,
      gender,
      summary,
      github,
      linkedin,
      interests,
      image: profileImage,
    }

    localStorage.setItem("userProfile", JSON.stringify(userData))

    // Navigate to the dashboard after a short delay to show loading state
    setTimeout(() => {
      router.push("/profile")
    }, 500)
  }

  // Format GitHub username (remove https://github.com/ if present)
  const formatGithubUsername = (input: string) => {
    let username = input
    if (username.includes("github.com/")) {
      username = username.split("github.com/")[1]
    }
    return username
  }

  // Format LinkedIn username (remove https://linkedin.com/in/ if present)
  const formatLinkedinUsername = (input: string) => {
    let username = input
    if (username.includes("linkedin.com/in/")) {
      username = username.split("linkedin.com/in/")[1]
    }
    return username
  }

  // Handle profile image upload
  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (event) => {
        if (event.target?.result) {
          setProfileImage(event.target.result as string)
        }
      }
      reader.readAsDataURL(file)
    }
  }

  return (
    <main className="min-h-screen py-12 relative">
      <AnimatedBackground />

      <div className="container mx-auto px-4 relative z-10">
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5 }}>
          <Card className="max-w-2xl mx-auto backdrop-blur-sm bg-white/90 border-none shadow-xl">
            <CardHeader>
              <CardTitle className="text-2xl text-center bg-clip-text text-transparent bg-gradient-to-r from-pink-600 to-purple-600">
                Tell us about yourself
              </CardTitle>
              <CardDescription className="text-center">
                Let's set up your profile to find the perfect developer match
              </CardDescription>
            </CardHeader>

            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-6">
                {/* Profile Image Upload */}
                <div className="flex flex-col items-center mb-4">
                  <div className="relative mb-4">
                    <div className="h-32 w-32 rounded-full overflow-hidden border-4 border-purple-200 shadow-lg">
                      <img
                        src={profileImage || "/placeholder.svg"}
                        alt="Profile"
                        className="h-full w-full object-cover"
                      />
                    </div>
                    <label
                      htmlFor="profile-image"
                      className="absolute bottom-0 right-0 bg-purple-600 text-white rounded-full p-2 cursor-pointer shadow-md hover:bg-purple-700 transition-colors"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      >
                        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                        <polyline points="17 8 12 3 7 8" />
                        <line x1="12" y1="3" x2="12" y2="15" />
                      </svg>
                    </label>
                    <input
                      id="profile-image"
                      type="file"
                      accept="image/*"
                      className="hidden"
                      onChange={handleImageUpload}
                    />
                  </div>
                  <p className="text-sm text-gray-500">Click the icon to upload a profile picture</p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="firstName">First Name</Label>
                    <Input
                      id="firstName"
                      value={firstName}
                      onChange={(e) => setFirstName(e.target.value)}
                      className="bg-black/50"
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="lastName">Last Name</Label>
                    <Input
                      id="lastName"
                      value={lastName}
                      onChange={(e) => setLastName(e.target.value)}
                      className="bg-black/50"
                      required
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="nickname">Nickname</Label>
                    <Input
                      id="nickname"
                      value={nickname}
                      onChange={(e) => setNickname(e.target.value)}
                      className="bg-white/50"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="age">Age</Label>
                    <Input
                      id="age"
                      type="number"
                      value={age}
                      onChange={(e) => setAge(e.target.value)}
                      className="bg-white/50"
                      required
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="gender">Gender</Label>
                  <Select value={gender} onValueChange={setGender} required>
                    <SelectTrigger className="bg-white/50">
                      <SelectValue placeholder="Select your gender" />
                    </SelectTrigger>
                    <SelectContent className="bg-gray-800 text-white border-gray-700">
                      <SelectItem value="male" className="hover:bg-gray-700 focus:bg-gray-700">
                        Male
                      </SelectItem>
                      <SelectItem value="female" className="hover:bg-gray-700 focus:bg-gray-700">
                        Female
                      </SelectItem>
                      <SelectItem value="non-binary" className="hover:bg-gray-700 focus:bg-gray-700">
                        Non-binary
                      </SelectItem>
                      <SelectItem value="other" className="hover:bg-gray-700 focus:bg-gray-700">
                        Other
                      </SelectItem>
                      <SelectItem value="prefer-not-to-say" className="hover:bg-gray-700 focus:bg-gray-700">
                        Prefer not to say
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="summary">Brief Summary</Label>
                  <Textarea
                    id="summary"
                    value={summary}
                    onChange={(e) => setSummary(e.target.value)}
                    placeholder="Tell us a bit about yourself, your experience, and what you're looking for..."
                    className="min-h-[100px] bg-white/50"
                    required
                  />
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="github" className="flex items-center gap-2">
                      <Github className="h-4 w-4" /> GitHub Username
                    </Label>
                    <Input
                      id="github"
                      value={github}
                      onChange={(e) => setGithub(formatGithubUsername(e.target.value))}
                      placeholder="yourusername"
                      className="bg-white/50"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="linkedin" className="flex items-center gap-2">
                      <Linkedin className="h-4 w-4" /> LinkedIn Username
                    </Label>
                    <Input
                      id="linkedin"
                      value={linkedin}
                      onChange={(e) => setLinkedin(formatLinkedinUsername(e.target.value))}
                      placeholder="yourusername"
                      className="bg-white/50"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="interests">Interests</Label>
                  <div className="flex space-x-2">
                    <Input
                      id="interests"
                      value={interest}
                      onChange={(e) => setInterest(e.target.value)}
                      onKeyDown={handleKeyDown}
                      placeholder="Add your tech interests (e.g., React, Python, AI)"
                      className="bg-white/50"
                    />
                    <Button type="button" onClick={handleAddInterest} variant="outline">
                      Add
                    </Button>
                  </div>

                  <div className="flex flex-wrap gap-2 mt-3">
                    {interests.map((item, index) => (
                      <Badge
                        key={index}
                        variant="secondary"
                        className="bg-gradient-to-r from-pink-100 to-purple-100 text-purple-800 pl-3 pr-2 py-1.5 flex items-center gap-1"
                      >
                        {item}
                        <button
                          type="button"
                          onClick={() => handleRemoveInterest(item)}
                          className="ml-1 rounded-full hover:bg-purple-200 p-0.5"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </Badge>
                    ))}
                    {interests.length === 0 && (
                      <p className="text-sm text-gray-500">Add some interests to help us find your perfect match</p>
                    )}
                  </div>
                </div>
              </form>
            </CardContent>

            <CardFooter>
              <Button
                onClick={handleSubmit}
                disabled={isSubmitting}
                className="w-full bg-gradient-to-r from-pink-500 to-purple-600 hover:from-pink-600 hover:to-purple-700"
              >
                {isSubmitting ? "Saving..." : "Continue"}
              </Button>
            </CardFooter>
          </Card>
        </motion.div>
      </div>
    </main>
  )
}
