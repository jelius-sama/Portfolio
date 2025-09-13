import { Fragment, useEffect } from "react"
import { StaticMetadata } from "@/contexts/metadata"
import { HeroSection } from "@/components/layout/hero-section"
import { AboutSection } from "@/components/layout/about-section"
import { SkillsSection } from "@/components/layout/skills-section"
import { ProjectsSection } from "@/components/layout/projects-section"
import { ExperienceSection } from "@/components/layout/experience-section"
import { ContactSection } from "@/components/layout/contact-section"

export default function Home() {
  useEffect(() => {
    const image = document.getElementById("slow-af") as HTMLImageElement | null;

    if (image) {
      const handleImageLoad = () => {
        const event = new CustomEvent("PageLoaded", {
          detail: { pathname: window.location.pathname },
        });
        window.dispatchEvent(event);
      };

      // If already loaded (from cache), fire immediately
      if (image.complete && image.naturalHeight !== 0) {
        handleImageLoad();
      } else {
        image.addEventListener("load", handleImageLoad);
      }

      return () => {
        image.removeEventListener("load", handleImageLoad);
      };
    }
  }, []);

  return (
    <Fragment>
      <StaticMetadata />

      <HeroSection />
      <AboutSection />
      <SkillsSection />
      <ExperienceSection />
      <ProjectsSection />
      <ContactSection />
    </Fragment>
  )
}
