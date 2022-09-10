# prl -> Easy sh command parallelization

![prl demo](demo.gif)

## Examples

```
# get file sizes for all zipped pdf files in the Downloads folder, using 6 parallel processes
prl -j 6 -cmd "du -h '{ ls ~/Downloads/*.zip | grep pdf }'"
# don't forget the the two quotes(') around the command - that way spaces etc in filenames are correctly handled
```

## Arguments

- j -> The number of concurrent processes to execute the command over
- cmd -> A string of the command to execute. 

The command string should have commands to parallelize over in parantheses -> {}
The results of these commands should result in a list of arguments, separated by newline(\n)
prl will execute the supplied command in parallel, substituting the arguments into the place of the parantheses
where the filename with parentheses is.

For example, 
```
prl -j 2 -cmd "ls {ls /home}"
```
Where /home contains:
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
- Many, many more tests
