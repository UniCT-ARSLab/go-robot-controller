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

	for true {
		time.Sleep(1 * time.Second)
	}

}
