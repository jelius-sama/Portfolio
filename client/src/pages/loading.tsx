import { useState, useEffect } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { RefreshCw, Loader } from "lucide-react"

export default function Loading() {
  const [progress, setProgress] = useState(0)
  const [currentTask, setCurrentTask] = useState("Initializing...")
  const [dots, setDots] = useState("")


  const tasks = [
    "Initializing system...",
    "Loading configuration files...",
    "Finalizing setup...",
  ]

  useEffect(() => {
    // Simulate loading progress
    const progressInterval = setInterval(() => {
      setProgress((prev) => {
        const newProgress = prev + Math.random() * 3
        if (newProgress >= 94) {
          clearInterval(progressInterval)
          setCurrentTask("Finalizing setup...")
          return 94
        }

        const taskIndex = Math.floor((newProgress / 100) * tasks.length)
        setCurrentTask(tasks[taskIndex] || tasks[tasks.length - 1])
        return newProgress
      })
    }, 150)

    // Animate dots
    const dotsInterval = setInterval(() => {
      setDots((prev) => {
        if (prev.length >= 3) return ""
        return prev + "."
      })
    }, 500)

    return () => {
      clearInterval(progressInterval)
      clearInterval(dotsInterval)
    }
  }, [])

  const getProgressColor = (progress: number) => {
    if (progress < 30) return "bg-red-400"
    if (progress < 70) return "bg-yellow-400"
    return "bg-green-400"
  }

  const getSystemStatus = (progress: number) => {
    if (progress < 25) return { status: "starting", color: "text-yellow-400", icon: "⚡" }
    if (progress < 50) return { status: "loading", color: "text-blue-400", icon: "🔄" }
    if (progress < 75) return { status: "processing", color: "text-orange-400", icon: "⚙️" }
    if (progress < 100) return { status: "finalizing", color: "text-green-400", icon: "✨" }
    return { status: "ready", color: "text-green-400", icon: "✅" }
  }

  const systemStatus = getSystemStatus(progress)

  return (
    <section className="max-w-6xl mx-auto pt-20 pb-12 px-4 sm:px-6">
      {/* Header */}
      <div className="text-center mb-8">
        <div className="w-20 h-20 mx-auto mb-6 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center">
          <Loader className="text-black animate-spin" size={32} />
        </div>
        <h1 className="text-3xl font-bold text-white mb-2">Loading</h1>
        <p className="text-gray-400 font-mono">Please wait while we prepare everything for you</p>
      </div>

      {/* Main Loading Terminal */}
      <TerminalWindow title="system-loader" className="mb-8">
        <div className="font-mono text-sm">
          <div className="text-orange-400 mb-4">$ ./portfolio-loader --verbose</div>

          {/* Progress Bar */}
          <div className="mb-6">
            <div className="flex justify-between text-xs mb-2">
              <span className="text-gray-400">Loading Progress</span>
              <span className="text-orange-400">{Math.round(progress)}%</span>
            </div>
            <div className="w-full bg-gray-800 rounded-full h-3 border border-orange-500/30">
              <div
                className={`h-3 rounded-full transition-all duration-300 ${getProgressColor(progress)}`}
                style={{ width: `${progress}%` }}
              ></div>
            </div>
          </div>

          {/* Current Task */}
          <div className="mb-6">
            <div className="flex items-center gap-2 text-yellow-400 mb-2">
              <RefreshCw className="animate-spin" size={14} />
              <span>
                {currentTask}
                {dots}
              </span>
            </div>
          </div>

          {/* System Status */}
          <div className="space-y-2 mb-6">
            <div className="text-orange-400 mb-2">$ systemctl status Portfolio.service</div>
            <div className={`${systemStatus.color} mb-2`}>
              {systemStatus.icon} Portfolio.service - Portfolio Application
            </div>
            <div className="text-gray-300 space-y-1 ml-4">
              <div>Loaded: loaded (/etc/systemd/system/Portfolio.service; enabled)</div>
              <div>Tasks: {Math.floor((progress / 100) * 3)}/3 completed</div>
            </div>
          </div>
        </div>
      </TerminalWindow>

      {/* Footer */}
      <div className="text-center text-gray-500 text-sm font-mono">
        <div>Loading System • Powered by Go</div>
        <div className="mt-2">
          {progress < 100 ? "Please wait while we prepare your experience..." : "Ready to proceed!"}
        </div>
      </div>
    </section>
  )
}
