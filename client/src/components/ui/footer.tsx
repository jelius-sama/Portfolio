import { useConfig } from "@/contexts/config"
import { type ReactNode } from "react"
import { cn } from "@/lib/utils"

export function Footer({ leading, trailing, className }: { leading?: ReactNode, trailing?: ReactNode, className?: string }) {
    const { app: { portfolio: me } } = useConfig()

    return (
        <div className={cn("text-center text-gray-500 text-sm font-mono", className)}>
            {leading}
            Copyright © {new Date().getFullYear()} {me.first_name + " " + me.last_name}. All Rights Reserved.
            {trailing}
        </div>

    )
}
