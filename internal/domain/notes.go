package domain

// ImportAction is how CommitImport resolves a zip entry that CheckImportConflicts previously
// reported as colliding with an existing note/file.
type ImportAction int

const (
	ImportActionSkip ImportAction = iota
	ImportActionOverwrite
	ImportActionRename
)

type ImportResolution struct {
	Path     string // conflicting destination path, as reported by CheckImportConflicts
	Action   ImportAction
	RenameTo string // vault-relative rename target; only used when Action == ImportActionRename
}
