import {Path} from "@/app/routing/Router.tsx"

const DEFAULT_NOTE_EXTENSION = ".md"

export function encodeNotePath(path: string): string {
    const displayPath = path.endsWith(DEFAULT_NOTE_EXTENSION)
        ? path.slice(0, -DEFAULT_NOTE_EXTENSION.length)
        : path
    return displayPath.split("/").map(encodeURIComponent).join("/")
}

export function decodeNotePath(raw: string): string {
    const decoded = raw.split("/").map(decodeURIComponent).join("/")
    const filename = decoded.split("/").pop() ?? decoded
    return filename.includes(".") ? decoded : decoded + DEFAULT_NOTE_EXTENSION
}

export function buildNotesUrl(vaultId: string | null, notePath: string | null): string {
    if (!vaultId) return Path.NotesPage
    if (!notePath) return `${Path.NotesPage}/${vaultId}`
    return `${Path.NotesPage}/${vaultId}/${encodeNotePath(notePath)}`
}
