package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/barasher/go-exiftool"
	"lukechampine.com/blake3"
)

// type thumbnail func(string, int) int



func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user's home directory:", err)
		os.Exit(1)
	}
	return usr.HomeDir
}

func getEnvOrFallback(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// func calculateHash(filePath string) string {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		os.Exit(1)
// 	}
// 	defer file.Close()

// 	hash := sha256.New()
// 	if _, err := io.Copy(hash, file); err != nil {
// 		fmt.Println("Error calculating hash:", err)
// 		os.Exit(1)
// 	}

// 	return fmt.Sprintf("%x", hash.Sum(nil))
// }


func calculateHash(filePath string) string {
	if chafaPreviewDebugTime == "1" {
		start = time.Now()
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	hash := blake3.New(256, nil)
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println("Error calculating hash:", err)
		os.Exit(1)
	}

	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("hash time: ",time.Since(start))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}





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
	if chafaPreviewDebugTime == "1" {
		start = time.Now()
	}
	var output = order_string{order, ""}
	// output[1] = output[1] + fmt.Sprint("test")

	// ch <- fmt.Sprint("test")
	// ch <- fmt.Sprint(exif_fmt(file, tags))
	output.content = output.content + fmt.Sprintln(sep1)
	// output.content = output.content + "test"
	output.content = output.content + exif_fmt(file, tags)
	ch <- output
	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("exif_fmt_gr time: ",time.Since(start))
	}
	// output := exif_fmt(file, tags)
	// gr_array[ar_index] = "test"
	// gr_array[1] = "test"
	// fmt.Println((*array)[ar_index])
	// fmt.Println(ar_index)
	// fmt.Println("testrgji")
}



var sep1 = "=================================================================="


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


func image_gr(filename string, width, height int, ch chan<- order_string, order int, wg *sync.WaitGroup, thumbnail_type string) {
	defer wg.Done()
	if chafaPreviewDebugTime == "1" {
		start = time.Now()
	}
	// gr_array[ar_index] = fmt.Sprintln(image(filename, width, height))
	// ch <- fmt.Sprint(image(filename, width, height))
	var output = order_string{order, ""}

	if thumbnail_type == "audio" {
		filename = thumbnail_music(filename)
	}

	// output.content = output.content + "test"
	output.content = output.content + image(filename, width, height)
	ch <- output
	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("image_gr time: ",time.Since(start))
	}
}

var time_output string
type order_string struct {
	order int
	content string
}
// var gr_array [2]string

func image_exif(image_file string, width, height int, file string, tags [][]string, thumbnail_type string) (string) {
	output := ""

	ch := make(chan order_string)
	// ch2 := make(chan string)

	// var gr_array [2]string


	var wg sync.WaitGroup

	wg.Add(2)
	go image_gr(image_file, width, height, ch, 0, &wg, thumbnail_type)
	go exif_fmt_gr(file, tags, ch, 1, &wg)


	go func() {
		wg.Wait()
		close(ch)
		// close(ch2)
	}()

	// gr_array[0] = "test0"
	// gr_array[1] = "test1"
	// output = output + fmt.Sprintln(gr_array[0])
	var temp_slice [20]string

	for result := range ch {
		temp_slice[result.order] = result.content
	}



	for _, val := range temp_slice {
		output = output + val
	}

	// output = output + fmt.Sprintln(sep1)
	// output = output + fmt.Sprintln(sep1)
	// output = output + fmt.Sprintln(gr_array[1])

	// close(ch)
	// close(ch2)

	return output
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}


func thumbnail_music(file string) string {
	if chafaPreviewDebugTime == "1" {
		start = time.Now()
	}
	// cache := filepath.Join(cacheFile, ".bmp")
	cache := thumbnail_cache + ".bmp"
	if !fileExists(cache) {
		// ffmpeg -i "$1" -an -c:v copy "${CACHE}.bmp"
		cmd := exec.Command("ffmpeg", "-y", "-hide_banner", "-loglevel", "error", "-nostats", "-i", file, "-an", "-c:v", "copy", cache)

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output), err)
			log.Fatal(err)
		}
	}



	// fmt.Println(string(output))
	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("thumbnail_music time: ",time.Since(start))
	}
	return cache
}



var userOpenFontRatio string
var chafaFmt []string
var chafaDither []string
var chafaColors []string



var start time.Time

var thumbnail_cache string
var chafaPreviewDebugTime string







func main() {
	chafaPreviewDebugTime = os.Getenv("LF_CHAFA_PREVIEW_DEBUG_TIME")

	if chafaPreviewDebugTime == "1" {
		start = time.Now()
	}

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


	for len(sep1) < width {
		sep1 = sep1 + sep1
	}


	defaultConfigBase := filepath.Join(getHomeDir(), ".config")
	configDir := getEnvOrFallback("XDG_CONFIG_HOME", defaultConfigBase)

	defaultCacheBase := filepath.Join(getHomeDir(), ".cache")
	cacheBase := getEnvOrFallback("XDG_CACHE_HOME", defaultCacheBase)

	lfCacheDir := filepath.Join(cacheBase, "lf")
	if _, err := os.Stat(lfCacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(lfCacheDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	// thumbnail_cache_dir := filepath.Join(lfCacheDir, "thumbnails")
	// if _, err := os.Stat(thumbnail_cache_dir); os.IsNotExist(err) {
	// 	err := os.MkdirAll(thumbnail_cache_dir, os.ModePerm)
	// 	if err != nil {
	// 		fmt.Println("Error creating directory:", err)
	// 		return
	// 	}
	// }



	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("init time: ",time.Since(start))
	}



	hash := calculateHash(file)


	// thumbnail_cache = filepath.Join(thumbnail_cache_dir, fmt.Sprintf("thumbnail.%s", hash))

	thumbnail_cache = filepath.Join(lfCacheDir, fmt.Sprintf("thumbnail.%s", hash))
	tmp := thumbnail_cache + configDir
	tmp = ""
	fmt.Print(tmp)


	

    switch ext {
    case ".bmp", ".jpg", ".jpeg", ".png", ".xpm", ".webp", ".tiff", ".gif", ".jfif", ".ico":
        // fmt.Println("It's an image file.")
		// fmt.Println(image(file, width, hight))
		// fmt.Println(width, hight)
		// fmt.Println(exif_fmt(file, image_tags))
		fmt.Print(image_exif(file, width, hight, file, image_tags, ""))
		// fmt.Println(exif_fmt(file))
    // case "Wednesday", "Thursday":
    //     fmt.Println("It's the middle of the week.")
    // case "Friday", "Saturday", "Sunday":
    //     fmt.Println("It's the end of the week.")
	case ".mp3", ".flac":
		// fmt.Println(exif_fmt(file, music_tags))
		fmt.Print(image_exif(file, width, hight, file, music_tags, "audio"))
		
    default:
        fmt.Println("sdf")
    }





	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("total time: ",time.Since(start))
		fmt.Println(sep1)
		fmt.Println(time_output)
	}

	


	
}
