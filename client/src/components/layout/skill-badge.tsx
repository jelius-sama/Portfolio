export type ValidLevel = "Beginner" | "Intermediate" | "Experienced" | "Expert"

interface SkillBadgeProps {
    skill: string
    level?: ValidLevel
}

export function SkillBadge({ skill, level = "Experienced" }: SkillBadgeProps) {
    const getLevelColor = (level: string) => {
        switch (level) {
            case "Expert":
                return "text-green-400"
            case "Experienced":
                return "text-orange-400"
            case "Intermediate":
                return "text-yellow-400"
            case "Beginner":
                return "text-blue-400"
            default:
                return "text-orange-400"
        }
    }

    return (
        <div className="flex items-center gap-3 p-3 bg-gray-800/50 rounded border border-orange-500/20 hover:border-orange-500/40 transition-colors">
            <div className="w-2 h-2 rounded-full bg-orange-400"></div>
            <div className="flex-1">
                <div className="text-white font-medium">{skill}</div>
                <div className={`text-sm ${getLevelColor(level)}`}>{level}</div>
            </div>
        </div>
    )
}
