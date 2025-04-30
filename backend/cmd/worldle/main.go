package main

import (
	"net/http"

	"github.com/MaarceloLuiz/worldle-replica/pkg/api"
	"github.com/sirupsen/logrus"
)

func main() {
	http.HandleFunc("/api/silhouette", api.SilhouetteHandler)
	http.HandleFunc("/api/newgame", api.NewGameHandler)

	logrus.Info("Starting server on :8080")
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
