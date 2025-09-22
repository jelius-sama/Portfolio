import { TerminalWindow } from "@/components/ui/terminal-window"
import { Github, ExternalLink, LinkIcon } from "lucide-react"
import { useConfig } from "@/contexts/config"
import { Link } from "react-router-dom"

const knownKeys = ["id", "title", "description", "image", "technologies", "codeUrl", "liveUrl"];

function ProjectLink({ url, label }: { url: string; label: string }) {
    const isRelative = url.startsWith("/");

    // Automatically detect type from label
    const type: "code" | "live" | "additional" =
        label === "Code" ? "code" : label === "Live Demo" ? "live" : "additional";

    const classes = {
        code: "flex items-center gap-2 px-3 py-2 bg-gray-800 text-gray-300 rounded border border-orange-500/30 hover:border-orange-500 hover:text-orange-400 transition-colors text-sm font-mono",
        live: "flex items-center gap-2 px-3 py-2 bg-orange-500 text-black rounded hover:bg-orange-400 transition-colors text-sm font-mono",
        additional: "flex items-center gap-2 px-3 py-2 bg-gray-800 text-gray-300 rounded border border-orange-500/30 hover:border-orange-500 hover:text-orange-400 transition-colors text-sm font-mono",
    };

    const icon = type === "live" ? <ExternalLink size={16} /> : <Github size={16} />;

    const content = <span className={classes[type]}>{type === "additional" ? isRelative ? <LinkIcon size={16} /> : <ExternalLink size={16} /> : icon}{label}</span>;

    return isRelative ? <Link to={url}>{content}</Link> : <a href={url} target="_blank" rel="noopener noreferrer">{content}</a>;
}

export function ProjectsSection() {
    const { app: { portfolio: me } } = useConfig()

    return (
        <section id="projects" className="max-w-6xl mx-auto py-20 px-4 sm:px-6">
            <h2 className="text-3xl font-bold text-center mb-12">
                <span className="text-orange-400">|</span> <span className="text-white">Projects</span>
            </h2>

            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
                {me.projects.map((project) => {
                    const specialKeys = Object.keys(project).filter(key => !knownKeys.includes(key));

                    return (
                        <TerminalWindow key={project.id} title={`project-${project.id}`} className="h-full">
                            <div className="flex flex-col h-full">
                                {/* Project Image/Mockup */}
                                {project.image.length > 0 && (
                                    <div className="mb-4 bg-gray-800 rounded border border-orange-500/20 overflow-hidden">
                                        <img
                                            src={project.image}
                                            alt={project.title}
                                            className="w-full aspect-video object-cover"
                                        />
                                    </div>
                                )}

                                {/* Project Info */}
                                <div className="flex-1">
                                    <h3 className="text-orange-400 font-bold text-xl mb-2 font-mono">{project.title}</h3>
                                    <p className="text-gray-300 text-sm mb-4 leading-relaxed">{project.description}</p>

                                    {/* Tech Stack */}
                                    <div className="mb-4">
                                        <div className="flex flex-wrap gap-2">
                                            {project.technologies.map((tech) => (
                                                <span
                                                    key={tech}
                                                    className="px-2 py-1 text-xs bg-gray-800 text-orange-400 rounded border border-orange-500/30 font-mono"
                                                >
                                                    {tech}
                                                </span>
                                            ))}
                                        </div>
                                    </div>
                                </div>

                                {/* Action Buttons */}
                                <div className="flex gap-3 mt-auto">
                                    {project.codeUrl && <ProjectLink url={project.codeUrl} label="Code" />}
                                    {project.liveUrl && <ProjectLink url={project.liveUrl} label="Live Demo" />}

                                    {specialKeys.map((key) => (
                                        <ProjectLink
                                            key={key}
                                            url={project[key as keyof typeof project] as string}
                                            label={key}
                                        />
                                    ))}
                                </div>
                            </div>
                        </TerminalWindow>
                    )
                })}
            </div>

            {/* Terminal Command for More Projects */}
            <TerminalWindow title="more-projects" className="mt-12 max-w-6xl mx-auto">
                <div className="font-mono text-sm">
                    <div className="text-orange-400 mb-2">$ ls projects/ --all</div>
                    <div className="text-gray-300 mb-2">Found {me.projects.length} projects in current directory</div>
                    <div className="text-orange-400 mt-4">$ echo "More projects coming soon..."</div>
                    <div className="text-gray-300">More projects coming soon...</div>
                </div>
            </TerminalWindow>
        </section>
    )
}
