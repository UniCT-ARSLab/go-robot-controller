package robot

import (
	"github.com/stianeikeland/go-rpio/v4"
	"periph.io/x/conn/v3/i2c"
)

type Robot struct {
	GPIOPin  rpio.Pin
	Device   i2c.Dev
	Location Position
}
