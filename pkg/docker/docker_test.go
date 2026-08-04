package docker

import (
	"errors"
	"reflect"
	"testing"
)

func TestRemoveImagesIgnoresMissingImages(t *testing.T) {
	var removed []string
	images := []string{"missing:latest", "present:latest"}
	err := removeImages(images, func(image string) ([]byte, error) {
		removed = append(removed, image)
		if image == "missing:latest" {
			return []byte("Error response from daemon: No such image: missing:latest"), errors.New("exit status 1")
		}
		return nil, nil
	})

	if err != nil {
		t.Fatalf("RemoveImages returned an error for an already-absent image: %v", err)
	}

	if len(removed) != 2 || removed[0] != images[0] || removed[1] != images[1] {
		t.Fatalf("removed images = %v, want %v", removed, images)
	}
}

func TestRemoveImagesReportsGenuineFailures(t *testing.T) {
	err := removeImages([]string{"postgres:16-alpine"}, func(image string) ([]byte, error) {
		if image == "postgres:16-alpine" {
			return []byte("permission denied"), errors.New("exit status 1")
		}
		return nil, nil
	})
	if err == nil {
		t.Fatal("RemoveImages returned nil for a genuine Docker failure")
	}
	want := "failed to remove postgres:16-alpine: permission denied"
	if err.Error() != want {
		t.Fatalf("RemoveImages error = %q, want %q", err, want)
	}
}

func TestDockerGroupUser(t *testing.T) {
	tests := []struct {
		name     string
		sudoUser string
		username string
		euid     int
		want     string
	}{
		{name: "uses sudo user", sudoUser: "ubuntu", username: "root", euid: 0, want: "ubuntu"},
		{name: "skips sudo root", sudoUser: "root", username: "root", euid: 0, want: ""},
		{name: "uses current user", username: "justin", euid: 1000, want: "justin"},
		{name: "skips root without sudo user", username: "root", euid: 0, want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := dockerGroupUser(test.sudoUser, test.username, test.euid); got != test.want {
				t.Fatalf("dockerGroupUser() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestAddUserToDockerGroup(t *testing.T) {
	var gotName string
	var gotArgs []string
	if _, err := addUserToDockerGroup("ubuntu", func(name string, args ...string) ([]byte, error) {
		gotName = name
		gotArgs = args
		return nil, nil
	}); err != nil {
		t.Fatalf("addUserToDockerGroup() error = %v", err)
	}

	if gotName != "sudo" {
		t.Fatalf("command = %q, want %q", gotName, "sudo")
	}
	wantArgs := []string{"usermod", "-aG", "docker", "ubuntu"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("arguments = %q, want %q", gotArgs, wantArgs)
	}
}

func TestAddUserToDockerGroupPropagatesError(t *testing.T) {
	wantErr := errors.New("runner failed")
	if _, err := addUserToDockerGroup("ubuntu", func(string, ...string) ([]byte, error) {
		return nil, wantErr
	}); err != wantErr {
		t.Fatalf("addUserToDockerGroup() error = %v, want %v", err, wantErr)
	}
}
