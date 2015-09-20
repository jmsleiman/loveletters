package main

import(
	"fmt"
	"math/rand"
	"time"
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
	
}