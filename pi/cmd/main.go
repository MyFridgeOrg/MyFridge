package main

import (
	"context"
	"fmt"

	"example.com/internal/camera"
	"example.com/internal/coordinator"
	"example.com/internal/sensor"
)

func main() {
	fmt.Print(banner)
	cam := camera.NewPiCamera(".")
	sensor := sensor.NewPIRMotionSensor("4")
	coord := coordinator.New(sensor, cam)
	ctx := context.Background()
	coord.Run(ctx)
}
