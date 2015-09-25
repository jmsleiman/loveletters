package main

import(
	"fmt"
	"math/rand"
	"time"
	//"io"
	"net/http"
	//"reflect"
	"encoding/json"
	
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	//"golang.org/x/net/websocket"
)

// Card is a struct that holds the information of a card type in the game
type Card struct{
	name string
	rank int
	description string
}

// Deck is a stack of card (with a defined order).
type Deck struct{
	cards Stack
}

// Hand is all the cards a player currently holds.
type Hand struct{
	cards[] Card
}

// Mgs is the message format passed between client and server.
type Msg struct{
	Name string `json:"name, omitempty"`
	Type string `json:"type, omitempty"` // "chat", "turn", "connect", "disconnect"
	Data string `json:"data, omitempty"`
}

// Room is a manager for a game of loveletters.
type Room struct{
	players []Player
	I chan Msg
}

// Player is a person that plays a game of loveletters.
type Player struct{
	name string
	score int
	O chan string
}

// Server holds all the information about the loveletter server.
type Server struct{
	rooms []Room
	players []Player
}

func (s *Server) handle(c *gin.Context) {
	s.wshandler(c.Writer, c.Request)
}

func (s *Server) findSuitableRoom() Room{
	fmt.Println("finding suitable room")
	return s.rooms[0] // dunno for now
}

func (s *Server) wshandler(w http.ResponseWriter, r *http.Request) {
	// so at this point we can have some fun, yeah
	
	
	// maybe take in a message at this point?
	// grab the room and grab the user id?
	// first message would be like:
	// {channel: channel, nick: myNick}
	room := s.findSuitableRoom()
	
	// if you have to create a room, make sure to give it an I:
	// I := make(chan string) // unbuffered channels
	
	player := Player{}
	s.players = append(s.players, player)
	
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	
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
				if m.Name == r.players[currentplayer].name{
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
	d: = []Card{
		{name:"Princess", rank: 8, description:"princess"},
		{name:"Countess", rank: 7, description:"countess"},
		{name:"King", rank: 6, description:"king"},
		{name:"Prince", rank: 5, description:"prince"},
		{name:"Prince", rank: 5, description:"prince"},
		{name:"Handmaid", rank: 4, description:"priest"},
		{name:"Handmaid", rank: 4, description:"priest"},
		{name:"Baron", rank: 3, description:"baron"},
		{name:"Baron", rank: 3, description:"baron"},
		{name:"Priest", rank: 2, description:"priest"},
		{name:"Priest", rank: 2, description:"priest"},
		{name:"Guard", rank: 1, description:"guard"},
		{name:"Guard", rank: 1, description:"guard"},
		{name:"Guard", rank: 1, description:"guard"},
		{name:"Guard", rank: 1, description:"guard"},
		{name:"Guard", rank: 1, description:"guard"},
	}
	
	rand.Seed(time.Now().UTC().UnixNano())
	
	for i := range d {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
	
	myhand := Hand{}
	
	deck := Deck{}
	
	for i := 0; i < len(d); i++{
		deck.cards.Push(d[i])
	}
	
	fmt.Println(deck.cards.Peek())
	fmt.Println(myhand)
	t := deck.cards.Pop().(Card)
	myhand.cards = append(myhand.cards, t)
	fmt.Println(deck.cards.Peek())
	fmt.Println(myhand)
	
	server := Server{}
	
	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		server.wshandler(c.Writer, c.Request)
	})

	r.Run("localhost:8080")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
