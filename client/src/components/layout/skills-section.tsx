import { SkillBadge, type ValidLevel } from "@/components/layout/skill-badge"
import { useConfig } from "@/contexts/config"

export function SkillsSection() {
    const { app: { portfolio: me } } = useConfig()

    return (
        <section id="skills" className="bg-gray-900/50">
            <div className="max-w-6xl mx-auto py-20 px-4 sm:px-6 ">
                <h2 className="text-3xl font-bold text-center mb-12">
                    <span className="text-orange-400">|</span> <span className="text-white">Skills & Technologies</span>
                </h2>

                <div className="text-gray-400 text-center mb-8 font-mono">Technologies and tools I work with:</div>

                <div className="grid md:grid-cols-3 gap-8">
                    <div>
                        <h3 className="text-orange-400 font-bold mb-4 font-mono">Languages</h3>
                        <div className="space-y-3">
                            {me.skills.languages.map((item) => (
                                <SkillBadge showLevel={false} key={item.skill} skill={item.skill} level={item.level as ValidLevel} />
                            ))}
                        </div>
                    </div>

                    <div>
                        <h3 className="text-orange-400 font-bold mb-4 font-mono">Frameworks</h3>
                        <div className="space-y-3">
                            {me.skills.frameworks.map((item) => (
                                <SkillBadge showLevel={false} key={item.skill} skill={item.skill} level={item.level as ValidLevel} />
                            ))}
                        </div>
                    </div>

                    <div>
                        <h3 className="text-orange-400 font-bold mb-4 font-mono">Others</h3>
                        <div className="space-y-3">
                            {me.skills.others.map((item) => (
                                <SkillBadge showLevel={false} key={item.skill} skill={item.skill} level={item.level as ValidLevel} />
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </section>
    )
}
