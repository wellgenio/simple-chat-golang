package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type IHandler interface {
	CreateRoom(c *gin.Context)
	JoinRoom(c *gin.Context)
	GetRooms(c *gin.Context)
	GetClients(c *gin.Context)
}

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{hub: h}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RoomsRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomsRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomsRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

type ClientRes struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	clients := make([]ClientRes, 0)
	roomID := c.Param("id_room")

	for _, c := range h.hub.Rooms[roomID].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, gin.H{"clients": clients})
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("id_room")
	clientID := c.Query("user_id")
	username := c.Query("username")

	cl := &Client{
		ID:       clientID,
		Username: username,
		RoomID:   roomID,
		Conn:     conn,
		Message:  make(chan *Message, 10),
	}

	msg := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- msg

	go cl.writeMessage()
	cl.readMessage(h.hub)
}
