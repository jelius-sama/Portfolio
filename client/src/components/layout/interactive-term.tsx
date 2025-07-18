import { useEffect, useState, useRef } from "react"
import { TerminalWindow } from "@/components/layout/terminal-window"
import { useConfig } from "@/contexts/config"
import { cn } from "@/lib/utils"

export function InteractiveTerminal() {
    const { app: { portfolio: me } } = useConfig()

    const hasMountedRef = useRef(false)
    const [command, setCommand] = useState("")
    const [output, setOutput] = useState<string[]>([])
    const inputRef = useRef<HTMLInputElement>(null)
    const scrollableContentRef = useRef<HTMLDivElement>(null) // Ref for the new scrollable div
    const neofetchRef = useRef<HTMLDivElement>(null);
    const inputWrapperRef = useRef<HTMLDivElement>(null);
    const initialMaxHeight = useRef<string | null>(null);
    const [maxHeight, setMaxHeight] = useState("auto");
    const [sudoPassword, setSudoPassword] = useState<string | null>(null);

    const verifySudoPassword = async (token: string): Promise<boolean> => {
        try {
            const response = await fetch("/api/sudo", {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.ok;
        } catch (error) {
            return false;
        }
    };

    const updateServer = async (token: string): Promise<boolean> => {
        try {
            const response = await fetch("/api/update_server", {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.ok;
        } catch {
            return false;
        }
    };

    const [awaitingPassword, setAwaitingPassword] = useState(false);
    const [pendingSudoCommand, setPendingSudoCommand] = useState<string | null>(null);

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
        if (
            neofetchRef.current &&
            inputWrapperRef.current &&
            !initialMaxHeight.current
        ) {
            const neofetchHeight = neofetchRef.current.offsetHeight;
            const inputHeight = inputWrapperRef.current.offsetHeight;
            const totalHeight = `${neofetchHeight + inputHeight}px`;

            initialMaxHeight.current = totalHeight;
            setMaxHeight(totalHeight); // Sets height, not maxHeight
        }
    }, []);

    const executeCommand = async (cmd: string) => {
        const lower = cmd.toLowerCase();

        // Handle sudo password input
        if (awaitingPassword) {
            setAwaitingPassword(false);
            setSudoPassword(cmd);

            const isValid = await verifySudoPassword(cmd || "");

            if (!isValid) {
                setOutput((prev) => [
                    ...prev,
                    "Sorry, try again.",
                ]);
                setPendingSudoCommand(null);
                return;
            }

            if (pendingSudoCommand === "help") {
                setOutput((prev) => [
                    ...prev,
                    "[sudo] password accepted",
                    "Super User help options:",
                    "- update-server: Updates Server to latest available version.",
                ]);
            } else if (pendingSudoCommand === "update-server") {
                setOutput((prev) => [...prev, "[sudo] password accepted"]);

                const success = await updateServer(sudoPassword || "");

                if (success) {
                    setOutput((prev) => [
                        ...prev,
                        "Server will be updated shortly...",
                    ]);
                } else {
                    setOutput((prev) => [
                        ...prev,
                        "Failed to schedule server update.",
                    ]);
                }
            } else {
                setOutput((prev) => [
                    ...prev,
                    "[sudo] password accepted",
                    `sudo: ${pendingSudoCommand}: command not found`,
                ]);
            }

            setPendingSudoCommand(null);
            return;
        }

        // Handle sudo prefix
        if (lower.startsWith("sudo ")) {
            const sudoCmd = cmd.slice(5).trim();
            setPendingSudoCommand(sudoCmd);
            setAwaitingPassword(true);
            return;
        }

        switch (lower) {
            case "update-server":
                setOutput((prev) => [
                    ...prev,
                    "Error: not enough privileges. Try using 'sudo update-server'.",
                ]);
                break;

            case "version":
                try {
                    const res = await fetch("/api/version");
                    if (!res.ok) throw new Error();
                    const data = await res.json();
                    setOutput((prev) => [...prev, `Version: ${data.version}`]);
                } catch {
                    setOutput((prev) => [...prev, "Failed to fetch version"]);
                }
                break;

            case "help":
                setOutput((prev) => [
                    ...prev,
                    "Usage:",
                    "- help: Show this help message",
                    "- version: Show the current server version",
                    "- clear: Clear the terminal output",
                    "- neofetch: Show system info",
                    "More commands will be implemented soon..."
                ]);
                break;

            case "clear":
                setOutput([]);
                if (neofetchRef.current) {
                    neofetchRef.current.innerHTML = "";
                }
                break;

            case "neofetch":
                setOutput((prev) => [
                    ...prev,
                    "Error: Not enough screen real estate. Try piping to less.",
                ]);
                break;

            case "neofetch | less":
                setOutput([]);
                if (neofetchRef.current) {
                    neofetchRef.current.innerHTML = `
                    <div class="text-orange-400 mb-2">$ neofetch | less</div>
                    <div class="flex w-full place-content-center gap-6 mb-6">
                        <div class="flex-shrink-0">
                            <div class="w-64 h-64 bg-[#1E1E1E] border border-orange-500/30 rounded flex items-center justify-center">
                                ${me.links["jelius.dev"]?.qr_code_link
                            ? `<img src="${me.links["jelius.dev"]?.qr_code_link}" alt="QR Code to portfolio page" class="w-58 h-58 rounded object-cover" />`
                            : `<span class="text-gray-500">QR Code not available</span>`
                        }
                            </div>
                        </div>
                    </div>
                `;
                }
                break;

            default:
                setOutput((prev) => [...prev, `Command not found: ${cmd}`]);
                break;
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            e.preventDefault();
            const trimmedCommand = command.trim();

            if (trimmedCommand) {
                if (!awaitingPassword) {
                    setOutput((prev) => [...prev, `$ ${trimmedCommand}`]);
                }
                void executeCommand(trimmedCommand);
            }

            setCommand("");
        }

        if (e.ctrlKey && e.key.toLowerCase() === "c") {
            e.preventDefault();
            setAwaitingPassword(false);
            setPendingSudoCommand(null);
            setOutput((prev) => [...prev, "^C"]);
            setCommand("");
        }
    };

    return (
        <TerminalWindow title="terminal">
            <div ref={scrollableContentRef} className="overflow-y-auto pr-2" style={{ height: `calc(${maxHeight} + calc(var(--spacing) * 6))` }}>
                <div id="neofetch" ref={neofetchRef} className="font-mono text-sm">
                    <div className="text-orange-400 mb-2">$ neofetch | less</div>

                    <div className="flex w-full place-content-center gap-6 mb-6">
                        <div className="flex-shrink-0">
                            <div className="w-64 h-64 bg-[#1E1E1E] border border-orange-500/30 rounded flex items-center justify-center">
                                {me.links["jelius.dev"]?.qr_code_link ? (
                                    <img
                                        src={me.links["jelius.dev"]?.qr_code_link}
                                        alt="QR Code to portfolio page"
                                        className="w-58 h-58 rounded object-cover"
                                    />
                                ) : (
                                    <span className="text-gray-500">QR Code not available</span>
                                )}
                            </div>
                        </div>
                    </div>
                </div>

                {/* Output History */}
                {output.map((line, index) => (
                    <div key={index} className={cn(line.startsWith("$ ") ? "text-orange-400" : "text-gray-300")}>
                        {line}
                    </div>
                ))}

                {/* Terminal Input */}
                <div ref={inputWrapperRef} className="flex items-center mt-4">
                    <span className="text-orange-400 mr-2">{awaitingPassword ? "[sudo] password for user:" : "$"}</span>
                    <input
                        ref={inputRef}
                        type="text"
                        value={command}
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
