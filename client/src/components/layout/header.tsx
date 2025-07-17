import { Link } from "react-router-dom"
import { Fragment } from "react"
// import { Sun, Moon } from "lucide-react"

// interface HeaderProps {
//     darkMode: boolean
//     toggleDarkMode: () => void
// }

export function Header() {
    const navItems = ["Home", "About", "Skills", "Experience", "Projects", "Contact", "Links"]

    return (
        <header className="fixed top-0 left-0 right-0 z-50 bg-gray-900/95 backdrop-blur-sm border-b border-orange-500/30">
            <div className="w-6xl mx-auto px-6 py-4">
                <div className="flex items-center justify-between">
                    <div className="text-orange-400 font-bold text-xl font-mono">
                        {"> "}
                        <span className="text-white">jelius.dev</span>
                    </div>

                    <nav className="hidden md:flex items-center gap-8">
                        {navItems.map((item, index) => (
                            <Fragment>
                                {item === "Links" ? (
                                    <Link
                                        key={index}
                                        to="/links"
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
                        ))}
                    </nav>

                    {/* <button onClick={toggleDarkMode} className="p-2 text-gray-400 hover:text-orange-400 transition-colors"> */}
                    {/*     {darkMode ? <Sun size={20} /> : <Moon size={20} />} */}
                    {/* </button> */}
                </div>
            </div>
        </header>
    )
}
