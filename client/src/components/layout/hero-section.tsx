import { TerminalWindow } from "@/components/ui/terminal-window"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"
import { useConfig } from "@/contexts/config"

export function HeroSection() {
    const { app: { portfolio: me } } = useConfig()

    const handleDownload = ({ url, filename }: { url: string, filename: string | null }) => {
        const link = document.createElement("a");
        link.href = url;
        link.download = filename || url.split("/").pop() || "download";
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <section id="home" className="min-h-screen w-full max-w-6xl text-center flex flex-col items-center justify-center mx-auto px-4 sm:px-6 pt-20">
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
                <p className="text-xl text-gray-300 font-mono">{"developer" in me ? me.developer + " Developer" : me.student + " Student"}</p>
            </div>

            <div className="flex flex-wrap gap-4 justify-center mb-8">
                <Dialog>
                    <DialogTrigger asChild={true}>
                        <button className="px-6 py-3 bg-orange-500 text-black font-medium rounded hover:bg-orange-400 transition-colors">
                            Download CV
                        </button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Download</DialogTitle>
                            <DialogDescription>
                                Choose your preferred file format for download: Markdown or PDF.
                            </DialogDescription>
                        </DialogHeader>

                        <div className="flex space-x-20 justify-center">
                            {/* MD button with note */}
                            <div className="flex flex-col items-start space-y-1">
                                <span className="text-sm text-gray-500">*Recommended</span>
                                <button onClick={() => handleDownload({ url: "https://jelius.dev/assets/README.md", filename: "CV.md" })} className="px-6 py-2 bg-orange-500 text-black font-medium rounded hover:bg-orange-400 transition-colors">Markdown</button>
                            </div>

                            {/* PDF button in its own flex-col to align properly */}
                            <div className="flex flex-col items-center justify-end">
                                <button onClick={() => handleDownload({ url: "https://jelius.dev/assets/README.pdf", filename: "CV.pdf" })} className="px-6 py-2 border border-orange-500 text-orange-400 font-medium rounded hover:bg-orange-500/10 transition-colors">PDF</button>
                            </div>
                        </div>
                    </DialogContent>
                </Dialog>

                <a href="#contact" className="px-6 py-3 border border-orange-500 text-orange-400 font-medium rounded hover:bg-orange-500/10 transition-colors">
                    Contact Info
                </a>
            </div>

            <TerminalWindow title="welcome" className="w-full text-left">
                <div className="font-mono text-sm">
                    <div className="text-orange-400 mb-2">$ whoami</div>
                    <div className="text-gray-300 mb-4">{me.introduction}</div>
                    <div className="text-orange-400 mb-2">$ cat skills.txt</div>
                    <div className="text-gray-300">{me.basic_skills}</div>
                </div>
            </TerminalWindow>
        </section>
    )
}
