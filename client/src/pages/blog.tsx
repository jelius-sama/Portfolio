import { useState, useEffect, Fragment } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { Calendar, Clock, ChevronLeft, ChevronRight, FileText, EyeIcon } from "lucide-react"
import type { Blog } from "@/types/blog"
import { useParams, Link, useLocation } from "react-router-dom"
import { Footer } from "@/components/ui/footer"
import { MarkdownRenderer } from "@/components/ui/markdown"
import { useConfig } from "@/contexts/config"
import { DynamicMetadata, PathBasedMetadata } from "@/contexts/metadata"
import { type StaticRoute } from "@/types/static.route"
import { formatViews, formatDate } from "@/lib/utils"

export function generateBlogMetadata(blog: Blog, fullPath: string): StaticRoute {
  const title = `${blog.title} | Jelius`;

  return {
    path: fullPath,
    title,
    meta: [
      { name: "description", content: blog.summary },
      { name: "application-name", content: "Jelius Basumatary" },
      { name: "robots", content: "index, follow" },
      { name: "format-detection", content: "telephone=no" },
      { name: "apple-mobile-web-app-capable", content: "yes" },
      { name: "apple-mobile-web-app-title", content: "Jelius Basumatary" },
      { name: "theme-color", content: "#000b11" },
      { name: "apple-mobile-web-app-status-bar-style", content: "default" },

      { property: "og:title", content: title },
      { property: "og:description", content: blog.summary },
      { property: "og:url", content: "https://jelius.dev" + fullPath },
      { property: "og:site_name", content: "Jelius Basumatary" },
      { property: "og:image", content: "/assets/jelius.jpg" },
      { property: "og:type", content: "article" },

      { name: "twitter:card", content: "summary" },
      { name: "twitter:site", content: "@jelius_sama" },
      { name: "twitter:creator", content: "@jelius_sama" },
      { name: "twitter:title", content: title },
      { name: "twitter:description", content: blog.summary },
      { name: "twitter:image", content: "/assets/jelius.jpg" },
    ],
    link: [
      { rel: "canonical", href: "https://jelius.dev" + fullPath }
    ]
  };
}

export default function BlogPostPage() {
  const { ssrData } = useConfig()
  const { pathname } = useLocation()

  const [blog, setBlog] = useState<Blog | null>(null)
  const [markdownContent, setMarkdownContent] = useState<string>("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { id } = useParams();

  useEffect(() => {
    (async () => {
      setLoading(true);
      setError(null);

      try {
        // Use SSR data if available and matches current path
        if (ssrData?.path === pathname) {
          if ("status" in ssrData && ssrData.status === 404) {
            throw new Error(`Blog with ID of "${id}" was not found!`)
          }
          setBlog(ssrData.api_resp);

          const markdownRes = await fetch(`/api/blog_file?id=${id}`);
          if (!markdownRes.ok) {
            throw new Error(`Failed to fetch markdown: ${markdownRes.status}`);
          }

          const markdownText = await markdownRes.text();
          setMarkdownContent(markdownText);
        } else {
          // Fetch blog and markdown in parallel
          const [blogRes, markdownRes] = await Promise.all([
            fetch(`/api/blog?id=${id}`),
            fetch(`/api/blog_file?id=${id}`)
          ]);

          if (!blogRes.ok) {
            if (blogRes.status === 404) {
              throw new Error(`Blog with ID of "${id}" was not found!`)
            } else {
              throw new Error(`Failed to fetch blog: ${blogRes.status}`);
            }
          }

          if (!markdownRes.ok) {
            throw new Error(`Failed to fetch markdown: ${markdownRes.status}`);
          }

          const blogData: Blog = await blogRes.json();
          const markdownText = await markdownRes.text();

          setBlog(blogData);
          setMarkdownContent(markdownText);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "An unexpected error occurred");
      } finally {
        setLoading(false);
      }
    })();
  }, [id, pathname]);

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
        <PathBasedMetadata paths={["*"]} />

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
      </Fragment>
    )
  }

  if (error) {
    return (
      <Fragment>
        <PathBasedMetadata paths={["*", "#not_found"]} />

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
      </Fragment>
    )
  }

  if (!blog) {
    return null
  }

  return (
    <Fragment>
      <DynamicMetadata currentMeta={pathname === ssrData?.path ? ssrData.metadata : generateBlogMetadata(blog, pathname)} />
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
              <div className="flex items-center gap-2">
                <EyeIcon size={16} className="text-orange-400" />
                <span>Views: {formatViews(blog.views)}</span>
              </div>

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
