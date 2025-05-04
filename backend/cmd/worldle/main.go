package main

import (
	"net/http"

	"github.com/MaarceloLuiz/worldle-replica/pkg/api"
	"github.com/sirupsen/logrus"
)

func main() {
	http.HandleFunc("/api/newgame", api.NewGameHandler)
	http.HandleFunc("/api/silhouette", api.SilhouetteHandler)
	http.HandleFunc("/api/territories", api.AllTerritoriesHandler)
	http.HandleFunc("/api/guess", api.GuessHandler)

	logrus.Info("Starting server on :8080")
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
