package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"

	"github.com/wellgenio/simple-chat-golang/internal/ws"
)

var r *gin.Engine

func InitRouter(handler ws.IHandler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.Use(secure.New(secure.Config{
		ContentSecurityPolicy: "connect-src 'self';",
	}))

	r.POST("/ws/create_room", handler.CreateRoom)
	r.GET("/ws/rooms", handler.GetRooms)
	r.GET("/ws/join_room/:id_room", handler.JoinRoom)
	r.GET("/ws/rooms/:id_room/clients", handler.GetClients)
}

func Start(addr string) error {
	return r.Run(addr)
}
