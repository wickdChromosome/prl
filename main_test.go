package main

import(

	"testing"
	"os/exec"
)


func TestLinesStringCount(t *testing.T) {
	/* Function to test the line counting function
	*/

	msg := "s1\ns2\n"
	if linesStringCount(msg) != 2 {
		t.Error("Error in linesStringCount()")
	}

}

func TestLs(t *testing.T) {
	/* Function to test running a simple ls in parallel
	*/

	test_command := "go build && ./prl -j 4 -cmd 'ls {for i in `seq 1 1000`; do echo prl; done}'"
        cmd := exec.Command("bash","-c",test_command)
                _,err := cmd.CombinedOutput()

	if err != nil {
		t.Error("Error returned when trying to run a simple ls in parallel.")
	}

}

