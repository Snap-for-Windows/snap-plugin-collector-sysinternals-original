package sysinternals

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

// SysinternalsCollector implementation.  Needed even if empty
type SysinternalsCollector struct{}

//stringInNamespace are the metrics that we are looking for
func stringInNamespace(givenString string) bool {
	availableMetrics := []string{"threadCount", "handleCount", "processorCount"}
	for _, metricName := range availableMetrics {
		if metricName == givenString {
			return true
		}
	}
	return false
}

// Unzip : unzip zip folders
// http://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
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

/*
* CollectMetrics collects metrics for testing.
* CollectMetrics() is be called by Snap when a task (which is collecting one+ of the metrics returned from the GetMetricTypes()) is started.
* Input: A slice of all the metric types being collected.
* Output: A slice (list) of the collected metrics as plugin.Metric with their values and an error if failure.
 */
func (SysinternalsCollector) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{} // Create a slice of MetricType objects. This is where the metrics requested by the task will be stored

	_, currentFilePath, _, _ := runtime.Caller(0) //get the current directory
	dirpath := path.Dir(currentFilePath)
	dirpath = strings.Replace(dirpath, "/", "\\", -1) //change the / to \ cause windows

	//http://stackoverflow.com/questions/6182369/exec-a-shell-command-in-go
	cmd := exec.Command("pslist", "/accepteula") // automatically accept eula if applicable, TODO use the "no banner" to remove banner and clean up code a bit?
	stdout, err := cmd.Output()

	//if the pslist exe does not exist
	if err != nil {
		println(err.Error()) //print out the error message
		fmt.Println("Attempting to Download Pslist")

		//http://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go
		out, err1 := os.Create("pstools.zip") //create a pstools.zip file to put data into
		if err1 != nil {
			err := os.Remove("pstools.zip") //remove the pstools zip file
			if err != nil {
				fmt.Println(err)
				return nil, fmt.Errorf("Unable to create unique folder for downloading PsTools")
			}
			//try one more time to download the tools
			out, err11 := os.Create("pstools.zip")
			if err11 != nil {
				return nil, fmt.Errorf("Unable to create a generic pstools file to download into")
			}
			defer out.Close()
		}
		defer out.Close()
		resp, err2 := http.Get("https://download.sysinternals.com/files/PSTools.zip") //get URL and connect to it
		if err2 != nil {
			return nil, fmt.Errorf("Unable to get to get to the URL to download PsTools")
		}
		defer resp.Body.Close()
		_, err3 := io.Copy(out, resp.Body) // actually download
		if err3 != nil {
			return nil, fmt.Errorf("Unable to download the file")
		}

		Unzip("pstools.zip", dirpath) //unzip the file and dump contents into directory
		fmt.Println("Successfully downloaded Pslist")
	}

	cmd = exec.Command("pslist")   //retry command after downloading sequence
	stdout, _ = cmd.Output()       //get the output of the command
	pslistOutput := string(stdout) // turn the output into a string

	var stringSlice []string
	stringSlice = strings.Split(pslistOutput, "\n")           //This splits the output into lines in a Slice
	stringSlice = append(stringSlice[:0], stringSlice[8:]...) //This removes the additional information about the computer (the first 7 lines), will be unnecessary if no banner works

	threadCount := 0
	handleCount := 0
	processCount := 0

	//Go row by row and parse each row
	//the last line of stringSlice is blank and so you must do len(stringSlice)-1 to avoid an out of bounds error
	for v := 0; v < len(stringSlice)-1; v++ {
		item := stringSlice[v]

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

	metricNames := make([]string, 0)
	for _, mt := range mts {
		metricNames = append(metricNames, mt.Namespace[len(mt.Namespace)-1].Value)
	}
	// Iterate through each of the metrics specified by task to collect
	for idx, mt := range mts {
		if _, err := mt.Config.GetString("test"); err == nil {
			continue
		}
		if mt.Namespace[len(mt.Namespace)-1].Value == "threadCount" {
			if val, err := mt.Config.GetInt("testint"); err == nil {
				mts[idx].Data = val
			} else {
				mts[idx].Data = threadCount
			}
			metrics = append(metrics, mts[idx])
		} else if mt.Namespace[len(mt.Namespace)-1].Value == "handleCount" {
			if val, err := mt.Config.GetInt("testint"); err == nil {
				mts[idx].Data = val
			} else {
				mts[idx].Data = handleCount
			}
			metrics = append(metrics, mts[idx])
		} else if mt.Namespace[len(mt.Namespace)-1].Value == "processCount" {
			if val, err := mt.Config.GetInt("testint"); err == nil {
				mts[idx].Data = val
			} else {
				mts[idx].Data = processCount
			}
			metrics = append(metrics, mts[idx])
		}
	}
	return metrics, nil
}

/*
 * GetMetricTypes returns a list of available metric types
 * GetMetricTypes() is called when this plugin is loaded in order to populate the "metric catalog" (where Snap
 * stores all of the available metrics for each plugin)
 * Input: Config info. This information comes from global Snap config settings
 * Output: A slice (list) of all plugin metrics, which are available to be collected by tasks
 */
func (SysinternalsCollector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	// slice to store list of all available perfmon metrics
	mts := []plugin.Metric{}

	mts = append(mts, plugin.Metric{
		Namespace: plugin.NewNamespace("intel", "sysinternals", "threadCount"),
		Version:   1,
	})

	mts = append(mts, plugin.Metric{
		Namespace: plugin.NewNamespace("intel", "sysinternals", "handleCount"),
		Version:   1,
	})
	mts = append(mts, plugin.Metric{
		Namespace: plugin.NewNamespace("intel", "sysinternals", "processCount"),
		Version:   1,
	})
	return mts, nil
}

/*
 * GetConfigPolicy() returns the config policy for this plugin
 *   A config policy allows users to provide configuration info to the plugin and is provided in the task. Here we define what kind of config info this plugin can take and/or needs.
 */
func (SysinternalsCollector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	// This rule is simply for unit testing, so I can pass in my own values for each metric rather than getting them from counters.go
	policy.AddNewFloatRule([]string{"random", "float"},
		"testfloat",
		false,
		plugin.SetMaxFloat(1000.0),
		plugin.SetMinFloat(0.0))

	// For now, assuming that perfmon has no configs. May need to add some if permissions becomes an issue.
	return *policy, nil
}
