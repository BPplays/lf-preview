package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/barasher/go-exiftool"
	"github.com/dhowden/tag"
	"github.com/kalafut/imohash"
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

var hash string = ""



// func calculateHash(filePath string) string {
// 	var hashstart time.Time
// 	if chafaPreviewDebugTime == "1" {
// 		hashstart = time.Now()
// 	}

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		os.Exit(1)
// 	}
// 	defer file.Close()

// 	hash := blake3.New(256, nil)
// 	if _, err := io.Copy(hash, file); err != nil {
// 		fmt.Println("Error calculating hash:", err)
// 		os.Exit(1)
// 	}


// 	output := limitStringToBytes(fmt.Sprintf("%x", hash.Sum(nil)), cache_byte_limit)

// 	if chafaPreviewDebugTime == "1" {
// 		time_output = time_output + fmt.Sprintln("hash time: ",time.Since(hashstart))
// 		// time_output = time_output + fmt.Sprintln("hash: ", fmt.Sprintf("%x", hash.Sum(nil)))
// 	}

// 	return output
// }


func get_hash() string {
	if hash == "" {
		hash = calculateHash(file)
	}

	return hash

}




// func calculateHash(filePath string) string {
// 	var hashstart time.Time
// 	if chafaPreviewDebugTime == "1" {
// 		hashstart = time.Now()
// 	}

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		os.Exit(1)
// 	}
// 	defer file.Close()

// 	key := make([]byte, 32)

// 	hash, err := highwayhash.New(key)

// 	if err != nil {
// 		fmt.Println(err)
// 		log.Fatal(err)
// 	}

// 	if _, err := io.Copy(hash, file); err != nil {
// 		fmt.Println("Error calculating hash:", err)
// 		os.Exit(1)
// 	}

// 	if chafaPreviewDebugTime == "1" {
// 		time_output = time_output + fmt.Sprintln("hash time: ",time.Since(hashstart))
// 	}

// 	output := limitStringToBytes(fmt.Sprintf("%x", hash.Sum(nil)), cache_byte_limit)

// 	return output
// }





// func calculateHash(filePath string) string {
// 	var hashstart time.Time
// 	if chafaPreviewDebugTime == "1" {
// 		hashstart = time.Now()
// 	}

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		os.Exit(1)
// 	}
// 	defer file.Close()

// 	// key := make([]byte, 32)

// 	hash := fnv.New128()

// 	hash2 := fnv.New128()

// 	file_data, err := os.ReadFile(filePath)
// 	if err != nil {
// 		fmt.Println("Error reading file:", err)
// 		log.Fatal(err)
// 	}

// 	hash.Write(file_data)
// 	hash2.Write(append(file_data, []byte("a")...))

// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	log.Fatal(err)
// 	// }

// 	if _, err := io.Copy(hash, file); err != nil {
// 		fmt.Println("Error calculating hash:", err)
// 		os.Exit(1)
// 	}

// 	if chafaPreviewDebugTime == "1" {
// 		time_output = time_output + fmt.Sprintln("hash time: ",time.Since(hashstart))
// 	}

// 	output := limitStringToBytes(fmt.Sprintf("%x", append(hash.Sum(nil), hash2.Sum(nil)...)), cache_byte_limit)

// 	return output
// }



func calculateHash(filePath string) string {
	var hashstart time.Time
	if chafaPreviewDebugTime == "1" {
		hashstart = time.Now()
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()


	hash := imohash.New()


	sum, err := hash.SumFile(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("hash time: ",time.Since(hashstart))
	}

	output := limitStringToBytes(fmt.Sprintf("%x", sum), cache_byte_limit)

	return output
}












func add_ext(file string, ext string, limit int) string {
	ext_bytes := []byte(ext)
	ext_len := len(ext_bytes)

	file_limit := limitStringToBytes(file, limit - ext_len)

	return file_limit+ext
}




func commandExists(command string) bool {
	cmd := exec.Command("which", command)
	err := cmd.Run()
	return err == nil
}








func get_folder_max_len(folder string) int {
	var i int
	if commandExists("getconf") {
		cmd := exec.Command("getconf", "NAME_MAX", folder)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output), err)
			log.Fatal(string(output), err)
		}
		cleanedString := strings.ReplaceAll(strings.ReplaceAll(string(output), " ", ""), "\n", "")
		i, err = strconv.Atoi(cleanedString)
		if err != nil {
			panic(err)
		}
	} else {
		i = 128
	}

	return i
}



