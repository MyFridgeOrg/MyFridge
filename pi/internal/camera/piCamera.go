package camera

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"example.com/internal/models"
)

type PiCamera struct {
	outputDir string
}

func NewPiCamera(outputDir string) *PiCamera {
	return &PiCamera{
		outputDir: outputDir,
	}
}

func (c *PiCamera) Capture(ctx context.Context) (models.SnapShot, error) {
	fileName := fmt.Sprintf("motion_.jpg")
	finalPath := filepath.Join(c.outputDir, fileName)

	cmd := exec.CommandContext(ctx, "rpicam-still",
		"-o", finalPath,
		"-t", "1",
		"--immediate",
	)

	if err := cmd.Run(); err != nil {
		return models.SnapShot{}, fmt.Errorf("CSI camera hardware failed execution: %w", err)
	}

	return models.SnapShot{
		EventID:    "test",
		FilePath:   finalPath,
		CapturedAt: time.Now(),
	}, nil
}
