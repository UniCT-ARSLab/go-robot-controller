package robot

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/arslab/robot_controller/utilities"
	"github.com/fatih/color"
	"github.com/stianeikeland/go-rpio/v4"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

type Connection struct {
	GPIOPin rpio.Pin
	Device  i2c.Dev
	Speed   int16
}

func NewConnection(gpioPIN uint16, i2cAddress uint16) *Connection {

	host.Init()
	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//open i2c
	b, erri2c := i2creg.Open("")
	if erri2c != nil {
		log.Fatal(erri2c)
		os.Exit(1)
	}
	//creo un "device" i2c usando il bus "b" (/etc/i2c)
	device := i2c.Dev{Addr: i2cAddress, Bus: b}

	connect := Connection{
		GPIOPin: rpio.Pin(gpioPIN),
		Device:  device,
	}

	return &connect

}

func (conn *Connection) Init() {
	conn.GPIOPin.Output()
	conn.GPIOPin.High()
	fmt.Println("Board initialised")
	log.Printf("[%s] %s", utilities.CreateColorString("Web Panel", color.FgHiCyan), "avaiable on port:"+strconv.Itoa(ws.Port))
}
