import { useState, useRef, useEffect } from "react"
import { ChevronDown } from "lucide-react"

interface Option {
    value: string
    label: string
}

interface SelectProps {
    value: string
    onChange: (value: string) => void
    options: Option[]
    disabled?: boolean
    className?: string
}

export function Select({ value, onChange, options, disabled = false, className = "" }: SelectProps) {
    const [isOpen, setIsOpen] = useState(false)
    const selectRef = useRef<HTMLDivElement>(null)

    const selectedOption = options.find((option) => option.value === value)

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (selectRef.current && !selectRef.current.contains(event.target as Node)) {
                setIsOpen(false)
            }
        }

        document.addEventListener("mousedown", handleClickOutside)
        return () => {
            document.removeEventListener("mousedown", handleClickOutside)
        }
    }, [])

    const handleSelect = (optionValue: string) => {
        onChange(optionValue)
        setIsOpen(false)
    }

    return (
        <div ref={selectRef} className={`relative ${className}`}>
            {/* Select Button */}
            <button
                type="button"
                onClick={() => !disabled && setIsOpen(!isOpen)}
                disabled={disabled}
                className={`
          w-full flex items-center justify-between px-4 py-2 
          bg-gray-800/50 border border-orange-500/30 rounded-lg
          text-white font-mono text-sm
          hover:border-orange-500 hover:bg-gray-800/70
          focus:outline-none focus:border-orange-500
          transition-all duration-200
          ${disabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"}
          ${isOpen ? "border-orange-500 bg-gray-800/70" : ""}
        `}
            >
                <span className="text-left">{selectedOption?.label || "Select..."}</span>
                <ChevronDown
                    className={`text-orange-400 transition-transform duration-200 ${isOpen ? "rotate-180" : ""}`}
                    size={16}
                />
            </button>

            {/* Dropdown Menu */}
            {isOpen && !disabled && (
                <div className="absolute top-full left-0 right-0 mt-1 z-50">
                    <div className="bg-gray-800 border border-orange-500/30 rounded-lg shadow-lg overflow-hidden">
                        {/* Terminal Header */}
                        <div className="flex items-center gap-2 px-3 py-2 bg-gray-900 border-b border-orange-500/30">
                            <div className="flex gap-1">
                                <div className="w-2 h-2 rounded-full bg-red-500"></div>
                                <div className="w-2 h-2 rounded-full bg-yellow-500"></div>
                                <div className="w-2 h-2 rounded-full bg-green-500"></div>
                            </div>
                            <span className="text-gray-400 text-xs font-mono">sort-options</span>
                        </div>

                        {/* Options */}
                        <div className="py-1">
                            {options.map((option) => (
                                <button
                                    key={option.value}
                                    type="button"
                                    onClick={() => handleSelect(option.value)}
                                    className={`
                    w-full px-4 py-2 text-left font-mono text-sm
                    hover:bg-gray-700 hover:text-orange-400
                    transition-colors duration-150
                    flex items-center gap-2
                    ${option.value === value ? "bg-gray-700 text-orange-400" : "text-gray-300"}
                  `}
                                >
                                    <span className="text-orange-400">$</span>
                                    <span>{option.label}</span>
                                    {option.value === value && <span className="ml-auto text-orange-400">●</span>}
                                </button>
                            ))}
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}
