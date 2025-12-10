import { Link, useLocation, useParams } from "react-router-dom"
import { Fragment, useState, type MouseEventHandler } from "react"
import { Menu, X } from "lucide-react"

const DISABLE_SECTION_LINKS = true;

function SmartLink({ type, to, href, item, index, className, ...rest }: { type: "link" | "a"; to?: string; href?: string; item: string; index: number; className?: string; onClick: MouseEventHandler<HTMLElement> }) {
    const Tag = type === "a" ? "a" : Link;
    const base = "flex items-center gap-3 p-3 bg-gray-800/50 border border-orange-500/20 rounded hover:border-orange-500 hover:bg-gray-800/70 transition-all duration-200 group";
    const number = String(index + 1).padStart(2, "0") + ".";

    const props = {
        ...(type === "a" ? { href } : { to }),
        className: `${base} ${className ?? ""}`,
        ...rest,
    };

    return (
        !(DISABLE_SECTION_LINKS && type === "a") && (
            // @ts-ignore
            <Tag {...props}>
                <span className="text-orange-400 group-hover:text-orange-300">{number}</span>
                <span className="text-gray-300 group-hover:text-white font-mono">{item}</span>
                <span className="ml-auto text-gray-500 group-hover:text-orange-400 text-xs">{type === "link" ? "page" : "section"}</span>
            </Tag>
        )
    );
};

export function Header() {
    const navItems = ["Home", "About", "Skills", "Experience", "Projects", "Contact", "Blogs", "Achievements", "Links"]
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
    const location = useLocation()
    const params = useParams()

    const handleNavClick = () => {
        setMobileMenuOpen(false)
    }

    return (
        <Fragment>
            <header className="fixed top-0 left-0 right-0 z-50 bg-gray-900/95 backdrop-blur-sm border-b border-orange-500/30">
                <div className="max-w-6xl mx-auto px-4 sm:px-6 py-4">
                    <div className="flex items-center justify-between">
                        <div className="text-orange-400 font-bold text-xl font-mono">
                            {"> "}
                            <span className="text-white">
                                jelius.dev
                                {location.pathname.startsWith("/blog/") && params?.id
                                    ? `/${params.id}`
                                    : ["/links", "/blogs", "/achievements"].includes(location.pathname)
                                        ? location.pathname
                                        : ""}
                            </span>
                        </div>

                        <nav className="hidden lg:flex items-center gap-8">
                            {location.pathname === "/" ? (
                                navItems.map((item, index) => item === "Links" || item === "Blogs" || item === "Achievements" ? (
                                    <Link key={index} to={"/" + item.toLowerCase()} className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm"> {item} </Link>
                                ) : !DISABLE_SECTION_LINKS && (
                                    <a key={index} href={`#${item.toLowerCase()}`} className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm">{item}</a>
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

                                    {location.pathname !== "/achievements" && (
                                        <Link to="/achievements" className="text-gray-300 hover:text-orange-400 transition-colors font-mono text-sm flex items-center gap-x-2">
                                            Achievements
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

                        <div className="lg:hidden flex items-center gap-2">
                            <button
                                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                                className="p-2 text-gray-400 hover:text-orange-400 transition-colors"
                            >
                                {mobileMenuOpen ? <X size={20} /> : <Menu size={20} />}
                            </button>
                        </div>
                    </div>
                </div>
            </header>

            {/* Mobile Navigation Menu */}
            {mobileMenuOpen && (
                <div className="fixed inset-0 z-40 lg:hidden">
                    {/* Backdrop */}
                    <div className="fixed inset-0 bg-gray-950/80 backdrop-blur-sm" onClick={() => setMobileMenuOpen(false)} />

                    {/* Mobile Menu */}
                    <div className="fixed top-20 left-4 right-4 bg-gray-900 rounded-lg border border-orange-500/30 overflow-hidden">
                        {/* Terminal Header */}
                        <div className="flex items-center gap-2 px-4 py-3 bg-gray-800 border-b border-orange-500/30">
                            <div className="flex gap-2">
                                <div className="w-3 h-3 rounded-full bg-red-500"></div>
                                <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                                <div className="w-3 h-3 rounded-full bg-green-500"></div>
                            </div>
                            <span className="text-gray-400 text-sm font-mono">navigation</span>
                        </div>

                        {/* Navigation Items */}
                        <div className="p-4">
                            <div className="font-mono text-sm">
                                <div className="text-orange-400 mb-4">$ ls /navigation/</div>
                                <div className="space-y-2">
                                    {location.pathname === "/" ? (
                                        navItems.map((item, index) => {
                                            if (DISABLE_SECTION_LINKS) index -= 6;
                                            return (
                                                <Fragment>
                                                    {item === "Links" || item === "Blogs" || item === "Achievements" ? (
                                                        <SmartLink type="link" key={index} index={index} item={item} onClick={handleNavClick} to={"/" + item.toLowerCase()} />
                                                    ) : (
                                                        <SmartLink type="a" key={index} index={index} item={item} onClick={handleNavClick} href={"#" + item.toLowerCase()} />
                                                    )}
                                                </Fragment>
                                            )
                                        })
                                    ) : (
                                        <Fragment>
                                            {[
                                                { path: "/", label: "Portfolio" },
                                                { path: "/blogs", label: "Blogs" },
                                                { path: "/achievements", label: "Achievements" },
                                                { path: "/links", label: "Links" },
                                            ]
                                                .filter((item) => item.path === "/" || location.pathname !== item.path)
                                                .map((item, index) => (
                                                    <SmartLink type="link" key={index} index={index} item={item.label} onClick={handleNavClick} to={item.path} />
                                                ))}
                                        </Fragment>
                                    )}
                                </div>

                                <div className="mt-4 pt-4 border-t border-orange-500/20">
                                    <div className="text-orange-400 mb-2">$ echo "Navigation ready"</div>
                                    <div className="text-gray-400 text-xs">Tap any item to navigate</div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </Fragment>
    )
}
