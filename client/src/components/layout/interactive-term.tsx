import { useEffect, useState, useRef } from "react"
import { TerminalWindow } from "@/components/ui/terminal-window"
import { useConfig } from "@/contexts/config"
import { useNavigate } from "react-router-dom"
import { cn } from "@/lib/utils"

type CommandHandler = {
    requiresSudo: boolean
    run: (
        stdin: () => Promise<string>,
        stdout: (line: string) => void,
        ctx: { sudoToken: string }
    ) => Promise<void>
}

export function InteractiveTerminal() {
    const { app: { portfolio: me } } = useConfig()

    const hasMountedRef = useRef(false)
    const inputRef = useRef<HTMLInputElement>(null)
    const scrollableContentRef = useRef<HTMLDivElement>(null)
    const neofetchRef = useRef<HTMLDivElement>(null)
    const inputWrapperRef = useRef<HTMLDivElement>(null)
    const initialMaxHeight = useRef<string | null>(null)

    const navigate = useNavigate()
    const [command, setCommand] = useState("")
    const [output, setOutput] = useState<string[]>([])
    const [maxHeight, setMaxHeight] = useState("auto")
    const [inputQueue, setInputQueue] = useState<((input: string) => void) | null>(null)
    const [pendingSudoCommand, setPendingSudoCommand] = useState<string | null>(null)
    const [pendingStdin, setPendingStdin] = useState<boolean>(false)

    const stdout = (line: string) => setOutput(prev => [...prev, line])

    const verifySudoPassword = async (token: string): Promise<boolean> => {
        try {
            const response = await fetch("/api/sudo", {
                method: "POST",
                headers: { Authorization: `Bearer ${token}` }
            })
            return response.ok
        } catch {
            return false
        }
    }

    const COMMANDS: Record<string, CommandHandler> = {
        "help": {
            requiresSudo: false,
            async run(_, stdout) {
                if (pendingSudoCommand === "help") {
                    stdout("Root-only commands:")
                    Object.entries(COMMANDS).forEach(([key, val]) => {
                        if (val.requiresSudo) stdout(`-> ${key}`)
                    })
                } else {
                    stdout("Usage:")
                    Object.entries(COMMANDS).forEach(([key, val]) => {
                        if (!val.requiresSudo) stdout(`-> ${key}`)
                    })

                }
            }
        },

        "version": {
            requiresSudo: false,
            async run(_, stdout) {
                try {
                    const res = await fetch("/api/version")
                    const data = await res.json()
                    stdout(`[INFO] Version: ${data.version}`)
                } catch {
                    stdout("[ERROR] Failed to fetch version")
                }
            }
        },

        "clear": {
            requiresSudo: false,
            async run() {
                setOutput([])
                if (neofetchRef.current) neofetchRef.current.innerHTML = ""
            }
        },

        "neofetch": {
            requiresSudo: false,
            async run(_, stdout) {
                stdout("[ERROR] Not enough screen real estate. Try piping to less.")
            }
        },

        "neofetch | less": {
            requiresSudo: false,
            async run() {
                setOutput([])
                if (neofetchRef.current) {
                    neofetchRef.current.innerHTML = `
                        <div class="text-orange-400 mb-2">$ neofetch | less</div>
                        <div class="flex w-full place-content-center gap-6 mb-6">
                            <div class="flex-shrink-0">
                                <div class="w-64 h-64 bg-[#1E1E1E] border border-orange-500/30 rounded flex items-center justify-center">
                                    ${me.links["jelius.dev"]?.qr_code_link
                            ? `<img src="${me.links["jelius.dev"]?.qr_code_link}" alt="QR Code" class="w-58 h-58 rounded object-cover" />`
                            : `<span class="text-gray-500">QR Code not available</span>`}
                                </div>
                            </div>
                        </div>
                    `
                }
            }
        },

        "update-server": {
            requiresSudo: true,
            async run(_, stdout, ctx) {
                const success = await fetch("/api/update_server", {
                    method: "POST",
                    headers: { Authorization: `Bearer ${ctx.sudoToken}` }
                }).then(r => r.ok).catch(() => false)

                stdout(success ? "[SUCCESS] Server will be updated shortly..." : "[ERROR] Failed to schedule server update.")
            }
        },

        "analytics": {
            requiresSudo: true,
            async run(_, stdout, ctx) {
                const response = await fetch("/api/authenticate", {
                    method: "POST",
                    headers: { Authorization: `Bearer ${ctx.sudoToken}` }
                })
                const success = response.ok

                stdout(success ? "[SUCCESS] You will be redirected soon..." : "[ERROR] Failed to authenticate.")
                if (success) {
                    navigate(`/analytics`)
                }
            }
        },

        "purge-cache": {
            requiresSudo: true,
            async run(stdin, stdout, ctx) {
                setPendingStdin(true)
                stdout("Enter relative paths(separated by comma):")
                try {
                    const rawDawg = await stdin()
                    const paths = rawDawg.split(",").map(p => p.trim());

                    const response = await fetch("/api/purge_cache", {
                        method: "POST",
                        headers: {
                            Authorization: `Bearer ${ctx.sudoToken}`,
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify({ paths }),
                    });

                    if (!response.ok) {
                        const error = await response.json();
                        console.error("Failed to purge cache:", error);
                        stdout("[ERROR] Failed to purge cache.")
                        return
                    }

                    stdout("[SUCCESS] Successfully to purge cache.")
                } catch (err) {
                    console.error("Error purging cache:", err);
                    stdout("[ERROR] Failed to purge cache.")
                } finally {
                    setPendingStdin(false)
                }
            }
        },

        "purge-all-cache": {
            requiresSudo: true,
            async run(_, stdout, ctx) {
                setPendingStdin(true)
                try {
                    const response = await fetch("/api/purge_all_cache", {
                        method: "POST",
                        headers: {
                            Authorization: `Bearer ${ctx.sudoToken}`,
                            "Content-Type": "application/json",
                        },
                    });

                    if (!response.ok) {
                        const error = await response.json();
                        console.error("Failed to purge all cache:", error);
                        stdout("[ERROR] Failed to purge all cache.")
                        return;
                    }

                    stdout("[SUCCESS] Successfully purged all cache.")
                } catch (err) {
                    console.error("Error purging all cache:", err);
                    stdout("[ERROR] Failed to purge all cache.")
                } finally {
                    setPendingStdin(false);
                }
            }
        },

        "post-blog": {
            requiresSudo: true,
            async run(stdin, stdout, ctx) {
                setPendingStdin(true)
                stdout("Enter blog title:")
                const title = await stdin()

                stdout("Enter blog summary:")
                const summary = await stdin()

                stdout("Select markdown file to upload...")

                setPendingStdin(false)

                const file = await new Promise<File | null>((resolve) => {
                    const input = document.createElement("input");
                    input.type = "file";
                    input.accept = ".md";

                    let resolved = false;

                    const handleChange = () => {
                        resolved = true;
                        cleanup();
                        resolve(input.files?.[0] || null);
                    };

                    const handleFocus = () => {
                        if (!resolved) {
                            // Give the file dialog a bit of time to trigger onchange before resolving
                            requestAnimationFrame(() => {
                                if (!resolved) {
                                    cleanup();
                                    resolve(null); // user cancelled
                                }
                            });
                        }
                    };

                    const cleanup = () => {
                        input.removeEventListener("change", handleChange);
                        window.removeEventListener("focus", handleFocus);
                    };

                    input.addEventListener("change", handleChange);
                    window.addEventListener("focus", handleFocus, { once: true });

                    input.click();
                });

                if (!file) {
                    stdout("[ERROR] File selection canceled.")
                    return
                }

                const formData = new FormData()
                formData.append("title", title)
                formData.append("summary", summary)
                formData.append("content_file", file)
                // TODO: Support `prequel`, `sequel` & `parts`.

                const headers = new Headers()
                headers.set("Authorization", `Bearer ${ctx.sudoToken}`)

                const success = await fetch("/api/blog", {
                    method: "POST",
                    headers,
                    body: formData
                }).then(r => r.ok).catch(() => false)

                stdout(success ? "[SUCCESS] Blog post uploaded successfully." : "[ERROR] Failed to upload blog post.")
            }
        },
    }

    const getStdin = (): (() => Promise<string>) => {
        return () => new Promise(resolve => {
            setInputQueue(() => {
                return (input: string) => {
                    setOutput(prev => [...prev, `-> ${input}`])
                    setInputQueue(null)
                    resolve(input)
                }
            })
        })
    }

    const executeCommand = async (raw: string) => {
        const cmd = raw.trim()
        if (!cmd) return

        if (inputQueue) {
            inputQueue(cmd)
            return
        }

        // Handling sudo flow
        if (pendingSudoCommand) {
            const verified = await verifySudoPassword(cmd)
            if (!verified) {
                stdout("[ERROR] Sorry, try again.")
                return
            }

            const handler = COMMANDS[pendingSudoCommand]
            await handler.run(getStdin(), stdout, { sudoToken: cmd })
            setPendingSudoCommand(null)
            return
        }

        stdout(`$ ${cmd}`)

        if (cmd.startsWith("sudo ")) {
            if (cmd.endsWith("help")) {
                setPendingSudoCommand("help")
                return
            }

            const baseCmd = cmd.slice(5).trim()
            const handler = COMMANDS[baseCmd]
            if (!handler) {
                stdout(`[ERROR] sudo: ${baseCmd}: command not found`)
                return
            }
            if (!handler.requiresSudo) {
                stdout(`[INFO] sudo: ${baseCmd}: command does not require sudo`)
                return
            }

            setPendingSudoCommand(baseCmd)
            return
        }

        const handler = COMMANDS[cmd]
        if (!handler) {
            stdout(`[ERROR] Command not found: ${cmd}`)
            return
        }

        if (handler.requiresSudo) {
            stdout(`[ERROR] not enough privileges. Try using 'sudo ${cmd}'.`)
            return
        }

        await handler.run(getStdin(), stdout, { sudoToken: "" })
    }

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            e.preventDefault()

            switch (command) {
                // INFO: Easter eggs
                case "music":
                    setCommand("")
                    sessionStorage.setItem('did_come_from_terminal', 'yes');
                    navigate("/music")
                    break

                default:
                    void executeCommand(command)
                    setCommand("")
                    break
            }
        }

        if (e.ctrlKey && e.key.toLowerCase() === "c") {
            e.preventDefault()
            setPendingSudoCommand(null)
            setInputQueue(null)
            stdout("^C")
            setCommand("")
        }
    }

    useEffect(() => {
        if (!hasMountedRef.current) {
            hasMountedRef.current = true
            return
        }

        if (scrollableContentRef.current) {
            scrollableContentRef.current.scrollTop = scrollableContentRef.current.scrollHeight
        }
    }, [output])

    useEffect(() => {
        if (neofetchRef.current && inputWrapperRef.current && !initialMaxHeight.current) {
            const totalHeight = `${neofetchRef.current.offsetHeight + inputWrapperRef.current.offsetHeight}px`
            initialMaxHeight.current = totalHeight
            setMaxHeight(totalHeight)
        }
    }, [])

    return (
        <TerminalWindow title="terminal">
            <div
                ref={scrollableContentRef}
                className="overflow-y-auto pr-2"
                style={{ height: `calc(${maxHeight} + calc(var(--spacing) * 6))` }}
            >
                <div id="neofetch" ref={neofetchRef} className="font-mono text-sm">
                    <div className="text-orange-400 mb-2">$ neofetch | less</div>
                    <div className="flex w-full place-content-center gap-6 mb-6">
                        <div className="flex-shrink-0">
                            <div className="w-64 h-64 bg-[#1E1E1E] border border-orange-500/30 rounded flex items-center justify-center">
                                {me.links["jelius.dev"]?.qr_code_link ? (
                                    <img
                                        src={me.links["jelius.dev"]?.qr_code_link}
                                        id="slow-af"
                                        alt="QR Code"
                                        className="w-58 h-58 rounded object-cover"
                                    />
                                ) : (
                                    <span className="text-gray-500">QR Code not available</span>
                                )}
                            </div>
                        </div>
                    </div>
                </div>

                {output.map((line, index) => (
                    <div
                        key={index}
                        className={cn(
                            line.startsWith("$ ")
                                ? "text-orange-400"
                                : line.startsWith("-> ")
                                    ? "text-cyan-400"
                                    : line.startsWith("[ERROR]")
                                        ? "text-red-400"
                                        : line.startsWith("[SUCCESS]")
                                            ? "text-green-400"
                                            : line.startsWith("[INFO]")
                                                ? "text-blue-400"
                                                : "text-gray-300"
                        )}
                    >
                        {line}
                    </div>
                ))}

                <div ref={inputWrapperRef} className="flex items-center mt-4">
                    {!pendingStdin && (
                        <span className="text-orange-400 mr-2">
                            {pendingSudoCommand ? "[sudo] password for user:" : "$"}
                        </span>
                    )}
                    <input
                        ref={inputRef}
                        type="text"
                        value={pendingSudoCommand ? "∙".repeat(command.length) : command}
                        onChange={(e) => setCommand(e.target.value)}
                        onKeyDown={handleKeyDown}
                        className="flex-1 bg-transparent text-white outline-none caret-orange-400"
                        spellCheck="false"
                        autoCapitalize="off"
                        autoComplete="off"
                        aria-label="Terminal command input"
                    />
                </div>
            </div>
        </TerminalWindow>
    )
}
