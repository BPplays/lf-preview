package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

const (
	/*
		Default inkscape binary path
	*/
	BINARY = "/usr/bin/inkscape"
)


type Converter struct {
	bin string
}

/*
	Return a new instance of the converter
*/
func svg_to_png_new() *Converter {
	var c Converter
	c.bin = BINARY
	return &c
}

/*
	Set custom inkscape binary path instead of the default
*/
func (c *Converter) SetBinary(b string) error {
	if len(b) == 0 {
		return fmt.Errorf("empty binary path")
	}
	c.bin = b
	return nil
}

/*
	Try to convert the input SVG to the PNG image
*/
func (c *Converter) Convert(in []byte) (out []byte, err error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(c.bin, "--export-type=png", "--export-filename=-", "--pipe")
	// cmd.Stdin = bytes.NewBuffer(in)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	pipe, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating pipe:", err)
		os.Exit(1)
	}

	go func() {
		defer pipe.Close() // Close the pipe when done

		// Write data to the command's standard input
		_, err := pipe.Write(in)
		if err != nil {
			fmt.Println("Error writing to pipe:", err)
			os.Exit(1)
		}
	}()

	if e := cmd.Run(); e != nil {
		err = fmt.Errorf("%s\nSTDERR:\n%s", e.Error(), stderr.String())
		return
	}

	if stdout.Len() == 0 {
		err = fmt.Errorf("got no data from inkscape")
		return
	}

	out = stdout.Bytes()

	return
}