package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/barasher/go-exiftool"
	"github.com/dhowden/tag"
	"github.com/kalafut/imohash"
	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
	"github.com/zeebo/blake3"
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

func calculatesha256Hash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		fmt.Println("Error calculating hash:", err)
		os.Exit(1)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

var hash_res string = ""


func blake3Hash(data []byte) string {
	var hashstartblake3 time.Time
	if debug_time {
		hashstartblake3 = time.Now()
	}

	hasher := blake3.New()

	hasher.Write(data)


	output := limitStringToBytes(fmt.Sprintf("%x", hasher.Sum(nil)), get_cache_byte_limit())

	if debug_time {
		time_output = time_output + fmt.Sprintln("blake3 hash time: ",time.Since(hashstartblake3))
		// time_output = time_output + fmt.Sprintln("hash: ", fmt.Sprintf("%x", hash.Sum(nil)))
	}

	return output
}



var hash_started bool = false

func get_hash() string {
	if hash_res == "" {
		if !hash_started {
			hash_started = true


			var hashstart time.Time


			if debug_time {
				hashstart = time.Now()
			}

			hash_res = limitStringToBytes(calculateHash(file)+blake3Hash([]byte(filepath.Base(file))), get_cache_byte_limit())

			if debug_time {
				time_output = time_output + fmt.Sprintln("hash time: ",time.Since(hashstart))
			}
		} else {
			for hash_res == "" {
				time.Sleep(80 * time.Microsecond)
			}
		}

	}

	return hash_res

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


	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// hasher := sparsehash.New(sha256.New)
	// hasher := sparsehash.New(highwayhash_hh)
	hasher := imohash.New()
	// hasher := imohash.NewCustom(10000, 64)


	sum, err := hasher.SumFile(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}



	output := limitStringToBytes(fmt.Sprintf("%x", sum), get_cache_byte_limit())

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


func clampToMax(value, max int) int {
    if value > max {
        return max
    }
    return value
}


func blocks_fmt(blocks []string) (string) {

	var builder strings.Builder
	len_blocks := len(blocks)

	for i, block := range blocks {
		if block != "" {
			builder.WriteString(block)

			if blocks[clampToMax(i+1, len_blocks-1)] != "" {
				if i < len_blocks-1 {
					builder.WriteString("\n")
				}
			}
		}

	}


	return builder.String()
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


		output = exif_fmt(exif, tags)



	}


	return output
}






