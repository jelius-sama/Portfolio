import { useState } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { Home, RefreshCw, AlertTriangle, Bug, Code, Mail, ExternalLink } from "lucide-react"

export default function ClientError() {
  const [errorCode] = useState(() => Math.random().toString(36).substr(2, 8).toUpperCase())
  const [timestamp] = useState(() => new Date().toISOString())
  const [retryCount, setRetryCount] = useState(0)
  const [diagnosticsRunning, setDiagnosticsRunning] = useState(false)

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

  const runDiagnostics = () => {
    setDiagnosticsRunning(true)
    setTimeout(() => {
      setDiagnosticsRunning(false)
    }, 3000)
  }

  const handleRetry = () => {
    setRetryCount((prev) => prev + 1)
    window.location.reload()
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
    <div className="min-h-screen bg-gray-950 text-white flex items-center justify-center px-6">
      <div className="max-w-3xl mx-auto w-full">
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
                <span className="font-bold">CRITICAL ERROR DETECTED</span>
              </div>
              <div className="text-gray-300 space-y-1">
                <div>
                  <span className="text-red-400">Error ID:</span> {errorCode}
                </div>
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

            <div className="text-orange-400 mb-2">$ echo "Error Stack Trace:"</div>
            <div className="bg-gray-900 border border-orange-500/20 rounded p-3 mb-4">
              <div className="text-red-300 space-y-1 text-xs">
                {errorDetails.stack.map((line, index) => (
                  <div key={index} className="font-mono">
                    {line}
                  </div>
                ))}
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
              <div>
                <span className="text-orange-400">Retry Count:</span> {retryCount}
              </div>
              <div>
                <span className="text-orange-400">Session ID:</span> sess_{Math.random().toString(36).substr(2, 12)}
              </div>
            </div>
          </div>
        </TerminalWindow>

        {/* System Diagnostics */}
        <TerminalWindow title="system-diagnostics" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-4">$ ./run-diagnostics.sh --client-side</div>

            {!diagnosticsRunning ? (
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-gray-300">JavaScript Engine:</span>
                  <span className="text-green-400">✓ Operational</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-300">DOM Rendering:</span>
                  <span className="text-yellow-400">⚠ Degraded</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-300">Network Connection:</span>
                  <span className="text-green-400">✓ Connected</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-300">Local Storage:</span>
                  <span className="text-green-400">✓ Available</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-300">Service Worker:</span>
                  <span className="text-red-400">✗ Error</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-300">React DevTools:</span>
                  <span className="text-gray-400">- Not Detected</span>
                </div>

                <button
                  onClick={runDiagnostics}
                  className="mt-4 flex items-center gap-2 px-4 py-2 bg-orange-500/20 border border-orange-500/50 rounded hover:bg-orange-500/30 transition-colors"
                >
                  <Code size={16} className="text-orange-400" />
                  <span className="text-orange-400">Run Full Diagnostics</span>
                </button>
              </div>
            ) : (
              <div className="space-y-2">
                <div className="text-yellow-400 mb-4">Running comprehensive diagnostics...</div>
                <div className="space-y-1 text-gray-300">
                  <div>→ Checking browser compatibility...</div>
                  <div>→ Validating JavaScript execution context...</div>
                  <div>→ Testing DOM manipulation capabilities...</div>
                  <div>→ Verifying network connectivity...</div>
                  <div>→ Analyzing error patterns...</div>
                  <div>→ Generating diagnostic report...</div>
                </div>
              </div>
            )}
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
                <div>Error ID: {errorCode} can be used for tracking purposes.</div>
                <div>If the issue persists, please contact support with this error ID.</div>
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
                onClick={handleRetry}
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
                href="mailto:support@yourname.dev?subject=Client Error Report&body=Error ID: {errorCode}%0ATimestamp: {timestamp}%0AUser Agent: {errorDetails.userAgent}"
                className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
              >
                <Mail className="text-orange-400 group-hover:text-orange-300" size={20} />
                <div>
                  <div className="text-white font-medium">Report This Error</div>
                  <div className="text-gray-400 text-xs">Send error details to support team</div>
                </div>
              </a>

              <a
                href="https://github.com/yourusername/portfolio/issues"
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
          <div>Error Handler v1.2.0 • Client-Side Error Boundary</div>
          <div className="mt-2">
            Error ID: {errorCode} • Timestamp: {new Date(timestamp).toLocaleString()}
          </div>
        </div>
      </div>
    </div>
  )
}
