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
	"strings"
	"regexp"
	"os"
	"log"
	"os/exec"
)

// Struct for storing the dynamic arg and its contents
type dynamic_arg struct {

	path string
	content []string

}

type command struct {

	cmd string
	result string

}

func linesStringCount(s string) int {
    n := strings.Count(s, "\n")
    if len(s) > 0 && !strings.HasSuffix(s, "\n") {
        n++
    }
    return n
}

func make_commands(cmd_string string, in_args []dynamic_arg) []string{

	// Inputs array retained command order,
	// so lets just do a simple substitution
	out_cmds := []string{}

	// Loop over all commands we are going to generate
	for i:=0;i<len(in_args[0].content)-1;i++ {

		new_cmd := cmd_string
		// Loop over all dynamic args for this command
		for j:=0;j<len(in_args);j++ {

			new_cmd = strings.Replace(new_cmd,in_args[j].path, in_args[j].content[i],-1)

		}
		out_cmds = append(out_cmds, new_cmd)
	}

	return out_cmds

}

func read_dynamic_args(in_cmd string) []dynamic_arg {
// Returns the dynamic args as a string 


	var inputs = []dynamic_arg {}
	// Find all non-static inputs
	re := regexp.MustCompile("\\{(.*?)\\}")
	match := re.FindAllString(in_cmd,-1)
	for _,arg := range match {

		this_path := arg[1:len(arg)-1]
		// Open file, read into string
		b, err := os.ReadFile(this_path)
		if err != nil {
			fmt.Print(err)
		}

		file_content := strings.Split(string(b),"\n")
		this_path = "{" + this_path + "}"
		inputs = append(inputs, dynamic_arg{path:this_path,content:file_content})

	}

	return inputs


}

func check_input(in_args []dynamic_arg) {
// Check input command sanity
// Also make sure that the input args 


	// For all lists of args, make sure to check that the number of
	// items is the same in all
	argcount := 0
	for i:=0;i<len(in_args);i++ {
		numlines := len(in_args[i].content)
		// If this is the only dynamic arg
		// or this is the first one, skip argc check
		argcount = numlines
		if i == 0 {
			continue
		}

		if argcount != len(in_args[i-1].content) {
			log.Fatal("ERROR: Number of input args in dynamic variables not the same")
		}
	}

	return


}

func exec_sh_worker(id int, commands <-chan string, results chan<-string) {

	for this_command := range commands {
		cmd := exec.Command("bash","-c",this_command)
		output,err := cmd.CombinedOutput()

		if err != nil {
			log.Fatal(err)
		}
		
		results <- string(output)
	}
}

func main() {

	// Parse arguments
	//################

	// Number of processes to exec in parallel
	numw := flag.Int("j", 4, "# of workers")
	// The command to execute and parallelize over
	cmd := flag.String("cmd", "ls", "Command to parallelize over")
	flag.Parse()

	// Extract dynamic args file names
	dynamic_args := read_dynamic_args(*cmd)

	// Input sanity checking for dynamic args
	check_input(dynamic_args)

	// Create shell command strings
	commands_list := make_commands(*cmd,dynamic_args)

	// Command execution
	// ###############
	numjobs := len(commands_list)

	// Make job and result channels
	jobs := make(chan string, numjobs)
	results := make(chan string, numjobs)
	
	// Make workers
	for w := 1; w <= *numw; w++ {
		go exec_sh_worker(w, jobs, results)
	}

	// Send commands to sh workers
	for j := 0; j < numjobs; j++ {
		jobs <- commands_list[j]
	}
	close(jobs)

	// Lets get the command output
	cmd_res := []string{}
	for a := 0; a < numjobs; a++ {
		cmd_res = append(cmd_res,<-results)
	}
	
	fmt.Println(cmd_res)



}
