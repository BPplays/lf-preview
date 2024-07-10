package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dhowden/tag"
)


func thumbnail_music(file string) *[]byte {
	// cache := get_thumbnail_cache_file(".bmp")

	var start time.Time
	if debug_time {
		start = time.Now()
	}
	// cache := filepath.Join(cacheFile, ".bmp")
	// cache := thumbnail_cache + ".bmp"
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	// if !fileExists(cache) {

	// }

	m, err := tag.ReadFrom(f)
	if err != nil {
		log.Fatal(err)
	}


	// fmt.Println(string(output))
	if debug_time {
		time_output = time_output + fmt.Sprintln("thumbnail_music time: ",time.Since(start))
	}
	return &m.Picture().Data
}