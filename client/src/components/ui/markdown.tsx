import React from "react"
import { useState } from "react"
import ReactMarkdown from "react-markdown"
import remarkGfm from "remark-gfm"
import remarkBreaks from "remark-breaks"
import { Copy, Check, ExternalLink, AlertTriangle, Info, CheckCircle, XCircle } from "lucide-react"

interface MarkdownRendererProps {
    content: string
    className?: string
}

function CodeBlock({ children, className }: { children: string; className?: string }) {
    const [copied, setCopied] = useState(false)
    const language = className?.replace("language-", "") || ""

    const handleCopy = async () => {
        try {
            await navigator.clipboard.writeText(children)
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        } catch (err) {
            console.error("Failed to copy code:", err)
        }
    }

    return (
        <div className="bg-gray-900 border border-orange-500/30 rounded-lg overflow-hidden my-6">
            {/* Terminal Header */}
            <div className="flex items-center justify-between px-4 py-2 bg-gray-800 border-b border-orange-500/30">
                <div className="flex items-center gap-2">
                    <div className="flex gap-1">
                        <div className="w-3 h-3 rounded-full bg-red-500"></div>
                        <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                        <div className="w-3 h-3 rounded-full bg-green-500"></div>
                    </div>
                    <span className="text-gray-400 text-sm font-mono">{language ? `${language}-code` : "code-block"}</span>
                </div>
                <button
                    onClick={handleCopy}
                    className="flex items-center gap-1 px-2 py-1 text-xs text-gray-400 hover:text-orange-400 transition-colors"
                >
                    {copied ? <Check size={12} /> : <Copy size={12} />}
                    {copied ? "Copied!" : "Copy"}
                </button>
            </div>
            {/* Code Content */}
            <div className="p-4 overflow-x-auto">
                <pre className="text-sm text-gray-300 font-mono">
                    <code>{children}</code>
                </pre>
            </div>
        </div>
    )
}

function Link({ href, children, title }: { href?: string; children: React.ReactNode; title?: string }) {
    if (!href) return <span>{children}</span>

    const isExternal = href.startsWith("http") || href.startsWith("//")
    const isEmail = href.startsWith("mailto:")
    const isAnchor = href.startsWith("#")

    return (
        <a
            href={href}
            title={title}
            target={isExternal ? "_blank" : undefined}
            rel={isExternal ? "noopener noreferrer" : undefined}
            className="text-orange-400 hover:text-orange-300 underline decoration-orange-500/50 hover:decoration-orange-300 transition-colors inline-flex items-center gap-1"
        >
            {children}
            {isExternal && <ExternalLink size={12} />}
            {isEmail && <span className="text-xs">✉</span>}
            {isAnchor && <span className="text-xs">#</span>}
        </a>
    )
}

function Blockquote({ children }: { children: React.ReactNode }) {
    // Check if this is a special blockquote (GitHub-style alerts)
    const childrenString = React.Children.toArray(children).join("")
    const alertMatch = childrenString.match(/^\[!(NOTE|WARNING|TIP|DANGER|INFO)\]/)

    if (alertMatch) {
        const type = alertMatch[1].toLowerCase()
        const getAlertStyle = (type: string) => {
            switch (type) {
                case "warning":
                    return {
                        border: "border-l-4 border-yellow-500",
                        bg: "bg-yellow-500/10",
                        icon: <AlertTriangle className="text-yellow-500" size={16} />,
                        label: "Warning",
                    }
                case "danger":
                    return {
                        border: "border-l-4 border-red-500",
                        bg: "bg-red-500/10",
                        icon: <XCircle className="text-red-500" size={16} />,
                        label: "Danger",
                    }
                case "tip":
                    return {
                        border: "border-l-4 border-green-500",
                        bg: "bg-green-500/10",
                        icon: <CheckCircle className="text-green-500" size={16} />,
                        label: "Tip",
                    }
                case "note":
                case "info":
                    return {
                        border: "border-l-4 border-blue-500",
                        bg: "bg-blue-500/10",
                        icon: <Info className="text-blue-500" size={16} />,
                        label: "Note",
                    }
                default:
                    return {
                        border: "border-l-4 border-orange-500",
                        bg: "bg-gray-800/30",
                        icon: null,
                        label: null,
                    }
            }
        }

        const style = getAlertStyle(type)

        return (
            <div className={`${style.border} ${style.bg} pl-4 py-3 my-4 rounded-r`}>
                {style.icon && (
                    <div className="flex items-center gap-2 mb-2 font-mono text-sm font-bold">
                        {style.icon}
                        {style.label}
                    </div>
                )}
                <div className="text-gray-300 [&>*:last-child]:mb-0">
                    {/* Remove the alert syntax from the content */}
                    {React.Children.map(children, (child) => {
                        if (React.isValidElement(child) && child.type === "p") {
                            //@ts-ignore
                            const content = React.Children.toArray(child.props.children)
                            const firstChild = content[0]
                            if (typeof firstChild === "string" && firstChild.includes("[!")) {
                                const cleanedContent = firstChild.replace(/^\[![^\]]+\]\s*/, "")
                                return React.cloneElement(child, {}, cleanedContent, ...content.slice(1))
                            }
                        }
                        return child
                    })}
                </div>
            </div>
        )
    }

    return (
        <div className="border-l-4 border-orange-500 bg-gray-800/30 pl-4 py-2 my-4 italic">
            <div className="text-gray-300">{children}</div>
        </div>
    )
}

