import { Fragment } from "react"
import { StaticMetadata } from "@/contexts/metadata"
import { Header } from "@/components/layout/header"
import { HeroSection } from "@/components/layout/hero-section"
import { AboutSection } from "@/components/layout/about-section"
import { SkillsSection } from "@/components/layout/skills-section"
import { ProjectsSection } from "@/components/layout/projects-section"
import { ExperienceSection } from "@/components/layout/experience-section"
import { ContactSection } from "@/components/layout/contact-section"

export default function Home() {

  return (
    <Fragment>
      <StaticMetadata />

      <div className="min-h-screen bg-gray-950 text-white">
        <Header />
        <HeroSection />
        <AboutSection />
        <SkillsSection />
        <ExperienceSection />
        <ProjectsSection />
        <ContactSection />
      </div>
    </Fragment>
  )
}
