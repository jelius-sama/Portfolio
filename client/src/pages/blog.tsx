import { useState, useEffect, Fragment } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { Calendar, Clock, ChevronLeft, ChevronRight, FileText } from "lucide-react"
import type { Blog } from "@/types/blog"
import { useParams, Link } from "react-router-dom"
import { Footer } from "@/components/ui/footer"
import { MarkdownRenderer } from "@/components/ui/markdown"

// TODO: Integrate with server to do SSR
export default function BlogPostPage() {
  const [blog, setBlog] = useState<Blog | null>(null)
  const [markdownContent, setMarkdownContent] = useState<string>("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { id } = useParams();

  useEffect(() => {
    (async () => {
      try {
        setLoading(true)
        setError(null)

        const [blogRes, markdownRes] = await Promise.all([
          fetch(`/api/blog?id=${id}`),
          fetch(`/api/blog_file?id=${id}`)
        ])

        if (!blogRes.ok) {
          throw new Error(`Failed to fetch blog: ${blogRes.status}`)
        }

        if (!markdownRes.ok) {
          throw new Error(`Failed to fetch markdown: ${markdownRes.status}`)
        }

        const blogData: Blog = await blogRes.json()
        const markdownText = await markdownRes.text()

        setBlog(blogData)
        setMarkdownContent(markdownText)
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred")
      } finally {
        setLoading(false)
      }
    })()
  }, [id])

  const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }
    return new Date(dateString).toLocaleDateString(undefined, options)
  }

  if (loading) {
    return (
      <section className="max-w-6xl mx-auto pt-20 pb-12 px-4 sm:px-6">
        <TerminalWindow title="loading">
          <div className="font-mono text-sm text-center py-8">
            <div className="text-orange-400 mb-2">$ curl -X GET /api/blog?id={id}</div>
            <div className="text-gray-300 mb-4">Fetching blog data...</div>
            <div className="text-orange-400 mb-2">$ curl -X GET /api/blog_file?id={id}</div>
            <div className="text-gray-300">Loading markdown content...</div>
          </div>
        </TerminalWindow>
      </section>
    )
  }

  if (error) {
    return (
      <section className="max-w-6xl mx-auto pt-20 pb-12 px-4 sm:px-6">
        <TerminalWindow title="error">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-2">$ curl -X GET /api/blog?id={id}</div>
            <div className="text-red-400 mb-4">Error: {error}</div>
            <div className="text-orange-400 mb-2">$ echo "Please try again later"</div>
            <div className="text-gray-300">Please try again later</div>
          </div>
        </TerminalWindow>
      </section>
    )
  }

  if (!blog) {
    return null
  }

  return (
    <Fragment>
      <div className="max-w-6xl mx-auto py-12 px-4 sm:px-6 mt-12">
        {/* Blog Header */}
        <div className="mb-8">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center">
            <FileText className="text-black" size={24} />
          </div>
          <h1 className="text-3xl font-bold text-white text-center mb-4">{blog.title}</h1>
          <p className="text-gray-400 text-center mb-6">{blog.summary}</p>
        </div>

        {/* Blog Meta */}
        <TerminalWindow title="blog-meta" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-2">$ cat blog-metadata.json</div>
            <div className="text-gray-300 space-y-1">
              <div className="flex items-center gap-2">
                <Calendar size={16} className="text-orange-400" />
                <span>Published: {formatDate(blog.createdAt)}</span>
              </div>
              {blog.createdAt !== blog.updatedAt && (
                <div className="flex items-center gap-2">
                  <Clock size={16} className="text-orange-400" />
                  <span>Updated: {formatDate(blog.updatedAt)}</span>
                </div>
              )}
              <div>ID: {blog.id}</div>
              {blog.parts.length > 0 && <div>Series: {blog.parts.length} parts</div>}
            </div>
          </div>
        </TerminalWindow>

        {/* Navigation (Prequel/Sequel) */}
        {(blog.prequelId || blog.sequelId) && (
          <TerminalWindow title="navigation" className="mb-8">
            <div className="font-mono text-sm">
              <div className="text-orange-400 mb-4">$ ls ../related-posts/</div>
              <div className="flex justify-between">
                {blog.prequelId ? (
                  <Link
                    to={`/blog/${blog.prequelId}`}
                    className="flex items-center gap-2 px-4 py-2 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200"
                  >
                    <ChevronLeft className="text-orange-400" size={16} />
                    <span className="text-gray-300">Previous Post</span>
                  </Link>
                ) : (
                  <div></div>
                )}
                {blog.sequelId && (
                  <Link
                    to={`/blog/${blog.sequelId}`}
                    className="flex items-center gap-2 px-4 py-2 bg-gray-800/50 border border-orange-500/30 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200"
                  >
                    <span className="text-gray-300">Next Post</span>
                    <ChevronRight className="text-orange-400" size={16} />
                  </Link>
                )}
              </div>
            </div>
          </TerminalWindow>
        )}

        {/* Blog Content */}
        <TerminalWindow title="blog-content">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-4">$ cat {blog.id}.md</div>
            <MarkdownRenderer content={markdownContent} />
          </div>
        </TerminalWindow>

        {/* Series Navigation */}
        {blog.parts.length > 0 && (
          <TerminalWindow title="series-navigation" className="mt-8">
            <div className="font-mono text-sm">
              <div className="text-orange-400 mb-4">$ ls series-parts/</div>
              <div className="text-gray-300 mb-4">This post is part of a {blog.parts.length}-part series:</div>
              <div className="grid gap-2">
                {blog.parts.map((partId, index) => (
                  <Link
                    key={partId}
                    to={`/blog/${partId}`}
                    className={`block px-3 py-2 rounded border transition-all duration-200 ${partId === blog.id
                      ? "bg-orange-500/20 border-orange-500 text-orange-400"
                      : "bg-gray-800/50 border-orange-500/30 text-gray-300 hover:border-orange-500 hover:bg-gray-800/70"
                      }`}
                  >
                    Part {index + 1}: {partId === blog.id ? blog.title : `Blog ${partId}`}
                  </Link>
                ))}
              </div>
            </div>
          </TerminalWindow>
        )}


        <Footer className=" mt-12" leading={
          <div>
            Share this post:{" "}
            <a
              href={`https://twitter.com/intent/tweet?url=${encodeURIComponent(
                `https://yourname.dev/blog/${blog.id}`,
              )}&text=${encodeURIComponent(blog.title)}`}
              target="_blank"
              rel="noopener noreferrer"
              className="text-orange-400 hover:underline"
            >
              Twitter
            </a>
          </div>

        } />
      </div>
    </Fragment>
  )
}
