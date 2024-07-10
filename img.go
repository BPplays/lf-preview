package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func chafa_image(image *[]byte, width, height int) (string) {


	cmd := exec.Command("chafa", fmt.Sprintf("--font-ratio=%s", userOpenFontRatio))
	cmd.Args = append(cmd.Args, chafaFmt...)
	cmd.Args = append(cmd.Args, chafaDither...)
	cmd.Args = append(cmd.Args, chafaColors...)
	cmd.Args = append(cmd.Args, "--color-space=din99d", "--scale=max", "-w", "9", "-O", "9", "-s", get_geometry(width, height), "--animate", "false")
	cmd.Args = append(cmd.Args, "--symbols", "block+border+space-wide+inverted+quad+extra+half+hhalf+vhalf")
	cmd.Args = append(cmd.Args, "--polite", "on", "--color-extractor=median")


	pipe, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating pipe:", err)
		os.Exit(1)
	}

	go func() {
		defer pipe.Close() // Close the pipe when done

		// Write data to the command's standard input
		_, err := pipe.Write(*image)
		if err != nil {
			fmt.Println("Error writing to pipe:", err)
			os.Exit(1)
		}
	}()


	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output), err)
		log.Fatal(string(output), err)
	}

	// fmt.Println(string(output))
	return string(output)
}

func isSVG(filename string) bool {
	// Check if the file extension is SVG
	if strings.HasSuffix(strings.ToLower(filename), ".svg") {
		return true
	}

	if strings.HasSuffix(strings.ToLower(filename), ".svgz") {
		return true
	}

	return false
}

func isSVGz(filename string) bool {
	// Check if the file extension is SVG

	return strings.HasSuffix(strings.ToLower(filename), ".svgz")
}

func svgz_to_svg(svgzData *[]byte) (*[]byte) {
	// Create a bytes reader from the input SVGZ data
	svgzReader := bytes.NewReader(*svgzData)

	// Create a gzip reader
	gzReader, err := gzip.NewReader(svgzReader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer gzReader.Close()

	// Read the decompressed SVG data into a byte slice
	svgData, err := io.ReadAll(gzReader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &svgData
}



func svg_to_png(input *[]byte) *[]byte {

	conv := svg_to_png_new()
	conv.SetBinary(findExecutableInPath("inkscape", "/usr/bin/inkscape"))

	output, err := conv.Convert(*input)
	if err != nil {
		fmt.Println(err)
		// os.Exit(1)
	}

	return &output
}


func image_gr(filename string, width, height int, ch chan<- order_string, order int, wg *sync.WaitGroup, thumbnail_type string) {
	defer wg.Done()
	var start time.Time
	var chafa_start time.Time
	if debug_time {
		start = time.Now()
	}


	cache := filepath.Join(get_thumbnail_cache_dir(), lfChafaPreviewFormat, file_font_ratio, get_geometry(width, height), limitStringToBytes(get_hash(), get_cache_byte_limit()))

	if !fileExists(filepath.Dir(cache)) {
		err := os.MkdirAll(filepath.Dir(cache), 0700)
		if err != nil {
			fmt.Println("Error Mkdir file:", err)
			log.Fatal(err)
		}
	}



	os.Chmod(filepath.Dir(cache), 0700)
	// gr_array[ar_index] = fmt.Sprintln(image(filename, width, height))
	// ch <- fmt.Sprint(image(filename, width, height))
	var output = order_string{order, ""}

	var image *[]byte

	// var err error

	var chafa_output string

	if fileExists(cache) {
		cache_data, err := os.ReadFile(cache)
		if err != nil {
			fmt.Println("Error reading file:", err)
			log.Fatal(err)
		}

		chafa_output = string(cache_data)
	} else {


		if thumbnail_type == "audio" {
			image = thumbnail_music(filename)

		} else if  thumbnail_type == "video" {
			vid_thumnr := vid_thm_new()

			vid_thumnr.SetBinary(findExecutableInPath("ffmpegthumbnailer", "ffmpegthumbnailer"))

			var err error
			image, err = vid_thumnr.vid_thm_Convert(filename)
			if err != nil {
				fmt.Println(err)
			}
			
		} else if thumbnail_type == "" {
			image_data, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
				log.Fatal(err)
			}

			if isSVG(filename) {
				if isSVGz(filename) {
					image_data = *(svgz_to_svg(&image_data))
				}

				image = svg_to_png(&image_data)
			} else {
				image = &image_data
			}
			
		}




		if debug_time {
			chafa_start = time.Now()
		}

		chafa_output = chafa_image(image, width, height)

		if debug_time {
			time_output = time_output + fmt.Sprintln("chafa time: ",time.Since(chafa_start))
		}

		err := os.WriteFile(cache, []byte(chafa_output), 0600)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			log.Fatal(err)
		}

	}








	// output.content = output.content + "test"



	// output.content = output.content + strings.TrimSuffix(chafa_output, "[?25h")
	output.content = output.content + chafa_output


	ch <- output
	if debug_time {
		time_output = time_output + fmt.Sprintln("image_gr time: ",time.Since(start))
	}
}