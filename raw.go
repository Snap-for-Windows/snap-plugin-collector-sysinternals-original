package main

import "os/exec"

import "os"
import "io"
import "fmt"
import "strings"
import "strconv"
import "bytes"

func main() {
	// run(exec.Command("cmd", "-Command", "netstat"))
	cmd := exec.Command("pslist")
	stdout, err := cmd.Output()
	fmt.Println("cmon!")
	if err != nil {
		println(err.Error())
		return
	}
	birdistheword := string(stdout)
	// print(birdistheword)

	row := strings.Fields(birdistheword)
	print(len(row))

	var stringSlice []string
	stringSlice = strings.Split(birdistheword, "\n")
	fmt.Println(strconv.Itoa(len(stringSlice)) + " what up")
	// i := 6
	stringSlice = append(stringSlice[:0], stringSlice[8:]...)
	fmt.Println(strconv.Itoa(len(stringSlice)) + " the second")

	// for index, element := range row {
	// 	print(element)
	// 	print(index)
	// }
	var buffer bytes.Buffer
	for _, theList := range stringSlice {
		buffer.WriteString(theList)
	}
	print(buffer.String())
	filename := "output.txt"
	fmt.Println("writing: " + filename)
	fo, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	n, err := io.WriteString(fo, buffer.String())
	if err != nil {
		fmt.Println(n, err)
	}
	fo.Close()

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
