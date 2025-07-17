import { TerminalWindow } from "@/components/layout/terminal-window"
import { LinkCard } from "@/components/layout/link-card"
import { Github, Linkedin, Twitter, Mail, FileText, Code, Rss, Youtube, Instagram, MoveLeftIcon } from "lucide-react"
import { Link } from "react-router-dom"
import { useConfig } from "@/contexts/config"
import { type JSX, Fragment } from "react"
import { StaticMetadata } from "@/contexts/metadata"

const iconMap: Record<string, JSX.Element> = {
  Github: <Github size={20} />,
  Linkedin: <Linkedin size={20} />,
  Code: <Code size={20} />,
  Rss: <Rss size={20} />,
  FileText: <FileText size={20} />,
  Twitter: <Twitter size={20} />,
  Youtube: <Youtube size={20} />,
  Instagram: <Instagram size={20} />,
  Mail: <Mail size={20} />
};

export default function Links() {
  const { app: { portfolio: me } } = useConfig()

  return (
    <Fragment>
      <StaticMetadata />

      <div className="min-h-screen bg-gray-950 text-white">
        {/* Header */}
        <header className="fixed top-0 left-0 right-0 z-50 bg-gray-900/95 backdrop-blur-sm border-b border-orange-500/30">
          <div className="w-6xl mx-auto px-6 py-4">
            <div className="flex items-center justify-between">
              <div className="text-orange-400 font-bold text-xl font-mono">
                {"> "}
                <span className="text-white">jelius.dev/links</span>
              </div>
              <Link to="/" className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm flex items-center gap-x-2">
                <MoveLeftIcon size={16} /> Back to Portfolio
              </Link>
            </div>
          </div>
        </header>

        {/* Main Content */}
        <main className="pt-20 pb-12 px-6">
          <div className="w-6xl mx-auto">
            {/* Profile Section */}
            <div className="text-center mb-8">
              <div className="w-24 h-24 mx-auto mb-4 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 p-1">
                <div className="w-full h-full rounded-full bg-gray-800 flex items-center justify-center overflow-hidden">
                  <img
                    src={me.pictures[0]}
                    alt="Profile"
                    className="w-full h-full object-cover rounded-full"
                  />
                </div>
              </div>
              <h1 className="text-2xl font-bold text-white mb-2">
                <span className="text-orange-400">@</span>{me.handle}
              </h1>
              <p className="text-gray-400 font-mono text-sm">{me.developer} Developer • {me.enthusiast} Enthusiast</p>
            </div>

            {/* Terminal Info */}
            <TerminalWindow title="user-info" className="mb-8">
              <div className="font-mono text-sm">
                <div className="text-orange-400 mb-2">$ whoami</div>
                <div className="text-gray-300 mb-4">{me.introduction}</div>
                <div className="text-orange-400 mb-2">$ ls social-links/</div>
                <div className="text-gray-300">Found {me.linktree.length} links • Click any link below to connect with me</div>
              </div>
            </TerminalWindow>

            {/* Links Section */}
            <TerminalWindow title="social-links" className="mb-8">
              <div className="space-y-3">
                {me.linktree.map((link, index) => (
                  <LinkCard
                    key={index}
                    title={link.title}
                    description={link.description}
                    url={link.url}
                    icon={iconMap[link.icon]}
                  />
                ))}
              </div>
            </TerminalWindow>

            {/* Terminal Footer */}
            <TerminalWindow title="stats" className="mb-8">
              <div className="font-mono text-sm">
                <div className="text-orange-400 mb-2">$ git log --stat</div>
                <div className="text-gray-300 space-y-1">
                  <div>• {me.linktree.length} social links configured</div>
                  <div>• Last updated: {new Date().toLocaleDateString()}</div>
                  <div>• Status: All systems operational</div>
                </div>

                <div className="text-orange-400 mt-4 mb-2">$ neofetch</div>

                {/* Neofetch Output */}
                <div className="flex gap-6 mb-6">
                  {/* Left side - QR Code */}
                  <div className="flex-shrink-0">
                    <div className="w-64 h-64 bg-gray-800 border border-orange-500/30 rounded flex items-center justify-center">
                      {me.links["jelius.dev/links"]?.qr_code_link ? (
                        <img
                          src={me.links["jelius.dev/links"]?.qr_code_link}
                          alt="QR Code to links page"
                          className="w-62 h-62 rounded object-cover"
                        />
                      ) : (
                        <span className="text-gray-500">QR Code not available</span>
                      )}

                    </div>
                  </div>

                  {/* Right side - System Info */}
                  <div className="flex-1 text-gray-300 space-y-1">
                    {me.neofetch.map((stat, index) => (
                      <div key={index}>
                        <span className="text-orange-400">{stat.label}:</span> {stat.value}
                      </div>
                    ))}
                  </div>
                </div>

                <div className="text-orange-400 mt-4 mb-2">$ echo "Thanks for visiting!"</div>
                <div className="text-gray-300">Thanks for visiting!</div>
              </div>
            </TerminalWindow>

            {/* Footer */}
            <div className="text-center text-gray-500 text-sm font-mono">
              <div className="mb-2">Built with ❤️ using terminal aesthetics</div>
              <div>Copyright © {new Date().getFullYear()} {me.first_name + " " + me.last_name}. All Rights Reserved.
              </div>
            </div>
          </div>
        </main>
      </div>
    </Fragment>
  )
}
