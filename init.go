package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
)


func Init() {
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

	file = os.Args[1]
	ext = path.Ext(file)
	// Subtract 2 from the parsed value
	width = arg2 - 2
	hight = arg3


	lfChafaPreviewFormat = os.Getenv("LF_CHAFA_PREVIEW_FORMAT")
	lfChafaPreviewFormatOverrideSixelRatio = os.Getenv("LF_CHAFA_PREVIEW_FORMAT_OVERRIDE_SIXEL_RATIO")
	lfChafaPreviewFormatOverrideKittyRatio = os.Getenv("LF_CHAFA_PREVIEW_FORMAT_OVERRIDE_KITTY_RATIO")
	fontRatio = os.Getenv("FONT_RATIO")
	chafaPreviewDither = os.Getenv("LF_CHAFA_PREVIEW_DITHER")
	chafaPreviewColors = os.Getenv("LF_CHAFA_PREVIEW_COLORS")



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
	configDir = getEnvOrFallback("XDG_CONFIG_HOME", defaultConfigBase)

	defaultCacheBase := filepath.Join(getHomeDir(), ".cache")
	cacheBase = getEnvOrFallback("XDG_CACHE_HOME", defaultCacheBase)

	lfCacheDir := filepath.Join(cacheBase, "lf")
	if _, err := os.Stat(lfCacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(lfCacheDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	thumbnail_cache_dir = filepath.Join(lfCacheDir, "thumbnails")
	if _, err := os.Stat(thumbnail_cache_dir); os.IsNotExist(err) {
		err := os.MkdirAll(thumbnail_cache_dir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}




	cache_byte_limit = get_folder_max_len(thumbnail_cache_dir)
}