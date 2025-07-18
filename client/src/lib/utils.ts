import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { useState, useEffect } from "react"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

export function timeAgo(dateString: string): string {
    const then = new Date(dateString)
    const now = new Date()
    const diff = Math.floor((now.getTime() - then.getTime()) / 1000) // in seconds

    if (diff < 60) return 'just now'
    if (diff < 3600) return `${Math.floor(diff / 60)} minutes ago`
    if (diff < 86400) return `${Math.floor(diff / 3600)} hours ago`
    return `${Math.floor(diff / 86400)} days ago`
}

export function useTimeAgo(dateString?: string): string | null {
    const [value, setValue] = useState<string | null>(() =>
        dateString ? timeAgo(dateString) : null
    )

    useEffect(() => {
        if (!dateString) return

        setValue(timeAgo(dateString))

        const interval = setInterval(() => {
            setValue(timeAgo(dateString))
        }, 60_000)

        return () => clearInterval(interval)
    }, [dateString])

    return value
}
