import { cn } from "@/lib/utils"

export type ValidLevel = "Beginner" | "Intermediate" | "Experienced" | "Expert"

interface SkillBadgeProps {
    skill: string
    level?: ValidLevel
    showLevel?: boolean
}

const levelStyles = {
    Expert: {
        dot: "bg-green-400",
        text: "text-green-400",
        border: "border-green-500/20 hover:border-green-500/40",
    },
    Experienced: {
        dot: "bg-orange-400",
        text: "text-orange-400",
        border: "border-orange-500/20 hover:border-orange-500/40",
    },
    Intermediate: {
        dot: "bg-yellow-400",
        text: "text-yellow-400",
        border: "border-yellow-500/20 hover:border-yellow-500/40",
    },
    Beginner: {
        dot: "bg-blue-400",
        text: "text-blue-400",
        border: "border-blue-500/20 hover:border-blue-500/40",
    },
} as const

export function SkillBadge({ skill, level = "Experienced", showLevel = true }: SkillBadgeProps) {
    const style = levelStyles[level]

    return (
        <div
            className={cn(
                "flex items-center gap-3 p-3 bg-gray-800/50 rounded border transition-colors",
                showLevel
                    ? "border-orange-500/20 hover:border-orange-500/40"
                    : style.border
            )}
        >
            <div className={cn("w-2 h-2 rounded-full", style.dot)} />
            <div className="flex-1">
                <div className="text-white font-medium">{skill}</div>
                {showLevel && <div className={cn("text-sm", style.text)}>{level}</div>}
            </div>
        </div>
    )
}
