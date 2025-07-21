import { TerminalWindow } from "@/components/ui/terminal-window"
import { useConfig } from "@/contexts/config"
import { Mail, Github, Linkedin, Twitter, LinkIcon } from "lucide-react"
import { type JSX } from "react"
import { Link } from "react-router-dom"
import { InteractiveTerminal } from "@/components/layout/interactive-term"
import { Footer } from "@/components/ui/footer"

const iconMap: Record<string, JSX.Element> = {
    Email: <Mail className="text-orange-400" size={20} />,
    GitHub: <Github size={18} />,
    LinkedIn: <Linkedin size={18} />,
    X: <Twitter size={18} />,
    "jelius.dev/links": <LinkIcon size={18} />
};

export function ContactSection() {
    const { app: { portfolio: me } } = useConfig()

    return (
        <section id="contact" className="py-20 px-6 bg-gray-900/50">
            <div className="max-w-6xl mx-auto">
                <h2 className="text-3xl font-bold text-center mb-12">
                    <span className="text-orange-400">|</span> <span className="text-white">Get in Touch</span>
                </h2>

                <div className="grid md:grid-cols-2 gap-8 mb-12">
                    <TerminalWindow title="contact-info">
                        <div className="space-y-4">
                            {/* Email */}
                            {"Email" in me.links && (
                                <a
                                    href={me.links.Email.link}
                                    className="flex items-center gap-3 text-gray-300"
                                >
                                    {iconMap.Email}
                                    <span>{me.links.Email.username}</span>
                                </a>
                            )}

                            {/* Social Media */}
                            <div className="text-gray-400 text-sm font-mono">Social Media Profiles</div>
                            <div className="space-y-3">
                                {Object.entries(me.links)
                                    .filter(([key]) => key !== "Email" && key in iconMap)
                                    .map(([key, value]) => (
                                        key === "jelius.dev/links" ? (
                                            <Link
                                                key={key}
                                                to={value.link}
                                                className="flex items-center gap-3 text-gray-300 hover:text-orange-400 cursor-pointer transition-colors"
                                            >
                                                {iconMap[key]}
                                                <span>Linktree</span>
                                            </Link>

                                        ) : (
                                            <a
                                                key={key}
                                                href={value.link}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className="flex items-center gap-3 text-gray-300 hover:text-orange-400 cursor-pointer transition-colors"
                                            >
                                                {iconMap[key]}
                                                <span>{key}</span>
                                            </a>
                                        )
                                    ))}
                            </div>
                        </div>
                    </TerminalWindow>

                    <InteractiveTerminal />
                </div>

                <TerminalWindow title="latest-logs" className="mb-8">
                    <div className="font-mono text-sm space-y-2">
                        {me.logs.map((log, index) => (
                            <div className="text-green-400" key={index}>{log}</div>
                        ))}
                    </div>
                </TerminalWindow>

                <Footer />
            </div>
        </section>
    )
}
