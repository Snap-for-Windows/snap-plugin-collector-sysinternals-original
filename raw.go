package main

import "os/exec"
import "strings"

func main() {
	// run(exec.Command("cmd", "-Command", "netstat"))
	cmd := exec.Command("pslist")
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}
	birdistheword := string(stdout)
	//print(birdistheword)

	row := strings.Split(birdistheword, "\n")
	print(len(row))

	for index, element := range row {
		print("jhfghljhgyftrdfdfcvjjliuytfrdgfcvhbjiuiyt")
		print(element)
		print(index)

	}

	//http://stackoverflow.com/questions/6182369/exec-a-shell-command-in-go
}

/*
Check for pstools in the folder
If it exists then try to update
else report error

Run tool(s)
Get all the input
Put it into rows by find the \n (vector)
Then sort into columns
*/
