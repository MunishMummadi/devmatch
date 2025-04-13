import { LoginForm } from "@/components/login-form"
import { FloatingShapes } from "@/components/floating-shapes"

export default function Home() {
  return (
    <main className="min-h-screen flex flex-col items-center justify-center p-4 relative overflow-hidden">
      <FloatingShapes />
      <div className="w-full max-w-md z-10">
        <div className="mb-8 flex items-center">
          <h1 className="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-indigo-600">
            DEVMatch.
          </h1>
        </div>
        <LoginForm />
      </div>
    </main>
  )
}
