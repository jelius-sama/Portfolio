import type React from "react"

interface TerminalWindowProps {
    title?: string
    children: React.ReactNode
    className?: string
}

export function TerminalWindow({ title = "terminal", children, className = "" }: TerminalWindowProps) {
    return (
        <div className={`bg-gray-900 rounded-lg border border-orange-500/30 overflow-hidden ${className}`}>
            <div className="flex items-center gap-2 px-4 py-3 bg-gray-800 border-b border-orange-500/30">
                <div className="flex gap-2">
                    <div className="w-3 h-3 rounded-full bg-red-500"></div>
                    <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                    <div className="w-3 h-3 rounded-full bg-green-500"></div>
                </div>
                <span className="text-gray-400 text-sm ml-2">{title}</span>
            </div>
            <div className="p-6">{children}</div>
        </div>
    )
}
