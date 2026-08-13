package githubdocs

import (
	"context"

	"go.redsock.ru/rerrors"
)

// walkTree recursively lists path (and every subdirectory under it) via client.listDir,
// accumulating every "file" entry into files and every "dir" path into folders. One HTTP call
// per directory level — fine for the small, fixed docs/public tree; there's no depth/fanout
// guard since this walks a hardcoded repo path, not user-controlled input.
func walkTree(ctx context.Context, c *client, path string) (files []contentsEntry, folders []string, err error) {
	entries, err := c.listDir(ctx, path)
	if err != nil {
		return nil, nil, rerrors.Wrap(err, "error listing github directory "+path)
	}

	for _, entry := range entries {
		switch entry.Type {
		case "file":
			files = append(files, entry)
		case "dir":
			folders = append(folders, entry.Path)

			var (
				subFiles   []contentsEntry
				subFolders []string
			)

			subFiles, subFolders, err = walkTree(ctx, c, entry.Path)
			if err != nil {
				return nil, nil, err
			}

			files = append(files, subFiles...)
			folders = append(folders, subFolders...)
		}
	}

	return files, folders, nil
}