func limitStringToBytes(input string, maxBytes int) string {
	// Ensure maxBytes is not negative
	if maxBytes <= 0 {
		return ""
	}

	// Convert string to a slice of bytes
	bytes := []byte(input)

	// Iterate through the string to get the substring within the byte limit
	for len(bytes) > maxBytes {
		_, size := utf8.DecodeLastRune(bytes)
		bytes = bytes[:len(bytes)-size]
	}

	// Convert the byte slice back to a string
	return string(bytes)
}


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


func exif_fmt(fileInfos []exiftool.FileMetadata, tags [][]string) (string) {
	// var start time.Time
	// if chafaPreviewDebugTime == "1" {
	// 	start = time.Now()
	// }

	// fileInfos := get_exif(file)

	// if chafaPreviewDebugTime == "1" {
	// 	time_output = time_output + fmt.Sprintln("get_exif time: ",time.Since(start))
	// }
	output := ""
	// cur := ""
	// if chafaPreviewDebugTime == "1" {
	// 	start = time.Now()
	// }
	for i, tag_small := range tags {
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
		if i != len(tags) - 1 {
			output = output + "\n"
		}

	}

	// if chafaPreviewDebugTime == "1" {
	// 	time_output = time_output + fmt.Sprintln("exif_fmt_loop time: ",time.Since(start))
	// }

	return output
}




func get_metadata(file string, tags [][]string) (string) {
	var output string

	cache := filepath.Join(metadata_cache_dir, add_ext(get_hash(), ".json", cache_byte_limit))

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


		output = exif_fmt(exif, tags)



	}


	return output
}






