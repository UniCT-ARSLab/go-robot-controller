package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/arslab/robot_controller/models"
	"github.com/arslab/robot_controller/robot"
	"github.com/arslab/robot_controller/utilities"
	_ "github.com/arslab/robot_controller/webserver/statik" //static file system
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rakyll/statik/fs"
	"gopkg.in/olahol/melody.v1"
)

//WebServer reppresents the WebServer and WebScoket server instance
type WebServer struct {
	Address       string
	Port          int
	Ssl           bool
	Router        *gin.Engine
	ServerSocket  *socketio.Server
	ServerSocketM *melody.Melody
}

var robotInstance *robot.Robot

//NewWebServer returns a new WebServer
func NewWebServer(robot *robot.Robot, address string, port int) *WebServer {
	robotInstance = robot
	serverSocket := newServerSocket()
	go func() {
		err := serverSocket.Serve()

		if err != nil {
			log.Fatal(err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New() // gin.Default()

	configCors := cors.DefaultConfig()
	//configCors.AllowOrigins = serverConfig.AllowOrigins
	configCors.AllowCredentials = true
	//router.Use(cors.New(configCors))
	router.Use(gin.Recovery())

	ws := WebServer{
		Address:       address,
		Port:          port,
		Ssl:           false,
		Router:        router,
		ServerSocket:  serverSocket,
		ServerSocketM: NewMelodyWebSocket(),
	}

	// Register REST
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	staticGroup := router.Group("/controller")
	staticGroup.StaticFS("/", statikFS)

	apiGroup := router.Group("/api")
	apiGroup.GET("/robot/position", func(context *gin.Context) { getRobotPosition(context) })
	apiGroup.POST("/robot/position", func(context *gin.Context) { setRobotPosition(context) })

	apiGroup.GET("/robot/speed", func(context *gin.Context) { getRobotSpeed(context) })
	apiGroup.POST("/robot/speed", func(context *gin.Context) { setRobotSpeed(context) })

	apiGroup.POST("/robot/move/distance", func(context *gin.Context) { robotForwardDistance(context) })
	apiGroup.POST("/robot/move/point", func(context *gin.Context) { robotForwardPoint(context) })

	apiGroup.POST("/robot/rotate/relative", func(context *gin.Context) { robotRelativeRotation(context) })
	apiGroup.POST("/robot/rotate/absolute", func(context *gin.Context) { robotAbsoluteRotation(context) })

	apiGroup.POST("/robot/motors/stop", func(context *gin.Context) { sendStop(context) })

	//apiGroup.GET("/system", func(context *gin.Context) { getSystemInformation(context) })

	router.GET("/socket.io/*any", gin.WrapH(serverSocket))
	router.POST("/socket.io/*any", gin.WrapH(serverSocket))

	router.GET("/ws", func(c *gin.Context) {
		ws.ServerSocketM.HandleRequest(c.Writer, c.Request)
	})

	router.GET("/", func(context *gin.Context) { context.Redirect(http.StatusMovedPermanently, "/controller") })

	if err != nil {
		log.Fatal(err)
	}
	return &ws
}

//Start starts the WebServer
func (ws *WebServer) Start() {
	go func() {
		log.Printf("[%s] %s", utilities.CreateColorString("WEB SERVER", color.FgHiBlue), "Avaiable on port:"+strconv.Itoa(ws.Port))
		err := ws.Router.Run(ws.Address + ":" + strconv.Itoa(ws.Port))
		log.Fatal(err.Error())
	}()

	go func() {
		for {
			//ws.listenInputChannel()
		}
	}()
}

/*

func (ws *WebServer) listenInputChannel() {

	select {
	case resp := <-ws.InputChannel:
		if serviceData, exists := (*ws.ServiceMap)[resp.ServiceName]; exists {
			//service exists
			protocolExists := false
			sendToWs := false

			for i := 0; i < len(serviceData.Protocols); i++ {
				protocolData := &serviceData.Protocols[i]
				if protocolData.Protocol.Type == resp.Protocol.Type && protocolData.Protocol.Server == resp.Protocol.Server && protocolData.Protocol.Port == resp.Protocol.Port {

					if protocolData.Err != nil && resp.Error != nil {
						if protocolData.Err.Error() != resp.Error.Error() {
							protocolData.Err = resp.Error
							sendToWs = true
						}
					} else {
						if protocolData.Err != resp.Error {
							protocolData.Err = resp.Error
							sendToWs = true
						}
					}
					protocolExists = true
				}

			}

			if !protocolExists {
				serviceData.Protocols = append(serviceData.Protocols, data.ProtocolData{
					Protocol: resp.Protocol,
					Err:      resp.Error,
				})

			}

			if sendToWs {
				ws.ServerSocket.BroadcastToRoom("/", ROOM_SERVICES_LISTENERS, EVENT_SERVICE_CHANGE, serviceData)
			}

			(*ws.ServiceMap)[resp.ServiceName] = serviceData

		}
	}

}
*/

func getRobotPosition(context *gin.Context) {

	postion := robotInstance.GetPosition()
	context.JSON(http.StatusOK, postion)
}

func setRobotPosition(context *gin.Context) {

	var position map[string]int16
	//value, _ := c.Request.GetBody()
	//fmt.Print(value)

	err := context.ShouldBindJSON(&position)
	if err == nil {
		newPosition := models.Position{
			X:     position["x"],
			Y:     position["y"],
			Angle: position["angle"],
		}
		robotInstance.SetPosition(newPosition)
		context.JSON(http.StatusOK, gin.H{"error": false})

	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func sendStop(context *gin.Context) {

	robotInstance.StopMotors()
	context.JSON(http.StatusOK, gin.H{"error": false})
}

func getRobotSpeed(context *gin.Context) {

	speed := robotInstance.Speed
	context.JSON(http.StatusOK, gin.H{"speed": speed})
}

func setRobotSpeed(context *gin.Context) {
	var json map[string]int16
	//value, _ := c.Request.GetBody()
	//fmt.Print(value)

	err := context.ShouldBindJSON(&json)
	if err == nil {
		newSpeed := json["speed"]
		errRobot := robotInstance.SetSpeed(newSpeed)
		if errRobot != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			context.JSON(http.StatusOK, gin.H{"error": false})
		}

	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func robotForwardDistance(context *gin.Context) {

	var json map[string]int16
	//value, _ := c.Request.GetBody()
	//fmt.Print(value)

	err := context.ShouldBindJSON(&json)
	if err == nil {
		distance := json["distance"]
		errRobot := robotInstance.ForwardDistance(distance)
		if errRobot != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			context.JSON(http.StatusOK, gin.H{"error": false})
		}

	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}

func robotForwardPoint(context *gin.Context) {

	var json map[string]int16
	//value, _ := c.Request.GetBody()
	//fmt.Print(value)

	err := context.ShouldBindJSON(&json)
	if err == nil {
		errRobot := robotInstance.ForwardToPoint(json["x"], json["y"])
		if errRobot != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			context.JSON(http.StatusOK, gin.H{"error": false})
		}

	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}

func robotRelativeRotation(context *gin.Context) {

	var json map[string]int16

	err := context.ShouldBindJSON(&json)
	if err == nil {
		errRobot := robotInstance.RelativeRotation(json["angle"])
		if errRobot != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			context.JSON(http.StatusOK, gin.H{"error": false})
		}

	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}

func robotAbsoluteRotation(context *gin.Context) {

	var json map[string]int16

	err := context.ShouldBindJSON(&json)
	if err == nil {
		errRobot := robotInstance.AbsoluteRotation(json["angle"])
		if errRobot != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			context.JSON(http.StatusOK, gin.H{"error": false})
		}

	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}

func newServerSocket() *socketio.Server {

	serverSocket := socketio.NewServer(nil)

	serverSocket.OnConnect("/", func(s socketio.Conn) error {
		s.Join("robot_controller")
		s.SetContext("")
		return nil
	})

	serverSocket.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	serverSocket.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	serverSocket.OnDisconnect("/", func(s socketio.Conn, reason string) {
		s.Close()
		//fmt.Println("closed", reason)
		//fmt.Println("connected remain:", serverSocket.Count())
	})

	return serverSocket
}

//NewMelodyWebSocket return melody server
func NewMelodyWebSocket() *melody.Melody {

	server := melody.New()

	server.HandleConnect(func(s *melody.Session) {

		log.Printf("[%s] %s", utilities.CreateColorString("WEB SOCKET", color.FgHiMagenta), "Client connected!")

		go func() {
			for true {
				if s.IsClosed() {
					//log.Printf("[%s] %s", utilities.CreateColorString("WEB SOCKET", color.FgHiMagenta), "Session closed!")
					return
				}
				position := robotInstance.GetPosition()
				wsMessage := models.WebSocketMessage{
					Command: "position",
					Payload: position,
				}
				message, err := json.Marshal(wsMessage)
				if err == nil {
					s.Write(message)
					//log.Printf("[%s] %s", utilities.CreateColorString("WEB SOCKET", color.FgHiMagenta), "Position sent!")
				}
				time.Sleep(time.Millisecond * 20)
			}
		}()
	})

	server.HandleDisconnect(func(s *melody.Session) {
		log.Printf("[%s] %s", utilities.CreateColorString("WEB SOCKET", color.FgHiMagenta), "Client disconnected!")
	})

	server.HandleMessage(func(s *melody.Session, msg []byte) {
		message := models.WebSocketMessage{}
		err := json.Unmarshal([]byte(msg), &message)
		if err != nil {
			s.Write([]byte("Error : Message Format"))
			log.Printf("[%s] %s", utilities.CreateColorString("WEB SOCKET", color.FgHiRed), err)
		} else {
			log.Printf("[%s] %s", utilities.CreateColorString("WEB SOCKET", color.FgHiMagenta), "Client sent command: "+message.Command)
			ManageWebSocketMessages(message)
			server.Broadcast([]byte(message.Command))
		}

	})

	return server
}

//ManageWebSocketMessages manage the websocket and socket.io messages
func ManageWebSocketMessages(msg models.WebSocketMessage) {

}

//GinMiddleware manage the cors
func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}
