import { Fragment, useState, useEffect } from "react"
import { StaticMetadata } from "@/contexts/metadata"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { Calendar, Clock, AwardIcon, EyeIcon } from "lucide-react"
import { Footer } from "@/components/ui/footer"
import { MarkdownRenderer } from "@/components/ui/markdown"
import { formatViews, formatDate } from "@/lib/utils"
import { useConfig } from "@/contexts/config"
import { useLocation } from "react-router-dom"

export default function Achievements() {
  const { ssrData, app: { portfolio: me } } = useConfig()
  const { pathname } = useLocation()

  const [stats, setStats] = useState<{ createdAt: number, updatedAt: number, views: number } | null>(null)
  const [markdownContent, setMarkdownContent] = useState<string>("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    (async () => {
      setLoading(true);
      setError(null);

      try {
        // Use SSR data if available and matches current path
        if (ssrData?.path === pathname) {
          if ("status" in ssrData && ssrData.status === 500) {
            throw new Error(`Something went wrong!`)
          }
          setStats(ssrData.api_resp);

          const markdownRes = await fetch(`/api/achievements_file`);
          if (!markdownRes.ok) {
            throw new Error(`Failed to fetch markdown: ${markdownRes.status}`);
          }

          const markdownText = await markdownRes.text();
          setMarkdownContent(markdownText);
        } else {
          // Fetch blog and markdown in parallel
          const [achievementsStats, markdownRes] = await Promise.all([
            fetch(`/api/achievements_stats`),
            fetch(`/api/achievements_file`)
          ]);

          if (!achievementsStats.ok) {
            throw new Error(`Failed to fetch achievements stats!`);
          }

          if (!markdownRes.ok) {
            throw new Error(`Failed to fetch markdown: ${markdownRes.status}`);
          }

          const data: { createdAt: number, updatedAt: number, views: number } = await achievementsStats.json();
          const markdownText = await markdownRes.text();

          setStats(data);
          setMarkdownContent(markdownText);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "An unexpected error occurred");
      } finally {
        setLoading(false);
      }
    })();
  }, [pathname]);

  useEffect(() => {
    if (!loading) {
      const event = new CustomEvent("PageLoaded", {
        detail: { pathname: window.location.pathname },
      });
      window.dispatchEvent(event);
    }
  }, [loading]);

  if (loading) {
    return (
      <Fragment>
        <StaticMetadata />

        <section className="max-w-6xl mx-auto pt-20 pb-12 px-4 sm:px-6">
          <TerminalWindow title="loading">
            <div className="font-mono text-sm text-center py-8">
              <div className="text-orange-400 mb-2">$ curl -X GET /api/achievements</div>
              <div className="text-gray-300 mb-4">Fetching achievements data...</div>
              <div className="text-orange-400 mb-2">$ curl -X GET /api/achievements_file</div>
              <div className="text-gray-300">Loading content...</div>
            </div>
          </TerminalWindow>
        </section>
      </Fragment>
    )
  }

  if (error) {
    return (
      <Fragment>
        <StaticMetadata />

        <section className="max-w-6xl mx-auto pt-20 pb-12 px-4 sm:px-6">
          <TerminalWindow title="error">
            <div className="font-mono text-sm">
              <div className="text-orange-400 mb-2">$ curl -X GET /api/achievements</div>
              <div className="text-red-400 mb-4">Error: {error}</div>
              <div className="text-orange-400 mb-2">$ echo "Please try again later"</div>
              <div className="text-gray-300">Please try again later</div>
            </div>
          </TerminalWindow>
        </section>
      </Fragment>
    )
  }

  if (!stats) {
    throw new Error("Something went wrong!")
  }

  return (
    <Fragment>
      <StaticMetadata />

      <div className="max-w-6xl mx-auto py-12 px-4 sm:px-6 mt-12">
        {/* Acheivement Header */}
        <div className="mb-8">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center">
            <AwardIcon absoluteStrokeWidth className="text-black" size={24} />
          </div>
          <h1 className="text-3xl font-bold text-white text-center mb-4">Achievements</h1>
          <p className="text-gray-400 text-center mb-6 flex flex-col">
            <span>A chronological record of my academic, professional, and personal achievements.</span>
            <span>This page will be continuously updated as I progress through my journey.</span>
          </p>
        </div>

        {/* Achievement Metadata */}
        <TerminalWindow title="metadata" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-2">$ cat metadata.json</div>
            <div className="text-gray-300 space-y-1">
              <div className="flex items-center gap-2">
                <Calendar size={16} className="text-orange-400" />
                <span>Published: {formatDate(stats.createdAt)}</span>
              </div>
              {stats.createdAt !== stats.updatedAt && (
                <div className="flex items-center gap-2">
                  <Clock size={16} className="text-orange-400" />
                  <span>Updated: {formatDate(stats.updatedAt)}</span>
                </div>
              )}
              <div className="flex items-center gap-2">
                <EyeIcon size={16} className="text-orange-400" />
                <span>Views: {formatViews(stats.views)}</span>
              </div>
            </div>
          </div>
        </TerminalWindow>

        {/* Blog Content */}
        <TerminalWindow title="achievements">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-4">$ cat achievements.md</div>
            <MarkdownRenderer content={markdownContent} />
          </div>
        </TerminalWindow>

        <Footer className="mt-12" leading={
          <div>
            Share this post:{" "}
            <a
              href={`https://twitter.com/intent/tweet?url=${encodeURIComponent(
                `https://${me.links["jelius.dev"].link}/achievements`,
              )}&text=Achievements`}
              target="_blank"
              rel="noopener noreferrer"
              className="text-orange-400 hover:underline"
            >
              X
            </a>
          </div>

        } />
      </div>

    </Fragment>
  )
}
