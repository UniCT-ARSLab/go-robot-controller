package webserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/arslab/robot_controller/utilities"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rakyll/statik/fs"

	_ "github.com/arslab/robot_controller/webserver/statik" //static file system
)

//WebServer reppresents the WebServer and WebScoket server instance
type WebServer struct {
	Address      string
	Port         int
	Ssl          bool
	Router       *gin.Engine
	ServerSocket *socketio.Server
}

//NewWebServer returns a new WebServer
func NewWebServer(serverConfPath string) *WebServer {

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
	router.Use(cors.New(configCors))
	router.Use(gin.Recovery())

	ws := WebServer{
		Address:      "0.0.0.0",
		Port:         9999,
		Ssl:          false,
		Router:       router,
		ServerSocket: serverSocket,
	}

	// Register REST
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	staticGroup := router.Group("/controller")
	staticGroup.StaticFS("/", statikFS)

	//apiGroup := router.Group("/api")
	//apiGroup.GET("/connection/status", func(context *gin.Context) { getServicesList(context, &ws) })
	//apiGroup.GET("/system", func(context *gin.Context) { getSystemInformation(context) })

	router.GET("/socket.io/*any", gin.WrapH(serverSocket))
	router.POST("/socket.io/*any", gin.WrapH(serverSocket))

	router.GET("/", func(context *gin.Context) { context.Redirect(http.StatusMovedPermanently, "/controller") })

	if err != nil {
		log.Fatal(err)
	}
	return &ws
}

//Start starts the WebServer
func (ws *WebServer) Start() {
	go func() {
		log.Printf("[%s] %s", utilities.CreateColorString("Web Panel", color.FgHiCyan), "avaiable on port:"+strconv.Itoa(ws.Port))
		err := ws.Router.Run(ws.Address + ":" + strconv.Itoa(ws.Port))
		log.Fatal(err.Error())
	}()

	go func() {
		for {
			//	ws.listenInputChannel()
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

func newServerSocket() *socketio.Server {

	serverSocket, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

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
