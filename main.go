package main

import (
	"fmt"

	"github.com/CodeGophercises/deck_of_cards/deck"
)

type Player struct {
	Name string
	hand []deck.Card
}

func (p *Player) show() {
	fmt.Println("Your cards are:")
	for _, c := range p.hand {
		fmt.Println(">>>", c.Name())
	}
}

// Just has some special rules
type Dealer Player

type Game struct {
	curCardIndex int
	NumDecks     int
	cards        []deck.Card
	players      []Player
	dealer       Dealer
}

// nd: Number of decks
func NewGame(nd int) *Game {
	cards := deck.NewMultiDeck(nd, deck.Shuffle)
	return &Game{
		NumDecks: nd,
		cards:    cards,
	}
}

func (g *Game) AddDealer(name string) {
	dealer := Dealer{
		Name: name,
	}
	g.dealer = dealer
}

func (g *Game) AddPlayer(name string) {
	g.players = append(g.players, Player{
		Name: name,
	})
}

func (g *Game) GetNextCard() deck.Card {
	if g.curCardIndex >= len(g.cards) {
		panic("No more cards")
	}
	c := g.cards[g.curCardIndex]
	g.curCardIndex += 1
	return c
}

func (g *Game) DealCards() {
	for i := 0; i < 2; i++ {
		// to players first
		for j := 0; j < len(g.players); j++ {
			c := g.GetNextCard()
			g.players[j].hand = append(g.players[j].hand, c)
		}
		// to dealer
		c := g.GetNextCard()
		g.dealer.hand = append(g.dealer.hand, c)
	}
}

func (g *Game) start() {
	g.DealCards()

	for _, player := range g.players {
		player.show()
		fmt.Print("Enter 0 to stand and 1 to hit:> ")
		var input int
		fmt.Scanf("%d\n", &input)
		for input != 0 {
			player.hand = append(player.hand, g.GetNextCard())
			player.show()
			fmt.Print("Enter 0 to stand and 1 to hit:> ")
			fmt.Scanf("%d\n", &input)
		}
	}

	// show dealer hand
	fmt.Println("Dealer cards: ")
	for _, c := range g.dealer.hand {
		fmt.Println(">>>", c.Name())
	}
}
func main() {
	game := NewGame(1)
	var dealer string
	fmt.Print("dealer name:> ")
	fmt.Scanf("%s\n", &dealer)
	game.AddDealer(dealer)
	var nPlayers int
	fmt.Print("Number of players in game:> ")
	fmt.Scanf("%d\n", &nPlayers)
	for i := 0; i < nPlayers; i++ {
		var player string
		fmt.Print("Player name:> ")
		fmt.Scanf("%s\n", &player)
		game.AddPlayer(player)
	}
	game.start()
}
