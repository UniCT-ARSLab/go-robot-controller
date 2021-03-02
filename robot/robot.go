package robot

import (
	"github.com/arslab/robot_controller/models/"

	"log"

	"github.com/arslab/robot_controller/utilities"
	"github.com/fatih/color"
)

type Robot struct {
	Connection Connection
	Position   models.Position
	Speed      int16
}

func NewRobot(gpioPIN int, i2cAddress uint16) (*Robot, error) {

	conn := NewConnection(gpioPIN, i2cAddress)

	robot := Robot{
		Connection: *conn,
		Position:   models.Position{X: 0, Y: 0, Angle: 0},
		Speed:      0,
	}

	if connError := robot.Connection.Init(); connError != nil {
		log.Printf("[%s] %s", utilities.CreateColorString("ROBOT", color.FgHiRed), "Connection Error!")
		return nil, connError
	}

	return &robot, nil
}

func (robot *Robot) SetSpeed(speed uint16) error {
	robot.Connection.ReceiveData()
}
