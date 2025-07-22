import { TerminalWindow } from "@/components/ui/terminal-window"
import { useConfig } from "@/contexts/config"

export interface ExperienceItem {
    title: string;
    subtitle: string;
    date: string;
    highlighted?: true;
    details?: string[];
}

export function ExperienceSection() {
    const { app: { portfolio: me } } = useConfig()

    return (
        <section id="experience" className="max-w-6xl mx-auto py-20 px-4 sm:px-6">
            <h2 className="text-3xl font-bold text-center mb-12">
                <span className="text-orange-400">|</span> <span className="text-white">Experience</span>
            </h2>

            <div className="text-gray-400 text-center mb-8 font-mono">My professional journey:</div>

            <TerminalWindow title="experience.log" className="mb-8">
                <div className="font-mono text-sm space-y-6">
                    {(me.experience as ExperienceItem[]).map((item, index) => (
                        <div
                            key={index}
                            className={item.highlighted ? "border-l-2 border-orange-500/30 pl-4" : ""}
                        >
                            <div className="text-orange-400 font-bold mb-2">{item.title}</div>
                            <div className="text-gray-400 mb-1">{item.subtitle}</div>
                            <div className="text-gray-500 mb-3">{item.date}</div>
                            {item.details && (
                                <div className="text-gray-300">
                                    {item.details.map((point, idx) => (
                                        <div key={idx}>• {point}</div>
                                    ))}
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            </TerminalWindow>
        </section>
    )
}
