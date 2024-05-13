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

type Money float64

type Player struct {
	Name string
	hand Hand
	bet  Money           // bet on current hand
	bets map[Money]Money // key:bet value:return on bet
}

func (p *Player) show() {
	fmt.Println()
	fmt.Printf("%s, Your cards are:\n", p.Name)
	for _, c := range p.hand {
		fmt.Println(">>>", c.Name())
	}
}

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
		bets: make(map[Money]Money),
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
	// Players place bets
	for i := 0; i < len(g.players); i++ {
		var bet float64
		fmt.Printf("%s bets:> ", g.players[i].Name)
		fmt.Scanf("%f\n", &bet)
		g.players[i].bet = Money(bet)
	}

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
	fmt.Println("\nFor dealer:")
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

func (g *Game) PlayerTurns() map[*Player]int {
	scores := make(map[*Player]int)
	for i, player := range g.players {
		player.show()
		score := player.hand.Score()
		if score == 21 {
			// blackjack
			fmt.Printf("Blackjack for %s. You get 150%% on your bet!\n", player.Name)
			g.players[i].bets[player.bet] = player.bet * 1.5
			continue
		}
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
			g.players[i].bets[player.bet] = 0
		} else {
			scores[&g.players[i]] = score
		}
	}
	return scores
}

func (g *Game) FindWinners(dealerScore int, playerScores map[*Player]int) {
	fmt.Println()
	if dealerScore > 21 {
		fmt.Println("Dealer gets busted.")
		for p, _ := range playerScores {
			fmt.Printf("%s gets double.\n", p.Name)
			p.bets[p.bet] = p.bet * 2
		}
	} else {
		for p, s := range playerScores {
			if s > dealerScore {
				fmt.Printf("%s gets double.\n", p.Name)
				p.bets[p.bet] = p.bet * 2
			} else if s < dealerScore {
				fmt.Printf("%s loses their bet!\n", p.Name)
				p.bets[p.bet] = 0
			} else {
				fmt.Printf("Tie for %s\n", p.Name)
				p.bets[p.bet] = p.bet
			}
		}
	}
}

func (g *Game) start() {
	g.DealCards()
	scores := g.PlayerTurns()

	// if players are left, continue otherwise end game
	if len(scores) == 0 {
		return
	}

	// show dealer hand
	dealerScore := g.ShowDealerHand()
	fmt.Printf("Dealer scores: %d\n", dealerScore)

	// determine winners
	g.FindWinners(dealerScore, scores)
}

func (g *Game) PrepareForNextRound() {
	g.dealer.hand = make([]deck.Card, 0)
	for i, _ := range g.players {
		g.players[i].hand = make([]deck.Card, 0)
		g.players[i].bet = Money(0)
	}
	return
}

func main() {
	numDecks := 1
	game := NewGame(numDecks)
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
		fmt.Println("\nPress q to quit the game or any other key to continue")
		fmt.Scanf("%s\n", &endGame)
		if strings.ToLower(endGame) == "q" {
			break
		}
		game.PrepareForNextRound()
	}

	fmt.Println()
	// Print money earned by players
	for _, player := range game.players {
		var bets, returns Money
		for b, r := range player.bets {
			bets += b
			returns += r
		}
		fmt.Printf("\nFor %s: Invested %f and made %f\n", player.Name, bets, returns)
	}
}
