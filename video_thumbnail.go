package main

import (
	"fmt"
	"log"
	"os/exec"
)

const (
	/*
		Default inkscape binary path
	*/
	th_BINARY = " ffmpegthumbnailer"
)




/*
	Return a new instance of the converter
*/
func vid_thm_new() *Converter {
	var c Converter
	c.bin = th_BINARY
	return &c
}



/*
	Try to convert the input SVG to the PNG image
*/
func (c *Converter) vid_thm_Convert(in string) (out *[]byte, err error) {
	// var stdout, stderr bytes.Buffer

	cmd := exec.Command(c.bin, "-q", "10", "-i", in, "-o", "/dev/stdout")
	// cmd.Stdin = bytes.NewBuffer(in)
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr



	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output), err)
		log.Fatal(string(output), err)
	}

	// if stdout.Len() == 0 {
	// 	err = fmt.Errorf("got no data from ffmpegthumbnailer")
	// 	return
	// }

	out = &output

	return
}