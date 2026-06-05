package coordinator

import "example.com/internal/models"

// camera captures a snapshot in response to a motion event.
type camera interface {
	Capture() (models.SnapShot, error)
}

// Coordinator orchestrates the motion-detection pipeline.
type Coordinator struct {
	//motionSensor
	cam camera
}

// New creates a Coordinator with the given dependencies.
func New(cam camera) *Coordinator {
	return &Coordinator{cam: cam}
}
