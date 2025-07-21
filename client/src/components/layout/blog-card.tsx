import { ArrowRight } from "lucide-react"
import type { Blog } from "@/types/blog"

interface BlogCardProps {
    blog: Blog
}

export function BlogCard({ blog }: BlogCardProps) {
    const formatDate = (dateString: string) => {
        const options: Intl.DateTimeFormatOptions = { year: "numeric", month: "long", day: "numeric" }
        return new Date(dateString).toLocaleDateString(undefined, options)
    }

    return (
        <a
            href={`/blogs/${blog.id}`
            }
            className="block p-4 bg-gray-800/50 border border-orange-500/30 rounded-lg hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
        >
            <div className="flex items-center justify-between mb-2" >
                <h3 className="text-white font-medium text-lg font-mono group-hover:text-orange-400 transition-colors" >
                    {blog.title}
                </h3>
                < ArrowRight className="text-gray-500 group-hover:text-orange-400 transition-colors" size={20} />
            </div>
            < p className="text-gray-400 text-sm mb-3 line-clamp-2" > {blog.summary} </p>
            < div className="flex justify-between text-xs text-gray-500 font-mono" >
                <span>Published: {formatDate(blog.createdAt)} </span>
                {blog.createdAt !== blog.updatedAt && <span>Updated: {formatDate(blog.updatedAt)} </span>}
            </div>
        </a>
    )
}
