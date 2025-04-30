package game

import (
	"errors"
	"sync"

	"github.com/MaarceloLuiz/worldle-replica/pkg/geography/silhouettes"
)

var ErrGameNotInitialized = errors.New("game not initialized")

type GameState struct {
	Country string
	Image   []byte
	Mutex   sync.RWMutex
}

var State = GameState{}

func StartNewGame() error {
	State.Mutex.Lock()
	defer State.Mutex.Unlock()

	country, err := silhouettes.GetRandomCountry()
	if err != nil {
		return err
	}

	img, err := silhouettes.FetchSilhouette(country)
	if err != nil {
		return err
	}

	State.Country = country
	State.Image = img
	return nil
}

func GetCurrentSilhouette() ([]byte, error) {
	State.Mutex.RLock()
	defer State.Mutex.RUnlock()

	if State.Image == nil {
		return nil, ErrGameNotInitialized
	}

	return State.Image, nil
}
