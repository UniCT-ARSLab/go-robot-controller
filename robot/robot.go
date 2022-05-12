package robot

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/arslab/robot_controller/models"
	"github.com/arslab/robot_controller/utilities"
	"github.com/brutella/can"
	"github.com/fatih/color"
)

const DEBUG_CAN = false

const (
	ID_ROBOT_POSITION       = 0x3E3
	ID_OTHER_ROBOT_POSITION = 0x3E5
	ID_ROBOT_SPEED          = 0x3E4
	ID_ROBOT_STATUS         = 0x402
	ID_MOTION_CMD           = 0x7F0
	ID_ST_CMD               = 0x710
	ID_OBST_MAP             = 0x70f
)

//Robot rappresents the logical Robot
type Robot struct {
	Connection             Connection
	StartPositionSetted    bool
	StartPosition          models.Position
	Position               models.Position
	Speed                  int16
	Stopped                bool
	CallbackPositionUpdate func(pos models.Position)
	Type                   string
	Color                  uint8
	StarterEnabled         bool
	TimerBattery           int16
}

//NewRobot return a new Robot instance
func NewRobot(networkInterface string) (*Robot, error) {

	conn := NewConnection(networkInterface)

	robot := Robot{
		Connection:          *conn,
		StartPositionSetted: false,
		Position:            models.Position{X: 0, Y: 0, Angle: 0},
		StartPosition:       models.Position{X: -1000, Y: -1000, Angle: -1000},
		Speed:               0,
		Stopped:             false,
		Type:                os.Getenv("ROBOT"),
		TimerBattery:        1200,
	}

	if connError := robot.Connection.Init(); connError != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), "Connection Error!")
		return nil, connError
	}

	robot.Connection.OnReceiveCallback(robot.onDataReceived)

	go func() {
		robot.Connection.Connect()
	}()

	go func() {

		for robot.TimerBattery > 0 {
			time.Sleep(time.Second)
			robot.TimerBattery--
		}
	}()

	log.Printf("Controller for robot %s started.", robot.Type)
	//robot.SetPosition(models.Position{X: 0, Y: 0, Angle: 0})
	return &robot, nil
}

