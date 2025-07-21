import type React from "react"
import { useState } from "react"
import { Copy, Check } from "lucide-react"

interface MarkdownRendererProps {
    content: string
    className?: string
}

interface CodeBlockProps {
    code: string
    language?: string
}

function CodeBlock({ code, language }: CodeBlockProps) {
    const [copied, setCopied] = useState(false)

    const handleCopy = async () => {
        try {
            await navigator.clipboard.writeText(code)
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
            <div className="p-4">
                <pre className="text-sm text-gray-300 font-mono overflow-x-auto">
                    <code>{code}</code>
                </pre>
            </div>
        </div>
    )
}

interface InlineCodeProps {
    children: string
}

function InlineCode({ children }: InlineCodeProps) {
    return (
        <code className="bg-gray-800 px-2 py-1 rounded text-orange-400 font-mono text-sm border border-orange-500/20">
            {children}
        </code>
    )
}

interface BlockquoteProps {
    children: React.ReactNode
}

function Blockquote({ children }: BlockquoteProps) {
    return (
        <div className="border-l-4 border-orange-500 bg-gray-800/30 pl-4 py-2 my-4 italic">
            <div className="text-gray-300">{children}</div>
        </div>
    )
}

interface TableProps {
    headers: string[]
    rows: string[][]
}

function Table({ headers, rows }: TableProps) {
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
                <table className="w-full font-mono text-sm">
                    <thead>
                        <tr className="bg-gray-800/50">
                            {headers.map((header, index) => (
                                <th key={index} className="px-4 py-2 text-left text-orange-400 border-b border-orange-500/30">
                                    {header}
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {rows.map((row, rowIndex) => (
                            <tr key={rowIndex} className="hover:bg-gray-800/30 transition-colors">
                                {row.map((cell, cellIndex) => (
                                    <td key={cellIndex} className="px-4 py-2 text-gray-300 border-b border-gray-700/50">
                                        {cell}
                                    </td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    )
}

export function MarkdownRenderer({ content, className = "" }: MarkdownRendererProps) {
    const parseMarkdown = (text: string): React.ReactNode[] => {
        const lines = text.split("\n")
        const elements: React.ReactNode[] = []
        let i = 0

        while (i < lines.length) {
            const line = lines[i]

            // Skip empty lines
            if (line.trim() === "") {
                i++
                continue
            }

            // Headers
            if (line.startsWith("# ")) {
                elements.push(
                    <h1 key={i} className="text-3xl font-bold text-orange-400 mb-6 mt-8 border-b border-orange-500/30 pb-2">
                        {line.slice(2)}
                    </h1>,
                )
                i++
                continue
            }

            if (line.startsWith("## ")) {
                elements.push(
                    <h2 key={i} className="text-2xl font-bold text-orange-400 mb-4 mt-6">
                        {line.slice(3)}
                    </h2>,
                )
                i++
                continue
            }

            if (line.startsWith("### ")) {
                elements.push(
                    <h3 key={i} className="text-xl font-bold text-orange-400 mb-3 mt-5">
                        {line.slice(4)}
                    </h3>,
                )
                i++
                continue
            }

            if (line.startsWith("#### ")) {
                elements.push(
                    <h4 key={i} className="text-lg font-bold text-orange-400 mb-2 mt-4">
                        {line.slice(5)}
                    </h4>,
                )
                i++
                continue
            }

            if (line.startsWith("##### ")) {
                elements.push(
                    <h5 key={i} className="text-base font-bold text-orange-400 mb-2 mt-3">
                        {line.slice(6)}
                    </h5>,
                )
                i++
                continue
            }

            if (line.startsWith("###### ")) {
                elements.push(
                    <h6 key={i} className="text-sm font-bold text-orange-400 mb-2 mt-3">
                        {line.slice(7)}
                    </h6>,
                )
                i++
                continue
            }

            // Code blocks
            if (line.startsWith("```")) {
                const language = line.slice(3).trim()
                const codeLines: string[] = []
                i++

                while (i < lines.length && !lines[i].startsWith("```")) {
                    codeLines.push(lines[i])
                    i++
                }

                elements.push(<CodeBlock key={i} code={codeLines.join("\n")} language={language || undefined} />)
                i++
                continue
            }

            // Blockquotes
            if (line.startsWith("> ")) {
                const quoteLines: string[] = []
                let j = i

                while (j < lines.length && lines[j].startsWith("> ")) {
                    quoteLines.push(lines[j].slice(2))
                    j++
                }

                elements.push(
                    <Blockquote key={i}>
                        {quoteLines.map((quoteLine, index) => (
                            <p key={index} className="mb-2 last:mb-0">
                                {parseInlineElements(quoteLine)}
                            </p>
                        ))}
                    </Blockquote>,
                )
                i = j
                continue
            }

            // Unordered lists
            if (line.startsWith("- ") || line.startsWith("* ") || line.startsWith("+ ")) {
                const listItems: string[] = []
                let j = i

                while (
                    j < lines.length &&
                    (lines[j].startsWith("- ") || lines[j].startsWith("* ") || lines[j].startsWith("+ "))
                ) {
                    listItems.push(lines[j].slice(2))
                    j++
                }

                elements.push(
                    <ul key={i} className="list-none space-y-2 my-4 ml-4">
                        {listItems.map((item, index) => (
                            <li key={index} className="flex items-start gap-2">
                                <span className="text-orange-400 mt-1">•</span>
                                <span className="text-gray-300">{parseInlineElements(item)}</span>
                            </li>
                        ))}
                    </ul>,
                )
                i = j
                continue
            }

            // Ordered lists
            if (/^\d+\. /.test(line)) {
                const listItems: string[] = []
                let j = i
                let counter = 1

                while (j < lines.length && new RegExp(`^${counter}\\. `).test(lines[j])) {
                    listItems.push(lines[j].replace(/^\d+\. /, ""))
                    counter++
                    j++
                }

                elements.push(
                    <ol key={i} className="list-none space-y-2 my-4 ml-4">
                        {listItems.map((item, index) => (
                            <li key={index} className="flex items-start gap-2">
                                <span className="text-orange-400 mt-1 font-mono">{index + 1}.</span>
                                <span className="text-gray-300">{parseInlineElements(item)}</span>
                            </li>
                        ))}
                    </ol>,
                )
                i = j
                continue
            }

            // Tables
            if (line.includes("|")) {
                const tableLines: string[] = []
                let j = i

                while (j < lines.length && lines[j].includes("|")) {
                    tableLines.push(lines[j])
                    j++
                }

                if (tableLines.length >= 2) {
                    const headers = tableLines[0]
                        .split("|")
                        .map((h) => h.trim())
                        .filter((h) => h !== "")
                    const rows = tableLines
                        .slice(2) // Skip header separator
                        .map((row) =>
                            row
                                .split("|")
                                .map((cell) => cell.trim())
                                .filter((cell) => cell !== ""),
                        )

                    elements.push(<Table key={i} headers={headers} rows={rows} />)
                    i = j
                    continue
                }
            }

            // Horizontal rule
            if (line.trim() === "---" || line.trim() === "***" || line.trim() === "___") {
                elements.push(<hr key={i} className="border-orange-500/30 my-8" />)
                i++
                continue
            }

            // Regular paragraphs
            const paragraphLines: string[] = []
            let j = i

            while (
                j < lines.length &&
                lines[j].trim() !== "" &&
                !lines[j].startsWith("#") &&
                !lines[j].startsWith("```") &&
                !lines[j].startsWith("> ") &&
                !lines[j].startsWith("- ") &&
                !lines[j].startsWith("* ") &&
                !lines[j].startsWith("+ ") &&
                !/^\d+\. /.test(lines[j]) &&
                !lines[j].includes("|") &&
                lines[j].trim() !== "---" &&
                lines[j].trim() !== "***" &&
                lines[j].trim() !== "___"
            ) {
                paragraphLines.push(lines[j])
                j++
            }

            if (paragraphLines.length > 0) {
                elements.push(
                    <p key={i} className="text-gray-300 mb-4 leading-relaxed">
                        {parseInlineElements(paragraphLines.join(" "))}
                    </p>,
                )
            }

            i = j
        }

        return elements
    }

    const parseInlineElements = (text: string): React.ReactNode => {
        // Handle inline code first
        const codeRegex = /`([^`]+)`/g
        const parts: React.ReactNode[] = []
        let lastIndex = 0
        let match

        while ((match = codeRegex.exec(text)) !== null) {
            // Add text before the code
            if (match.index > lastIndex) {
                const beforeText = text.slice(lastIndex, match.index)
                parts.push(parseOtherInlineElements(beforeText))
            }

            // Add the inline code
            parts.push(<InlineCode key={match.index}>{match[1]}</InlineCode>)
            lastIndex = match.index + match[0].length
        }

        // Add remaining text
        if (lastIndex < text.length) {
            parts.push(parseOtherInlineElements(text.slice(lastIndex)))
        }

        return parts.length === 1 ? parts[0] : parts
    }

    const parseOtherInlineElements = (text: string): React.ReactNode => {
        // Bold
        text = text.replace(/\*\*(.*?)\*\*/g, '<strong class="text-white font-bold">$1</strong>')
        text = text.replace(/__(.*?)__/g, '<strong class="text-white font-bold">$1</strong>')

        // Italic
        text = text.replace(/\*(.*?)\*/g, '<em class="text-gray-200 italic">$1</em>')
        text = text.replace(/_(.*?)_/g, '<em class="text-gray-200 italic">$1</em>')

        // Strikethrough
        text = text.replace(/~~(.*?)~~/g, '<del class="text-gray-400 line-through">$1</del>')

        // Links
        const linkRegex = /\[([^\]]+)\]\(([^)]+)\)/g;
        text = text.replace(linkRegex, (_, linkText, url) => {
            const isExternal = url.startsWith("http") || url.startsWith("//")
            return `<a href="${url}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ""
                } class="text-orange-400 hover:text-orange-300 underline decoration-orange-500/50 hover:decoration-orange-300 transition-colors">${linkText}${isExternal
                    ? ' <svg class="inline w-3 h-3 ml-1" fill="currentColor" viewBox="0 0 20 20"><path d="M11 3a1 1 0 100 2h2.586l-6.293 6.293a1 1 0 101.414 1.414L15 6.414V9a1 1 0 102 0V4a1 1 0 00-1-1h-5z"></path><path d="M5 5a2 2 0 00-2 2v6a2 2 0 002 2h6a2 2 0 002-2v-2a1 1 0 10-2 0v2H5V7h2a1 1 0 000-2H5z"></path></svg>'
                    : ""
                }</a>`
        })

        return <span dangerouslySetInnerHTML={{ __html: text }} />
    }

    return <div className={`terminal-markdown ${className}`}>{parseMarkdown(content)}</div>
}
