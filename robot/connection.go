package robot

import (
	"bytes"
	"encoding/binary"
	"log"
	"strconv"
	"time"

	"github.com/arslab/robot_controller/utilities"
	"github.com/fatih/color"
	"github.com/stianeikeland/go-rpio/v4"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

type Connection struct {
	GPIONumPin int
	GPIOPin    rpio.Pin
	I2CAddress uint16
	Device     i2c.Dev
	Speed      int16
}

func NewConnection(gpioPIN int, i2cAddress uint16) *Connection {

	connect := Connection{
		GPIONumPin: gpioPIN,
		GPIOPin:    rpio.Pin(gpioPIN),
		I2CAddress: i2cAddress,
		//Device:     device,
	}

	return &connect
}

func (conn *Connection) Init() error {
	host.Init()
	if _, err := driverreg.Init(); err != nil {
		return err
	}

	conn.GPIOPin.Output()
	conn.GPIOPin.High()

	host.Init()
	if _, err := driverreg.Init(); err != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiRed), err)
		return err
	}

	b, erri2c := i2creg.Open("")
	if erri2c != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiRed), erri2c)
		return erri2c
	}

	conn.Device = i2c.Dev{Addr: conn.I2CAddress, Bus: b}

	log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgYellow), "Connection Initialised on GPIO Pin :"+strconv.Itoa(conn.GPIONumPin))
	return nil
}

func (conn *Connection) Reset() {
	log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgYellow), "Resetting the connection...")
	conn.GPIOPin.Low()
	time.Sleep(time.Millisecond * 100)
	conn.GPIOPin.High()
	time.Sleep(time.Second * 2)
	log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgYellow), "Connection Resetted")
}

func (conn *Connection) SendData(payload interface{}, register byte) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, payload)
	if err != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiRed), err)
		return err
	}

	write := append([]byte{register}, buf.Bytes()...)

	if err := conn.Device.Tx(write, nil); err != nil {
		return err
	}

	return nil
}

func (conn *Connection) ReceiveData(read []byte, register byte) error {
	write := []byte{register}
	if err := conn.Device.Tx(write, read); err != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiRed), err)
		return err
	}
	return nil
}

func (conn *Connection) SendReceiveData(payload interface{}, read []byte, register byte) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, payload)
	if err != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiRed), err)
		return err
	}

	write := append([]byte{register}, buf.Bytes()...)

	if err := conn.Device.Tx(write, read); err != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiRed), err)
		return err
	}

	return nil
}
