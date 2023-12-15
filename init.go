package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
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

	for len(sep1) < width {
		sep1 = sep1 + sep1
	}
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







func Init() {

	gr_initall()

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


func get_geometry() string {

	if geometry == "" {
		geometry = fmt.Sprintf("%dx%d", width, hight)
	}



	return geometry
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
	
		lfCacheDir = filepath.Join(cacheBase, "lf")
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
		if chafaPreviewDebugTime == "1" {
			start = time.Now()
		}
		
		// cache_byte_limit = get_folder_max_len(get_thumbnail_cache_dir())
		cache_byte_limit = 200

		if chafaPreviewDebugTime == "1" {
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





