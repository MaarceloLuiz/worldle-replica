package api

import (
	"net/http"

	"github.com/MaarceloLuiz/worldle-replica/pkg/game"
)

func SilhouetteHandler(w http.ResponseWriter, r *http.Request) {
	silhouette, err := game.GetCurrentSilhouette()
	if err != nil {
		http.Error(w, "Failed to fetch silhouette", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(silhouette)
}

func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	err := game.StartNewGame()
	if err != nil {
		http.Error(w, "Failed to start new game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New game started"))
}
