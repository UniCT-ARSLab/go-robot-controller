package robot

import (
	//"bytes"
	//"encoding/binary"

	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/signal"

	//"time"

	"github.com/arslab/robot_controller/utilities"
	"github.com/brutella/can"
	"github.com/fatih/color"
)

//Connection is the interface between the logical Robot and the i2C Bus
type Connection struct {
	Interface string
	Bus       *can.Bus
	OnReceive func(data can.Frame)
}

// func handleCANFrame(frm can.Frame) {
// 	data := frm.Data[:] //trimSuffix(frm.Data[:], 0x00)
// 	//length := fmt.Sprintf("[%x]", frm.Length)

// 	log.Printf("%s : [%x], Length Data : %d\n", "ID", frm.ID, len(data))
// 	if frm.ID == 0x3E3 {
// 		//position
// 		posX := data[0:1]
// 		posY := data[2:3]
// 		angle := data[4:5]
// 		log.Printf("%s : [X : %d, Y : %d, A : %d]\n", "Position", posX, posY, angle)
// 	}
// 	if frm.ID == 0x3E7 {
// 		wheel := data[0:1]
// 		speed := data[2:3]
// 		tSpeed := data[4:5]
// 		log.Printf("%s : [Wheel : %d, Speed : %d, TargetSpeed : %d]\n", "Position", wheel, speed, tSpeed)
// 	}

// 	if frm.ID == 0x3E4 {
// 		linearSpeed := data[0:3]
// 		log.Printf("%s : [%d]\n", "Linear Speed", linearSpeed)
// 	}

// 	if frm.ID == 0x70F {
// 		if len(data) > 0 {
// 			obstacle := data[0]
// 			validity := data[1]
// 			angleStart := data[2:3]
// 			angleEnd := data[4:5]
// 			log.Printf("%s : [Obstacle : %d, Validity : %d, AngleStart : %d, AngleEnd : %d]\n", "Obstacle", obstacle, validity, angleStart, angleEnd)
// 		} else {
// 			log.Printf("%s : [No Obstacle]\n", "Obstacle")
// 		}

// 	}
// 	//log.Printf("%-3s %-4x %-3s % -24X '%s'\n", "can0", frm.ID, length, data, printableString(data[:]))
// }

//NewConnection return a new I2C Connection specifying the GPIO Pin and the Address of the device
func NewConnection(networkInterface string) *Connection {

	bus, err := can.NewBusForInterfaceWithName(networkInterface)
	if err != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiRed), err)
		os.Exit(1)
	}

	connection := Connection{
		Interface: networkInterface,
		Bus:       bus,
	}

	return &connection
}

func (conn *Connection) OnReceiveCallback(cb func(data can.Frame)) {
	conn.OnReceive = cb
	conn.Bus.SubscribeFunc(conn.OnReceive)
}

//Init initialise the CAN connection
func (conn *Connection) Init() error {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	go func() {

		select {
		case <-c:
			log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiYellow), "Connection Closed on Interface:"+conn.Interface)
			conn.Bus.Disconnect()
			os.Exit(1)
		}
	}()

	log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgYellow), "Connection Initialised on Interface:"+conn.Interface)
	return nil
}

func (conn *Connection) Disconnect() {
	conn.Bus.Disconnect()
}

func (conn *Connection) Connect() {
	go func() {
		conn.Bus.ConnectAndPublish()
	}()
}

//SendData allows to send data through the bus
func (conn *Connection) SendData(payload interface{}, id uint32) error {

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, payload)
	if err != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiRed), err)
		return err
	}
	err = nil

	var pl [8]uint8 // = buf.Bytes()
	for i, v := range buf.Bytes() {
		pl[i] = v
	}

	//log.Println(pl)

	frm := can.Frame{
		Length: 8,
		ID:     id,
		Data:   pl,
	}
	err = conn.Bus.Publish(frm)
	if err != nil {
		log.Println("Errore nella publish")
		log.Printf("[%s] %s", utilities.CreateColorString("CONNECTION", color.FgHiRed), err)

		return err
	}

	return nil
}
