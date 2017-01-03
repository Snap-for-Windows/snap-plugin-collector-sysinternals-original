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

	// row := strings.Fields(birdistheword)
	// print(len(row))

	var stringSlice []string
	//This splits the output into lines
	stringSlice = strings.Split(birdistheword, "\n")
	fmt.Println(strconv.Itoa(len(stringSlice)) + " what up")
	//This removes the additional information about the computer (the first 7 lines)
	stringSlice = append(stringSlice[:0], stringSlice[8:]...)
	fmt.Println(strconv.Itoa(len(stringSlice)) + " the second")

	// var finalSlice []map[string]string
	
	// for _, Selement := range stringSlice {
		
	// 	tempSlice := strings.Split(Selement, "  ")
	// 	for _, element := range tempSlice {
	// 		stringMap := make(map[string]string)
	// 		fmt.Println(element)// stringMap["Name"] = element[1].String()
	// 		// stringMap["Pid"] = element[2]
	// 		// stringMap["Pri"] = element[3]
	// 		finalSlice = append(finalSlice, stringMap)
	// 	}
	// }

	//this is used to split the stuff in the slice into different slices
	//This does not work!!
	//It is expected that after the Name of the process there is 2 spaces
	//Otherwise there may not be two spaces

	//So first get up to two spaces and then get single spaces
	for _, Selement := range stringSlice {
		nameSplit := strings.Index(Selement, "  ")
		SelementName := Selement[0:nameSplit]
		fmt.Println(SelementName)

		Selement = Selement[nameSplit:len(Selement)]
		Selement = strings.Trim(Selement, " ")
		// fmt.Println(Selement)
		number := strings.Index(Selement, " ")
		stringMap := make(map[string]string)

		if number != 0 && len(Selement) != 0 {
			//do this in a while loop (instead of if loop)
			stringMap["Name"] = SelementName

			stringMap["Pid"] = Selement[0:1]
			Selement = Selement[number:len(Selement)]
			Selement = strings.Trim(Selement, " ")

			stringMap["Pri"] = Selement[0:1]
			Selement = Selement[number:len(Selement)]
			Selement = strings.Trim(Selement, " ")

			stringMap["Thd"] = Selement[0:1]
			Selement = Selement[number:len(Selement)]
			Selement = strings.Trim(Selement, " ")

			stringMap["Hnd"] = Selement[0:1]
			Selement = Selement[number:len(Selement)]
			Selement = strings.Trim(Selement, " ")

			stringMap["Priv"] = Selement[0:1]
			Selement = Selement[number:len(Selement)]
			Selement = strings.Trim(Selement, " ")
		}
		// for key, value := range stringMap {
		// 	fmt.Println("Key: ", key, "Value: ", value)
		// }
	}

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