func exif_fmt_gr(file string, tags [][]string, ch chan<- order_string, order int, wg *sync.WaitGroup) {
	defer wg.Done()
	var start time.Time
	if chafaPreviewDebugTime == "1" {
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
	"PictureBitsPerPixel":                "Picture Bits Per Pixel",
	"FileModifyDate":                "File Modify Date",
	"FileAccessDate":                "File Access Date",
	"PictureDescription":                "Picture Description",
	"Directory":                "Directory",
	"TrackNumber":                "TrackNumber",
	"Duration":                "Duration",
	"Date":                "Date",
	"FileTypeExtension":                "File Type Extension",
	"FileSize":                "File Size",
	"SampleRate":                "Sample Rate",
	"FileName":                "File Name",
	"FileType":                "File Type",
	"Album":                "Album",
	"Artist":                "Artist",
	"Comment":                "Comment",
	"ImageSize":                "Image Size",
	"YCbCrSubSampling": "Y Cb Cr Sub Sampling",
	"BitsPerSample": "Bits Per Sample",
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



func chafa_image(image []byte, width, height int) (string) {
	geometry := fmt.Sprintf("%dx%d", width, height)

	cmd := exec.Command("chafa", fmt.Sprintf("--font-ratio=%s", userOpenFontRatio))
	cmd.Args = append(cmd.Args, chafaFmt...)
	cmd.Args = append(cmd.Args, chafaDither...)
	cmd.Args = append(cmd.Args, chafaColors...)
	cmd.Args = append(cmd.Args, "--color-space=din99d", "--scale=max", "-w", "9", "-O", "9", "-s", geometry, "--animate", "false", "--symbols", "all")

	pipe, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating pipe:", err)
		os.Exit(1)
	}

	go func() {
		defer pipe.Close() // Close the pipe when done

		// Write data to the command's standard input
		_, err := pipe.Write(image)
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


func image_gr(filename string, width, height int, ch chan<- order_string, order int, wg *sync.WaitGroup, thumbnail_type string) {
	defer wg.Done()
	var start time.Time
	var chafa_start time.Time
	if chafaPreviewDebugTime == "1" {
		start = time.Now()
	}
	// gr_array[ar_index] = fmt.Sprintln(image(filename, width, height))
	// ch <- fmt.Sprint(image(filename, width, height))
	var output = order_string{order, ""}

	var image []byte

	var err error

	if thumbnail_type == "audio" {
		image = thumbnail_music(filename)
	} else if thumbnail_type == "" {
		image, err = os.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading file:", err)
			log.Fatal(err)
		}
	}

	// output.content = output.content + "test"

	if chafaPreviewDebugTime == "1" {
		chafa_start = time.Now()
	}

	output.content = output.content + chafa_image(image, width, height)

	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("chafa time: ",time.Since(chafa_start))
	}
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





func get_thumbnail_cache_file(ext string) string {
	return filepath.Join(thumbnail_cache_dir, add_ext(get_hash(), ext, cache_byte_limit))
}






// func thumbnail_music(file string) string {
// 	cache := get_thumbnail_cache_file(".bmp")

// 	var start time.Time
// 	if chafaPreviewDebugTime == "1" {
// 		start = time.Now()
// 	}
// 	// cache := filepath.Join(cacheFile, ".bmp")
// 	// cache := thumbnail_cache + ".bmp"
// 	if !fileExists(cache) {
// 		// ffmpeg -i "$1" -an -c:v copy "${CACHE}.bmp"
// 		cmd := exec.Command("ffmpeg", "-y", "-hide_banner", "-loglevel", "error", "-nostats", "-i", file, "-an", "-c:v", "copy", cache)

// 		output, err := cmd.CombinedOutput()
// 		if err != nil {
// 			fmt.Println(string(output), err)
// 			log.Fatal(err)
// 		}
// 	}



// 	// fmt.Println(string(output))
// 	if chafaPreviewDebugTime == "1" {
// 		time_output = time_output + fmt.Sprintln("thumbnail_music time: ",time.Since(start))
// 	}
// 	return cache
// }




func thumbnail_music(file string) []byte {
	// cache := get_thumbnail_cache_file(".bmp")

	var start time.Time
	if chafaPreviewDebugTime == "1" {
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
	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("thumbnail_music time: ",time.Since(start))
	}
	return m.Picture().Data
}








func read_file(file string) string {
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}


func file_size_mb(file_path string) float64 {
	// Open the file
	file, err := os.Open(file_path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		log.Fatal(err)
	}
	defer file.Close()

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file information:", err)
		log.Fatal(err)
	}

	// Calculate file size in megabytes
	fileSizeInMB := float64(fileInfo.Size()) / (1 << 20) // 1 MB = 1 << 20 bytes

	return fileSizeInMB
}

var file_mb float64 = -1

func get_file_mb() float64 {
	if file_mb == -1 {
		file_mb = file_size_mb(file)
	}

	return file_mb
}




var userOpenFontRatio string
var chafaFmt []string
var chafaDither []string
var chafaColors []string



// var start time.Time

// var thumbnail_cache string
var metadata_cache_dir string

var thumbnail_cache_dir string
var chafaPreviewDebugTime string



var cache_byte_limit int
var file string

var ext string

var width int
var hight int


var configDir string
var cacheBase string
var lfCacheDir string

var lfChafaPreviewFormat string
var lfChafaPreviewFormatOverrideSixelRatio string
var lfChafaPreviewFormatOverrideKittyRatio string
var fontRatio string
var chafaPreviewDither string
var chafaPreviewColors string





func main() {

	chafaPreviewDebugTime = os.Getenv("LF_CHAFA_PREVIEW_DEBUG_TIME")
	var prgstart time.Time
	if chafaPreviewDebugTime == "1" {
		prgstart = time.Now()
	}


	Init()


	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("init time: ",time.Since(prgstart))
	}




	


	// thumbnail_cache = filepath.Join(thumbnail_cache_dir, fmt.Sprintf("thumbnail.%s", hash))
	// thumbnail_cache = filepath.Join(thumbnail_cache_dir, hash)


	// thumbnail_cache = filepath.Join(lfCacheDir, fmt.Sprintf("thumbnail.%s", hash))

	tmp := configDir
	tmp = ""
	fmt.Print(tmp)


	

    switch ext {
    case ".bmp", ".jpg", ".jpeg", ".png", ".xpm", ".webp", ".tiff", ".gif", ".jfif", ".ico":

		
		if get_file_mb() > 100 {
			fmt.Print("file to big to preview")
		} else {
			fmt.Print(image_exif(file, width, hight, file, image_tags, ""))
		}
	// case ".mp3", ".flac", ".ogg":
	case ".wav", ".mp3", ".flac", ".m4a", ".wma", ".ape", ".ac3", ".ogg", ".spx", ".opus", ".mka":
		// fmt.Println(exif_fmt(file, music_tags))
		
		fmt.Print(image_exif(file, width, hight, file, music_tags, "audio"))
		
    default:
        // fmt.Println("sdf")
		if get_file_mb() > 0.1 {
			fmt.Print("file to big to preview")
		} else {
			fmt.Print(read_file(file))
		}
		
    }





	if chafaPreviewDebugTime == "1" {
		time_output = time_output + fmt.Sprintln("total time: ",time.Since(prgstart))
		fmt.Println(sep1)
		fmt.Println(time_output)
	}

	


	
}
