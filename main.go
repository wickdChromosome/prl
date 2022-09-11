package main

/*

This module is used to parallelize sh commands.

It allows you to enter some command, and parallelize inputs on 
a pool of processes.

All non-static arguments (those in {}) must have the same amount of members to iterate over.

Everything else in the cmd string will just get copied over.

Try running w/ --dry-run to see the commands that will be executed.

*/



import (

	"fmt"
	"flag"
	"strings"
	"regexp"
	"log"
	"os/exec"
	"github.com/schollz/progressbar/v3"

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
	/* Counts the number of new line
	   characters in the input string
	*/

	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
	    n++
	}
	return n
}


func make_commands(cmd_string string, in_args []dynamic_arg) []string{
	/* Creates an array of command strings which will be ran 
	   in parallel
	*/

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
	/* Returns the dynamic args as a string 
	*/

	var inputs = []dynamic_arg {}
	// Find all non-static inputs
	re := regexp.MustCompile("\\{(.*?)\\}")
	match := re.FindAllString(in_cmd,-1)
	if len(match) == 0 {
		flag.PrintDefaults()
		log.Fatal("Error: no dynamic arguments supplied, nothing to parallelize over")
	}

	for _,arg := range match {

		this_command := arg[1:len(arg)-1]

		cmd := exec.Command("bash","-c",this_command)
		output,err := cmd.CombinedOutput()

		if err != nil {
			fmt.Print(err)
		}

		command_output := strings.Split(string(output),"\n")
		this_command = "{" + this_command + "}"
		inputs = append(inputs, dynamic_arg{path:this_command,content:command_output})

	}

	return inputs


}

func check_input(in_args []dynamic_arg) {
	/* Check input command sanity
	   
	*/


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
			flag.PrintDefaults()
			log.Fatal("ERROR: Number of input args in dynamic variables not the same")
		}
	}

	return


}

func exec_sh_worker(id int, commands <-chan string, results chan<-string) {
	/* Shell execution worker
	*/

	for this_command := range commands {
		cmd := exec.Command("bash","-c",this_command)
		output,err := cmd.CombinedOutput()

		if err != nil {
			log.Fatal(err)
		}
		results <- string(output)
	}
}

func show_output(cmds_out []string) {
	/* Shows the output stats of the executed command 
	*/

	for _,this_cmd := range cmds_out {

		// Print a command divider
		fmt.Println("\n==================================================")

		// Print the command results
		fmt.Println(this_cmd)

	}

}

func main() {

	// Parse arguments
	//################

	// Number of processes to exec in parallel
	numw := flag.Int("j", 4, "# of workers")
	// The command to execute and parallelize over
	cmd := flag.String("cmd", "", "Command to parallelize over")

	// Allow a dry run, where the generated commands get output
	// to stdout
	dry_run := flag.Bool("dry-run", false, "Print out generated commands, but dont execute them")

	// Progbar info string
	progbar_str := flag.String("progbar-string", "Executing commands...", "Progress bar info string")

	// Silent mode
	silent := flag.Bool("s", false, "Don't print out command results")

	flag.Parse()

	// Extract dynamic args file names
	dynamic_args := read_dynamic_args(*cmd)

	// Input sanity checking for dynamic args
	check_input(dynamic_args)

	// Create shell command strings
	commands_list := make_commands(*cmd,dynamic_args)

	// If dry run, just print commands and exit
	if *dry_run {
		if !*silent {
			show_output(commands_list)
		}
		return
	}

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

	// Make a progress bar
	bar := progressbar.NewOptions(numjobs,
	    progressbar.OptionEnableColorCodes(true),
	    progressbar.OptionSetWidth(15),
	    progressbar.OptionSetDescription("[green]"+*progbar_str+"[reset] "),
	    progressbar.OptionSetTheme(progressbar.Theme{
	    Saucer:        "[green]=[reset]",
	    SaucerHead:    "[green]>[reset]",
	    SaucerPadding: " ",
	    BarStart:      "[",
	    BarEnd:        "]",
	}))

	// Lets get the command output
	cmd_res := []string{}
	for a := 0; a < numjobs; a++ {
		cmd_res = append(cmd_res,<-results)

		// Update progress bar
		bar.Add(1)
	}

	// Add a new line, just to prettify things
	fmt.Println("")

	// If silent, don't print anything and just exit
	if *silent {
		return
	}

	// Lets prettify the output and show it
	show_output(cmd_res)



}
