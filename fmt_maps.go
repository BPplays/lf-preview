package main


var video_tags = [][]string{
	{"Duration", "FileSize"},
	{"ImageSize", "VideoFrameRate"},
	{"VideoCodecID", "MIMEType"},
	{"Megapixels"},
}

var music_tags = [][]string{
	{"Title", "Duration"},
	{"Genre", "Album", "Artist", "Composer", "Date"},
	{"SampleRate", "Channels", "MIMEType"},
}




var image_tags = [][]string{
	{"ImageSize", "Megapixels", "FileSize"},
	{"MIMEType", "ColorSpace", "ColorPrimaries", "Compression"},
	{"BitDepth", "BitsPerSample", "YCbCrSubSampling", "ChromaFormat"},
}



// replacing FileType with MIMEType





var exif_key_map = map[string]string{
	"Title":                "Title",
	"Genre":                "Genre",
	"Composer":                "Composer",
	"PictureBitsPerPixel":                "Picture Bits Per Pixel",
	"FileModifyDate":                "File Modify Date",
	"FileAccessDate":                "File Access Date",
	"PictureDescription":                "Picture Description",
	"Directory":                "Directory",
	"TrackNumber":                "Track Number",
	"Duration":                "Duration",
	"Date":                "Date",
	"FileTypeExtension":                "File Type Extension",
	"FileSize":                "File Size",
	"SampleRate":                "Sample Rate",
	"FileName":                "File Name",
	"FileType":                "File Type",
	// "MIMEType":                "MIME Type",
	"MIMEType":                "Media Type",
	"Album":                "Album",
	"Artist":                "Artist",
	"Comment":                "Comment",
	"ImageSize":                "Image Size",
	"YCbCrSubSampling": "Y Cb Cr Sub Sampling",
	"BitsPerSample": "Bits Per Sample",
	"ColorSpace": "Color Space",
	"BitDepth": "Bit Depth",
	"ChromaFormat": "Chroma Format",
	"ColorPrimaries": "Color Primaries",
}


