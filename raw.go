package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Unzip : unzip zip folders
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	ex, err := os.ex()
	if err != nil {
		panic(err)
	}
	exPath := path.Dir(ex)

	// run(exec.Command("cmd", "-Command", "netstat"))
	cmd := exec.Command("pslist")
	stdout, err := cmd.Output()
	if err != nil {
		println(err.Error())
		fmt.Println("Attempting to Download Pslist")

		out, err1 := os.Create("pstools.zip")
		if err1 != nil {
			fmt.Println("Unable to create file to download")
		}
		defer out.Close()
		resp, err2 := http.Get("https://download.sysinternals.com/files/PSTools.zip")
		if err2 != nil {
			fmt.Println("Unable to get to URL")
		}
		defer resp.Body.Close()
		location, err3 := io.Copy(out, resp.Body)
		if err3 != nil {
			fmt.Println("Unable to download")
		}
		Unzip(location, "%PATH")
		return
	}
	birdistheword := string(stdout)
	// print(birdistheword)

	// row := strings.Fields(birdistheword)
	// print(len(row))

	var stringSlice []string
	//This splits the output into lines
	stringSlice = strings.Split(birdistheword, "\n")
	// fmt.Println(strconv.Itoa(len(stringSlice)) + " what up")
	//This removes the additional information about the computer (the first 7 lines)
	stringSlice = append(stringSlice[:0], stringSlice[8:]...)
	// fmt.Println(strconv.Itoa(len(stringSlice)) + " the second")

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

	// get how many elements are in the slice
	// slicenumber := len(stringSlice)
	// fmt.Println(slicenumber)

	// nameCount := slicenumber
	threadCount := 0
	handleCount := 0
	processCount := 0

	//Go row by row and parse each row
	//the last line of stringSlice is blank and so you must do len(stringSlice)-1 to avoid an error
	for v := 0; v < len(stringSlice)-1; v++ {
		item := stringSlice[v]
		if v == len(stringSlice)-2 {
			fmt.Println(item)
		}
		processCount++

		//get process name and remove it
		nameSplit := strings.Index(item, "  ") //this gets the number of characters a name is
		item = item[nameSplit:]                //This gets from the name to the end of item
		item = strings.Trim(item, " ")         //this will trim the whitespace at the beginning

		//get PID and remove
		pidSplit := strings.Index(item, " ")
		item = item[pidSplit:]
		item = strings.Trim(item, " ")

		//get Priority and remove
		priSplit := strings.Index(item, " ")
		item = item[priSplit:]
		item = strings.Trim(item, " ")

		//get thread and add to overall thread count and then remove
		thdSplit := strings.Index(item, " ")
		thdCount, _ := strconv.Atoi(item[0:thdSplit]) //get the string that is the thread number and convert it
		threadCount += thdCount                       //add the thread count to the master variable
		item = item[thdSplit:]
		item = strings.Trim(item, " ")

		// fmt.Println(item)
		//get handle and add to overall handle count and then remove
		hndSplit := strings.Index(item, " ")
		hndCount, _ := strconv.Atoi(item[0:hndSplit]) //get the string that is the handle number and convert it
		handleCount += hndCount                       //add the handle count to the master variable
		item = item[hndSplit:]
		item = strings.Trim(item, " ")
	}

	fmt.Println(threadCount)
	fmt.Println(processCount)
	fmt.Println(handleCount)

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
