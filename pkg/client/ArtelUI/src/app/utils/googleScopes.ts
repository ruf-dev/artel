export const SCOPE_INFO: Record<string, string> = {
    "spreadsheets": "Read and write your Google Sheets spreadsheets",
    "spreadsheets.readonly": "Read-only access to your Google Sheets",
    "drive.file": "Access files Artel creates or opens in Google Drive",
    "drive.readonly": "Read-only access to your Google Drive files",
    "gmail.readonly": "Read your Gmail messages",
    "calendar.readonly": "Read your Google Calendar events",
    "calendar": "Read and write your Google Calendar events",
}

export function trimScope(s: string): string {
    const idx = s.lastIndexOf("/")
    return idx >= 0 ? s.slice(idx + 1) : s
}

export function parseScopeList(scopes: string): string[] {
    return scopes.split(/[ ,]+/).filter(Boolean)
}
