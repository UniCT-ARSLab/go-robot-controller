package main

import (
	"os"
	"time"

	"github.com/arslab/robot_controller/robot"
	"github.com/arslab/robot_controller/webserver"
)

func main() {

	robot, err := robot.NewRobot(4, 0x34)
	if err != nil {
		os.Exit(1)
	}

	webServer := webserver.NewWebServer(robot, "0.0.0.0", 9998)

	webServer.Start()

	// go func() {
	// 	for true {
	// 		time.Sleep(1 * time.Second)
	// 		pos := robot.GetPosition()
	// 		fmt.Println("Posizione 1:", pos)

	// 	}
	// }()

	for true {
		time.Sleep(1 * time.Second)
	}

	// pos := robot.GetPosition()
	// fmt.Println("Posizione:", pos)
	// robot.SetSpeed(200)
	// robot.ForwardDistance(100)

	// go func() {
	// 	for true {
	// 		time.Sleep(4 * time.Second)
	// 		fmt.Println("RESET POSITION:")
	// 		robot.SetPosition(models.Position{X: 10, Y: 10, Angle: 0})

	// 	}
	// }()

}
