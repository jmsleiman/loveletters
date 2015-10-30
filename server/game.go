package main

import(
	"fmt"
	"math/rand"
	"time"
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
}

func setup(){
	d := []Card{
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
	/*
	myhand := Hand{}
	*/
	deck := Deck{}
	
	for i := 0; i < len(d); i++{
		deck.cards.Push(d[i])
	}
	
	fmt.Println(deck.cards)
	/*
	fmt.Println(myhand)
	t := deck.cards.Pop().(Card)
	myhand.cards = append(myhand.cards, t)
	fmt.Println(deck.cards.Peek())
	fmt.Println(myhand)*/
}