func (robot *Robot) onDataReceived(frm can.Frame) {
	data := frm.Data[:]

	switch frm.ID {
	case ID_ROBOT_POSITION:
		//position

		var posX int16
		var posY int16
		var angle int16

		buf := bytes.NewBuffer(data[:2])
		binary.Read(buf, binary.LittleEndian, &posX)
		buf = bytes.NewBuffer(data[2:4])
		binary.Read(buf, binary.LittleEndian, &posY)
		buf = bytes.NewBuffer(data[4:6])
		binary.Read(buf, binary.LittleEndian, &angle)
		angle = angle / 100.0

		if robot.StartPosition.X < -999 {
			robot.StartPosition.X = posX
			robot.StartPosition.Y = posY
			robot.StartPosition.Angle = angle
		}

		robot.Position.X = posX
		robot.Position.Y = posY
		robot.Position.Angle = angle
		if DEBUG_CAN {
			log.Printf("%s : [X : %d, Y : %d, A : %d]\n", "Position", posX, posY, angle)
		}
	case ID_ROBOT_SPEED:
		var speed int16
		buf := bytes.NewBuffer(data[:4])
		binary.Read(buf, binary.LittleEndian, &speed)
		robot.Speed = speed
		if DEBUG_CAN {
			log.Printf("%s : [%d]\n", "Linear Speed", speed)
		}
	case ID_ROBOT_STATUS:
		var status int16
		buf := bytes.NewBuffer(data[2:4])
		binary.Read(buf, binary.LittleEndian, &status)
		if DEBUG_CAN {
			log.Printf("%s : [%d]\n", "Status", status)
		}
	case ID_OBST_MAP:
		var obstacle_number uint8
		var valid uint8
		var angleStart int16
		var angleEnd int16
		var distance int16

		buf := bytes.NewBuffer(data[:1])
		binary.Read(buf, binary.LittleEndian, &obstacle_number)
		buf = bytes.NewBuffer(data[1:2])
		binary.Read(buf, binary.LittleEndian, &valid)
		buf = bytes.NewBuffer(data[2:4])
		binary.Read(buf, binary.LittleEndian, &angleStart)
		buf = bytes.NewBuffer(data[4:6])
		binary.Read(buf, binary.LittleEndian, &angleEnd)
		buf = bytes.NewBuffer(data[6:8])
		binary.Read(buf, binary.LittleEndian, &distance)
		if DEBUG_CAN {
			log.Printf("%s : Number: [%d], Valid: [%d], AStart: [%d], AEnd: [%d], Distance: [%d]\n", "Obstacle map", obstacle_number, valid, angleStart, angleEnd, distance)
		}

		//default:
		//	log.Printf("%s : [%x]\n", "UNKNOWN ID CAN", frm.ID)
	}

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

// //UpdatePosition logical position from I2C board
// func (robot *Robot) UpdatePosition() error {

// 	read := make([]byte, 6)
// 	if err := robot.Connection.ReceiveData(read, REG_GET_POSITION); err != nil {
// 		printError("UpdatePosition Error!")
// 		return err
// 	}
// 	//numBytes := []byte{read[0], read[1]}
// 	var x int16
// 	buf := bytes.NewBuffer(read[:2])
// 	binary.Read(buf, binary.LittleEndian, &x)

// 	var y int16
// 	buf = bytes.NewBuffer(read[2:4])
// 	binary.Read(buf, binary.LittleEndian, &y)

// 	var a int16
// 	buf = bytes.NewBuffer(read[4:6])
// 	binary.Read(buf, binary.LittleEndian, &a)
// 	a = a / 100.0

// 	if !robot.StartPositionSetted {
// 		robot.StartPositionSetted = true

// 		robot.StartPosition = models.Position{
// 			X:     x,
// 			Y:     y,
// 			Angle: a,
// 		}
// 		printInfo("Start Position: X: " + strconv.Itoa(int(robot.StartPosition.X)) + ", Y: " + strconv.Itoa(int(robot.StartPosition.Y)) + ", Angle: " + strconv.Itoa(int(robot.StartPosition.Angle)))

// 	}

// 	// robot.Position.X = (x - robot.LastBoardPosition.X) + robot.Position.X
// 	// robot.Position.Y = (y - robot.LastBoardPosition.Y) + robot.Position.Y
// 	// robot.Position.Angle = (a - robot.LastBoardPosition.Angle) + robot.Position.Angle

// 	deltaX := (x - robot.StartPosition.X)
// 	deltaY := (y - robot.StartPosition.Y)

// 	radiands := float64(robot.StartPosition.Angle) * (math.Pi / 180)

// 	robot.LastBoardPosition.X = x
// 	robot.LastBoardPosition.Y = y
// 	robot.LastBoardPosition.Angle = a

// 	//println("dX", deltaX, "dY", deltaY)

// 	newX := int16(float64(deltaX)*math.Cos(radiands) + float64(deltaY)*math.Sin(radiands))
// 	newY := int16(float64(-deltaX)*math.Sin(radiands) + float64(deltaY)*math.Cos(radiands))
// 	newA := a - robot.StartPosition.Angle

// 	if newA < -180 {
// 		newA = 360 + newA
// 	} else if newA > 180 {
// 		newA = 360 - newA
// 	}

// 	robot.Position.X = newX
// 	robot.Position.Y = newY
// 	robot.Position.Angle = newA

// 	//println("LogicX:", robot.Position.X, "LogicY:", robot.Position.Y, "LogicAngle:", robot.Position.Angle)

// 	//printInfo("Position updated")
// 	if robot.CallbackPositionUpdate != nil {
// 		robot.CallbackPositionUpdate(robot.Position)
// 	}
// 	return nil
// }

//SetPosition set the position on the I2C board
func (robot *Robot) SetPosition(p models.Position) error {

	motionCMD := models.MotionCommand{
		CMD:     models.MC_SET_POSITION,
		PARAM_1: p.X,
		PARAM_2: p.Y,
		PARAM_3: p.Angle,
	}

	err := robot.Connection.SendData(motionCMD, ID_MOTION_CMD)

	if err == nil {
		log.Printf("[%s] %s : X: %d, Y: %d, Angle: %d", utilities.CreateColorString("ROBOT", color.FgHiCyan), "Position changed", p.X, p.Y, p.Angle)
		return nil
	} else {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), err)
		return err
	}

}

//SetSpeed of the Robot
func (robot *Robot) SetSpeed(speed int16) error {
	motionCMD := models.MotionCommand{
		CMD:     models.MC_SET_SPEED,
		PARAM_1: speed,
	}

	err := robot.Connection.SendData(motionCMD, ID_MOTION_CMD)

	if err != nil {
		log.Printf("[%s] %s : Speed: %d", utilities.CreateColorString("ROBOT", color.FgHiCyan), "Set Speed", speed)
		return nil
	} else {
		return err
	}
}

