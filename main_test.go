package main

import (
	"bytes"     //// manipulate raw byte data
	"io/ioutil" // read data from files ?? deprecated 1.16
	"os"

	//// delete files
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
	resultFile = "test1.md.html"
)

// unit
func TestParseContent(t *testing.T) {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	result := parseContent(input)

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	// which compares two slices of bytes.
	if !bytes.Equal(expected, result) {
		t.Logf("golden\n%s\n", expected)
		t.Logf("result\n%s\n", result)
		t.Error("result doesn't match expected")
	}
}

// integrated
func TestRun(t *testing.T) {
	if err := run(inputFile); err != nil {
		t.Fatal(err)
	}

	result, err := ioutil.ReadFile(resultFile) // test1.md.html
	if err != nil {
		t.Fatal(err)
	}

	expected, err := ioutil.ReadFile(goldenFile) //./testdata/test1.md.html
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden\n%s\n", expected)
		t.Logf("result\n%s\n", result)
		t.Error("result doesn't match golden")
	}

	os.Remove(resultFile) // test1.md.html
}
