import { createRoot } from "react-dom/client"
import { MarkdownRenderer } from "./renderer"

if (!(window as any).__mdRendererLoaded) {
    ; (window as any).__mdRendererLoaded = true

    const roots = new WeakMap<Element, ReturnType<typeof createRoot>>()

    function mountAll() {
        document.querySelectorAll<HTMLElement>("[data-markdown-root]").forEach((el) => {
            if (roots.has(el)) return
            const src = el.querySelector('script[data-markdown-source]')
            if (!src?.textContent) return
            const content = JSON.parse(src.textContent) as string
            const root = createRoot(el)
            root.render(<MarkdownRenderer content={content} />)
            roots.set(el, root)
        })
    }

    function unmount(el: Element) {
        const root = roots.get(el)
        if (root) { root.unmount(); roots.delete(el) }
    }

    document.addEventListener("DOMContentLoaded", mountAll)
    document.addEventListener("htmx:afterSettle", mountAll)
    document.addEventListener("htmx:beforeCleanupElement", (e) => unmount(e.target as Element))
}
