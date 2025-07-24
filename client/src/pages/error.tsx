import { useState } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { useConfig } from "@/contexts/config"
import { Home, RefreshCw, AlertTriangle, Bug, Mail, ExternalLink } from "lucide-react"

export default function ClientError({ error }: { error?: Error }) {
  const { app } = useConfig()
  const [timestamp] = useState(() => new Date().toISOString())

  if (error) {
    console.log("Caught an error:", error)
  }

  const errorDetails = {
    type: "CLIENT_ERROR",
    message: "An unexpected error occurred in the client application",
    stack: [
      "at ComponentDidCatch (Portfolio.tsx:45:12)",
      "at ErrorBoundary.render (ErrorBoundary.tsx:23:8)",
      "at ReactDOM.render (react-dom.js:1234:16)",
      "at Application.start (main.tsx:12:4)",
    ],
    userAgent: typeof navigator !== "undefined" ? navigator.userAgent : "Unknown",
    url: typeof window !== "undefined" ? window.location.href : "Unknown",
  }

  const troubleshootingSteps = [
    "Clear your browser cache and cookies",
    "Disable browser extensions temporarily",
    "Try using an incognito/private browsing window",
    "Check your internet connection",
    "Update your browser to the latest version",
    "Try accessing the site from a different device",
  ]

  return (
    <section className="max-w-6xl mx-auto pt-20 pb-12 px-4 sm:px-6">
      {/* Header */}
      <div className="text-center mb-8">
        <div className="w-20 h-20 mx-auto mb-6 rounded-full bg-gradient-to-br from-red-500 to-red-600 flex items-center justify-center">
          <AlertTriangle className="text-white" size={32} />
        </div>
        <h1 className="text-3xl font-bold text-white mb-2">Client Error</h1>
        <p className="text-gray-400 font-mono">Something went wrong in the browser</p>
      </div>

      {/* Error Details */}
      <TerminalWindow title="error-report" className="mb-8">
        <div className="font-mono text-sm">
          <div className="text-orange-400 mb-4">$ cat /var/log/client-error.log</div>
          <div className="bg-red-500/10 border border-red-500/30 rounded p-4 mb-4">
            <div className="flex items-center gap-2 text-red-400 mb-2">
              <Bug size={16} />
              <span className="font-bold">ERROR DETECTED</span>
            </div>
            <div className="text-gray-300 space-y-1">
              <div>
                <span className="text-red-400">Timestamp:</span> {timestamp}
              </div>
              <div>
                <span className="text-red-400">Type:</span> {errorDetails.type}
              </div>
              <div>
                <span className="text-red-400">Message:</span> {errorDetails.message}
              </div>
            </div>
          </div>

          <div className="text-orange-400 mb-2">$ browser-info --details</div>
          <div className="text-gray-300 space-y-1">
            <div>
              <span className="text-orange-400">URL:</span> {errorDetails.url}
            </div>
            <div>
              <span className="text-orange-400">User Agent:</span> {errorDetails.userAgent.substring(0, 80)}...
            </div>
            {localStorage.getItem("session_id") && (
              <div>
                <span className="text-orange-400">Session ID:</span> {localStorage.getItem("session_id")}
              </div>
            )}
          </div>
        </div>
      </TerminalWindow>

      {/* Troubleshooting */}
      <TerminalWindow title="troubleshooting-guide" className="mb-8">
        <div className="font-mono text-sm">
          <div className="text-orange-400 mb-4">$ cat /usr/share/docs/troubleshooting.md</div>
          <div className="text-gray-300 mb-4">Try these steps to resolve the issue:</div>
          <div className="space-y-2">
            {troubleshootingSteps.map((step, index) => (
              <div key={index} className="flex items-start gap-3">
                <span className="text-orange-400 mt-1">{index + 1}.</span>
                <span className="text-gray-300">{step}</span>
              </div>
            ))}
          </div>

          <div className="mt-6 p-3 bg-blue-500/10 border border-blue-500/30 rounded">
            <div className="flex items-center gap-2 text-blue-400 mb-2">
              <AlertTriangle size={16} />
              <span className="font-bold">Developer Note</span>
            </div>
            <div className="text-gray-300 text-xs">
              <div>This error has been automatically reported to our monitoring system.</div>
            </div>
          </div>
        </div>
      </TerminalWindow>

      {/* Actions */}
      <TerminalWindow title="recovery-actions" className="mb-8">
        <div className="font-mono text-sm">
          <div className="text-orange-400 mb-4">$ ls /recovery/actions/</div>
          <div className="space-y-3">
            <button
              onClick={window.location.reload}
              className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group w-full"
            >
              <RefreshCw className="text-orange-400 group-hover:text-orange-300" size={20} />
              <div className="text-left">
                <div className="text-white font-medium">Retry Application</div>
                <div className="text-gray-400 text-xs">Reload the page and try again</div>
              </div>
            </button>

            <a
              href="/"
              className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
            >
              <Home className="text-orange-400 group-hover:text-orange-300" size={20} />
              <div>
                <div className="text-white font-medium">Return to Homepage</div>
                <div className="text-gray-400 text-xs">Go back to the main portfolio</div>
              </div>
            </a>

            <a
              href={`mailto:work@${app.portfolio.links.Email.domain}?subject=Client Error Report&body=Timestamp: ${timestamp}%0AUser Agent: ${errorDetails.userAgent}`}
              className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
            >
              <Mail className="text-orange-400 group-hover:text-orange-300" size={20} />
              <div>
                <div className="text-white font-medium">Report This Error</div>
                <div className="text-gray-400 text-xs">Send error details to support team</div>
              </div>
            </a>

            <a
              href={`${app.portfolio.links.GitHub.link}/Portfolio/issues`}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
            >
              <ExternalLink className="text-orange-400 group-hover:text-orange-300" size={20} />
              <div>
                <div className="text-white font-medium">Report on GitHub</div>
                <div className="text-gray-400 text-xs">Create an issue in the repository</div>
              </div>
            </a>
          </div>
        </div>
      </TerminalWindow>

      {/* Footer */}
      <div className="text-center text-gray-500 text-sm font-mono">
        <div className="mt-2">
          Timestamp: {new Date(timestamp).toLocaleString()}
        </div>
      </div>
    </section>
  )
}
