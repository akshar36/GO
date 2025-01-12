package main

import (
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	d := newDeck()

	if len(d) != 56 {
		t.Errorf("Expected deck of length 56, but got %v", len(d))
	}

	if d[0] != "Ace of Spades" {
		t.Errorf("Expected first card to be Ace of Spades, but got %v", d[0])
	}

	if d[len(d)-1] != "Joker" {
		t.Errorf("Expected last card to be Joker, but got %v", d[len(d)-1])
	}
}

func TestSaveToFileAndNewDeckFromFile(t *testing.T) {
	os.Remove("_decktesting")
	deck := newDeck()
	deck.saveDeckToFile("_decktesting")
	loadedDeck := readDeckFromFile("_decktesting")

	if len(loadedDeck) != 56 {
		t.Errorf("Expected loaded deck to be of length 56, but got %v", len(loadedDeck))
	}
	os.Remove("_decktesting")
}
