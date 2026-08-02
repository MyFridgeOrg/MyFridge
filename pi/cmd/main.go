package main

import (
	"context"

	"example.com/internal/camera"
	"example.com/internal/coordinator"
)

func main() {
	cam := camera.NewPiCamera(".")
	coord := coordinator.New(cam)
	ctx := context.Background()
	coord.Run(ctx)
}
