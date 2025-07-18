import { TerminalWindow } from "@/components/layout/terminal-window"
import { Home, Mail, RefreshCw, AlertTriangle } from "lucide-react"

/** 
 * TODO: Work Needed 
*/
export default function ServerErrorPage() {
  const troubleshootingSteps = [
    "Check server logs for detailed error information",
    "Verify database connectivity and status",
    "Restart application services if necessary",
    "Contact system administrator if issue persists",
  ]

  return (
    <div className="min-h-screen bg-gray-950 text-white flex items-center justify-center px-6">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="text-6xl font-bold text-red-400 mb-4 font-mono">500</div>
          <h1 className="text-2xl font-bold text-white mb-2">Internal Server Error</h1>
          <p className="text-gray-400 font-mono">Something went wrong on our end</p>
        </div>

        {/* Terminal Error Output */}
        <TerminalWindow title="error-500" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-2">$ systemctl status portfolio.service</div>
            <div className="text-red-400 mb-2">● portfolio.service - Portfolio Application</div>
            <div className="text-gray-300 space-y-1 mb-4">
              <div> Loaded: loaded (/etc/systemd/system/portfolio.service; enabled)</div>
              <div>
                {" "}
                Active: <span className="text-red-400">failed (Result: exit-code)</span>
              </div>
              <div> Process: 2048 ExecStart=/usr/bin/portfolio-app (code=exited, status=1)</div>
              <div> Main PID: 2048 (code=exited, status=1)</div>
            </div>

            <div className="text-orange-400 mb-2">$ tail -f /var/log/application/error.log</div>
            <div className="text-red-400 space-y-1 mb-4">
              <div>[{new Date().toISOString()}] ERROR: Unhandled exception in request handler</div>
              <div>[{new Date().toISOString()}] ERROR: Database connection timeout after 30s</div>
              <div>[{new Date().toISOString()}] ERROR: Failed to process request - Internal server error</div>
            </div>

            <div className="text-orange-400 mb-2">$ ps aux | grep portfolio</div>
            <div className="text-gray-300 mb-4">No running processes found for 'portfolio'</div>

            <div className="text-orange-400 mb-2">$ systemctl restart portfolio.service</div>
            <div className="text-yellow-400 mb-4">Attempting to restart service...</div>

            <div className="text-orange-400 mb-2">$ echo "Our team has been notified and is working on a fix"</div>
            <div className="text-gray-300">Our team has been notified and is working on a fix</div>
          </div>
        </TerminalWindow>

        {/* System Status */}
        <TerminalWindow title="system-status" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-4">$ curl -s https://status.yourname.dev/api/status</div>
            <div className="space-y-2">
              <div className="flex justify-between">
                <span className="text-gray-300">Web Server:</span>
                <span className="text-red-400">● DOWN</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">Database:</span>
                <span className="text-yellow-400">⚠ DEGRADED</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">CDN:</span>
                <span className="text-green-400">● OPERATIONAL</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">Load Balancer:</span>
                <span className="text-green-400">● OPERATIONAL</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">Monitoring:</span>
                <span className="text-green-400">● OPERATIONAL</span>
              </div>
            </div>

            <div className="mt-4 p-3 bg-yellow-500/10 border border-yellow-500/30 rounded">
              <div className="flex items-center gap-2 text-yellow-400 mb-2">
                <AlertTriangle size={16} />
                <span className="font-bold">Incident Report</span>
              </div>
              <div className="text-gray-300 text-xs">
                <div>Incident ID: INC-{Math.random().toString(36).substr(2, 8).toUpperCase()}</div>
                <div>Started: {new Date().toLocaleString()}</div>
                <div>Status: Investigating</div>
                <div>ETA: 15-30 minutes</div>
              </div>
            </div>
          </div>
        </TerminalWindow>

        {/* Troubleshooting */}
        <TerminalWindow title="troubleshooting" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-4">$ cat /usr/share/docs/troubleshooting.md</div>
            <div className="text-gray-300 mb-4">While we work on fixing this issue, you can try:</div>
            <div className="space-y-2">
              {troubleshootingSteps.map((step, index) => (
                <div key={index} className="flex items-start gap-3">
                  <span className="text-orange-400 mt-1">•</span>
                  <span className="text-gray-300">{step}</span>
                </div>
              ))}
            </div>

            <div className="mt-6 space-y-3">
              <a
                href="/"
                className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
              >
                <Home className="text-orange-400 group-hover:text-orange-300" size={20} />
                <div>
                  <div className="text-white font-medium">Return to Homepage</div>
                  <div className="text-gray-400 text-xs">Go back to the main site</div>
                </div>
              </a>

              <a
                href="mailto:admin@yourname.dev?subject=Server Error Report"
                className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
              >
                <Mail className="text-orange-400 group-hover:text-orange-300" size={20} />
                <div>
                  <div className="text-white font-medium">Report This Error</div>
                  <div className="text-gray-400 text-xs">Send error details to administrator</div>
                </div>
              </a>

              <button
                onClick={() => window.location.reload()}
                className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group w-full"
              >
                <RefreshCw className="text-orange-400 group-hover:text-orange-300" size={20} />
                <div>
                  <div className="text-white font-medium">Retry Request</div>
                  <div className="text-gray-400 text-xs">Refresh the page to try again</div>
                </div>
              </button>
            </div>
          </div>
        </TerminalWindow>

        {/* Footer */}
        <div className="text-center text-gray-500 text-sm font-mono">
          <div>Error Code: HTTP 500 • Server: nginx/1.24.0 • Instance: i-0a1b2c3d4e5f6g7h8</div>
          <div className="mt-2">Request ID: {Math.random().toString(36).substr(2, 16).toUpperCase()}</div>
        </div>
      </div>
    </div>
  )
}
