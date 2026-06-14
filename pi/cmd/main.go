package main

import (
	"context"

	"example.com/internal/camera"
	"example.com/internal/coordinator"
)

func main() {
	cam := camera.NewLogitechCamera("/dev/video0", ".")
	coord := coordinator.New(cam)
	ctx := context.Background()
	coord.Run(ctx)
}
