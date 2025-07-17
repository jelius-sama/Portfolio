import type React from "react"
import { ExternalLink } from "lucide-react"

interface LinkCardProps {
    title: string
    description?: string
    url: string
    icon?: React.ReactNode
}

export function LinkCard({ title, description, url, icon }: LinkCardProps) {
    return (
        <a
            href={url}
            target="_blank"
            rel="noopener noreferrer"
            className="block w-full p-4 bg-gray-800/50 border border-orange-500/30 rounded-lg hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group"
        >
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                    {icon && <div className="text-orange-400 group-hover:text-orange-300">{icon}</div>}
                    <div>
                        <div className="text-white font-medium font-mono">{title}</div>
                        {description && <div className="text-gray-400 text-sm">{description}</div>}
                    </div>
                </div>
                <ExternalLink className="text-gray-500 group-hover:text-orange-400 transition-colors" size={18} />
            </div>
        </a>
    )
}
