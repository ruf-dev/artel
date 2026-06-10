package couchdb

// SetupStatus reports which initial setup steps have been applied to a CouchDB instance.
type SetupStatus struct {
	ClusterModeEnabled      bool
	UsersDbInitialized      bool
	ReplicatorDbInitialized bool
}

// Public document models returned by LiveSyncClient methods.

type NoteEntry struct {
	Path  string
	Mtime int64
}

type NoteDoc struct {
	Id      string
	Rev     string
	Content string
	Mtime   int64
	Ctime   int64
	Size    int64
	Deleted bool
}

type FileEntry struct {
	Path     string
	Mtime    int64
	MimeType string
}

type FileDoc struct {
	Id       string
	Rev      string
	RawBytes []byte
	MimeType string
	Mtime    int64
	Ctime    int64
	Size     int64
	Deleted  bool
}

// Internal CouchDB document shapes used for JSON marshaling.

type docScan struct {
	Type    string `json:"type"`
	Deleted bool   `json:"deleted"`
	Mtime   int64  `json:"mtime"`
}

type docFull struct {
	Id       string   `json:"_id"`
	Rev      string   `json:"_rev"`
	Type     string   `json:"type"`
	Data     string   `json:"data"`
	Children []string `json:"children"`
	Mtime    int64    `json:"mtime"`
	Ctime    int64    `json:"ctime"`
	Size     int64    `json:"size"`
	Deleted  bool     `json:"deleted"`
}

type docRev struct {
	Rev string `json:"_rev"`
}

type noteWrite struct {
	Id       string   `json:"_id"`
	Rev      string   `json:"_rev,omitempty"`
	Data     string   `json:"data"`
	Children []string `json:"children"`
	Mtime    int64    `json:"mtime"`
	Ctime    int64    `json:"ctime"`
	Size     int64    `json:"size"`
	Type     string   `json:"type"`
}

type docDelete struct {
	Id       string   `json:"_id"`
	Rev      string   `json:"_rev"`
	Children []string `json:"children"`
	Type     string   `json:"type"`
	Deleted  bool     `json:"deleted"`
}

type leafDoc struct {
	Data string `json:"data"`
}

type docType struct {
	Type string `json:"type"`
}

type fileWrite struct {
	Id       string   `json:"_id"`
	Data     string   `json:"data"`
	Children []string `json:"children"`
	Mtime    int64    `json:"mtime"`
	Ctime    int64    `json:"ctime"`
	Size     int64    `json:"size"`
	Type     string   `json:"type"`
}
