"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { MessageSquare, Heart, User, LogOut } from "lucide-react"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"

export function DashboardHeader() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState("discover")
  const [userInitial, setUserInitial] = useState("U")

  useEffect(() => {
    // Get user initial from localStorage
    const savedUserData = localStorage.getItem("userProfile")
    if (savedUserData) {
      const userData = JSON.parse(savedUserData)
      if (userData.firstName) {
        setUserInitial(userData.firstName.charAt(0).toUpperCase())
      }
    }

    // Set active tab based on current path
    const path = window.location.pathname
    if (path.includes("/favorites")) {
      setActiveTab("favorites")
    } else if (path.includes("/chats")) {
      setActiveTab("chats")
    } else if (path.includes("/profile")) {
      setActiveTab("profile")
    } else {
      setActiveTab("discover")
    }
  }, [])

  return (
    <header className="backdrop-blur-sm bg-white/80 shadow-sm sticky top-0 z-20">
      <div className="container mx-auto px-4 py-3 flex items-center justify-between">
        <Link
          href="/dashboard"
          className="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-pink-600 to-purple-600"
        >
          DEVMatch
        </Link>

        <div className="flex items-center space-x-2">
          <Button
            variant={activeTab === "discover" ? "default" : "ghost"}
            size="sm"
            onClick={() => {
              setActiveTab("discover")
              router.push("/dashboard")
            }}
            className={activeTab === "discover" ? "bg-gradient-to-r from-pink-500 to-purple-600" : ""}
          >
            <User className="h-4 w-4 mr-2" />
            Discover
          </Button>

          <Button
            variant={activeTab === "favorites" ? "default" : "ghost"}
            size="sm"
            onClick={() => {
              setActiveTab("favorites")
              router.push("/dashboard/favorites")
            }}
            className={activeTab === "favorites" ? "bg-gradient-to-r from-pink-500 to-purple-600" : ""}
          >
            <Heart className="h-4 w-4 mr-2" />
            Favorites
          </Button>

          <Button
            variant={activeTab === "chats" ? "default" : "ghost"}
            size="sm"
            onClick={() => {
              setActiveTab("chats")
              router.push("/dashboard/chats")
            }}
            className={activeTab === "chats" ? "bg-gradient-to-r from-pink-500 to-purple-600" : ""}
          >
            <MessageSquare className="h-4 w-4 mr-2" />
            Chats
          </Button>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="rounded-full">
                <div className="h-8 w-8 rounded-full bg-gradient-to-r from-pink-500 to-purple-600 flex items-center justify-center text-white font-medium">
                  {userInitial}
                </div>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="backdrop-blur-sm bg-white/90">
              <DropdownMenuItem onClick={() => router.push("/profile")}>
                <User className="mr-2 h-4 w-4" />
                <span>Profile</span>
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => router.push("/")}>
                <LogOut className="mr-2 h-4 w-4" />
                <span>Logout</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  )
}
