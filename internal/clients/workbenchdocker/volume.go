package workbenchdocker

import (
	"context"

	"github.com/docker/docker/api/types/volume"

	"go.redsock.ru/rerrors"
)

// CreateVolume creates a named Docker volume, labeled for operational visibility.
func (c *Client) CreateVolume(ctx context.Context, name string) error {
	labels := map[string]string{
		workbenchLabelKey: workbenchLabelValue,
	}

	createOptions := volume.CreateOptions{
		Name:   name,
		Labels: labels,
	}

	_, err := c.cli.VolumeCreate(ctx, createOptions)
	if err != nil {
		return rerrors.Wrap(err, "creating workbench volume")
	}

	return nil
}

// RemoveVolume force-removes a named Docker volume.
func (c *Client) RemoveVolume(ctx context.Context, name string) error {
	const force = true

	err := c.cli.VolumeRemove(ctx, name, force)
	if err != nil {
		return rerrors.Wrap(err, "removing workbench volume")
	}

	return nil
}
