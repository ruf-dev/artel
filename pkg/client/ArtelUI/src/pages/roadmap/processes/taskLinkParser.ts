// Pure parsing helpers for the roadmap graph. A "task link" is a line inside a Trello card
// comment written by AddTaskLinkDialog (or by hand) in the form "<Relation>: <url>" — see
// dialogs/AddTaskLinkDialog for the writer side. roadmapGraph.ts's expandNode() is the reader
// side: it runs every loaded card's comments through parseTaskLinks() to discover edges.

export type RoadmapRelation = "depends_on" | "blocks" | "blocked_by" | "related_to"

export interface ParsedTaskLink {
    relation: RoadmapRelation
    url: string
    /** Trello card shortLink extracted from `url`, if it matched the trello.com/c/<shortLink> shape. */
    shortLink?: string
}

const LINK_LINE_RE = /^(Depends on|Blocks|Blocked by|Related to):\s*(\S+)/i

const RELATION_BY_LABEL: Record<string, RoadmapRelation> = {
    "depends on": "depends_on",
    "blocks": "blocks",
    "blocked by": "blocked_by",
    "related to": "related_to",
}

const TRELLO_CARD_URL_RE = /trello\.com\/c\/([a-zA-Z0-9]+)/

/** Pulls the shortLink segment out of a Trello card URL (e.g. "https://trello.com/c/abC123/9-x"
 * -> "abC123"). Returns undefined for anything that isn't a recognizable Trello card URL. */
export function extractTrelloShortLink(url: string): string | undefined {
    const match = TRELLO_CARD_URL_RE.exec(url)
    return match?.[1]
}

/** Scans a comment body line-by-line for "<Relation>: <url>" markers and returns every one it
 * recognizes. Unrelated lines (the rest of a human comment) are ignored. */
export function parseTaskLinks(text: string): ParsedTaskLink[] {
    const links: ParsedTaskLink[] = []

    for (const rawLine of text.split("\n")) {
        const match = LINK_LINE_RE.exec(rawLine.trim())
        if (!match) continue

        const relation = RELATION_BY_LABEL[match[1].toLowerCase()]
        if (!relation) continue

        const url = match[2]
        links.push({relation, url, shortLink: extractTrelloShortLink(url)})
    }

    return links
}
