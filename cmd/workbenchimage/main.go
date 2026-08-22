// SPDX-License-Identifier: AGPL-3.0-only

// Command workbenchimage prints the content-hash tag the backend's EnsureImage
// (internal/clients/workbenchdocker/image.go) will look for on the Docker daemon, so a local
// `docker build` can be tagged to match it exactly — see `make workbench-image`.
package main

import (
	"fmt"
	"os"

	"github.com/ruf-dev/artel/internal/clients/workbenchdocker"
)

func main() {
	tag, err := workbenchdocker.WorkbenchImageTag()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err = fmt.Fprintln(os.Stdout, tag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
