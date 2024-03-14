package main

import (
	"bytes"
	"fmt"
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

	cmd := exec.Command(c.bin, "-o", "-", "-f", "-")
	cmd.Stdin = bytes.NewBuffer(in)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

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