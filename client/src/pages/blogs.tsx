import { useState, useEffect, useRef, Fragment } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { BlogCard } from "@/components/layout/blog-card"
import type { Blog } from "@/types/blog"
import { Rss } from "lucide-react"
import { Footer } from "@/components/ui/footer"
import { StaticMetadata } from "@/contexts/metadata"
import { Select } from "@/components/ui/select"

const BLOGS_PER_PAGE = 5
const IS_RSS_IMPLEMENTED = false

interface BlogsResponse {
  blogs: Blog[] | null
  total: number
  page: number
  totalPages: number
}

// TODO: Use tanstack query instead of fetch for client side caching
export default function BlogListPage() {
  const [sortOrder, setSortOrder] = useState<"newest" | "oldest" | "updated">("newest")
  const [blogs, setBlogs] = useState<Blog[]>([])
  const [currentPage, setCurrentPage] = useState(1)
  const [totalBlogs, setTotalBlogs] = useState(0)
  const [totalPages, setTotalPages] = useState(0)
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const loadMoreRef = useRef<HTMLDivElement>(null)

  const sortOptions = [
    { value: "newest", label: "Sort by: Newest" },
    { value: "oldest", label: "Sort by: Oldest" },
    { value: "updated", label: "Sort by: Last Updated" },
  ]

  const fetchBlogs = async (page: number, isNewSort = false) => {
    try {
      if (page === 1) {
        setLoading(true)
      } else {
        setLoadingMore(true)
      }
      setError(null)

      let sortBy = "created_at"
      let order = "desc"

      if (sortOrder === "oldest") {
        order = "asc"
      } else if (sortOrder === "updated") {
        sortBy = "updated_at"
        order = "desc"
      }

      const res = await fetch(`/api/blogs?page=${page}&limit=${BLOGS_PER_PAGE}&sortBy=${sortBy}&order=${order}`)

      if (!res.ok) {
        throw new Error("Failed to fetch blogs!")
      }

      const data: BlogsResponse = await res.json()

      if (page === 1 || isNewSort) {
        // First page or new sort - replace all blogs
        setBlogs(data.blogs || [])
      } else {
        // Subsequent pages - append to existing blogs
        setBlogs(prev => [...prev, ...(data.blogs ?? [])])
      }

      setTotalBlogs(data.total)
      setTotalPages(data.totalPages)
      setCurrentPage(page)

    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred")
    } finally {
      setLoading(false)
      setLoadingMore(false)
    }
  }

  // Initial load and sort changes
  useEffect(() => {
    setCurrentPage(1)
    fetchBlogs(1, true)
  }, [sortOrder])

  // Infinite scroll observer
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && currentPage < totalPages && !loading && !loadingMore) {
          fetchBlogs(currentPage + 1)
        }
      },
      { threshold: 1.0 }
    )

    if (loadMoreRef.current) {
      observer.observe(loadMoreRef.current)
    }

    return () => {
      if (loadMoreRef.current) {
        observer.unobserve(loadMoreRef.current)
      }
    }
  }, [currentPage, totalPages, loading, loadingMore])

  const hasMoreBlogs = currentPage < totalPages

  return (
    <Fragment>
      <StaticMetadata />
      <section className="max-w-6xl mx-auto py-12 mt-12 px-4 sm:px-6">
        {/* Header */}
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
            <div className="text-gray-300 mb-4">
              {loading && blogs.length === 0 ? "Loading articles..." : `Found ${totalBlogs} articles.`}
            </div>
            <div className="text-orange-400 mb-2">$ cat /etc/blog-config.conf</div>
            <div className="text-gray-300">
              <div>Sort by: {sortOrder}</div>
              <div>
                Displaying: {blogs.length} of {totalBlogs}
              </div>
              {totalPages > 1 && (
                <div>
                  Loaded pages: {currentPage} of {totalPages}
                </div>
              )}
            </div>
          </div>
        </TerminalWindow>

        {/* Sorting Options */}
        <div className="flex justify-end mb-6">
          <Select
            value={sortOrder}
            onChange={(value) => setSortOrder(value as "newest" | "oldest" | "updated")}
            options={sortOptions}
            disabled={loading}
            className="w-64"
          />
        </div>

        {/* Blog List */}
        <TerminalWindow title="blog-posts">
          <div className="space-y-6">
            {error && (
              <div className="text-center text-red-400 font-mono py-4">
                Error: {error}
              </div>
            )}

            {blogs.map((blog) => (
              <BlogCard key={blog.id} blog={blog} />
            ))}

            {hasMoreBlogs && (
              <div ref={loadMoreRef} className="text-center text-gray-500 font-mono py-4">
                {loadingMore ? "Loading more posts..." : "Scroll down for more..."}
              </div>
            )}

            {!hasMoreBlogs && blogs.length > 0 && (
              <div className="text-center text-gray-500 font-mono py-4">End of posts.</div>
            )}

            {blogs.length === 0 && !loading && !error && (
              <div className="text-center text-gray-500 font-mono py-4">No blog posts found.</div>
            )}

            {loading && blogs.length === 0 && (
              <div className="text-center text-gray-500 font-mono py-4">Loading posts...</div>
            )}
          </div>
        </TerminalWindow>

        <Footer
          className="mt-12"
          leading={IS_RSS_IMPLEMENTED && (
            <div>
              Read more on my RSS feed:{" "}
              <a href="/rss.xml" className="text-orange-400 hover:underline">
                /rss.xml
              </a>
            </div>
          )}
        />
      </section>
    </Fragment>
  )
}
