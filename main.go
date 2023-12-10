package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"

	"github.com/barasher/go-exiftool"
)

func get_exif(file string) ([]exiftool.FileMetadata) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		// fmt.Printf("Error when intializing: %v\n", err)
		// return "", err
		log.Fatal(err)
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


func exif_fmt(file string) (string) {
	fileInfos := get_exif(file)
	output := ""
	// cur := ""
	for _, fileInfo := range fileInfos {
		output = output + fmt.Sprintln("fileInfos")
		for k, v := range fileInfo.Fields {
			output = output + fmt.Sprintf("[%v] %v\n", k, v)
		}
	}

	return output
}

// music_tags=(
//     "-Title -Duration"
//     "-Genre -Album -Artist -Composer -Date"
//     "-SampleRate -Channels -FileType"
// )

// video_tags=(
//     "-Duration"
//     "-ImageSize -FileSize"
//     "-VideoCodecID -FileType"
// 	"-Megapixels"
// )

// image_tags=(
//     "-ImageSize -Megapixels -FileSize"
//     "-FileType -ColorSpace -Compression"
// 	# " "
// 	"-BitsPerSample -YCbCrSubSampling"
// )




var exif_key_map = map[string]string{
	"alma":                "\uF31D",
}



func image(filename string, width, height int) (string) {
	geometry := fmt.Sprintf("%dx%d", width, height)

	cmd := exec.Command("chafa", filename, fmt.Sprintf("--font-ratio=%s", userOpenFontRatio))
	cmd.Args = append(cmd.Args, chafaFmt...)
	cmd.Args = append(cmd.Args, chafaDither...)
	cmd.Args = append(cmd.Args, chafaColors...)
	cmd.Args = append(cmd.Args, "--color-space=din99d", "--scale=max", "-w", "9", "-O", "9", "-s", geometry, "--animate", "false")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(output))
	return string(output)
}


func image_gr(filename string, width, height int, ch chan<-string, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- fmt.Sprint(image(filename, width, height))
}





// func image_exif(filename string, width, height int) (string) {

// 	ch1 := make(chan string)
// 	ch2 := make(chan string)

// 	// defer ch1.Close
// 	// defer ch2.Close

// 	var wg sync.WaitGroup

// 	go image_gr(filename, width, height, ch1, &wg)
// }








var userOpenFontRatio string
var chafaFmt []string
var chafaDither []string
var chafaColors []string

var music_tags [][]string



func main() {
	arg2, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error parsing argument:", err)
		return
	}

	arg3, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Error parsing argument:", err)
		return
	}

	file := os.Args[1]
	ext := path.Ext(file)
	// Subtract 2 from the parsed value
	width := arg2 - 2
	hight := arg3


	lfChafaPreviewFormat := os.Getenv("LF_CHAFA_PREVIEW_FORMAT")
	lfChafaPreviewFormatOverrideSixelRatio := os.Getenv("LF_CHAFA_PREVIEW_FORMAT_OVERRIDE_SIXEL_RATIO")
	lfChafaPreviewFormatOverrideKittyRatio := os.Getenv("LF_CHAFA_PREVIEW_FORMAT_OVERRIDE_KITTY_RATIO")
	fontRatio := os.Getenv("FONT_RATIO")
	chafaPreviewDither := os.Getenv("LF_CHAFA_PREVIEW_DITHER")
	chafaPreviewColors := os.Getenv("LF_CHAFA_PREVIEW_COLORS")

	defaultUserOpenFontRatio := "1/2"

	if lfChafaPreviewFormat == "sixel" {
		defaultUserOpenFontRatio = "1/1"

		if lfChafaPreviewFormatOverrideSixelRatio != "1" {
			defaultUserOpenFontRatio = "1/1"
		}
	}

	if lfChafaPreviewFormat == "kitty" {
		if lfChafaPreviewFormatOverrideKittyRatio != "1" {
			defaultUserOpenFontRatio = "100/225"
		}
	}

	userOpenFontRatio = fontRatio
	if userOpenFontRatio == "" {
		userOpenFontRatio = defaultUserOpenFontRatio
	}

	chafaFmt = []string{}
	if lfChafaPreviewFormat != "" {
		chafaFmt = append(chafaFmt, "-f", lfChafaPreviewFormat)
	}

	chafaDither = []string{}
	if chafaPreviewDither != "" {
		chafaDither = append(chafaDither, fmt.Sprintf("--dither=%s", chafaPreviewDither))
	}

	chafaColors = []string{"--colors=full"}
	if chafaPreviewColors != "" {
		chafaColors[0] = fmt.Sprintf("--colors=%s", chafaPreviewColors)
	}







	

    switch ext {
    case ".bmp", ".jpg", ".jpeg", ".png", ".xpm", ".webp", ".tiff", ".gif", ".jfif", ".ico":
        // fmt.Println("It's an image file.")
		// fmt.Println(image(file, width, hight))
		fmt.Println(width, hight)
		fmt.Println(exif_fmt(file))
    // case "Wednesday", "Thursday":
    //     fmt.Println("It's the middle of the week.")
    // case "Friday", "Saturday", "Sunday":
    //     fmt.Println("It's the end of the week.")
	case ".mp3", ".flac":
		fmt.Println(exif_fmt(file))
    default:
        fmt.Println("sdf")
    }










	
}
