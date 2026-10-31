package camera

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"
)

type mockRunner struct {
	called bool
	name   string
	args   []string
	err    error
}

func (m *mockRunner) Run(_ context.Context, name string, args ...string) error {
	m.called = true
	m.name = name
	m.args = args
	return m.err
}

func TestNewPiCamera(t *testing.T) {
	cam := NewPiCamera("/tmp/captures")

	if cam.outputDir != "/tmp/captures" {
		t.Errorf("expected outputDir '/tmp/captures', got %q", cam.outputDir)
	}
	if _, ok := cam.runner.(execRunner); !ok {
		t.Errorf("expected runner to be execRunner, got %T", cam.runner)
	}
}

func TestPiCamera_Capture(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRunner{}
		cam := &PiCamera{
			outputDir: "/tmp/captures",
			runner:    mock,
		}

		before := time.Now()
		snapshot, err := cam.Capture(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if !mock.called {
			t.Fatal("expected runner to be called")
		}
		if mock.name != "rpicam-still" {
			t.Errorf("expected command 'rpicam-still', got %q", mock.name)
		}
		if snapshot.EventID == "" {
			t.Error("expected EventID to be set")
		}
		if snapshot.CapturedAt.Before(before) {
			t.Errorf("expected CapturedAt after %v, got %v", before, snapshot.CapturedAt)
		}
	})

	t.Run("file path", func(t *testing.T) {
		mock := &mockRunner{}
		cam := &PiCamera{
			outputDir: "/tmp/captures",
			runner:    mock,
		}

		snapshot, err := cam.Capture(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if dir := filepath.Dir(snapshot.FilePath); dir != "/tmp/captures" {
			t.Errorf("expected file in /tmp/captures, got %q", dir)
		}
		name := filepath.Base(snapshot.FilePath)
		if !regexp.MustCompile(`^motion_capture_\d+\.jpg$`).MatchString(name) {
			t.Errorf("expected file name like motion_capture_<timestamp>.jpg, got %q", name)
		}
	})

	t.Run("command arguments", func(t *testing.T) {
		mock := &mockRunner{}
		cam := &PiCamera{
			outputDir: "/tmp/captures",
			runner:    mock,
		}

		snapshot, err := cam.Capture(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		want := []string{"-o", snapshot.FilePath, "-t", "5000", "--immediate"}
		if !slices.Equal(mock.args, want) {
			t.Errorf("expected args %q, got %q", want, mock.args)
		}
	})

	t.Run("runner error", func(t *testing.T) {
		runErr := errors.New("no cameras available")
		mock := &mockRunner{err: runErr}
		cam := &PiCamera{
			outputDir: "/tmp/captures",
			runner:    mock,
		}

		snapshot, err := cam.Capture(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, runErr) {
			t.Errorf("expected error to wrap %v, got: %v", runErr, err)
		}
		if !strings.Contains(err.Error(), "CSI camera hardware failed execution") {
			t.Errorf("expected wrapped error message, got: %v", err)
		}
		if snapshot.FilePath != "" || snapshot.EventID != "" || !snapshot.CapturedAt.IsZero() {
			t.Errorf("expected empty snapshot on error, got %+v", snapshot)
		}
	})
}
