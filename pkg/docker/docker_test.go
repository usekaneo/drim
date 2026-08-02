package docker

import (
	"errors"
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
