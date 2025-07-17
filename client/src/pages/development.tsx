import { TerminalWindow } from "@/components/layout/terminal-window"
import { Github, Mail, Wrench } from "lucide-react"
import { suggestions } from "@/pages/not-found"
import { Link, useLocation } from "react-router-dom"
import { useConfig } from "@/contexts/config"
import { StaticMetadata } from "@/contexts/metadata"
import { Fragment } from "react"

export default function DevelopmentPage() {
  const location = useLocation()
  const { app: { portfolio: me } } = useConfig()

  return (
    <Fragment>
      <StaticMetadata />

      <div className="min-h-screen flex items-center justify-center px-6">
        <div className="w-6xl mx-auto px-6 py-4">
          {/* Header */}
          <div className="text-center my-8">
            <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center">
              <Wrench className="text-black" size={24} />
            </div>
            <h1 className="text-3xl font-bold text-white mb-2">Under Development</h1>
            <p className="text-gray-400 font-mono">This page is currently being built</p>
          </div>

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
              <a
                href={me.links.GitHub.link}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
              >
                <Github className="text-orange-400 group-hover:text-orange-300" size={20} />
                <div>
                  <div className="text-white font-medium font-mono">Follow Development</div>
                  <div className="text-gray-400 text-sm">Watch progress on GitHub</div>
                </div>
              </a>

              <a
                href={`${me.links.Email.link}?subject=Development Page Feedback`}
                className="flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
              >
                <Mail className="text-orange-400 group-hover:text-orange-300" size={20} />
                <div>
                  <div className="text-white font-medium font-mono">Send Feedback</div>
                  <div className="text-gray-400 text-sm">Share ideas or suggestions</div>
                </div>
              </a>
            </div>
          </TerminalWindow>

          {/* Footer */}
          <div className="text-center text-gray-500 text-sm font-mono">
            <div>Development Environment • Go 1.24.5</div>
            <div className="mt-2">Last commit: {Math.random().toString(36).substr(2, 7)} • Branch: master{location.pathname}</div>
          </div>
        </div>
      </div>
    </Fragment>
  )
}
