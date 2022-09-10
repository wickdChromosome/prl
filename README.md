# prl
A simple tool for concurrent shell command execution

## Examples

Supply arguments in a pre-made file
```
# zip every file in paths.txt, with 5 workers in parallel
go build
./prl -j 5 -cmd "zip -r {paths.txt}.zip {paths.txt}"
```

Supply arguments on the go, using temp files
```
temp_f=$(mktemp) && ls test/*.zip > $temp_f && prl -j 6 -cmd "du -sh {$(printf $temp_f)}" && rm $temp_f
```

## Arguments

- j -> The number of concurrent processes to execute the command over
- cmd -> A string of the command to execute. 

The command string should have filenames with arguments to parallelize over in parantheses -> {}
These arguments inside should be separated by new line(\n).
prl will execute the supplied command in parallel, substituting the arguments inside the file into
where the filename with parentheses is.

For example, if a file, paths.txt is supplied like this:
```
./prl -j 2 -cmd "ls {paths.txt}"
```
Where paths.txt contains:
```
/home/user1
/home/user2
```
The commands executed concurrently will be:
```
ls /home/user1
ls /home/user2
```
But this can be used for any shell command.

## TODO
- Better logging, where all the output is captured and sorted by command
- Better on-the-go command execution. Maybe using sh's <() operator? Temp files are still a bit messy
- Many, many more tests
- Add support for spaces in filenames
