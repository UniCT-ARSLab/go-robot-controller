package robot

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	Connection             Connection
	Position               models.Position // L
	LastBoardPosition      models.Position // U
	Speed                  int16
	Stopped                bool
	CallbackPositionUpdate func(pos models.Position)
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
			time.Sleep(500 * time.Millisecond)
		}
	}()
	robot.SetPosition(models.Position{X: 0, Y: 0, Angle: 0})
	return &robot, nil
}

func printError(s string) {
	log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), s)
}

func printInfo(s string) {
	log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiCyan), s)
}

//SetCallbackUpadetePosition set the callback position update function
func (robot *Robot) SetCallbackUpadetePosition(cb func(position models.Position)) {
	robot.CallbackPositionUpdate = cb
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

	//fmt.Println("board:", x, "last:", robot.LastBoardPosition.X, "logical:", robot.Position.X)

	//robot.Position.X = (robot.Position.X + x) - (robot.LastBoardPosition.X)
	//robot.Position.Y = (robot.Position.Y + y) - (robot.LastBoardPosition.Y)
	//robot.Position.Angle = (robot.Position.Angle + a) - (robot.LastBoardPosition.Angle)

	//printInfo("Position updated")
	if robot.CallbackPositionUpdate != nil {
		robot.CallbackPositionUpdate(robot.Position)
	}
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

	robot.Position.X = p.X
	robot.Position.Y = p.Y
	robot.Position.Angle = p.Angle

	log.Printf("[%s] %s : X: %d, Y: %d, Angle: %d", utilities.CreateColorString("ROBOT", color.FgHiCyan), "Position changed", p.X, p.Y, p.Angle)
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
	robot.Speed = speed
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
	printInfo("Relative Rotated by " + strconv.Itoa(int(degree)))
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
