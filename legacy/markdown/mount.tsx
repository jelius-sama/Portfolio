/*!
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright (c) 2026 Jelius Basumatary
 */

import { createRoot, type Root } from "react-dom/client"
import { MarkdownRenderer } from "./renderer"

if (!(window as any).__mdRendererLoaded) {
    (window as any).__mdRendererLoaded = true

    const roots = new WeakMap<Element, Root>()

    function mountOne(rootEl: HTMLElement) {
        if (roots.has(rootEl)) return

        const src = rootEl.querySelector<HTMLScriptElement>('script[data-markdown-source]')
        if (!src?.textContent) return
        const content = JSON.parse(src.textContent) as string

        rootEl.querySelector('[data-markdown-loading]')?.remove()

        const mount = rootEl.querySelector<HTMLElement>('[data-markdown-mount]')
        if (!mount) return // markup out of date — nothing to mount into

        const root = createRoot(mount)
        root.render(<MarkdownRenderer content={content} />)
        roots.set(rootEl, root)
    }

    function mountAll() {
        document.querySelectorAll<HTMLElement>("[data-markdown-root]").forEach(mountOne)
    }

    function unmount(el: Element) {
        const root = roots.get(el)
        if (root) {
            root.unmount()
            roots.delete(el)
        }
    }

    // Run immediately: whether this executes during initial parse (article
    // already precedes it in the markup) or is force-run by htmx after an
    // SPA swap (the whole fragment, article included, lands in the DOM
    // before htmx re-executes the script), the target already exists by
    // the time this line runs — no dependency on DOMContentLoaded, which
    // may have fired long before this script was even fetched.
    mountAll()

    document.addEventListener("htmx:afterSettle", mountAll)
    document.addEventListener("htmx:beforeCleanupElement", (e) => unmount(e.target as Element))
}
