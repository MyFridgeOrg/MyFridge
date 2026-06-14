package camera

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"example.com/internal/models"
)

// commandRunner abstracts command execution so it can be swapped in tests.
type commandRunner interface {
	Run(ctx context.Context, name string, args ...string) error
}

// execRunner is the real implementation that runs OS commands.
type execRunner struct{}

func (e execRunner) Run(ctx context.Context, name string, args ...string) error {
	return exec.CommandContext(ctx, name, args...).Run()
}

type LogitechCamera struct {
	devicePath string // e.g., "/dev/video0"
	outputDir  string // e.g., "/home/gustav/captures"
	runner     commandRunner
}

// NewLogitechCamera instantiates our physical webcam driver
func NewLogitechCamera(devicePath, outputDir string) *LogitechCamera {
	return &LogitechCamera{
		devicePath: devicePath,
		outputDir:  outputDir,
		runner:     execRunner{},
	}
}

func (c *LogitechCamera) Capture(ctx context.Context) (models.SnapShot, error) {
	fileName := fmt.Sprintf("motion.jpg")
	finalPath := filepath.Join(c.outputDir, fileName)

	err := c.runner.Run(ctx, "fswebcam",
		"-d", c.devicePath,
		"-r", "1920x1080",
		"--no-banner",
		finalPath,
	)

	if err != nil {
		return models.SnapShot{}, fmt.Errorf("camera hardware failure executing fswebcam: %w", err)
	}

	return models.SnapShot{
		EventID:    "MotionEvent test",
		FilePath:   finalPath,
		CapturedAt: time.Now(),
	}, nil
}
