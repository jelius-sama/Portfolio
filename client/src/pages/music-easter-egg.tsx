import { Fragment, useLayoutEffect, useState, lazy } from "react"
import { PathBasedMetadata } from "@/contexts/metadata"
import { useLocation } from "react-router-dom"
import { Music, Heart } from 'lucide-react';
import { Footer } from "@/components/ui/footer"
import { TerminalWindow } from "@/components/ui/terminal-window"

const NotFoundPage = lazy(() => import("@/pages/not-found"))
const MusicPlayer = lazy(() => import("@/components/layout/music-player"))

export default function Achievements() {
  const location = useLocation()
  const [shouldAllow, setShouldAllow] = useState(false)

  useLayoutEffect(() => {
    const result = sessionStorage.getItem("did_come_from_terminal")
    setShouldAllow(result === "yes")
  }, [location])

  return shouldAllow ? (
    <Fragment>
      <PathBasedMetadata paths={["*", "#music_easter_egg"]} />

      <section className="max-w-6xl mx-auto pt-20 px-4 sm:px-6 font-mono flex flex-col w-full h-screen">
        <h2 className="text-white text-2xl font-bold">Easter egg</h2>
        <p className="text-muted-foreground text-md mb-6">My top 4 favorite songs</p>

        <div className="w-full flex flex-col gap-y-4">
          <MusicPlayer
            title="だから僕は音楽を辞めた"
            artist="ヨルシカ"
            albumArt="/assets/easter eggs/だから僕は音楽を辞めた.jpg"
            audioSrc="/assets/easter eggs/だから僕は音楽を辞めた.mp3"
          />

          <MusicPlayer
            title="アンコール"
            artist="YOASOBI"
            albumArt="/assets/easter eggs/encore.jpg"
            audioSrc="/assets/easter eggs/アンコール.mp3"
          />

          <MusicPlayer
            title="Ref：rain"
            artist="Aimer"
            albumArt="/assets/easter eggs/ref_rain.jpg"
            audioSrc="/assets/easter eggs/ref_rain.mp3"
          />

          <MusicPlayer
            title="ハレハレヤ | Harehare Ya - mix ver."
            artist="Kityod x keita x sou"
            albumArt="/assets/easter eggs/harehareya.jpg"
            audioSrc="/assets/easter eggs/Harehare Ya mix ver.mp3"
          />
        </div>

        <Footer className="mt-auto max-w-6xl mx-auto px-4 sm:px-6 py-8" leading={<EasterEggFooter />} trailing={<p className="mt-1 text-gray-600">Thanks for finding this hidden gem 🎵</p>} />
      </section>
    </Fragment>
  ) : <NotFoundPage />
}

function EasterEggFooter() {
  return (
    <div className="my-8">
      {/* Main Footer Content */}
      <div className="text-center space-y-2">
        {/* Music Icon and Message */}
        <div className="flex items-start justify-center gap-2 text-gray-300">
          <Music className="text-orange-400" size={20} />
          <span className="font-mono text-sm">
            Music is the soundtrack to life's debugging sessions
          </span>
        </div>

        {/* Made with Love */}
        <TerminalWindow className="gap-2 text-gray-400 text-sm" title="">
          <div className="font-mono text-xs text-left">
            <div className="text-orange-400 mb-1">$ echo "easter_egg_discovered"</div>
            <div className="text-gray-300">
              "The best code is like good music - it tells a story"
            </div>
          </div>

          <span className="flex flex-row gap-2 items-center justify-center mt-4">
            Made with
            <Heart className="text-red-400 fill-current" size={16} />
          </span>
        </TerminalWindow>
      </div>
    </div>
  );
}