function Table({ children }: { children: React.ReactNode }) {
    return (
        <div className="overflow-x-auto my-6">
            <div className="bg-gray-900 border border-orange-500/30 rounded-lg overflow-hidden">
                {/* Terminal Header */}
                <div className="flex items-center gap-2 px-4 py-2 bg-gray-800 border-b border-orange-500/30">
                    <div className="flex gap-1">
                        <div className="w-3 h-3 rounded-full bg-red-500"></div>
                        <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                        <div className="w-3 h-3 rounded-full bg-green-500"></div>
                    </div>
                    <span className="text-gray-400 text-sm font-mono">data-table</span>
                </div>
                {/* Table Content */}
                <div className="overflow-x-auto">
                    <table className="w-full font-mono text-sm">{children}</table>
                </div>
            </div>
        </div>
    )
}

function TaskListItem({ children, checked }: { children: React.ReactNode; checked?: boolean }) {
    return (
        <li className="flex items-start gap-3 list-none">
            <div className="mt-1">
                {checked ? (
                    <div className="w-4 h-4 bg-orange-500 rounded border border-orange-500 flex items-center justify-center">
                        <Check size={10} className="text-black" />
                    </div>
                ) : (
                    <div className="w-4 h-4 border border-orange-500/50 rounded"></div>
                )}
            </div>
            <span className="text-gray-300 flex-1">{children}</span>
        </li>
    )
}

