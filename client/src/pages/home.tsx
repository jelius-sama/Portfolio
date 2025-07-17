import { Fragment } from "react"
import { StaticMetadata } from "@/contexts/metadata"
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

      <script
        type="application/ld+json"
        data-aria-hidden="true"
        aria-hidden="true"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify({
            "@context": "https://schema.org",
            "@type": "Person",
            name: "Jelius",
            url: "https://jelius.dev",
            sameAs: [
              "https://www.linkedin.com/in/jelius-basumatary-485044339/",
              "https://github.com/jelius-sama"
            ]
          }),
        }}
      />

      <HeroSection />
      <AboutSection />
      <SkillsSection />
      <ExperienceSection />
      <ProjectsSection />
      <ContactSection />
    </Fragment>
  )
}
