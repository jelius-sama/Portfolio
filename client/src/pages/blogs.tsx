import { useState, useEffect, useRef, useMemo, Fragment } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { BlogCard } from "@/components/layout/blog-card"
import type { Blog } from "@/types/blog"
import { Rss } from "lucide-react"
import { Footer } from "@/components/ui/footer"
import { StaticMetadata } from "@/contexts/metadata"

// Hardcoded blog data
const ALL_BLOGS: Blog[] = [
  {
    id: "chat-system",
    title: "Building a Real-Time Chat System from Scratch",
    summary: "A deep dive into TCP/IP networking, covering server and client implementation.",
    markdownUrl: "/blogs/chat-system.md",
    createdAt: "2025-08-04T10:00:00Z",
    updatedAt: "2025-08-04T10:00:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
  {
    id: "udp-networking",
    title: "Understanding UDP Networking: Building a Simple Server-Client Communication System",
    summary: "Exploring UDP for lightweight, connectionless communication.",
    markdownUrl: "/blogs/udp-networking.md",
    createdAt: "2025-07-08T14:30:00Z",
    updatedAt: "2025-07-08T14:30:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
  {
    id: "ping-utility",
    title: "Building a Ping Utility from Scratch",
    summary: "Recreating the classic network diagnostic tool using raw sockets.",
    markdownUrl: "/blogs/ping-utility.md",
    createdAt: "2025-05-28T09:15:00Z",
    updatedAt: "2025-05-28T09:15:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
  {
    id: "http-server",
    title: "Building a Simple HTTP Server from Scratch",
    summary: "Understanding the basics of HTTP and how to serve web content.",
    markdownUrl: "/blogs/http-server.md",
    createdAt: "2025-05-25T11:00:00Z",
    updatedAt: "2025-05-25T11:00:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
  {
    id: "docker-basics",
    title: "Docker for Beginners: Containerizing Your Applications",
    summary: "A gentle introduction to Docker concepts and basic commands.",
    markdownUrl: "/blogs/docker-basics.md",
    createdAt: "2025-04-10T16:00:00Z",
    updatedAt: "2025-04-10T16:00:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
  {
    id: "react-hooks",
    title: "Mastering React Hooks: useState, useEffect, and Beyond",
    summary: "Deep dive into essential React Hooks for better state management.",
    markdownUrl: "/blogs/react-hooks.md",
    createdAt: "2025-03-15T08:45:00Z",
    updatedAt: "2025-03-15T08:45:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
  {
    id: "go-concurrency",
    title: "Go Concurrency Patterns: Goroutines and Channels",
    summary: "Understanding Go's powerful concurrency model for efficient programming.",
    markdownUrl: "/blogs/go-concurrency.md",
    createdAt: "2025-02-20T13:00:00Z",
    updatedAt: "2025-02-20T13:00:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
  {
    id: "tailwind-css",
    title: "Styling with Tailwind CSS: A Utility-First Approach",
    summary: "How to rapidly build modern UIs with Tailwind CSS.",
    markdownUrl: "/blogs/tailwind-css.md",
    createdAt: "2025-01-05T10:00:00Z",
    updatedAt: "2025-01-05T10:00:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
  {
    id: "database-design",
    title: "Relational Database Design Best Practices",
    summary: "Tips and tricks for designing robust and scalable databases.",
    markdownUrl: "/blogs/database-design.md",
    createdAt: "2024-12-10T14:00:00Z",
    updatedAt: "2024-12-10T14:00:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
  {
    id: "api-security",
    title: "Securing Your APIs: A Guide to Best Practices",
    summary: "Essential steps to protect your API endpoints from common vulnerabilities.",
    markdownUrl: "/blogs/api-security.md",
    createdAt: "2024-11-22T09:00:00Z",
    updatedAt: "2024-11-22T09:00:00Z",
    prequelId: null,
    sequelId: null,
    parts: [],
  },
]

const BLOGS_PER_PAGE = 5

// TODO: Implement search feature
export default function BlogListPage() {
  const [sortOrder, setSortOrder] = useState<"newest" | "oldest" | "updated">("newest")
  const [visibleBlogsCount, setVisibleBlogsCount] = useState(BLOGS_PER_PAGE)
  const loadMoreRef = useRef<HTMLDivElement>(null)

  const sortedBlogs = useMemo(() => {
    const sorted = [...ALL_BLOGS]
    if (sortOrder === "newest") {
      sorted.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    } else if (sortOrder === "oldest") {
      sorted.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime())
    } else if (sortOrder === "updated") {
      sorted.sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
    }
    return sorted
  }, [sortOrder])

  const displayedBlogs = useMemo(() => {
    return sortedBlogs.slice(0, visibleBlogsCount)
  }, [sortedBlogs, visibleBlogsCount])

  const hasMoreBlogs = visibleBlogsCount < ALL_BLOGS.length

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMoreBlogs) {
          setVisibleBlogsCount((prevCount) => prevCount + BLOGS_PER_PAGE)
        }
      },
      { threshold: 1.0 },
    )

    if (loadMoreRef.current) {
      observer.observe(loadMoreRef.current)
    }

    return () => {
      if (loadMoreRef.current) {
        observer.unobserve(loadMoreRef.current)
      }
    }
  }, [hasMoreBlogs])

  return (
    <Fragment>
      <StaticMetadata />
      <div className="max-w-6xl mx-auto py-12 px-6 mt-12">
        {/* Page Title */}
        <div className="text-center mb-8">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center">
            <Rss className="text-black" size={24} />
          </div>
          <h1 className="text-3xl font-bold text-white mb-2">My Blog Posts</h1>
          <p className="text-gray-400 font-mono">Explore my thoughts on development and technology</p>
        </div>

        {/* Terminal Info */}
        <TerminalWindow title="blog-info" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-2">$ ls /var/www/blogs/posts/</div>
            <div className="text-gray-300 mb-4">Found {ALL_BLOGS.length} articles.</div>
            <div className="text-orange-400 mb-2">$ cat /etc/blog-config.conf</div>
            <div className="text-gray-300">
              <div>Sort by: {sortOrder}</div>
              <div>
                Displaying: {displayedBlogs.length} of {ALL_BLOGS.length}
              </div>
            </div>
          </div>
        </TerminalWindow>

        {/* Sorting Options */}
        <div className="flex justify-end mb-6">
          <div className="relative">
            <select
              value={sortOrder}
              onChange={(e) => setSortOrder(e.target.value as "newest" | "oldest" | "updated")}
              className="appearance-none bg-gray-800/50 border border-orange-500/30 text-white px-4 py-2 rounded-lg pr-8 focus:outline-none focus:border-orange-500 font-mono text-sm"
            >
              <option value="newest">Sort by: Newest</option>
              <option value="oldest">Sort by: Oldest</option>
              <option value="updated">Sort by: Last Updated</option>
            </select>
            <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-400">
              <svg className="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                <path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z" />
              </svg>
            </div>
          </div>
        </div>

        {/* Blog List */}
        <TerminalWindow title="blog-posts">
          <div className="space-y-6">
            {displayedBlogs.map((blog) => (
              <BlogCard key={blog.id} blog={blog} />
            ))}
            {hasMoreBlogs && (
              <div ref={loadMoreRef} className="text-center text-gray-500 font-mono py-4">
                Loading more posts...
              </div>
            )}
            {!hasMoreBlogs && displayedBlogs.length > 0 && (
              <div className="text-center text-gray-500 font-mono py-4">End of posts.</div>
            )}
            {displayedBlogs.length === 0 && (
              <div className="text-center text-gray-500 font-mono py-4">No blog posts found.</div>
            )}
          </div>
        </TerminalWindow>

        <Footer className="mt-12" leading={
          <div>
            Read more on my RSS feed:{" "}
            <a href="/rss.xml" className="text-orange-400 hover:underline">
              /rss.xml
            </a>
          </div>
        }
        />
      </div>
    </Fragment>
  )
}
