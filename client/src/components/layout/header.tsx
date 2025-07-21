import { Link, useLocation, useParams } from "react-router-dom"
import { Fragment } from "react"

export function Header() {
    const navItems = ["Home", "About", "Skills", "Experience", "Projects", "Contact", "Blogs", "Links"]
    const location = useLocation()
    const params = useParams()

    return (
        <header className="fixed top-0 left-0 right-0 z-50 bg-gray-900/95 backdrop-blur-sm border-b border-orange-500/30">
            <div className="max-w-6xl mx-auto py-4">
                <div className="flex items-center justify-between">
                    <div className="text-orange-400 font-bold text-xl font-mono">
                        {"> "}
                        <span className="text-white">
                            jelius.dev
                            {location.pathname.startsWith("/blog/") && params?.id
                                ? `/${params.id}`
                                : ["/links", "/blogs"].includes(location.pathname)
                                    ? location.pathname
                                    : ""}
                        </span>
                    </div>

                    <nav className="hidden md:flex items-center gap-8">
                        {location.pathname === "/" ? (
                            navItems.map((item, index) => (
                                <Fragment>
                                    {item === "Links" || item === "Blogs" ? (
                                        <Link
                                            key={index}
                                            to={"/" + item.toLowerCase()}
                                            className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm"
                                        >
                                            {item}
                                        </Link>

                                    ) : (
                                        <a
                                            key={index}
                                            href={`#${item.toLowerCase()}`}
                                            className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm"
                                        >
                                            {item}
                                        </a>

                                    )}
                                </Fragment>
                            ))
                        ) : (
                            <Fragment>
                                <Link to="/" className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm flex items-center gap-x-2">
                                    Portfolio
                                </Link>

                                {location.pathname !== "/blogs" && (
                                    <Link to="/blogs" className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm flex items-center gap-x-2">
                                        Blogs
                                    </Link>
                                )}

                                {location.pathname !== "/links" && (
                                    <Link to="/links" className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm flex items-center gap-x-2">
                                        Links
                                    </Link>
                                )}
                            </Fragment>
                        )}
                    </nav>
                </div>
            </div>
        </header>
    )
}
