package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
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

	"github.com/kalafut/imohash"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/pflag"
	"github.com/zeebo/blake3"
)

// type thumbnail func(string, int) int





var ignored_folders = map[string]bool{
	".ssh": true,
	"ssh": true,
}

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



var hash_mutex sync.Mutex

func get_hash() string {
	if hash_res == "" {
		hash_mutex.Lock()
		defer hash_mutex.Unlock()
	}

	if hash_res == "" {
		var hashstart time.Time


		if debug_time {
			hashstart = time.Now()
		}

		hash_res = limitStringToBytes(calculateHash(file)+blake3Hash([]byte(filepath.Base(file))), get_cache_byte_limit())

		if debug_time {
			time_output = time_output + fmt.Sprintln("hash time: ",time.Since(hashstart))
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








var sep1 = ""


// music_tags=(
//     "-Title -Duration"
//     "-Genre -Album -Artist -Composer -Date"
//     "-SampleRate -Channels -FileType"
// )

// video_tags=(
//     "-Duration -FileSize"
//     "-ImageSize -VideoFrameRate"
//     "-VideoCodecID -FileType"
// 	"-Megapixels"
// )

// image_tags=(
//     "-ImageSize -Megapixels -FileSize"
//     "-FileType -ColorSpace -Compression"
// 	# " "
// 	"-BitsPerSample -YCbCrSubSampling"
// )





func findExecutableInPath(executable string, default_path string) (string) {
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, path := range paths {
		fullPath := path + string(os.PathSeparator) + executable
		_, err := exec.LookPath(fullPath)
		if err == nil {
			return fullPath
		}
	}
	return default_path
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
	// fileSizeInMB := float64(fileInfo.Size()) / (1 << 20) // 1 MB = 1 << 20 bytes
	fileSizeInMB := float64(fileInfo.Size()) * math.Pow10(-6)

	return fileSizeInMB
}

var file_mb float64 = -1

func get_file_mb() float64 {
	if file_mb == -1 {
		file_mb = file_size_mb(file)
	}

	return file_mb
}



func unhideCursor() {
	cmd := exec.Command("tput", "cnorm")
	cmd.Stdout = os.Stdout
	cmd.Run()
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
var height int


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

var dbg_print_exif bool

var no_info bool






func getBaseFolder(filePath string) string {
	folder := filepath.Dir(filePath)
	baseFolder := filepath.Base(folder)
	return baseFolder
}




var disable_compat bool


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

	pflag.BoolVarP(&no_info, "no_info", "n", false, "no info section")

	pflag.Parse()






	debug_time = stringNumberToBool(os.Getenv("LF_CHAFA_PREVIEW_DEBUG_TIME"))
	disable_compat = stringNumberToBool(os.Getenv("LF_CHAFA_PREVIEW_DISABLE_COMPAT"))
	dbg_print_exif = stringNumberToBool(os.Getenv("LF_CHAFA_PREVIEW_DEBUG_EXIF_PRINT"))
	var prgstart time.Time
	if debug_time {
		prgstart = time.Now()
	}

	disable_wordwrap := os.Getenv("LF_CHAFA_PREVIEW_DISABLE_WORDWRAP")

	Init()


	// if chafaPreviewDebugTime == "1" {
	// 	time_output = time_output + fmt.Sprintln("init time: ",time.Since(prgstart))
	// }

	mime.AddExtensionType(".webp", "image/webp")
	mime.AddExtensionType(".avif", "image/avif")
	mime.AddExtensionType(".avifs", "image/avif")
	mime.AddExtensionType(".jxl", "image/jxl")


	preview_output := ""



	// thumbnail_cache = filepath.Join(thumbnail_cache_dir, fmt.Sprintf("thumbnail.%s", hash))
	// thumbnail_cache = filepath.Join(thumbnail_cache_dir, hash)


	// thumbnail_cache = filepath.Join(lfCacheDir, fmt.Sprintf("thumbnail.%s", hash))

	// tmp := configDir
	// tmp = ""
	// fmt.Print(tmp)


	mimetmp := mime.TypeByExtension(ext)
	tmpt := strings.Split(mimetmp, ";")
	mime_sl := strings.Split(tmpt[0], "/")
	mime_top := mime_sl[0]

	switch mime_top {
	// case ".bmp", ".jpg", ".jpeg", ".png", ".xpm", ".webp", ".tiff", ".gif", ".jfif", ".ico", ".svg", ".svgz":
	case "image":


		if get_file_mb() > 100 {
			preview_output = "file to big to preview"
		} else {
			preview_output = image_exif(file, width, height, file, image_tags, "image")
		}
	// case ".mp3", ".flac", ".ogg":
	// case ".wav", ".mp3", ".flac", ".m4a", ".wma", ".ape", ".ac3", ".ogg", ".spx", ".opus", ".mka":
	case "audio":
		// fmt.Println(exif_fmt(file, music_tags))
		// get_hash()

		preview_output = image_exif(file, width, height, file, music_tags, "audio")


	// case ".mkv", ".mp4", ".webm", ".avi", ".mts", ".m2ts", ".mov", ".flv":
	case "video":
		preview_output = image_exif(file, width, height, file, video_tags, "video")

	default:
		// fmt.Println("sdf")

		if get_file_mb() > 0.1 {
			preview_output = fmt.Sprintf("file to big to preview\n%v mb", get_file_mb())
		} else {
			// fmt.Println(getBaseFolder(file))
			if ignored_folders[getBaseFolder(file)] {
				var sb strings.Builder
				for i := range ignored_folders {
					sb.WriteString("\"")
					sb.WriteString(i)
					sb.WriteString("\"")
					sb.WriteString(" ")
				}
				preview_output = fmt.Sprintln("file in ignored folders list:", sb.String())

			} else {
				preview_output = read_file(file)
			}

			if disable_wordwrap != "1" {
				preview_output = wordwrap.WrapString(preview_output, uint(width))
				preview_output = char_wrap(preview_output, width)
			}

		}

	}





	if debug_time {
		time_output = time_output + fmt.Sprintln("total time: ",time.Since(prgstart))
		preview_output += sep1 + "\n"
		preview_output += time_output
		preview_output += fmt.Sprint(mime_sl)
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

	// unhideCursor()

}
