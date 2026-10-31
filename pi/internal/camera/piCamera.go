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
	uuid := fmt.Sprintf("%d", time.Now().UnixNano())
	fileName := fmt.Sprintf("motion_capture_%v.jpg", uuid)
	finalPath := filepath.Join(c.outputDir, fileName)

	err := exec.CommandContext(ctx, "rpicam-still",
		"-o", finalPath,
		"-t", "5000",
		"--immediate",
	).Run()

	if err != nil {
		return models.SnapShot{}, fmt.Errorf("CSI camera hardware failed execution: %w", err)
	}

	return models.SnapShot{
		EventID:    "test",
		FilePath:   finalPath,
		CapturedAt: time.Now(),
	}, nil
}
