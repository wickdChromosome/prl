package main

import(

	"testing"
)


func TestLinesStringCount(t *testing.T) {
// Function to test running a simple ls
// in parallel

	msg := "s1\ns2\n"
	if linesStringCount(msg) != 2 {
		t.Error("Error in linesStringCount()")
	}

}
