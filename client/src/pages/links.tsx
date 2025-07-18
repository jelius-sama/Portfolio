import { TerminalWindow } from "@/components/layout/terminal-window"
import { LinkCard } from "@/components/layout/link-card"
import { Github, Linkedin, Twitter, Mail, FileText, Code, Rss, Youtube, Instagram } from "lucide-react"
import { useConfig } from "@/contexts/config"
import { type JSX, Fragment } from "react"
import { StaticMetadata } from "@/contexts/metadata"
import { useQuery } from "@tanstack/react-query"
import { toast } from "sonner"

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

  const { isPending, error, data } = useQuery({
    queryKey: ['neofetch'],
    queryFn: async () => {
      const res = await fetch('/api/neofetch')
      if (!res.ok) {
        const errMsg = "Failed to fetch netfetch command!"
        toast.error(errMsg)
        throw new Error(errMsg)
      }
      return res.json() as Promise<Array<{ label: string, value: string }>>
    }
  })

  const neofetch: undefined | Array<{ label: string, value: string }> = isPending ? undefined : error ? me.neofetch : data && data

  return (
    <Fragment>
      <StaticMetadata />

      <script
        type="application/ld+json"
        data-aria-hidden="true"
        aria-hidden="true"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify({
            "@context": "https://schema.org",
            "@type": "Person",
            name: "Jelius",
            url: "https://jelius.dev/links",
            sameAs: [
              "https://www.linkedin.com/in/jelius-basumatary-485044339/",
              "https://github.com/jelius-sama"
            ]
          }),
        }}
      />

      {/* Main Content */}
      <main className="pt-20 pb-12 px-6">
        <div className="max-w-6xl mx-auto">
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
            <p className="text-gray-400 font-mono text-sm">{"developer" in me ? me.developer + " Developer" : me.student + " Student"} • {me.enthusiast} Enthusiast</p>
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

              <div className="text-orange-400 mt-4 mb-2 sm:block hidden">$ neofetch</div>
              <div className="text-orange-400 mt-4 mb-2 sm:hidden block">$ neofetch | less</div>

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
                {neofetch && (
                  <div className="sm:block flex-1 text-gray-300 space-y-1 hidden">
                    {neofetch.map((stat, index) => (
                      <div key={index}>
                        <span className="text-orange-400">{stat.label}:</span> {stat.value}
                      </div>
                    ))}
                  </div>
                )}
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
    </Fragment>
  )
}
