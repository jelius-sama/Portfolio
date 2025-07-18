import { TerminalWindow } from "@/components/layout/terminal-window"
import { Github, Mail, Wrench } from "lucide-react"
import { suggestions } from "@/pages/not-found"
import { Link } from "react-router-dom"
import { useConfig } from "@/contexts/config"
import { StaticMetadata } from "@/contexts/metadata"
import { Fragment } from "react"
import { useQuery } from '@tanstack/react-query'
import { toast } from "sonner"
import { Skeleton } from "@/components/ui/skeleton"
import { useTimeAgo } from "@/lib/utils"

type CommitResponse = {
  sha: string
  commit: {
    message: string
    author: {
      name: string
      email: string
      date: string
    }
  }
  html_url: string
}

export default function DevelopmentPage() {
  const { app: { portfolio: me } } = useConfig()

  const { isPending, error, data } = useQuery({
    queryKey: ['latestCommit'],
    queryFn: async () => {
      const res = await fetch('/api/latest_commit?repo=Portfolio&branch=main')
      if (!res.ok) {
        const errMsg = "Failed to fetch last commit!"
        toast.error(errMsg)
        throw new Error(errMsg)
      }
      return res.json() as Promise<CommitResponse>
    }
  })

  const timeAgo = useTimeAgo(data?.commit.author.date)

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
          <div className="flex flex-col justify-center items-center text-gray-500 text-sm font-mono">
            <div>Development Environment • Go 1.24.5</div>
            {isPending ? (
              <div className="mt-2 flex items-center justify-center">Last commit: <Skeleton className="h-3 w-[100px] mx-1" /> • Branch: Portfolio/main</div>
            ) : error ? (
              <div className="mt-2">{error.message}</div>
            ) : data && (
              <div className="mt-2">Last commit: {data.sha.slice(0, 7)} • {timeAgo} • Branch: Portfolio/main</div>
            )}
          </div>
        </div>
      </div>
    </Fragment>
  )
}
