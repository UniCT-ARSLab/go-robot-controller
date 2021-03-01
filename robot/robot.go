package robot

import (
	"github.com/arslab/robot_controller/models/"
)

type Robot struct {
	Connection Connection
	Position   models.Position
	Speed      int16
}
