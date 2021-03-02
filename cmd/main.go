package main

import (
	"fmt"
	"os"
	"time"

	"github.com/arslab/robot_controller/models"
	"github.com/arslab/robot_controller/robot"
)

func main() {

	robot, err := robot.NewRobot(4, 0x34)
	if err != nil {
		os.Exit(1)
	}

	robot.SetPosition(models.Position{X: 0, Y: 0, Angle: 0})
	pos := robot.GetPosition()
	fmt.Println("Posizione:", pos)
	robot.SetSpeed(200)
	robot.ForwardDistance(100)

	go func() {
		for true {
			time.Sleep(1 * time.Second)
			pos = robot.GetPosition()
			fmt.Println("Posizione 1:", pos)

		}
	}()

	// go func() {
	// 	for true {
	// 		time.Sleep(4 * time.Second)
	// 		fmt.Println("RESET POSITION:")
	// 		robot.SetPosition(models.Position{X: 10, Y: 10, Angle: 0})

	// 	}
	// }()

	for true {
		time.Sleep(1 * time.Second)
	}

}
