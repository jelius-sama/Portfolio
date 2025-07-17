import { TerminalWindow } from "@/components/layout/terminal-window"
import { useConfig } from "@/contexts/config"
import { toast } from "sonner"

export function HeroSection() {
    const { app: { portfolio: me } } = useConfig()

    return (
        <section id="home" className="min-h-screen flex items-center justify-center px-6 pt-20">
            <div className="max-w-4xl mx-auto text-center">
                <div className="mb-8">
                    <div className="w-32 h-32 mx-auto mb-6 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 p-1">
                        <div className="w-full h-full rounded-full bg-gray-800 flex items-center justify-center">
                            {me.pictures.length > 0 ? (
                                <img src={me.pictures[0]} className="rounded-full overflow-hidden" />
                            ) : (
                                <span className="text-orange-400 text-4xl font-mono">{"{ }"}</span>
                            )}
                        </div>
                    </div>
                    <div className="text-gray-400 text-lg mb-2">Hello, I'm</div>
                    <h1 className="text-5xl md:text-6xl font-bold text-white mb-4">
                        <span className="text-orange-400">{me.first_name}</span> {me.last_name}
                    </h1>
                    <p className="text-xl text-gray-300 font-mono">{me.developer} Developer</p>
                </div>

                <div className="flex flex-wrap gap-4 justify-center mb-8">
                    <button onClick={() => toast.info("CV Not Available")} className="px-6 py-3 bg-orange-500 text-black font-medium rounded hover:bg-orange-400 transition-colors">
                        Download CV
                    </button>
                    <a href="#contact" className="px-6 py-3 border border-orange-500 text-orange-400 font-medium rounded hover:bg-orange-500/10 transition-colors">
                        Contact Info
                    </a>
                </div>

                <TerminalWindow title="welcome" className="max-w-2xl mx-auto text-left">
                    <div className="font-mono text-sm">
                        <div className="text-orange-400 mb-2">$ whoami</div>
                        <div className="text-gray-300 mb-4">{me.introduction}</div>
                        <div className="text-orange-400 mb-2">$ cat skills.txt</div>
                        <div className="text-gray-300">{me.basic_skills}</div>
                    </div>
                </TerminalWindow>
            </div>
        </section>
    )
}
