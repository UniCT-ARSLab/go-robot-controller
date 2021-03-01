package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

type i2cMessage struct {
	Command byte
	Val1    int16
	Val2    int16
	Val3    int16
	End     byte
}

var pin = rpio.Pin(4)
var device i2c.Dev

func BoardInit() {
	pin.Output()
	pin.High()
	fmt.Println("Board initialised")
}

func BoardReset() {
	fmt.Println("Board Resetting...")
	pin.Low()
	time.Sleep(time.Millisecond * 100)
	pin.High()
	time.Sleep(time.Second * 2)
	fmt.Println("Board Resetted")
}

func OpenI2C() {
	host.Init()
	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
	}

	//apro i2c
	b, erri2c := i2creg.Open("")
	if erri2c != nil {
		log.Fatal(erri2c)
		os.Exit(1)
	}
	//creo un "device" i2c usando il bus "b" (/etc/i2c)
	device = i2c.Dev{Addr: 0x34, Bus: b}
}

func SendCommand(message i2cMessage) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, message)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	write := append([]byte{0x60}, buf.Bytes()...)

	if err := device.Tx(write, nil); err != nil {
		log.Fatal(err)
		BoardReset()
		os.Exit(1)
	}
}

func SetSpeed(speed int16) {

	//velocità

	toSend := i2cMessage{
		Command: 0x8c,
		Val1:    speed,
		Val2:    0,
		Val3:    0,
		End:     0,
	}

	fmt.Printf("Set Speed: %d\n", speed)
	SendCommand(toSend)
}

func ForwardDistance(distance int16) {
	//movimento
	toSend := i2cMessage{
		Command: 0x85,
		Val1:    distance,
		Val2:    0,
		Val3:    0,
		End:     0,
	}

	fmt.Printf("Foward: %d\n", distance)

	SendCommand(toSend)
}

func RelativeRotate(angle int16) {
	//movimento
	toSend := i2cMessage{
		Command: 0x88,
		Val1:    angle,
		Val2:    0,
		Val3:    0,
		End:     0,
	}
	fmt.Printf("Rotation: %d\n", angle)
	SendCommand(toSend)
}

func GetPosition() {
	// posizione
	write := []byte{0x01}
	read := make([]byte, 6)
	if err := device.Tx(write, read); err != nil {
		log.Fatal(err)
		BoardReset()
		os.Exit(1)
	}
	//numBytes := []byte{read[0], read[1]}
	var x int16
	buf := bytes.NewBuffer(read[:2])
	binary.Read(buf, binary.LittleEndian, &x)

	var y int16
	buf = bytes.NewBuffer(read[2:4])
	binary.Read(buf, binary.LittleEndian, &y)

	var a int16
	buf = bytes.NewBuffer(read[4:6])
	binary.Read(buf, binary.LittleEndian, &a)
	a = a / 100.0

	//y := binary.LittleEndian.Uint16(read[2:4])
	//a := binary.LittleEndian.Uint16(read[4:6])
	fmt.Printf("Position: {X:%d, Y:%d, ANGLE:%d}\n", x, y, a)

}

func main() {

	rpio.Open()
	defer rpio.Close()
	BoardInit()
	OpenI2C()

	GetPosition()
	SetSpeed(200)
	ForwardDistance(100)

}
