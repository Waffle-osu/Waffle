package main

import (
	"fmt"
	"io/ioutil"
)

func BeatmapImporter(songsDir string) {
	beatmapSetDirectories, err := ioutil.ReadDir(songsDir)

	if err != nil {
		fmt.Printf("Failed to import beatmaps: %s", err.Error())
	}

	for _, directory := range beatmapSetDirectories {
		
	}

	fmt.Printf("Beatmapsets Found: %d\n", len(beatmapSetDirectories))
}