//ForwardDistance move the robot about the given millimeters
func (robot *Robot) ForwardDistance(distance int16) error {

	// if robot.Stopped {
	// 	printError("The robot id Stopped")
	// 	return errors.New("The robot id Stopped")
	// }

	motionCMD := models.MotionCommand{
		CMD:     models.MC_FW_TO_DISTANCE,
		PARAM_1: distance,
	}

	err := robot.Connection.SendData(motionCMD, ID_MOTION_CMD)

	if err != nil {
		log.Printf("[%s] %s : Distance: %d", utilities.CreateColorString("ROBOT", color.FgHiCyan), "Forward Distance", distance)
		return nil
	} else {
		return err
	}
}

func (robot *Robot) StopMotors() error {

	motionCMD := models.MotionCommand{
		CMD: models.MC_STOP,
	}

	err := robot.Connection.SendData(motionCMD, ID_MOTION_CMD)

	if err == nil {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiCyan), "Motors Stopped")
		return nil
	} else {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), err)
		return err
	}
}

func (robot *Robot) Align(colorIn uint8) error {
	var cmd uint8
	if robot.Type == "piccolo" {
		cmd = models.ST_ALIGN_PICCOLO
	} else {
		cmd = models.ST_ALIGN_GRANDE
	}

	robot.Color = colorIn

	motionCMD := models.StrategyCommand{
		CMD:   cmd,
		FLAGS: colorIn,
	}

	err := robot.Connection.SendData(motionCMD, ID_ST_CMD)

	if err == nil {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiCyan), "Aligning")
		return nil
	} else {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), err)
		return err
	}
}

func (robot *Robot) ToggleStarter(enable bool) error {
	var cmd uint8
	if enable {
		cmd = models.ST_ENABLE_STARTER
	} else {
		cmd = models.ST_DISABLE_STARTER
	}

	robot.StarterEnabled = enable

	motionCMD := models.StrategyCommand{
		CMD: cmd,
	}

	err := robot.Connection.SendData(motionCMD, ID_ST_CMD)

	if err == nil {
		log.Printf("[%s] %s %t", utilities.CreateColorString("ROBOT", color.FgHiCyan), "Starter Toggled", enable)
		return nil
	} else {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), err)
		return err
	}
}

func (robot *Robot) ResetBoard() error {

	_, err := exec.Command("openocd", "-f", "board/st_nucleo_f4.cfg", "-c", "init", "-c", "reset", "-c", "exit").Output()
	log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiCyan), "Board Reset")
	if err != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), err)
		return err
	}
	return nil
}

//ForwardToPoint move the robot to the defined point
func (robot *Robot) ForwardToPoint(x int16, y int16) error {

	// if robot.Stopped {
	// 	printError("The robot id Stopped")
	// 	return errors.New("The robot id Stopped")
	// }

	// if err := robot.Connection.SendData(
	// 	models.I2CMessage{
	// 		Command: REG_FW_POINT,
	// 		Val1:    x,
	// 		Val2:    y,
	// 		Val3:    0,
	// 		End:     0,
	// 	},
	// 	0x60); err != nil {
	// 	printError("ForwardToPoint Error!")
	// 	return err
	// }
	// printInfo("Moved to X:" + strconv.Itoa(int(x)) + " Y:" + strconv.Itoa(int(y)))
	return nil
}

//RelativeRotation rotate the robot about the given degrees
func (robot *Robot) RelativeRotation(degree int16) error {

	// if robot.Stopped {
	// 	printError("The robot id Stopped")
	// 	return errors.New("The robot id Stopped")
	// }

	// if err := robot.Connection.SendData(
	// 	models.I2CMessage{
	// 		Command: REG_REL_ROTATION,
	// 		Val1:    degree,
	// 		Val2:    0,
	// 		Val3:    0,
	// 		End:     0,
	// 	},
	// 	0x60); err != nil {
	// 	printError("RelativeRotation Error!")
	// 	return err
	// }
	// printInfo("Relative Rotated by " + strconv.Itoa(int(degree)))
	return nil
}

//AbsoluteRotation rotate the robot about the given degrees
func (robot *Robot) AbsoluteRotation(degree int16) error {

	// if robot.Stopped {
	// 	printError("The robot id Stopped")
	// 	return errors.New("The robot id Stopped")
	// }

	// if err := robot.Connection.SendData(
	// 	models.I2CMessage{
	// 		Command: REG_ABS_ROTATION,
	// 		Val1:    degree,
	// 		Val2:    0,
	// 		Val3:    0,
	// 		End:     0,
	// 	},
	// 	0x60); err != nil {
	// 	printError("AbsoluteRotation Error!")
	// 	return err
	// }

	return nil
}
