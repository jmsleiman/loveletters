package main

import(
	"errors"
	"math/rand"
	"time"
	"fmt"
	"strconv"
	log "github.com/Sirupsen/logrus"
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

func (d* Deck) IsEmpty() (bool){
	return d.cards.IsEmpty()
}

func (d* Deck) dealCard() (Card, error){
	// need to check if the stack is empty. if so, return error
	dealtCard, err := d.cards.Pop()
	if(err != nil){
		return Card{}, errors.New("No cards can be dealt")
	} else{
		return dealtCard.(Card), nil
	}
}

type DiscardPile struct{
	cards Stack
}

func (d* DiscardPile) add (c Card){
	d.cards.Push(c)
}

func (d* DiscardPile) Peek()(c Card){
	return d.cards.Peek().(Card)
}

type Move struct{
	Using int
	Target int
	Action string
}

// Player is a person that plays a game of loveletters.
type Player struct{
	name string
	hand []Card
	discardPile []Card
	score int
	protected bool
	active bool
}

func (p* Player) currentCardRank()(int){
	return p.hand[0].rank
}

func (p* Player) addCard(c Card){
	p.hand = append(p.hand, c)
}

func (p* Player) discardCard(n int)(Card){
	tmp := p.hand[n]
	
	if(tmp.rank == 8){
		// player just lost the game!
		p.active = false
	}
	
	p.hand = append(p.hand[:n], p.hand[n+1:]...)
	return tmp
}

type Round struct{
	players []*Player
	turnIndicator int
	deck Deck
	setAside Card
	discards DiscardPile
	finished bool
}

/*
 * 8 : If you discard the Princess (8), you lose the game.
 * 7 : If you have the Countess and either the King (6) or Prince (5) you must discard the Countess. You can discard the Countess (7) regardless.
 * 6 : Trade the card in your handw ith the card held by another player of your choice. If you have no options, this card does nothing.
 * 5 : When you discard the Prince (5), choose one player still in the round (including yourself). That play discards his or her hand (not applying its effect), and draws a new card. If the deck is empty, they take the card that was removed at the start of the round.
 * 4 : Immunity from other players' cards.
 * 3 : When discarded, choose one other player. Secretly compare hands. Lower hand is knocked out. Tie, nothing happens. If all are protected by the Handmaid (4), nothing happens.
 * 2 : You can look at one other player's hand. Do not reveal the hand.
 * 1 : Choose a player and name a card (other than Guard (1)). If that player has that card, they are knocked out of the round. If all are protected, this card does nothing.
 * 
 * */

func createDeck() (Deck){
	d := []Card{
		{name:"Princess", rank: 8, description:"princess"},
		{name:"Countess", rank: 7, description:"countess"},
		{name:"King", rank: 6, description:"king"},
		{name:"Prince", rank: 5, description:"prince"},
		{name:"Prince", rank: 5, description:"prince"},
		{name:"Handmaid", rank: 4, description:"handmaid"},
		{name:"Handmaid", rank: 4, description:"handmaid"},
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
	/*
	myhand := Hand{}
	*/
	deck := Deck{}
	
	for i := 0; i < len(d); i++{
		deck.cards.Push(d[i])
	}
	
	return deck
}

func (round* Round) playerDraws(p* Player){
	outOfCards := round.drawCard(p)
	if(outOfCards == true){
		round.highestCardWins()
	}
}

func (round* Round) playCard(p1* Player, cardNum int, p2* Player, options int) (bool, error){
	validTargets := round.identifyPossibleTargets()
	
	if(p2 != nil){
		targetValid := false
		for _, person := range validTargets{
			if(p2 == person){
				targetValid = true
				break
			}
		}
		
		if(!targetValid){
			return false, errors.New("Could not target that player")
		}
		
	}
	
	cardobj := p1.discardCard(cardNum)
	
	switch cardobj.rank{
		case 1:
			fmt.Println("Played a Guard")
			fmt.Println("Select the target player")
			target := ""
			accusation := ""
			fmt.Scanln(&target)
			fmt.Println("Choose an accusation")
			fmt.Scanln(&accusation)
			
			// remember, can't accuse of being a guard!!
			// remember, can't accuse a protected target!!
			// if there's no valid target, card fizzles
			
			enemy_rank, _ := strconv.Atoi(accusation)
			ti, _ := strconv.Atoi(target)
			
			if(round.players[ti].hand[0].rank == enemy_rank){
				fmt.Println("Player accused " + target + " of holding the " + accusation + " card... Busted!")
				round.players[ti].active = false
			} else {
				fmt.Println("Player accused " + target + " of holding the " + accusation + " card... No luck!")
			}
			
		case 2:
			fmt.Println("Played a Priest")
			fmt.Println("Select the target player")
			target := ""
			fmt.Scanln(&target)
			ti, _ := strconv.Atoi(target)
			
			// if no valid targets, then the player can just discard
			
			for{
				
				if ( round.players[ti].protected ){
					fmt.Println("That player is protected. Choose a new player.")
					fmt.Scanln(&target)
					ti, _ = strconv.Atoi(target)
				} else {
					break
				}
			}
			
			fmt.Print("Player has: ")
			fmt.Println(round.players[ti].hand[0])
			
		case 3:
			fmt.Println("Played a Baron")
			fmt.Println("Select the target player")
			target := ""
			fmt.Scanln(&target)
			ti, _ := strconv.Atoi(target)
			
			// if there's no valid target, card fizzles
			
			for{
				if ( round.players[ti].protected ){
					fmt.Println("That player is protected. Choose a new player.")
					fmt.Scanln(&target)
					ti, _ = strconv.Atoi(target)
				} else {
					break
				}
			}
			
			fmt.Print("Comparing hands with Target... ")
			fmt.Print(p1.hand[0])
			fmt.Print(" vs ")
			fmt.Print(round.players[ti].hand[0])
			fmt.Print("... winner: ")
			
			if (p1.hand[0].rank == round.players[ti].hand[0].rank){
				fmt.Print("... a tie!")
			} else if (p1.hand[0].rank > round.players[ti].hand[0].rank){
				fmt.Print("... you!")
				round.players[ti].active = false
			} else{
				fmt.Print("... not you :(")
				p1.active = false
			}
			
			
			
		case 4:
			fmt.Println("Played a Handmaid")
			fmt.Println("Protected!!")
			p1.protected = true
			
		case 5:
			fmt.Println("Played a Prince")
			fmt.Println("Select the target player")
			target := ""
			fmt.Scanln(&target)
			ti, _ := strconv.Atoi(target)
			
			for{
				if ( round.players[ti].protected ){
					fmt.Println("That player is protected. Choose a new player.")
					fmt.Scanln(&target)
					ti, _ = strconv.Atoi(target)
				} else {
					break
				}
			}
			
			round.players[ti].discardCard(0)
			outOfCards := round.drawCard(round.players[ti])
			
			if(outOfCards == true){
				round.highestCardWins()
			}
			
		case 6:
			fmt.Println("Played the King")
			fmt.Println("Select the target player")
			target := ""
			fmt.Scanln(&target)
			ti, _ := strconv.Atoi(target)
			
			for{
				if ( round.players[ti].protected ){
					fmt.Println("That player is protected. Choose a new player.")
					fmt.Scanln(&target)
					ti, _ = strconv.Atoi(target)
				} else {
					break
				}
			}
			
			tmp := p1.hand
			p1.hand = round.players[ti].hand
			round.players[ti].hand = tmp
			
			log.WithFields(log.Fields{
				"player1": p1.hand,
				"player2": round.players[ti].hand,
			}).Info("players' hands")
			
			
		case 7:
			fmt.Println("Played the Countess")
			
			
		case 8:
			fmt.Println("Played the Princess")
			fmt.Println("You lose!!")
			p1.active = false
			
	}
	
	return true, nil
	
}

func (round* Round) playerTurn(p* Player){
// 	playedCard := ""
// 	fmt.Println("Select the card you'd like to play: ")
// 	fmt.Scanln(&playedCard)
//	cardnum, _ := strconv.Atoi(playedCard)
	
	// remember, if the player has the 7 and the 5, 6, or 8, she must discard the 7 immediately, and that's her move!
	
	
}

func (r* Round) identifyPossibleTargets() ([] *Player){
	s := make([]*Player, 0)
	
	for _, p := range r.players{
		if(p.active && !p.protected){
			s = append(s, p)
		}
	}
	
	log.Print("the size of the array was: ")
	log.Println(len(s))
	
	return s
	
}

func (r* Round) drawCard(p* Player) (bool){
	drawnCard, err := r.deck.dealCard()
	
	log.WithFields(log.Fields{
		"player": p.name,
		"player-hand": p.hand,
		"drawn-card": drawnCard,
	}).Info("player has drawn")
	
	if err != nil{
		log.WithFields(log.Fields{
			"player": p.name,
		}).Info("no more draws can be made")
		
		r.highestCardWins()
		return false
		
		// then there's no more drawing to be had!
	}else{
		p.addCard(drawnCard)
		return true
	}
}

func (r* Round) dealHands(){
	tmp, err := r.deck.dealCard()
	
	if(err != nil){
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Dealing hands failed!!")
	} else{
		r.setAside = tmp
		for _, player := range r.players{
			r.drawCard(player)
		}
	}
}

func (r* Round) highestCardWins(){
	// it's possible to have a tie in the score. if so, make sure the player with the highest sum in their graveyard is the winner.
	
	highScore := -1
	highPlayer := &Player{}
	
	for _, p := range r.players{
		if p.currentCardRank() > highScore {
			highScore = p.currentCardRank()
			highPlayer = p
		} else {
			p.active = false
		}
	}
	
	// remember if there's a tie, it's the sum of the played cards that decide
	
	if !r.over(){
		log.Println("game didn't end properly?")
	} else {
		highPlayer.score += 1
	}
	
}

func (r* Round) over() (bool){
	log.Println("testing if the round is over...")
	tally := 0
	
	for _, p := range r.players{
		if(p.active == true){
			tally = tally + 1
		}
	}
	
	if(tally > 1){
		r.finished = false
	} else {
		r.finished = true
	}
	
	if(r.finished){
		log.Println("the round is over!...")
		return true
	} else {
		log.Println("the round is not!...")
		return false
	}
	
}

	/*
	fmt.Println(myhand)
	t := deck.cards.Pop().(Card)
	myhand.cards = append(myhand.cards, t)
	fmt.Println(deck.cards.Peek())
	fmt.Println(myhand)*/






