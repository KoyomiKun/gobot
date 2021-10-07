package texas

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	CardFormat = "花色:%s 数值:%c\n"
)

var (
	colors = []string{"Spade", "Heart", "Club", "Diamond"}
)

type Card struct {
	Num        byte
	ColorIndex int
}

func (c Card) String() string {
	return fmt.Sprintf(CardFormat, colors[c.ColorIndex], c.Num)
}

type Cards struct {
	Content  []Card
	TopIndex int
}

func NewCards() *Cards {
	cards := make([]Card, 0, 52)
	for i := 1; i <= 10; i++ {
		for j := 0; j < 4; j++ {
			cards = append(cards, Card{
				Num:        byte(i) + '0',
				ColorIndex: j,
			})
		}
	}
	for i := 11; i <= 13; i++ {
		for j := 0; j < 4; j++ {
			cards = append(cards, Card{
				Num:        byte(i) - 11 + 'J',
				ColorIndex: j,
			})
		}
	}
	return &Cards{
		TopIndex: len(cards) - 1,
		Content:  cards,
	}
}

func (c *Cards) Shuffle() {
	c.TopIndex = len(c.Content) - 1
	rand.Seed(time.Now().UnixNano())
	for i := range c.Content {
		j := rand.Intn(i + 1)
		c.Content[i], c.Content[j] = c.Content[j], c.Content[i]
	}
}

func (c *Cards) Pop() *Card {
	card := c.Content[c.TopIndex]
	c.TopIndex--
	return &card
}
