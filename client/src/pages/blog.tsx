import { useState, useEffect, Fragment } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { Calendar, Clock, ChevronLeft, ChevronRight, FileText } from "lucide-react"
import type { Blog } from "@/types/blog"
import { useParams, Link } from "react-router-dom"
import { Footer } from "@/components/ui/footer"

// TODO: Integrate with server to do SSR
export default function BlogPostPage() {
  const [blog, setBlog] = useState<Blog | null>(null)
  const [markdownContent, setMarkdownContent] = useState<string>("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { id } = useParams();

  useEffect(() => {
    const fetchBlogData = async () => {
      try {
        setLoading(true)
        setError(null)

        // First, fetch the blog details
        const blogResponse = await fetch(`/api/blog?id=${id}`)
        if (!blogResponse.ok) {
          throw new Error(`Failed to fetch blog: ${blogResponse.status}`)
        }
        const blogData: Blog = await blogResponse.json()
        setBlog(blogData)

        // Then, fetch the markdown file
        const markdownResponse = await fetch(`/api/blog_file?id=${blogData.id}`)
        if (!markdownResponse.ok) {
          throw new Error(`Failed to fetch markdown: ${markdownResponse.status}`)
        }
        const markdownText = await markdownResponse.text()
        setMarkdownContent(markdownText)
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred")
      } finally {
        setLoading(false)
      }
    }

    fetchBlogData()
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

  // Simple markdown renderer (you might want to use a proper markdown library)
  const renderMarkdown = (content: string) => {
    return content.split("\n").map((line, index) => {
      // Headers
      if (line.startsWith("# ")) {
        return (
          <h1 key={index} className="text-2xl font-bold text-orange-400 mb-4 mt-6">
            {line.slice(2)}
          </h1>
        )
      }
      if (line.startsWith("## ")) {
        return (
          <h2 key={index} className="text-xl font-bold text-orange-400 mb-3 mt-5">
            {line.slice(3)}
          </h2>
        )
      }
      if (line.startsWith("### ")) {
        return (
          <h3 key={index} className="text-lg font-bold text-orange-400 mb-2 mt-4">
            {line.slice(4)}
          </h3>
        )
      }

      // Code blocks
      if (line.startsWith("```")) {
        return (
          <div key={index} className="bg-gray-800 border border-orange-500/30 rounded p-4 my-4 font-mono text-sm">
            <div className="text-orange-400 mb-2">$ code</div>
          </div>
        )
      }

      // Inline code
      if (line.includes("`")) {
        const parts = line.split("`")
        return (
          <p key={index} className="text-gray-300 mb-3 leading-relaxed">
            {parts.map((part, i) =>
              i % 2 === 0 ? (
                part
              ) : (
                <code key={i} className="bg-gray-800 px-2 py-1 rounded text-orange-400 font-mono text-sm">
                  {part}
                </code>
              ),
            )}
          </p>
        )
      }

      // Lists
      if (line.startsWith("- ")) {
        return (
          <li key={index} className="text-gray-300 mb-1 ml-4">
            {line.slice(2)}
          </li>
        )
      }

      // Empty lines
      if (line.trim() === "") {
        return <br key={index} />
      }

      // Regular paragraphs
      return (
        <p key={index} className="text-gray-300 mb-3 leading-relaxed">
          {line}
        </p>
      )
    })
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-950 text-white">
        <header className="fixed top-0 left-0 right-0 z-50 bg-gray-900/95 backdrop-blur-sm border-b border-orange-500/30">
          <div className="max-w-4xl mx-auto px-6 py-4">
            <div className="flex items-center justify-between">
              <div className="text-orange-400 font-bold text-xl font-mono">
                {"> "}
                <span className="text-white">blog/{id}</span>
              </div>
              <Link to="/blogs" className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm">
                ← Back to Blogs
              </Link>
            </div>
          </div>
        </header>

        <main className="pt-20 pb-12 px-6">
          <div className="max-w-4xl mx-auto">
            <TerminalWindow title="loading">
              <div className="font-mono text-sm text-center py-8">
                <div className="text-orange-400 mb-2">$ curl -X GET /api/blog?id={id}</div>
                <div className="text-gray-300 mb-4">Fetching blog data...</div>
                <div className="text-orange-400 mb-2">$ curl -X GET /api/blog_file?id={id}</div>
                <div className="text-gray-300">Loading markdown content...</div>
              </div>
            </TerminalWindow>
          </div>
        </main>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-950 text-white">
        <header className="fixed top-0 left-0 right-0 z-50 bg-gray-900/95 backdrop-blur-sm border-b border-orange-500/30">
          <div className="max-w-4xl mx-auto px-6 py-4">
            <div className="flex items-center justify-between">
              <div className="text-orange-400 font-bold text-xl font-mono">
                {"> "}
                <span className="text-white">error</span>
              </div>
              <Link to="/blogs" className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm">
                ← Back to Blogs
              </Link>
            </div>
          </div>
        </header>

        <main className="pt-20 pb-12 px-6">
          <div className="max-w-4xl mx-auto">
            <TerminalWindow title="error">
              <div className="font-mono text-sm">
                <div className="text-orange-400 mb-2">$ curl -X GET /api/blog?id={id}</div>
                <div className="text-red-400 mb-4">Error: {error}</div>
                <div className="text-orange-400 mb-2">$ echo "Please try again later"</div>
                <div className="text-gray-300">Please try again later</div>
              </div>
            </TerminalWindow>
          </div>
        </main>
      </div>
    )
  }

  if (!blog) {
    return null
  }

  return (
    <Fragment>
      <div className="max-w-6xl mx-auto py-12 px-6 mt-12">
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
            <div className="prose prose-invert max-w-none">{renderMarkdown(markdownContent)}</div>
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
