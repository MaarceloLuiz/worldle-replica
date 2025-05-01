package api

import (
	"encoding/json"
	"net/http"

	"github.com/MaarceloLuiz/worldle-replica/pkg/game"
	geography "github.com/MaarceloLuiz/worldle-replica/pkg/geography/territories"
)

func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	err := game.StartNewGame()
	if err != nil {
		http.Error(w, "Failed to start new game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New game started"))
}

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

func AllTerritoriesHandler(w http.ResponseWriter, r *http.Request) {
	territories, err := geography.GetFormattedTerritoryNames()
	if err != nil {
		http.Error(w, "Failed to fetch territories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(territories)
	if err != nil {
		http.Error(w, "Failed to encode territories", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
