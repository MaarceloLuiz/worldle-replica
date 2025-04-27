package main

import (
	"fmt"
	"os"

	"github.com/MaarceloLuiz/worldle-replica/pkg/geography/silhouettes"
)

func main() {
	img, err := silhouettes.FetchSilhouette()
	if err != nil {
		fmt.Println("Error fetching silhouette:", err)
		return
	}

	if err := os.WriteFile("output.png", img, 0644); err != nil {
		fmt.Println("could not write file:", err)
		return
	}
	fmt.Println("OK — wrote out.png (size:", len(img), "bytes)")
}
