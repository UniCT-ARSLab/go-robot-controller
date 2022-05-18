package main

import (
	//"os"

	"os"
	"os/signal"
	"time"

	//"github.com/arslab/robot_controller/robot"

	"github.com/arslab/robot_controller/robot"
	"github.com/arslab/robot_controller/webserver"
)

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	time.Sleep(time.Second * 5)

	robot, err := robot.NewRobot("can0")
	if err != nil {
		os.Exit(1)
	}

	webServer := webserver.NewWebServer(robot, "0.0.0.0", 9998)

	webServer.Start()

	for true {
		time.Sleep(1 * time.Second)
	}

}
