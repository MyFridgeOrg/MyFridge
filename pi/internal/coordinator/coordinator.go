package coordinator

import (
	"context"
	"log"

	"example.com/internal/models"
)

// camera captures a snapshot in response to a motion event.
type camera interface {
	Capture(ctx context.Context) (models.SnapShot, error)
}

// Coordinator orchestrates the motion-detection pipeline.
type Coordinator struct {
	//motionSensor
	cam camera
}

// New creates a Coordinator with the given dependencies.
func New(cam camera) *Coordinator {
	return &Coordinator{
		cam: cam,
	}
}

func (c *Coordinator) Run(ctx context.Context) {
	snapshot, err := c.cam.Capture(ctx)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	log.Printf("Successfully captured an image stored at: %s", snapshot.FilePath)
}
