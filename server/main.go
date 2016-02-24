package main

import(
	"fmt"
	"strconv"
)

func main(){
	server := Server{[]Room{}, []PlayerConnection{}}
	fmt.Println(server) 
	
	deck := createDeck()
	
	// so time to try this out
	
	playerOne := Player{
		name: "Joseph",
		hand: []Card{},
		discardPile: []Card{},
		score: 0,
		protected: false,
		active: true,
	}
	
	playerTwo := Player{
		name: "Juliette",
		hand: []Card{},
		discardPile: []Card{},
		score: 0,
		protected: false,
		active: true,
	}
	
	round := Round{
		players: make([]*Player, 2),
		turnIndicator: 0,
		deck: deck,
		discards: DiscardPile{},
		finished: false,
	}
	
	round.players[0] = &playerOne
	round.players[1] = &playerTwo
	
	round.dealHands()
	
 	for !deck.IsEmpty(){
		// for each player, if not eliminated
		if round.over(){
			break
		}
		for _, p := range round.players{
			if(p.active && !round.over()){
				// so start by making the player draw a card
				gameOver := !round.drawCard(p)
				
				// then you inform the player of the valid targets
				targets := round.identifyPossibleTargets()
				
				// if possible, you make the player choose a card
				playedCard := ""
				fmt.Println("Select the card you'd like to play: ")
				fmt.Scanln(&playedCard)
				cardnum, _ := strconv.Atoi(playedCard)
				
				// for loop to test and see if the card is actually in their hand
				// switch cases to determine if they need extra parameters
				// for loop to make sure they chose the right player to target
				// should be able to just test it and not necessarily play it
				round.playCard(p, cardnum, nil, 0)
			}
		}
	}
}