func exif_fmt_gr(file string, tags [][]string, ch chan<- order_string, order int, wg *sync.WaitGroup) {
	defer wg.Done()
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



var sep1 = ""


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



func chafa_image(image *[]byte, width, height int) (string) {


	cmd := exec.Command("chafa", fmt.Sprintf("--font-ratio=%s", userOpenFontRatio))
	cmd.Args = append(cmd.Args, chafaFmt...)
	cmd.Args = append(cmd.Args, chafaDither...)
	cmd.Args = append(cmd.Args, chafaColors...)
	cmd.Args = append(cmd.Args, "--color-space=din99d", "--scale=max", "-w", "9", "-O", "9", "-s", get_geometry(width, height), "--animate", "false")
	cmd.Args = append(cmd.Args, "--symbols", "block+border+space-wide+inverted+quad+extra+half+hhalf+vhalf")


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



















func image_gr(filename string, width, height int, ch chan<- order_string, order int, wg *sync.WaitGroup, thumbnail_type string) {
	defer wg.Done()
	var start time.Time
	var chafa_start time.Time
	if debug_time {
		start = time.Now()
	}


	cache := filepath.Join(get_thumbnail_cache_dir(), file_font_ratio, get_geometry(width, height), limitStringToBytes(get_hash(), get_cache_byte_limit()))

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
		} else if thumbnail_type == "" {
			image_data, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
				log.Fatal(err)
			}
			image = &image_data
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



	output.content = output.content + chafa_output


	ch <- output
	if debug_time {
		time_output = time_output + fmt.Sprintln("image_gr time: ",time.Since(start))
	}
}

var time_output string
type order_string struct {
	order int
	content string
}
// var gr_array [2]string



func countRune(s string, r rune) int {
    count := 0
    for _, c := range s {
        if c == r {
            count++
        }
    }
    return count
}








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

	lines := 0
	for result := range ch {
		lines += countRune(result.content, '\n')
		temp_slice[result.order] = result.content
	}



	if lines > height {




		new_height := height-countRune(temp_slice[1], '\n')

		if new_height < 2 {
			temp_slice[0] = ""
		} else {
			ch2 := make(chan order_string)

			wg.Add(1)
			go image_gr(image_file, width, new_height, ch2, 0, &wg, thumbnail_type)
			go func() {
				wg.Wait()
				close(ch2)
			}()

			for result := range ch2 {
				temp_slice[result.order] = result.content
			}
		}
		// fmt.Println(lines)
		// fmt.Printf("h: %v, nh: %v\n", height, new_height)


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
	return filepath.Join(get_thumbnail_cache_dir(), add_ext(get_hash(), ext, get_cache_byte_limit()))
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





func char_wrap(s string, limit int) string {
	var result strings.Builder

	var rune_sl []rune
	var diff float64

	var aj_limit float64

	var fl_len float64

	var int_aj_limit int



	// if len(rune_sl) < limit {
	// 	return s
	// }

	string_split := strings.Split(s, "\n")


	// var inc int64 = 0

	for _, str := range string_split {
		rune_sl = []rune(str)

		for {
				// inc += 1
				// fmt.Println(inc)

				fl_len = float64(len(rune_sl))

				diff = float64(runewidth.StringWidth(string(rune_sl))) / fl_len

				if diff > 0 {
					aj_limit = float64(limit) / diff
					int_aj_limit = int(math.Floor(aj_limit))
				}


				if len(rune_sl) <= int_aj_limit {
					result.WriteString(string(rune_sl))
					result.WriteString("\n")
					break
				}

				// diff = runewidth.StringWidth(str) - len(rune_sl)
				// fmt.Printf("len: %v\n", len(rune_sl))
				// fmt.Printf("fl_len: %v\n", fl_len)
				// fmt.Printf("aj_limit: %v\n", aj_limit)
				// fmt.Printf("int_aj_limit: %v\n", int_aj_limit)
				// fmt.Printf("diff: %v\n", diff)



				result.WriteString(string(rune_sl[:int_aj_limit-1]))
				result.WriteString("⏎\n")

				rune_sl = rune_sl[int_aj_limit-1:]
		}



	}

	return result.String()
	}













var userOpenFontRatio string
var chafaFmt []string
var chafaDither []string
var chafaColors []string



var file_font_ratio string

// var start time.Time

// var thumbnail_cache string
var metadata_cache_dir string = ""

var thumbnail_cache_dir string = ""
var debug_time bool



var cache_byte_limit int = -1
var file string

var ext string

var width int
var hight int


var configDir string = ""
var cacheBase string
var lfCacheDir string = ""

var lfChafaPreviewFormat string = ""
var lfChafaPreviewFormatOverrideSixelRatio string = ""
var lfChafaPreviewFormatOverrideKittyRatio string = ""
var fontRatio string = ""
var chafaPreviewDither string = ""
var chafaPreviewColors string = ""


func stringNumberToBool(strNumber string) bool {
	if strNumber == "" {
		return false
	}
	// Parse the string to an integer
	intValue, err := strconv.Atoi(strNumber)
	if err != nil {
		// Handle the error (e.g., invalid string format)
		fmt.Println("Error:", err)
		return false
	}

	// Convert the integer to a boolean
	return intValue != 0
}


func main() {

	cpuprofile := os.Getenv("LF_CHAFA_PREVIEW_DEBUG_CPUPROF")
	memprofile := os.Getenv("LF_CHAFA_PREVIEW_DEBUG_MEMPROF")

	flag.Parse()
    if cpuprofile != "" {
        f, err := os.Create(cpuprofile)
        if err != nil {
            log.Fatal("could not create CPU profile: ", err)
        }
        defer f.Close() // error handling omitted for example
        if err := pprof.StartCPUProfile(f); err != nil {
            log.Fatal("could not start CPU profile: ", err)
        }
        defer pprof.StopCPUProfile()
    }



	debug_time = stringNumberToBool(os.Getenv("LF_CHAFA_PREVIEW_DEBUG_TIME"))
	var prgstart time.Time
	if debug_time {
		prgstart = time.Now()
	}

	disable_wordwrap := os.Getenv("LF_CHAFA_PREVIEW_DISABLE_WORDWRAP")

	Init()


	// if chafaPreviewDebugTime == "1" {
	// 	time_output = time_output + fmt.Sprintln("init time: ",time.Since(prgstart))
	// }



	preview_output := ""



	// thumbnail_cache = filepath.Join(thumbnail_cache_dir, fmt.Sprintf("thumbnail.%s", hash))
	// thumbnail_cache = filepath.Join(thumbnail_cache_dir, hash)


	// thumbnail_cache = filepath.Join(lfCacheDir, fmt.Sprintf("thumbnail.%s", hash))

	tmp := configDir
	tmp = ""
	fmt.Print(tmp)




    switch ext {
    case ".bmp", ".jpg", ".jpeg", ".png", ".xpm", ".webp", ".tiff", ".gif", ".jfif", ".ico":


		if get_file_mb() > 100 {
			preview_output = "file to big to preview"
		} else {
			preview_output = image_exif(file, width, hight, file, image_tags, "")
		}
	// case ".mp3", ".flac", ".ogg":
	case ".wav", ".mp3", ".flac", ".m4a", ".wma", ".ape", ".ac3", ".ogg", ".spx", ".opus", ".mka":
		// fmt.Println(exif_fmt(file, music_tags))
		// get_hash()

		preview_output = image_exif(file, width, hight, file, music_tags, "audio")

    default:
        // fmt.Println("sdf")
		if get_file_mb() > 0.1 {
			preview_output = "file to big to preview"
		} else {
			preview_output = read_file(file)

			if disable_wordwrap != "1" {
				preview_output = wordwrap.WrapString(preview_output, uint(width))
				preview_output = char_wrap(preview_output, width)
			}
		}

    }





	if debug_time {
		time_output = time_output + fmt.Sprintln("total time: ",time.Since(prgstart))
		preview_output = preview_output + sep1 + "\n"
		preview_output = preview_output + time_output
	}



	if get_print_output() {
		fmt.Print(preview_output)
	}


    if memprofile != "" {
        f, err := os.Create(memprofile)
        if err != nil {
            log.Fatal("could not create memory profile: ", err)
        }
        defer f.Close() // error handling omitted for example
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal("could not write memory profile: ", err)
        }
    }



}
