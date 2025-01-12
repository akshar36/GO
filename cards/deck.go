package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type deck []string //like typedef (not exactly, for OOP alternative)

func newDeck() deck {
	cards := deck{}
	cardSuits := []string{"Spades", "Diamonds", "Clubs", "Hearts"}
	cardValues := []string{"Ace ", "Two ", "Three ", "Four ", "Five ", "Six ", "Seven ", "Eight ", "Nine ", "Ten ", "Jack ", "Queen ", "King "}

	for _, suit := range cardSuits {
		for _, value := range cardValues {
			cards = append(cards, value+"of "+suit)
		}
	}

	count := 4
	for count > 0 {
		cards = append(cards, "Joker")
		count--
	}
	return cards
}

func (d deck) print() {
	for i, card := range d {
		fmt.Println(i, card)
	}
}

func (d deck) deal(handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

func (d deck) saveDeckToFile(fileName string) {
	err:= os.WriteFile(fileName, []byte(strings.Join([]string(d), ",")), 0666)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func readDeckFromFile(fileName string) (deck) {
	byteOp, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	strOp := string(byteOp)
	return deck(strings.Split(strOp,","))
}

func (d deck) shuffle(){
	seed := rand.NewSource(time.Now().UnixNano())
	rg := rand.New(seed)

	for i := range d{
		newPosition := rg.Intn(len(d) - 1)
		d[i], d[newPosition] = d[newPosition], d[i]
	}

}