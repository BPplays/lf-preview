package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-runewidth"
)

func init1(wg *sync.WaitGroup) {
	defer wg.Done()
	arg2, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error parsing argument:", err)
		return
	}
	width = arg2 - 2
}

func init2(wg *sync.WaitGroup) {
	defer wg.Done()
	arg3, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Error parsing argument:", err)
		return
	}
	hight = arg3
}


func init3(wg *sync.WaitGroup) {
	defer wg.Done()

	file = os.Args[1]
	ext = path.Ext(file)
}


func init4(wg *sync.WaitGroup) {
	defer wg.Done()

	lfChafaPreviewFormat = os.Getenv("LF_CHAFA_PREVIEW_FORMAT")

	chafaFmt = []string{}
	if lfChafaPreviewFormat != "" {
		chafaFmt = append(chafaFmt, "-f", lfChafaPreviewFormat)
	}
}


func init5(wg *sync.WaitGroup) {
	defer wg.Done()

	lfChafaPreviewFormatOverrideSixelRatio = os.Getenv("LF_CHAFA_PREVIEW_FORMAT_OVERRIDE_SIXEL_RATIO")
	lfChafaPreviewFormatOverrideKittyRatio = os.Getenv("LF_CHAFA_PREVIEW_FORMAT_OVERRIDE_KITTY_RATIO")

	fontRatio = os.Getenv("FONT_RATIO")
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

	file_font_ratio = strings.Replace(userOpenFontRatio, "/", "x", -1)
}


func init6(wg *sync.WaitGroup) {
	defer wg.Done()

	chafaPreviewDither = os.Getenv("LF_CHAFA_PREVIEW_DITHER")

	chafaDither = []string{}
	if chafaPreviewDither != "" {
		chafaDither = append(chafaDither, fmt.Sprintf("--dither=%s", chafaPreviewDither))
	}
}


func init7(wg *sync.WaitGroup) {
	defer wg.Done()

	chafaPreviewColors = os.Getenv("LF_CHAFA_PREVIEW_COLORS")

	chafaColors = []string{"--colors=full"}
	if chafaPreviewColors != "" {
		chafaColors[0] = fmt.Sprintf("--colors=%s", chafaPreviewColors)
	}
}


// func init4(wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	lfChafaPreviewFormatOverrideSixelRatio = os.Getenv("LF_CHAFA_PREVIEW_FORMAT_OVERRIDE_SIXEL_RATIO")
// }

func init8(wg *sync.WaitGroup) {
	defer wg.Done()

	// for len(sep1) < width {
	// 	sep1 = sep1 + sep1
	// }
	sep1 = strings.Repeat("=", width) 
}




var init_functions = []func(wg *sync.WaitGroup){
	init1,
	init2,
	init3,
	init4,
	init5,
	init6,
	init7,
	init8,
}

func gr_initall() {
	var wg sync.WaitGroup

	// var start [64]time.Time



	for _, fn := range init_functions {
		// if chafaPreviewDebugTime == "1" {
		// 	start[i] = time.Now()
		// }
		wg.Add(1)
		fn(&wg)
	}

	go func() {
		wg.Wait()
	}()
	

}

const (
	baseChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ!@#$%^&*(){}[]:;'\\/é¥"
)




// func intToBase(n *big.Int, base *big.Int) string {
// 	var result strings.Builder
// 	// base := big.NewInt(62)

// 	zero := big.NewInt(0)
// 	for n.Cmp(zero) > 0 {
// 		remainder := new(big.Int)
// 		remainder.Mod(n, base)
// 		result.WriteByte(baseChars[remainder.Int64()])
// 		n.Div(n, base)
// 	}

// 	// Reverse the result string
// 	runes := []rune(result.String())
// 	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
// 		runes[i], runes[j] = runes[j], runes[i]
// 	}

// 	return string(runes)
// }





