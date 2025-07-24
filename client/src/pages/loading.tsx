import { useState, useEffect } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { Home, RefreshCw, Loader } from "lucide-react"

export function LoadingPage() {
  const [progress, setProgress] = useState(0)
  const [currentTask, setCurrentTask] = useState("Initializing...")
  const [dots, setDots] = useState("")

  const tasks = [
    "Initializing system...",
    "Loading configuration files...",
    "Establishing database connection...",
    "Fetching user data...",
    "Rendering components...",
    "Optimizing assets...",
    "Finalizing setup...",
  ]

  useEffect(() => {
    // Simulate loading progress
    const progressInterval = setInterval(() => {
      setProgress((prev) => {
        const newProgress = prev + Math.random() * 3
        if (newProgress >= 100) {
          clearInterval(progressInterval)
          setCurrentTask("Loading complete!")
          return 100
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
    <div className="min-h-screen bg-gray-950 text-white flex items-center justify-center px-6">
      <div className="max-w-2xl mx-auto w-full">
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
              <div className="text-orange-400 mb-2">$ systemctl status portfolio.service</div>
              <div className={`${systemStatus.color} mb-2`}>
                {systemStatus.icon} portfolio.service - Portfolio Application
              </div>
              <div className="text-gray-300 space-y-1 ml-4">
                <div>Loaded: loaded (/etc/systemd/system/portfolio.service; enabled)</div>
                <div>
                  Active: <span className={systemStatus.color}>{systemStatus.status}</span> since{" "}
                  {new Date().toLocaleTimeString()}
                </div>
                <div>Memory: {Math.round(45 + (progress / 100) * 20)}M</div>
                <div>CPU: {((progress / 100) * 15).toFixed(1)}%</div>
                <div>Tasks: {Math.floor((progress / 100) * 7)}/7 completed</div>
              </div>
            </div>

            {/* Loading Log */}
            <div className="text-orange-400 mb-2">$ tail -f /var/log/portfolio-loader.log</div>
            <div className="bg-gray-900 border border-orange-500/20 rounded p-3 max-h-32 overflow-y-auto">
              <div className="text-gray-300 space-y-1 text-xs">
                <div>[{new Date().toLocaleTimeString()}] INFO: Starting portfolio application</div>
                <div>[{new Date().toLocaleTimeString()}] INFO: Loading React components...</div>
                <div>[{new Date().toLocaleTimeString()}] INFO: Initializing terminal interface</div>
                <div>[{new Date().toLocaleTimeString()}] INFO: Connecting to analytics service</div>
                {progress > 30 && <div>[{new Date().toLocaleTimeString()}] INFO: Assets optimization in progress</div>}
                {progress > 60 && <div>[{new Date().toLocaleTimeString()}] INFO: Preparing user interface</div>}
                {progress > 90 && <div>[{new Date().toLocaleTimeString()}] INFO: Final checks completed</div>}
                {progress >= 100 && (
                  <div className="text-green-400">[{new Date().toLocaleTimeString()}] SUCCESS: Portfolio ready!</div>
                )}
              </div>
            </div>
          </div>
        </TerminalWindow>

        {/* System Resources */}
        <TerminalWindow title="system-resources" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-4">$ htop --summary</div>
            <div className="grid grid-cols-2 gap-6">
              <div>
                <div className="text-gray-400 mb-2">CPU Usage</div>
                <div className="flex items-center gap-2">
                  <div className="w-20 bg-gray-800 rounded-full h-2">
                    <div
                      className="bg-blue-400 h-2 rounded-full transition-all duration-300"
                      style={{ width: `${Math.min(progress * 0.8, 80)}%` }}
                    ></div>
                  </div>
                  <span className="text-blue-400 text-xs">{Math.round(progress * 0.8)}%</span>
                </div>
              </div>
              <div>
                <div className="text-gray-400 mb-2">Memory</div>
                <div className="flex items-center gap-2">
                  <div className="w-20 bg-gray-800 rounded-full h-2">
                    <div
                      className="bg-green-400 h-2 rounded-full transition-all duration-300"
                      style={{ width: `${Math.min(progress * 0.6, 60)}%` }}
                    ></div>
                  </div>
                  <span className="text-green-400 text-xs">{Math.round(progress * 0.6)}%</span>
                </div>
              </div>
              <div>
                <div className="text-gray-400 mb-2">Network</div>
                <div className="flex items-center gap-2">
                  <div className="w-20 bg-gray-800 rounded-full h-2">
                    <div
                      className="bg-purple-400 h-2 rounded-full transition-all duration-300"
                      style={{ width: `${Math.min(progress * 0.4, 40)}%` }}
                    ></div>
                  </div>
                  <span className="text-purple-400 text-xs">{Math.round(progress * 0.4)}%</span>
                </div>
              </div>
              <div>
                <div className="text-gray-400 mb-2">Disk I/O</div>
                <div className="flex items-center gap-2">
                  <div className="w-20 bg-gray-800 rounded-full h-2">
                    <div
                      className="bg-yellow-400 h-2 rounded-full transition-all duration-300"
                      style={{ width: `${Math.min(progress * 0.3, 30)}%` }}
                    ></div>
                  </div>
                  <span className="text-yellow-400 text-xs">{Math.round(progress * 0.3)}%</span>
                </div>
              </div>
            </div>
          </div>
        </TerminalWindow>

        {/* Service Status */}
        <TerminalWindow title="service-status" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-4">$ docker ps --format "table"</div>
            <div className="space-y-2">
              <div className="grid grid-cols-4 gap-4 text-xs text-gray-400 border-b border-orange-500/20 pb-2">
                <span>SERVICE</span>
                <span>STATUS</span>
                <span>UPTIME</span>
                <span>HEALTH</span>
              </div>
              <div className="grid grid-cols-4 gap-4 text-xs">
                <span className="text-gray-300">portfolio-web</span>
                <span className={progress > 20 ? "text-green-400" : "text-yellow-400"}>
                  {progress > 20 ? "Running" : "Starting"}
                </span>
                <span className="text-gray-400">{Math.floor(progress / 10)}s</span>
                <span className={progress > 20 ? "text-green-400" : "text-yellow-400"}>
                  {progress > 20 ? "Healthy" : "Pending"}
                </span>
              </div>
              <div className="grid grid-cols-4 gap-4 text-xs">
                <span className="text-gray-300">analytics-db</span>
                <span className={progress > 40 ? "text-green-400" : "text-yellow-400"}>
                  {progress > 40 ? "Running" : "Starting"}
                </span>
                <span className="text-gray-400">{Math.floor(progress / 15)}s</span>
                <span className={progress > 40 ? "text-green-400" : "text-yellow-400"}>
                  {progress > 40 ? "Healthy" : "Pending"}
                </span>
              </div>
              <div className="grid grid-cols-4 gap-4 text-xs">
                <span className="text-gray-300">blog-service</span>
                <span className={progress > 60 ? "text-green-400" : "text-yellow-400"}>
                  {progress > 60 ? "Running" : "Starting"}
                </span>
                <span className="text-gray-400">{Math.floor(progress / 20)}s</span>
                <span className={progress > 60 ? "text-green-400" : "text-yellow-400"}>
                  {progress > 60 ? "Healthy" : "Pending"}
                </span>
              </div>
            </div>
          </div>
        </TerminalWindow>

        {/* Navigation */}
        {progress >= 100 && (
          <TerminalWindow title="ready" className="mb-8">
            <div className="font-mono text-sm">
              <div className="text-orange-400 mb-4">$ echo "System ready! Available actions:"</div>
              <div className="space-y-3">
                <a
                  href="/"
                  className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
                >
                  <Home className="text-orange-400 group-hover:text-orange-300" size={20} />
                  <div>
                    <div className="text-white font-medium">Continue to Portfolio</div>
                    <div className="text-gray-400 text-xs">Access the main application</div>
                  </div>
                </a>
                <button
                  onClick={() => window.location.reload()}
                  className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group w-full"
                >
                  <RefreshCw className="text-orange-400 group-hover:text-orange-300" size={20} />
                  <div>
                    <div className="text-white font-medium">Reload Application</div>
                    <div className="text-gray-400 text-xs">Restart the loading process</div>
                  </div>
                </button>
              </div>
            </div>
          </TerminalWindow>
        )}

        {/* Footer */}
        <div className="text-center text-gray-500 text-sm font-mono">
          <div>Loading System v2.1.0 • Powered by React & Next.js</div>
          <div className="mt-2">
            {progress < 100 ? "Please wait while we prepare your experience..." : "Ready to proceed!"}
          </div>
        </div>
      </div>
    </div>
  )
}
