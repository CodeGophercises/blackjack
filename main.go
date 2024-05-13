package main

import (
	"fmt"
	"strings"

	"github.com/CodeGophercises/blackjack/scoring"
	"github.com/CodeGophercises/deck_of_cards/deck"
)

type Hand []deck.Card

func (h *Hand) Score() int {
	score := 0
	for _, c := range *h {
		score += scoring.GetCardScore(c)
	}
	return score
}

type Player struct {
	Name string
	hand Hand
}

func (p *Player) show() {
	fmt.Printf("%s, Your cards are:\n", p.Name)
	for _, c := range p.hand {
		fmt.Println(">>>", c.Name())
	}
}

// Just has some special rules
type Dealer Player

func (d *Dealer) show() {
	fmt.Printf("Dealer %s has:\n", d.Name)
	for _, c := range d.hand {
		fmt.Println(">>>", c.Name())
	}
}

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

// TODO: Add soft 17 rule
func (g *Game) ShowDealerHand() int {
	fmt.Println("For dealer:")
	score := g.dealer.hand.Score()
	for score <= 16 {
		c := g.GetNextCard()
		g.dealer.hand = append(g.dealer.hand, c)
		score += scoring.GetCardScore(c)
	}
	fmt.Println(">Cards:")
	for _, c := range g.dealer.hand {
		fmt.Println(">>", c.Name())
	}
	return score

}
func (g *Game) start() {
	g.DealCards()
	scores := make(map[string]int)
	for _, player := range g.players {
		player.show()
		score := player.hand.Score()
		fmt.Print("Enter 0 to stand and 1 to hit:> ")
		var input int
		fmt.Scanf("%d\n", &input)
		for input != 0 {
			c := g.GetNextCard()
			player.hand = append(player.hand, c)
			score += scoring.GetCardScore(c)
			if score > 21 {
				break
			}
			player.show()
			fmt.Print("Enter 0 to stand and 1 to hit:> ")
			fmt.Scanf("%d\n", &input)
		}
		fmt.Printf("%s's hand scores: %d\n", player.Name, score)
		if score > 21 {
			// busted
			fmt.Printf("Busted. You lose your bet %s\n", player.Name)
		} else {
			scores[player.Name] = score
		}
	}

	// if players are left, continue otherwise end game
	if len(scores) == 0 {
		return
	}

	// show dealer hand
	dealerScore := g.ShowDealerHand()
	fmt.Printf("Dealer scores: %d\n", dealerScore)
	// determine winners
	if dealerScore > 21 {
		fmt.Println("Dealer gets busted.")
		for p, _ := range scores {
			fmt.Printf("%s gets double.\n", p)
		}
	} else {
		for p, s := range scores {
			if s > dealerScore {
				fmt.Printf("%s gets double.\n", p)
			} else if s < dealerScore {
				fmt.Printf("%s loses their bet!\n", p)
			} else {
				fmt.Printf("Tie for %s\n", p)
			}
		}
	}
}

func (g *Game) PrepareForNextRound() {
	g.dealer.hand = make([]deck.Card, 0)
	for i, _ := range g.players {
		g.players[i].hand = make([]deck.Card, 0)
	}
	return
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
	fmt.Println()
	var endGame string
	for {
		game.start()
		fmt.Println("Press q to quit the game or any other key to continue")
		fmt.Scanf("%s\n", &endGame)
		if strings.ToLower(endGame) == "q" {
			break
		}
		game.PrepareForNextRound()
	}
}
