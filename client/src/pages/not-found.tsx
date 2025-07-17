import { TerminalWindow } from "@/components/layout/terminal-window"
import { Home, LinkIcon } from "lucide-react"
import { Fragment } from "react"
import { StaticMetadata } from "@/contexts/metadata"
import { useLocation, Link } from "react-router-dom"

export const suggestions = [
  { title: "Portfolio", description: "Go back to the main page", url: "/", icon: <Home size={20} /> },
  { title: "Links", description: "Find all my social links", url: "/links", icon: <LinkIcon size={20} /> },
]

export default function NotFoundPage() {
  const location = useLocation()

  return (
    <Fragment>
      <StaticMetadata />

      <main className="flex items-center justify-center px-6 py-12">
        <div className="max-w-6xl mx-auto px-6 py-4">
          <div className="text-center my-8">
            <div className="text-6xl font-bold text-orange-400 mb-4 font-mono">404</div>
            <h1 className="text-2xl font-bold text-white mb-2">Page Not Found</h1>
            <p className="text-gray-400 font-mono">The requested resource could not be located</p>
          </div>

          {/* ASCII Art */}
          <TerminalWindow title="ascii-art" className="mb-8">
            <div className="font-mono text-xs text-orange-400 text-center">
              <pre>{`
    ╔══════════════════════════════════════╗
    ║                                      ║
    ║    ┌─┐┌─┐┌─┐┌─┐  ┌┐┌┌─┐┌┬┐           ║
    ║    ├─┘├─┤│ ┬├┤   ││││ │ │            ║
    ║    ┴  ┴ ┴└─┘└─┘  ┘└┘└─┘ ┴            ║
    ║              ┌─┐┌─┐┬ ┬┌┐┌┌┬┐         ║
    ║              ├┤ │ ││ ││││ ││         ║
    ║              └  └─┘└─┘┘└┘─┴┘         ║
    ║                                      ║
    ╚══════════════════════════════════════╝
            `}</pre>
            </div>
          </TerminalWindow>

          {/* Terminal Error Output */}
          <TerminalWindow title="error-404" className="mb-8">
            <div className="font-mono text-sm">
              <div className="text-orange-400 mb-2">$ ls -la {location.pathname}</div>
              <div className="text-red-400 mb-4">ls: cannot access '{location.pathname}': No such file or directory</div>

              <div className="text-orange-400 mb-2">$ find / -name "{location.pathname}" 2&gt;/dev/null</div>
              <div className="text-gray-300 mb-4">Search completed. No matching files found.</div>

              <div className="text-orange-400 mb-2">$ cat /var/log/jelius.dev/error.log | tail -1</div>
              <div className="text-yellow-400 mb-4">
                {new Date().toISOString()} [error] 404: File not found - The page you're looking for has moved or never
                existed
              </div>

              <div className="text-orange-400 mb-2">$ systemctl status page-finder.service</div>
              <div className="text-red-400 mb-2">● page-finder.service - Page Location Service</div>
              <div className="text-gray-300 space-y-1 mb-4">
                <div> Loaded: loaded (/etc/systemd/system/page-finder.service; enabled)</div>
                <div>
                  {" "}
                  Active: <span className="text-red-400">failed (Result: exit-code)</span>
                </div>
                <div> Process: 1337 ExecStart=/usr/bin/find-page (code=exited, status=404)</div>
              </div>

              <div className="text-orange-400 mb-2">$ echo "Don't worry, let's get you back on track!"</div>
              <div className="text-gray-300">Don't worry, let's get you back on track!</div>
            </div>
          </TerminalWindow>

          {/* Navigation Suggestions */}
          <TerminalWindow title="suggested-paths" className="mb-8">
            <div className="space-y-3">
              <div className="text-orange-400 font-mono text-sm mb-4">$ ls /available/routes</div>
              {suggestions.map((suggestion, index) => (
                <Link
                  key={index}
                  to={suggestion.url}
                  className="block p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
                >
                  <div className="flex items-center gap-3">
                    <div className="text-orange-400 group-hover:text-orange-300">{suggestion.icon}</div>
                    <div>
                      <div className="text-white font-medium font-mono">{suggestion.title}</div>
                      <div className="text-gray-400 text-sm">{suggestion.description}</div>
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          </TerminalWindow>

          {/* Footer */}
          <div className="text-center text-gray-500 text-sm font-mono">
            <div>Error Code: HTTP 404 • Powered by Go 1.24.5 (net/http)</div>
            <div className="mt-2">If you believe this is an error, please contact the administrator</div>
          </div>
        </div>
      </main>
    </Fragment>
  )
}
