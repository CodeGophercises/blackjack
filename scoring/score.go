package scoring

import (
	"fmt"

	"github.com/CodeGophercises/deck_of_cards/deck"
)

func GetCardScore(c deck.Card) int {
	var score int
	switch {
	case c.Rank >= deck.Two && c.Rank <= deck.Ten:
		score = int(c.Rank)
	case c.Rank >= deck.Jack:
		score = 10
	default: // Ace
		fmt.Printf("For %s, choose score: 1 or 11 > ", c.Name())
		fmt.Scanf("%d\n", &score)
	}
	return score
}
