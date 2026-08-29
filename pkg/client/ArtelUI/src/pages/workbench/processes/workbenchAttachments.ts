// Builds the user-message text sent to the agent: when files are attached, a single
// preamble line naming them (vault-relative paths, as returned by NotesAPI.ListNotes),
// then a blank line, then the user's text. No attachments => text unchanged. The
// agent reads the named files via its existing read_file MCP tool.
export function withAttachmentsPreamble(text: string, paths: string[]): string {
    return paths.length ? `[Attached vault files: ${paths.join(", ")}]\n\n${text}` : text
}

// Inverse of withAttachmentsPreamble: reconstructs the exact preamble that would
// have been built for `paths` and strips it off `text` as a literal prefix, for
// display purposes (e.g. UserMessageBubble rendering a sent message's caption
// without the mechanical "[Attached vault files: ...]" line the model actually
// received). Exact-prefix based rather than a generic regex strip, since `paths`
// is already known — safe and reversible. Returns `text` unchanged when there are
// no paths, or when the reconstructed preamble isn't a literal prefix of `text`
// (e.g. a legacy row persisted before this feature existed).
export function stripAttachmentsPreamble(text: string, paths: string[]): string {
    if (paths.length === 0) return text
    const preamble = withAttachmentsPreamble("", paths)
    return text.startsWith(preamble) ? text.slice(preamble.length) : text
}
