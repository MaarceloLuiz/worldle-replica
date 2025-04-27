package main

import (
	"fmt"

	"github.com/MaarceloLuiz/worldle-replica/pkg/geography/silhouettes"
)

func main() {
	response, err := silhouettes.FetchSilhouette()
	if err != nil {
		fmt.Println("Error fetching silhouette:", err)
	}

	fmt.Println(string(response))
}
