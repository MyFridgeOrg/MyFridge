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

type motionSensor interface {
	// WaitForMotion blocks until a motion event is detected.
	WaitForMotion(ctx context.Context) <-chan struct{}
}

// Coordinator orchestrates the motion-detection pipeline.
type Coordinator struct {
	motionSensor motionSensor
	cam          camera
}

// New creates a Coordinator with the given dependencies.
func New(motionSensor motionSensor, cam camera) *Coordinator {
	return &Coordinator{
		motionSensor: motionSensor,
		cam:          cam,
	}
}

func (c *Coordinator) Run(ctx context.Context) {
	events := c.motionSensor.WaitForMotion(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-events:
			snapshot, err := c.cam.Capture(ctx)
			if err != nil {
				log.Fatalf("ERROR: %v", err)
			}
			log.Printf("Successfully captured an image stored at: %s", snapshot.FilePath)
		}
	}
}
