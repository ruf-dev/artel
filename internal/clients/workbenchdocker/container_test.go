package workbenchdocker

import (
	"reflect"
	"testing"

	"github.com/docker/go-connections/nat"
)

// TestContainerNetworkName exercises containerNetworkName, the seed for the dedicated
// per-container Docker network CreateContainer creates (see its doc comment) — a plain suffix
// derivation from CreateOpts.Name, mirrored here so a future change to the naming scheme has to
// touch this test too.
func TestContainerNetworkName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "typical workbench name",
			in:   "workbench-vault123-user456",
			want: "workbench-vault123-user456-net",
		},
		{
			name: "empty name",
			in:   "",
			want: "-net",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containerNetworkName(tt.in)

			if got != tt.want {
				t.Fatalf("containerNetworkName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestNewWorkbenchHostConfig exercises newWorkbenchHostConfig's isolation-hardening and mount
// fields directly, without needing a live Docker daemon — see CreateContainer's own doc comment
// on why these fields (CapDrop, SecurityOpt, PidsLimit) were added.
func TestNewWorkbenchHostConfig(t *testing.T) {
	portBindings := nat.PortMap{
		bridgeNatPort: []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "20000"}},
	}

	got := newWorkbenchHostConfig("workbench-vault123-user456", portBindings)

	if len(got.Mounts) != 1 {
		t.Fatalf("expected exactly one mount, got %d", len(got.Mounts))
	}

	mount := got.Mounts[0]
	if mount.Target != homeMountPath {
		t.Errorf("mount target = %q, want %q", mount.Target, homeMountPath)
	}

	if mount.Source != "workbench-vault123-user456" {
		t.Errorf("mount source = %q, want %q", mount.Source, "workbench-vault123-user456")
	}

	wantCapDrop := []string{"ALL"}
	if !reflect.DeepEqual([]string(got.CapDrop), wantCapDrop) {
		t.Errorf("CapDrop = %v, want %v", got.CapDrop, wantCapDrop)
	}

	wantSecurityOpt := []string{"no-new-privileges"}
	if !reflect.DeepEqual(got.SecurityOpt, wantSecurityOpt) {
		t.Errorf("SecurityOpt = %v, want %v", got.SecurityOpt, wantSecurityOpt)
	}

	if got.PidsLimit == nil {
		t.Fatal("PidsLimit is nil, want a set limit")
	}

	if *got.PidsLimit != 512 {
		t.Errorf("PidsLimit = %d, want 512", *got.PidsLimit)
	}

	if !reflect.DeepEqual(got.PortBindings, portBindings) {
		t.Errorf("PortBindings = %v, want %v", got.PortBindings, portBindings)
	}

	if _, ok := got.Tmpfs[envDropDir]; !ok {
		t.Errorf("Tmpfs missing envDropDir entry %q: %v", envDropDir, got.Tmpfs)
	}
}

// TestFilterRemovableNetworks exercises the exclusion logic RemoveContainer relies on to never
// delete the old shared workbenchNetworkName even if it's found attached to a container created
// before the per-container-network change — see filterRemovableNetworks's doc comment.
func TestFilterRemovableNetworks(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "only the dedicated per-container network",
			in:   []string{"workbench-vault123-user456-net"},
			want: []string{"workbench-vault123-user456-net"},
		},
		{
			name: "shared network is excluded",
			in:   []string{workbenchNetworkName},
			want: []string{},
		},
		{
			name: "shared network alongside a dedicated one — only the dedicated one is removable",
			in:   []string{workbenchNetworkName, "workbench-vault123-user456-net"},
			want: []string{"workbench-vault123-user456-net"},
		},
		{
			name: "no networks",
			in:   []string{},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterRemovableNetworks(tt.in)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("filterRemovableNetworks(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// TestHomeMountPath pins homeMountPath's value — the const workspaceMountPath was renamed to, per
// the move from a bare /workspace mount to a vault subdirectory under the container's $HOME (see
// client.go's doc comment). A regression here breaks every file already written into workbench
// volumes under the old path.
func TestHomeMountPath(t *testing.T) {
	const want = "/root/vault"

	if homeMountPath != want {
		t.Fatalf("homeMountPath = %q, want %q", homeMountPath, want)
	}
}
