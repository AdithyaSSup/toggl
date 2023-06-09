package dao

import (
	"OnlineDeck/pkg/models"
	"context"
	"errors"
	"github.com/google/uuid"
	"math/rand"
)

var ErrUUIDGeneration = errors.New("failed to generate unique UUID")
var ErrDeckNotFound = errors.New("decks does not exist for the given id")
var ErrInvalidDraw = errors.New(" insufficient cards available in the deck ")
var ErrInvalidUUID = errors.New(" id is not a valid uuid ")

type DeckDao struct {
	Decks map[uuid.UUID]*models.Deck
}

// NewDeckDao returns a new instance of DeckDao.
func NewDeckDao(decks map[uuid.UUID]*models.Deck) *DeckDao {
	return &DeckDao{Decks: decks}
}

// Create creates a new deck with the given cards and shuffle option.
func (d *DeckDao) Create(ctx context.Context, cards []models.Card, shuffle bool) (*models.Deck, error) {

	idr, err := d.CreateUUID(ctx, 0) // idr holds the reference to UUID
	if err != nil {
		return nil, err
	}

	id := *idr // id contains the actual id

	d.Decks[id] = &models.Deck{
		ID:    id,
		Cards: cards,
	}

	if shuffle {
		return d.Shuffle(ctx, d.Decks[id])
	}

	return d.Decks[id], err
}

// Get returns the deck with the given ID, or an error if it does not exist.
func (d *DeckDao) Get(_ context.Context, id string) (*models.Deck, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	value, exists := d.Decks[uid]
	if exists {
		return value, nil
	}

	return nil, ErrDeckNotFound
}

// Draw : removes the given number of cards from the deck with the given ID.
func (d *DeckDao) Draw(_ context.Context, deck *models.Deck, count int) ([]models.Card, error) {

	cards := deck.Cards

	if len(cards) < count {
		return nil, ErrInvalidDraw
	}

	res := make([]models.Card, count)
	for i := 0; i < count; i++ {
		res[i] = cards[len(cards)-1]
		cards = cards[:len(cards)-1]
	}
	// update the deck with the cards slice which has removed elements
	deck.Cards = cards

	return res, nil

}

// Shuffle shuffles the deck with the given ID.
func (d *DeckDao) Shuffle(_ context.Context, deck *models.Deck) (*models.Deck, error) {

	cards := deck.Cards
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	deck.Cards = cards
	deck.Shuffled = true

	return deck, nil
}

// CreateUUID generates a UUID which is always unique
func (d *DeckDao) CreateUUID(ctx context.Context, count int) (*uuid.UUID, error) {
	if count > 10 {
		return nil, ErrUUIDGeneration
	}

	id := uuid.New()

	_, exists := d.Decks[id]
	if exists {
		return d.CreateUUID(ctx, count+1)
	}
	return &id, nil
}
