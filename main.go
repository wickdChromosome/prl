package main

/*

This module is used to parallelize sh commands.

It allows you to enter some command, and parallelize inputs on 
a pool of processes.

You can enter either:
- a single, static argument
- the path to a \n separated list of arguments

All non-static arguments must have the same amount of members to iterate over.
Static arguments just get copied over in every process.

Static arguments can be included as the actual command string. List of arguments
should be supplied as {/command/file/path}

*/



import (

	"fmt"
	"flag"
	_"strings"
	"regexp"

)

func linesStringCount(s string) int {
    n := strings.Count(s, "\n")
    if len(s) > 0 && !strings.HasSuffix(s, "\n") {
        n++
    }
    return n
}

func read_dynamic_args(in_cmd string) [][]string {
// Check input command sanity
// Also make sure that the input args 


	// For all lists of args, make sure to check that the number of
	// items is the same in all

	// Find all non-static inputs
	re := regexp.MustCompile("\\{(.*?)\\}")
	match := re.FindAllString(in_cmd,-1)
	for _,arg := range match {

		this_path := arg[1:len(arg)-1]
		// Open file, read into string

		// Count number of args in there


	}




	return


}


func check_input(in_cmd string) {
// Check input command sanity
// Also make sure that the input args 


	// For all lists of args, make sure to check that the number of
	// items is the same in all

	// Find all non-static inputs
	re := regexp.MustCompile("\\{(.*?)\\}")
	match := re.FindAllString(in_cmd,-1)
	for _,arg := range match {

		this_path := arg[1:len(arg)-1]
		fmt.Println(this_path)

	}

	return


}

func main() {

	// Parse arguments
	//################

	// Number of processes to exec in parallel
	num := flag.Int("j", 4, "# of iterations")

	// The command to execute and parallelize over
	cmd := flag.String("cmd", "ls", "Command to parallelize over")

	flag.Parse()

	// Extract dynamic args file names
	read_dynamic_args(*cmd)

	// Input sanity checking
	check_input(*cmd)

	// Create shell command strings

	// Command execution
	// ###############

	// Make worker pool
	//for i:=0;i<*num;i++ {
	//	fmt.Println(i)
	//}

	// Execute commands 




}
