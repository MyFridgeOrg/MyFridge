package sensor

import (
	"context"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type PIRMotionSensor struct {
	pin gpio.PinIO
}

func NewPIRMotionSensor(pinName string) *PIRMotionSensor {
	if _, err := host.Init(); err != nil {
		return nil
	}

	pin := gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatalf("Failed to find pin %s", pinName)
	}

	if err := pin.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		log.Fatal(err)
	}
	return &PIRMotionSensor{
		pin: pin,
	}
}

func (p *PIRMotionSensor) Read() gpio.Level {
	return p.pin.Read()
}

func (p *PIRMotionSensor) WaitForMotion(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if p.pin.WaitForEdge(100 * time.Millisecond) {
				if p.Read() == gpio.High {
					select {
					case ch <- struct{}{}: // send event to channel
					case <-ctx.Done(): // if noone is there to receive from the channel, return
						return
					}
				}
			}
			select {
			case ch <- <-ctx.Done():
				return
			default:
			}
		}
	}()
	return ch
}
