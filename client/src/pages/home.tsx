import { Fragment } from "react"
import { StaticMetadata } from "@/contexts/metadata"

export default function Home() {
  return (
    <Fragment>
      <StaticMetadata />

      <p>Hello, World!</p>
    </Fragment>
  )
}
