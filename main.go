package main

import (
	"bytes"
	"flag"
	"fmt"
	"io" // use io.Writer interface
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"time"

	// "path/filepath"

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
func run(filename string, out io.Writer, skipPreview bool) error {

	// reads the content of the input md file into a slice of bytes
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// responsible for converting md to html (bf uses here)
	htmlData := parseContent(input)

	// make tmp and check for errors
	temp, err := ioutil.TempFile("", "mdv-*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	// outName := fmt.Sprintf("%s.html", filepath.Base(filename)) // file.md.html
	outName := temp.Name() ///tmp/mdv-2227813109.html

	// Fprintln(w io.Writer, a ...any) (n int, err error)
	// prints the remaining arguments to that interface.
	fmt.Fprintln(out, outName)

	// returns a potential error when writing the HTML file,
	// !!! which the run() also returns as its error.
	// return saveHTML(outName, htmlData)

	// check for the error instead of directly returning it as the function
	// now continues to preview the file.
	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	// clean temporary
	// another benefit of using the run function;
	// since it returns a value, instead of relying on os.Exit to exit,
	// can safely use the defer statement to clean up the resources.
	defer os.Remove(outName)

	return preview(outName)
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

// uses the os/exec package to execute a separate process
// in this case, a command that opens a default application
// based on the given file
// uses exec.LookPath to locate the executable in the $PATH
// and executes it, passing the extra parameters
// and the temporary file name as arguments.
func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// define executable based on os
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("os not supported")
	}

	// append filename to parameters slice
	cParams = append(cParams, fname)

	// locate executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// By deleting the file automatically, introduce a small race condition:
	// the browser may not have time to open the file before it gets deleted.
	// can solve this in different ways, but to keep things simple,
	// add a small delay to the preview() before it returns,

	// return exec.Command(cPath, cParams...).Run()
	err = exec.Command(cPath, cParams...).Run()

	// give the browser some time to open the file before deleting it​
	// ! adding a delay isn’t a recommended long-term solution.
	// can update this function to clean up resources using a signal [somewhere]
	time.Sleep(3 * time.Second)

	return err

}

func main() {
	// check if flag has been set and use it as input to the run function.
	// otherwise, return the usage information to the user and terminate the program.
	// finally, check the error return value from the run function
	// 		and exiting with an error message in case it isn’t nil.
	filename := flag.String("file", "", "md file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	// main uses run() return value to decide whether to exit with an error code.
	// run() itself uses saveHTML() return value
	// run(filename string, out io.Writer) error
	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
