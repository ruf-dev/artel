// Builds the user-message text sent to the agent: when files are attached, a single
// preamble line naming them (vault-relative paths, as returned by NotesAPI.ListNotes),
// then a blank line, then the user's text. No attachments => text unchanged. The
// agent reads the named files via its existing read_file MCP tool.
export function withAttachmentsPreamble(text: string, paths: string[]): string {
    return paths.length ? `[Attached vault files: ${paths.join(", ")}]\n\n${text}` : text
}
