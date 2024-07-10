package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/barasher/go-exiftool"
)

func get_exif(file string) ([]exiftool.FileMetadata) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		// fmt.Printf("Error when intializing: %v\n", err)
		// return "", err
		fmt.Println("get_exif", err)
		log.Fatal("get_exif", err)
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(file)


	// for _, fileInfo := range fileInfos {
	// 	if fileInfo.Err != nil {
	// 		fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
	// 		continue
	// 	}

	// 	for k, v := range fileInfo.Fields {
	// 		fmt.Printf("[%v] %v\n", k, v)
	// 	}
	// }
	return fileInfos
}

func exif_fmt_gr(file string, tags [][]string, ch chan<- order_string, order int, wg *sync.WaitGroup) {
	defer wg.Done()

	if no_info {
		return
	}

	var start time.Time
	if debug_time {
		start = time.Now()
	}
	var output = order_string{order, ""}
	// output[1] = output[1] + fmt.Sprint("test")

	// ch <- fmt.Sprint("test")
	// ch <- fmt.Sprint(exif_fmt(file, tags))
	output.content = output.content + fmt.Sprintln(sep1)
	// output.content = output.content + "test"
	output.content = output.content + get_metadata(file, tags)
	output.content = output.content + fmt.Sprintln(sep1)
	ch <- output
	if debug_time {
		time_output = time_output + fmt.Sprintln("exif_fmt_gr time: ",time.Since(start))
	}
	// output := exif_fmt(file, tags)
	// gr_array[ar_index] = "test"
	// gr_array[1] = "test"
	// fmt.Println((*array)[ar_index])
	// fmt.Println(ar_index)
	// fmt.Println("testrgji")
}


func get_metadata(file string, tags [][]string) (string) {
	var output string

	cache := filepath.Join(get_metadata_cache_dir(), add_ext(get_hash(), ".json", get_cache_byte_limit()))

	if fileExists(cache) {
		cache_data, err := os.ReadFile(cache)
		if err != nil {
			fmt.Println("Error reading file:", err)
			log.Fatal(err)
		}

		var exif []exiftool.FileMetadata

		err = json.Unmarshal(cache_data, &exif)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			log.Fatal()
		}

		output = exif_fmt(exif, tags)
	} else {
		exif := get_exif(file)

		jsonData, err := json.Marshal(exif)
		if err != nil {
			fmt.Println("Error marshalling to JSON:", err)
			log.Fatal(err)
		}


		err = os.WriteFile(cache, jsonData, 0600)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			log.Fatal(err)
		}

		if dbg_print_exif {
			fmt.Println(exif)
		}
		output = exif_fmt(exif, tags)



	}


	return output
}


func exif_fmt(fileInfos []exiftool.FileMetadata, tags [][]string) (string) {

	var blocks []string
	var builder strings.Builder




	var tag_name string
	var tag_val string

	var ok bool
	var val any

	for _, tag_block := range tags {
		for _, tag := range tag_block {

			for _, fileInfo := range fileInfos {
				// output = output + fmt.Sprintln("fileInfos")
				val, ok = fileInfo.Fields[tag]
				// If the key exists
				if ok {

					tag_name = tag
					tag_val, ok = exif_key_map[tag]
					// If the key exists
					if ok {
						tag_name = tag_val
					}
					builder.WriteString(fmt.Sprintf("%v: %v\n", tag_name, val))
				}

			}

		}


		blocks = append(blocks, builder.String())
		builder.Reset()

	}


	return blocks_fmt(blocks)
}

