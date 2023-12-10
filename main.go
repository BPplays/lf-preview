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


func exif_fmt(file string, tags [][]string) (string) {
	fileInfos := get_exif(file)
	output := ""
	// cur := ""
	for _, tag_small := range tags {
		for _, tag := range tag_small {
			for _, fileInfo := range fileInfos {
				// output = output + fmt.Sprintln("fileInfos")
				val, ok := fileInfo.Fields[tag]
				// If the key exists
				if ok {
					tag_name := tag
					tag_val, ok := exif_key_map[tag]
					// If the key exists
					if ok {
						tag_name = tag_val
					}
					output = output + fmt.Sprintf("%v: %v\n", tag_name, val)
				}
					
			}

		}
		output = output + "\n"
	}

	return output
}


func exif_fmt_gr(file string, tags [][]string, ch chan<- order_string, order int, wg *sync.WaitGroup) {
	defer wg.Done()
	var output = order_string{order, ""}
	// output[1] = output[1] + fmt.Sprint("test")

	// ch <- fmt.Sprint("test")
	// ch <- fmt.Sprint(exif_fmt(file, tags))
	output.content = output.content + "test"
	output.content = output.content + exif_fmt(file, tags)
	ch <- output
	// output := exif_fmt(file, tags)
	// gr_array[ar_index] = "test"
	// gr_array[1] = "test"
	// fmt.Println((*array)[ar_index])
	// fmt.Println(ar_index)
	// fmt.Println("testrgji")
}



var sep1 = "============================================================================="


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

var music_tags = [][]string{
	{"Title", "Duration"},
	{"Genre", "Album", "Artist", "Composer", "Date"},
	{"SampleRate", "Channels", "FileType"},
}




var image_tags = [][]string{
	{"ImageSize", "Megapixels", "FileSize"},
	{"FileType", "ColorSpace", "Compression"},
	{"BitsPerSample", "YCbCrSubSampling"},
}



var exif_key_map = map[string]string{
	"Title":                "Title",
	"Genre":                "Genre",
	"Composer":                "Composer",
	"PictureBitsPerPixel":                "PictureBitsPerPixel",
	"FileModifyDate":                "FileModifyDate",
	"FileAccessDate":                "FileAccessDate",
	"PictureDescription":                "PictureDescription",
	"Directory":                "Directory",
	"TrackNumber":                "TrackNumber",
	"Duration":                "Duration",
	"Date":                "Date",
	"FileTypeExtension":                "FileTypeExtension",
	"FileSize":                "FileSize",
	"SampleRate":                "SampleRate",
	"FileName":                "FileName",
	"FileType":                "FileType",
	"Album":                "Album",
	"Artist":                "Artist",
	"Comment":                "Comment",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
	// "Genre":                "Genre",
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


func image_gr(filename string, width, height int, ch chan<- order_string, order int, wg *sync.WaitGroup) {
	defer wg.Done()
	// gr_array[ar_index] = fmt.Sprintln(image(filename, width, height))
	// ch <- fmt.Sprint(image(filename, width, height))
	var output = order_string{order, ""}

	// output.content = output.content + "test"
	output.content = output.content + image(filename, width, height)
	ch <- output
}


type order_string struct {
	order int
	content string
}
// var gr_array [2]string

func image_exif(image_file string, width, height int, file string, tags [][]string) (string) {
	output := ""

	ch := make(chan order_string, 2)
	// ch2 := make(chan string)

	// var gr_array [2]string


	var wg sync.WaitGroup

	wg.Add(2)
	go image_gr(image_file, width, height, ch, 0, &wg)
	go exif_fmt_gr(file, tags, ch, 1, &wg)

	go func() {
		wg.Wait()
		close(ch)
		// close(ch2)
	}()

	// gr_array[0] = "test0"
	// gr_array[1] = "test1"
	// output = output + fmt.Sprintln(gr_array[0])
	var temp_slice []string

	for result := range ch {
		temp_slice[result.order] = result.content
	}

	for _, val := range temp_slice {
		output = output + val
	}

	output = output + fmt.Sprintln(sep1)
	output = output + fmt.Sprintln(sep1)
	// output = output + fmt.Sprintln(gr_array[1])

	// close(ch)
	// close(ch2)

	return output
}








var userOpenFontRatio string
var chafaFmt []string
var chafaDither []string
var chafaColors []string





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
		// fmt.Println(width, hight)
		fmt.Println(image_exif(file, width, hight, file, image_tags))
		//! fmt.Println(exif_fmt(file))
    // case "Wednesday", "Thursday":
    //     fmt.Println("It's the middle of the week.")
    // case "Friday", "Saturday", "Sunday":
    //     fmt.Println("It's the end of the week.")
	case ".mp3", ".flac":
		fmt.Println(exif_fmt(file, music_tags))
		
    default:
        fmt.Println("sdf")
    }










	
}
