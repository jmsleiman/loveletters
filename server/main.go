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

type Card struct{
	name string
	rank uint16
	description string
}

type Deck struct{
	cards Stack
}

type Hand struct{
	cards[] Card
}

type Msg struct{
	Name string `json:"name"`
	Type string `json:"type"` // "chat", "turn", "connect", "disconnect"
	Data string `json:"data"`
}

type Room struct{
	players []Player
	I chan Msg
}

type Player struct{
	name string
	score uint8
	O chan string
}

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
	var d = []Card{}
	d = append(d, Card{name:"Princess", rank: 8, description:"princess"})
	d = append(d, Card{name:"Countess", rank: 7, description:"countess"})
	d = append(d, Card{name:"King", rank: 6, description:"king"})
	d = append(d, Card{name:"Prince", rank: 5, description:"prince"})
	d = append(d, Card{name:"Prince", rank: 5, description:"prince"})
	d = append(d, Card{name:"Handmaid", rank: 4, description:"priest"})
	d = append(d, Card{name:"Handmaid", rank: 4, description:"priest"})
	d = append(d, Card{name:"Baron", rank: 3, description:"baron"})
	d = append(d, Card{name:"Baron", rank: 3, description:"baron"})
	d = append(d, Card{name:"Priest", rank: 2, description:"priest"})
	d = append(d, Card{name:"Priest", rank: 2, description:"priest"})
	d = append(d, Card{name:"Guard", rank: 1, description:"guard"})
	d = append(d, Card{name:"Guard", rank: 1, description:"guard"})
	d = append(d, Card{name:"Guard", rank: 1, description:"guard"})
	d = append(d, Card{name:"Guard", rank: 1, description:"guard"})
	d = append(d, Card{name:"Guard", rank: 1, description:"guard"})
	
	rand.Seed(time.Now().UTC().UnixNano())
	
	for i := range d {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
	
	var myhand = Hand{}
	
	deck := Deck{}
	
	for i := 0; i < len(d); i++{
		deck.cards.Push(d[i])
	}
	
	fmt.Println(deck.cards.Peek())
	fmt.Println(myhand)
	var t = deck.cards.Pop().(Card)
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