export function MarkdownRenderer({ content, className = "" }: MarkdownRendererProps) {
    return (
        <div className={`terminal-markdown prose prose-invert max-w-none ${className}`}>
            <ReactMarkdown
                remarkPlugins={[remarkGfm, remarkBreaks]}
                components={{
                    // Headers
                    h1: ({ children, ...props }) => (
                        <h1 className="text-3xl font-bold text-orange-400 mb-6 mt-8 border-b border-orange-500/30 pb-2" {...props}>
                            {children}
                        </h1>
                    ),
                    h2: ({ children, ...props }) => (
                        <h2 className="text-2xl font-bold text-orange-400 mb-4 mt-6" {...props}>
                            {children}
                        </h2>
                    ),
                    h3: ({ children, ...props }) => (
                        <h3 className="text-xl font-bold text-orange-400 mb-3 mt-5" {...props}>
                            {children}
                        </h3>
                    ),
                    h4: ({ children, ...props }) => (
                        <h4 className="text-lg font-bold text-orange-400 mb-2 mt-4" {...props}>
                            {children}
                        </h4>
                    ),
                    h5: ({ children, ...props }) => (
                        <h5 className="text-base font-bold text-orange-400 mb-2 mt-3" {...props}>
                            {children}
                        </h5>
                    ),
                    h6: ({ children, ...props }) => (
                        <h6 className="text-sm font-bold text-orange-400 mb-2 mt-3" {...props}>
                            {children}
                        </h6>
                    ),

                    // Paragraphs
                    p: ({ children, ...props }) => (
                        <p className="text-gray-300 mb-4 leading-relaxed" {...props}>
                            {children}
                        </p>
                    ),

                    // Links
                    a: ({ href, children, title, ...props }) => (
                        <Link href={href} title={title} {...props}>
                            {children}
                        </Link>
                    ),

                    // Code - Enhanced detection for inline vs block code
                    //@ts-ignore
                    code: ({ children, className, inline, node, ...props }) => {
                        const codeContent = String(children).replace(/\n$/, "")

                        // Check if this is definitely a fenced code block
                        const isFencedCodeBlock = className && className.startsWith("language-")

                        // Check if this is inside a pre element (another indicator of code block)
                        //@ts-ignore
                        const isInsidePreElement = node?.parent && node.parent.tagName === "pre"

                        // Check if content has multiple lines (likely a code block)
                        const hasMultipleLines = codeContent.includes("\n")

                        // If any of these conditions are true, it's a code block
                        const isCodeBlock = isFencedCodeBlock || isInsidePreElement || hasMultipleLines

                        // If it's explicitly marked as inline OR none of the block indicators are present
                        if (inline === true || (!isCodeBlock && !isFencedCodeBlock)) {
                            return (
                                <code className="bg-gray-800 px-2 py-1 rounded text-orange-400 font-mono text-sm border border-orange-500/20">
                                    {children}
                                </code>
                            )
                        }

                        // Otherwise, it's a code block - use terminal window
                        return (
                            <CodeBlock className={className} {...props}>
                                {codeContent}
                            </CodeBlock>
                        )
                    },

                    // Pre - handle pre elements that contain code blocks
                    pre: ({ children }) => {
                        // If pre contains a code element, let the code component handle the styling
                        return <>{children}</>
                    },

                    // Lists
                    ul: ({ children, ...props }) => (
                        <ul className="list-none space-y-2 my-4 ml-4" {...props}>
                            {children}
                        </ul>
                    ),
                    ol: ({ children, ...props }) => (
                        <ol className="list-none space-y-2 my-4 ml-4" {...props}>
                            {children}
                        </ol>
                    ),

                    //@ts-ignore
                    li: ({ children, checked, ...props }) => {
                        // Handle task list items
                        if (typeof checked === "boolean") {
                            return (
                                <TaskListItem checked={checked} {...props}>
                                    {children}
                                </TaskListItem>
                            )
                        }

                        // Regular list items
                        //@ts-ignore
                        const isOrderedList = props.node?.parent?.tagName === "ol"
                        //@ts-ignore
                        const index = props.index || 0

                        return (
                            <li className="flex items-start gap-2" {...props}>
                                <span className="text-orange-400 mt-1 font-mono">{isOrderedList ? `${index + 1}.` : "•"}</span>
                                <span className="text-gray-300">{children}</span>
                            </li>
                        )
                    },

                    // Blockquotes
                    blockquote: ({ children, ...props }) => <Blockquote {...props}>{children}</Blockquote>,

                    // Tables
                    table: ({ children, ...props }) => <Table {...props}>{children}</Table>,
                    thead: ({ children, ...props }) => <thead {...props}>{children}</thead>,
                    tbody: ({ children, ...props }) => <tbody {...props}>{children}</tbody>,
                    tr: ({ children, ...props }) => (
                        <tr className="hover:bg-gray-800/30 transition-colors" {...props}>
                            {children}
                        </tr>
                    ),
                    th: ({ children, ...props }) => (
                        <th className="px-4 py-3 text-left text-orange-400 border-b border-orange-500/30" {...props}>
                            {children}
                        </th>
                    ),
                    td: ({ children, ...props }) => (
                        <td className="px-4 py-3 text-gray-300 border-b border-gray-700/50" {...props}>
                            {children}
                        </td>
                    ),

                    // Text formatting
                    strong: ({ children, ...props }) => (
                        <strong className="text-white font-bold" {...props}>
                            {children}
                        </strong>
                    ),
                    em: ({ children, ...props }) => (
                        <em className="text-gray-200 italic" {...props}>
                            {children}
                        </em>
                    ),
                    del: ({ children, ...props }) => (
                        <del className="text-gray-400 line-through" {...props}>
                            {children}
                        </del>
                    ),

                    // Horizontal rule
                    hr: ({ ...props }) => <hr className="border-orange-500/30 my-8" {...props} />,

                    // Images
                    img: ({ src, alt, title, ...props }) => (
                        <div className="my-6">
                            <div className="bg-gray-900 border border-orange-500/30 rounded-lg overflow-hidden">
                                <div className="flex items-center gap-2 px-4 py-2 bg-gray-800 border-b border-orange-500/30">
                                    <div className="flex gap-1">
                                        <div className="w-3 h-3 rounded-full bg-red-500"></div>
                                        <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                                        <div className="w-3 h-3 rounded-full bg-green-500"></div>
                                    </div>
                                    <span className="text-gray-400 text-sm font-mono">{alt || "image"}</span>
                                </div>
                                <div className="p-4">
                                    <img
                                        src={src || "/placeholder.svg"}
                                        alt={alt}
                                        title={title}
                                        className="max-w-full h-auto rounded border border-orange-500/20"
                                        {...props}
                                    />
                                    {title && <div className="text-gray-400 text-sm mt-2 font-mono text-center">{title}</div>}
                                </div>
                            </div>
                        </div>
                    ),
                }}
            >
                {content}
            </ReactMarkdown>
        </div>
    )
}
