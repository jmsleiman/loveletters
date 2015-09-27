package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	//"golang.org/x/net/websocket"
	//"io"
	//"reflect"
)

// Mgs is the message format passed between client and server.
type Msg struct{
	Name string `json:"name"`
	Type string `json:"type"` // "chat", "turn", "connect", "disconnect"
	Data string `json:"data"`
}

// Room is a manager for a game of loveletters.
type Room struct{
	players []PlayerConnection
	I chan Msg
}

type PlayerConnection struct{
	O chan string
	player Player
}

// Server holds all the information about the loveletter server.
type Server struct{
	rooms []Room
	players []PlayerConnection
}

func (s *Server) handle(c *gin.Context) {
	s.wshandler(c.Writer, c.Request)
}

func (s *Server) findSuitableRoom() Room{
	fmt.Println("finding suitable room")
	return s.rooms[0] // dunno for now
}

func (s *Server) wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	
	// so at this point we can have some fun, yeah
	fmt.Println("Expecting a message...")
	_, msg, err := conn.ReadMessage()
	
	res := Msg{}
	
	json.Unmarshal(msg, &res)
	
	fmt.Println(err)
	fmt.Println(res)
	fmt.Println(res.Name)
	
	// maybe take in a message at this point?
	// grab the room and grab the user id?
	// first message would be like:
	// {channel: channel, nick: myNick}
	room := s.findSuitableRoom()
	
	// if you have to create a room, make sure to give it an I:
	// I := make(chan string) // unbuffered channels
	
	player := PlayerConnection{}
	s.players = append(s.players, player)
	
	//make new channels
	O := make(chan string) // more of a basket than a tube
	player.O = O
	
	// assign O to player
	// assign I to room
	
	go func() {
		// loop over channel’s contents
		// quit only if the channel closes
		for s := range O {
			conn.WriteMessage(websocket.TextMessage, []byte(s))
		}
	}()
	
	go func(){
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			
			res := Msg{}
			json.Unmarshal(msg, &res)
			
			room.I <- res
		}
	}()
	
	// keep note: there's a limit to how many players can join a room
	// make([]Room, 4) when creating the room?
	// and make them shut when the players start playing?
	
	room.players = append(room.players, player)
	
}

// Loop does the game loop
// slap this in a goroutine?
func (r Room) Loop(){
	currentplayer := 0
	for m := range r.I{//m is a Msg
		switch m.Type{
			case "turn":
				if m.Name == r.players[currentplayer].player.name{
					fmt.Println("uhhhh")
				}
			case "chat":
				//chat stuf
			case "connect":
				// do we need to do this?
				// idek, its just an example
				// i intend to reuse this code...
		
			case "disconnect":
				//probably need to do this right?
		}
		currentplayer = (currentplayer+1)%len(r.players)
	}
}

func main() {
	server := Server{}
	
	r := gin.Default()
	r.LoadHTMLFiles("index_debug.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index_debug.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		server.wshandler(c.Writer, c.Request)
	})
	
	// we also need to figure out how to loop the game itself, so...

	r.Run("localhost:8080")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
