package robot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/arslab/robot_controller/models"
	"github.com/arslab/robot_controller/utilities"
	"github.com/fatih/color"
)

const (
	REG_GET_POSITION = 0x01
	REG_SET_POSITION = 0x84
	REG_SET_SPEED    = 0x8c
	REG_FW_DISTANCE  = 0x85
	REG_FW_POINT     = 0x85
	REG_REL_ROTATION = 0x88
	REG_ABS_ROTATION = 0x89
)

//Robot rappresents the logical Robot
type Robot struct {
	Connection        Connection
	Position          models.Position // L
	LastBoardPosition models.Position // U
	Speed             int16
	Stopped           bool
}

//NewRobot return a new Robot instance
func NewRobot(gpioPIN int, i2cAddress uint16) (*Robot, error) {

	conn := NewConnection(gpioPIN, i2cAddress)

	robot := Robot{
		Connection:        *conn,
		Position:          models.Position{X: 0, Y: 0, Angle: 0},
		Speed:             0,
		LastBoardPosition: models.Position{X: 0, Y: 0, Angle: 0},
		Stopped:           false,
	}

	if connError := robot.Connection.Init(); connError != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), "Connection Error!")
		return nil, connError
	}

	go func() {
		for true {
			if err := robot.UpdatePosition(); err != nil {
				robot.Stopped = true
				robot.Connection.Reset()
				robot.Stopped = false
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return &robot, nil
}

func printError(s string) {
	log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), s)
}

func printInfo(s string) {
	log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiCyan), s)
}

//GetPosition returns the position of the Robot
func (robot *Robot) GetPosition() models.Position {

	return robot.Position
}

//UpdatePosition logical position from I2C board
func (robot *Robot) UpdatePosition() error {

	read := make([]byte, 6)
	if err := robot.Connection.ReceiveData(read, REG_GET_POSITION); err != nil {
		printError("UpdatePosition Error!")
		return err
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

	robot.Position.X = (x - robot.LastBoardPosition.X) + robot.Position.X
	robot.Position.Y = (y - robot.LastBoardPosition.Y) + robot.Position.Y
	robot.Position.Angle = (a - robot.LastBoardPosition.Angle) + robot.Position.Angle

	robot.LastBoardPosition.X = x
	robot.LastBoardPosition.Y = y
	robot.LastBoardPosition.Angle = a

	fmt.Println("board:", x, "last:", robot.LastBoardPosition.X, "logical:", robot.Position.X)

	//robot.Position.X = (robot.Position.X + x) - (robot.LastBoardPosition.X)
	//robot.Position.Y = (robot.Position.Y + y) - (robot.LastBoardPosition.Y)
	//robot.Position.Angle = (robot.Position.Angle + a) - (robot.LastBoardPosition.Angle)

	//printInfo("Position updated")
	return nil
}

//SetPosition set the position on the I2C board
func (robot *Robot) SetPosition(p models.Position) error {

	// if err := robot.Connection.SendData(
	// 	models.I2CMessage{
	// 		Command: REG_SET_POSITION,
	// 		Val1:    p.X,
	// 		Val2:    p.Y,
	// 		Val3:    p.Angle,
	// 		End:     0,
	// 	},
	// 	0x60); err != nil {
	// 	printError("SetPosition Error!")
	// 	return err
	// }

	// robot.LastBoardPosition.X = robot.LastBoardPosition.X + robot.Position.X
	// robot.LastBoardPosition.Y = robot.LastBoardPosition.Y + robot.Position.Y
	// robot.LastBoardPosition.Angle = robot.LastBoardPosition.Angle + robot.Position.Angle

	robot.Position.X = p.X
	robot.Position.Y = p.Y
	robot.Position.Angle = p.Angle

	printInfo("Position changed")
	return nil
}

//SetSpeed of the Robot
func (robot *Robot) SetSpeed(speed int16) error {
	if err := robot.Connection.SendData(
		models.I2CMessage{
			Command: REG_SET_SPEED,
			Val1:    speed,
			Val2:    0,
			Val3:    0,
			End:     0,
		},
		0x60); err != nil {
		printError("SetSpeed Error!")
		return err
	}

	printInfo("Speed setted at " + strconv.Itoa(int(speed)))
	return nil
}

//ForwardDistance move the robot about the given millimeters
func (robot *Robot) ForwardDistance(distance int16) error {

	if robot.Stopped {
		printError("The robot id Stopped")
		return errors.New("The robot id Stopped")
	}
	if err := robot.Connection.SendData(
		models.I2CMessage{
			Command: REG_FW_DISTANCE,
			Val1:    distance,
			Val2:    0,
			Val3:    0,
			End:     0,
		},
		0x60); err != nil {
		printError("ForwardDistance Error!")
		return err
	}

	printInfo("Moved by " + strconv.Itoa(int(distance)))
	return nil
}

//ForwardToPoint move the robot to the defined point
func (robot *Robot) ForwardToPoint(x int16, y int16) error {

	if robot.Stopped {
		printError("The robot id Stopped")
		return errors.New("The robot id Stopped")
	}

	if err := robot.Connection.SendData(
		models.I2CMessage{
			Command: REG_FW_POINT,
			Val1:    x,
			Val2:    y,
			Val3:    0,
			End:     0,
		},
		0x60); err != nil {
		printError("ForwardToPoint Error!")
		return err
	}
	printInfo("Moved to X:" + strconv.Itoa(int(x)) + " Y:" + strconv.Itoa(int(y)))
	return nil
}

//RelativeRotation rotate the robot about the given degrees
func (robot *Robot) RelativeRotation(degree int16) error {

	if robot.Stopped {
		printError("The robot id Stopped")
		return errors.New("The robot id Stopped")
	}

	if err := robot.Connection.SendData(
		models.I2CMessage{
			Command: REG_REL_ROTATION,
			Val1:    degree,
			Val2:    0,
			Val3:    0,
			End:     0,
		},
		0x60); err != nil {
		printError("RelativeRotation Error!")
		return err
	}

	return nil
}

//AbsoluteRotation rotate the robot about the given degrees
func (robot *Robot) AbsoluteRotation(degree int16) error {

	if robot.Stopped {
		printError("The robot id Stopped")
		return errors.New("The robot id Stopped")
	}

	if err := robot.Connection.SendData(
		models.I2CMessage{
			Command: REG_ABS_ROTATION,
			Val1:    degree,
			Val2:    0,
			Val3:    0,
			End:     0,
		},
		0x60); err != nil {
		printError("AbsoluteRotation Error!")
		return err
	}

	return nil
}