func intToBase(n int64, base int64) string {
	var result strings.Builder
	// var result string = ""

	var rune_sl []rune

	var minLimit int64 = 1



	// for _, char := range baseChars {
	// 	rune_sl = append(rune_sl, char)
	// }

	rune_sl = []rune(baseChars)

	if base > int64(len(rune_sl)) {
		base = int64(len(rune_sl))
	}

	if base < minLimit {
		log.Fatal("fuck")
	}


	for n > 0 {
		remainder := n % base
		result.WriteRune(rune_sl[remainder])
		n /= base
	}

	// Reverse the result string
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// func intToBase(n int64, base int64) string {

// 	var minLimit int64 = 1
// 	var maxLimit int64 = 100

// 	if base < minLimit || base > maxLimit {
// 		log.Fatal("fuck")
// 	}
// 	var result strings.Builder
// 	// base := int64(62)

// 	for n > 0 {
// 		remainder := n % base
// 		result.WriteByte(baseChars[remainder])
// 		n /= base
// 	}

// 	// Reverse the result string
// 	runes := []rune(result.String())
// 	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
// 		runes[i], runes[j] = runes[j], runes[i]
// 	}

// 	return string(runes)
	
// }

func hw_test() {
	output := ""
    for i := 1; i <= width; i++ {
        // output += fmt.Sprint(i)
		output += fmt.Sprint(intToBase(int64(i), 99964))
    }
	for i := width; i <= width+50; i++ {
        // output += fmt.Sprint("+", i-width)
		output += fmt.Sprint(intToBase(int64(i-width), 99964))
    }
	output += "\n"
    for i := 2; i <= hight; i++ {
        output += fmt.Sprint(intToBase(int64(i), 99964), "   |   width: ", width, " hight: ", hight, "\n")
    }
	for i := hight; i <= hight+50; i++ {
        output += fmt.Sprint(intToBase(int64(i+1-hight), 99964), "\n")
    }
	fmt.Print(output)
	os.Exit(0)
}

func len_test() {

	str := "にっぽん"

	// Calculate the length of the string in terminal characters
	// length := utf8.RuneCountInString(str)
	width := runewidth.StringWidth(str)

	fmt.Println(str)

	fmt.Println("12345678")

	fmt.Printf("Length of the string: %d\n", width)

	os.Exit(0)
}


// var term_height int
// var term_width int

func Init() {
	var start time.Time
	if debug_time {
		start = time.Now()
	}

	debug_hw_test := os.Getenv("LF_CHAFA_PREVIEW_DEBUG_HW_TEST")
	debug_len_test := os.Getenv("LF_CHAFA_PREVIEW_DEBUG_LEN_TEST")

	// var err error

	// term_width, term_height, err = terminal.GetSize(int(os.Stdin.Fd()))
	// if err != nil {
	// 	fmt.Print(err)
	// 	log.Fatal(err)
	// }

	gr_initall()

	if debug_hw_test == "1" {
		hw_test()
	}

	if debug_len_test == "1" {
		len_test()
	}

	if debug_time {
		time_output = time_output + fmt.Sprintln("init time: ",time.Since(start))
	}
}









var geometry string = ""




func get_thumbnail_cache_dir() string {

	if thumbnail_cache_dir == "" {
		thumbnail_cache_dir = filepath.Join(get_lfCacheDir(), "thumbnails")
		if _, err := os.Stat(thumbnail_cache_dir); os.IsNotExist(err) {
			err := os.MkdirAll(thumbnail_cache_dir, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				log.Fatal(err)
			}
		}
	}



	return thumbnail_cache_dir
}


func get_geometry(width, height int) string {
	return fmt.Sprintf("%dx%d", width, height)
}

var preview_print_output bool
var set_preview_print_output bool

func get_print_output() bool {
	if !set_preview_print_output {
		if os.Getenv("LF_CHAFA_PREVIEW_PRINT_OUTPUT") != "0" {
			preview_print_output = true
			set_preview_print_output = true
		} else {
			preview_print_output = false
		}
	}


	return preview_print_output
}





func get_lfCacheDir() string {
	if lfCacheDir == "" {
		defaultCacheBase := filepath.Join(getHomeDir(), ".cache")
		cacheBase = getEnvOrFallback("XDG_CACHE_HOME", defaultCacheBase)
	
		lfCacheDir = filepath.Join(cacheBase, "lf-preview")
		if _, err := os.Stat(lfCacheDir); os.IsNotExist(err) {
			err := os.MkdirAll(lfCacheDir, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				log.Fatal(err)
			}
		}
	}


	return lfCacheDir
}





func get_metadata_cache_dir() string {
	if metadata_cache_dir == "" {
		metadata_cache_dir = filepath.Join(get_lfCacheDir(), "metadata", "v2")
		if _, err := os.Stat(metadata_cache_dir); os.IsNotExist(err) {
			err := os.MkdirAll(metadata_cache_dir, 0700)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				log.Fatal()
			}
		}
	}


	return metadata_cache_dir
}






func get_cache_byte_limit() int {

	if cache_byte_limit == -1 {
		var start time.Time
		if debug_time {
			start = time.Now()
		}

		// cache_byte_limit = get_folder_max_len(get_thumbnail_cache_dir())
		cache_byte_limit = 200

		if debug_time {
			time_output = time_output + fmt.Sprintln("get_cache_byte_limit time: ",time.Since(start))
		}

	}



	return cache_byte_limit
}

func get_configDir() string {
	if configDir == "" {
		defaultConfigBase := filepath.Join(getHomeDir(), ".config")
		configDir = getEnvOrFallback("XDG_CONFIG_HOME", defaultConfigBase)
	}


	return configDir
}





