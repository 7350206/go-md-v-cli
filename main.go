package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	bm "github.com/microcosm-cc/bluemonday"
	bf "github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
	<html>
		<head>
			<meta http-equiv="content-type" content="text/html; charset=utf-8">
			<title>Markdown Preview Tool</title>
		</head>
		<body>
		`
	footer = `
		</body>
	</html>`
)

// coordinate the execution of the remaining functions
// returns a potential error
// main uses the return value to decide whether to exit with an error code.
func run(filename string) error {
	// reads the content of the input md file into a slice of bytes
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// responsible for converting md to html (bf uses here)
	htmlData := parseContent(input)

	outName := fmt.Sprintf("%s.html", filepath.Base(filename)) // file.md.html

	fmt.Println(outName)

	// returns a potential error when writing the HTML file,
	// !!! which the run() also returns as its error.
	return saveHTML(outName, htmlData)
	// return nil
}

func parseContent(input []byte) []byte {
	// Blackfriday has various options and plugins to customize the results
	// Run([]byte) parses md using most common extensions, such as rendering tables and code blocks.
	// Use the input param as input to Blackfriday and pass its returned content to Bluemonday

	// generated a valid block of HTML will constitut the body of the page
	output := bf.Run(input)
	body := bm.UGCPolicy().SanitizeBytes(output)

	// combine this body with the header and footer const to generate the complete html content.
	// use a buffer of bytes `bytes.Buffer` to join all the HTML parts

	// create buffer to write to file
	var buffer bytes.Buffer

	// write html to buffer
	buffer.WriteString(header) // string
	buffer.Write(body)         // []byte
	buffer.WriteString(footer)

	return buffer.Bytes() // [] byte

}

func saveHTML(outFName string, data []byte) error {
	// write bytes to file
	return ioutil.WriteFile(outFName, data, 0644) //  readable and writable by the owner, readonly by anyone else.
}

func main() {
	// check if flag has been set and use it as input to the run function.
	// otherwise, return the usage information to the user and terminate the program.
	// finally, check the error return value from the run function
	// 		and exiting with an error message in case it isn’t nil.
	filename := flag.String("file", "", "md file to preview")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	// main uses run() return value to decide whether to exit with an error code.
	// run() itself uses saveHTML() return value
	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
