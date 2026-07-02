import { useEffect, useState } from "react"

export type Theme = "dark" | "light"

const STORAGE_KEY = "artel.theme"

function applyTheme(theme: Theme) {
    if (theme === "light") {
        document.documentElement.setAttribute("data-theme", "light")
    } else {
        document.documentElement.removeAttribute("data-theme")
    }
}

export function useTheme() {
    const [theme, setTheme] = useState<Theme>(
        () => (localStorage.getItem(STORAGE_KEY) === "light" ? "light" : "dark")
    )

    useEffect(() => {
        applyTheme(theme)
        localStorage.setItem(STORAGE_KEY, theme)
    }, [theme])

    function toggleTheme() {
        setTheme(t => (t === "light" ? "dark" : "light"))
    }

    return { theme, toggleTheme }
}
