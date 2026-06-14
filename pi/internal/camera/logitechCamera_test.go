package camera

import (
	"context"
	"errors"
	"strings"
	"testing"
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

func TestLogitechCamera_Capture(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRunner{}
		cam := &LogitechCamera{
			devicePath: "/dev/video",
			outputDir:  "/tmp/captures",
			runner:     mock,
		}

		snapshot, err := cam.Capture(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if !mock.called {
			t.Fatal("expected runner to be called")
		}
		if mock.name != "fswebcam" {
			t.Errorf("expected command 'fswebcam', got %q", mock.name)
		}
		if snapshot.EventID == "" {
			t.Error("expected EventID to be set")
		}
		if !strings.HasPrefix(snapshot.FilePath, "/tmp/captures/") {
			t.Errorf("expected FilePath under /tmp/captures/, got %q", snapshot.FilePath)
		}
		if snapshot.CapturedAt.IsZero() {
			t.Error("expected CapturedAt to be set")
		}
	})

	t.Run("runner error", func(t *testing.T) {
		mock := &mockRunner{err: errors.New("device not found")}
		cam := &LogitechCamera{
			devicePath: "/dev/video0",
			outputDir:  "/tmp/captures",
			runner:     mock,
		}
		_, err := cam.Capture(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "camera hardware failure") {
			t.Errorf("expected wrapped error message, got: %v", err)
		}

	})

}
