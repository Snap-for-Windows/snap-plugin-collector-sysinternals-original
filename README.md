# Snap Collector Plugin - SysInternals
This Snap plugin collects a few metrics using pslist.exe (# of processes, threads, handles).  The script will also automatically download the executable from the Microsoft site.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [Known Issues](#known-issues)
6. [License](#license-and-authors)
7. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements 
* [golang 1.7+](https://golang.org/dl/) (needed only for building as code is written in Go)

### Installation
To build this plugin through the Snap github repository:

1. Download the zip file from github (or clone the repo)

2. Set up Golang on your machine

  a. Download Golang
  
  b. Install Golang
  
  c. Set up environmnet variables (both GO and GOPATH)
  
3. Download Glide (to be able to download packages)

  a. Download Glide
  
  b. Set an environment path for glide (or just move the exe file around)
  
4. Set up the dependencies through Glide

  ```
  glide.exe install
  ```

  5. Build Snap through Go (you cannot be in the scripts folder!)

  ```
  go build scripts\build_all.sh
  ```

6. Download the sysinternals plugin

7. Install dependencies for the plugin using glide

8. Build the project through Go

  ```
  go build github.com\Snap-for-Windows\snap-plugin-collector-sysinternals
  ```

### Configuration and Usage
1. Run snapteld.exe with the appropriate flags (example below)

  ```
  snapteld.exe --plugin-trust 0 --log-level 1
  ```

2. Load the plugin through snaptel.exe
 ```
 snaptel.exe plugin load snap-plugin-collector-sysinternals.exe
 Plugin loaded
 Name: sysinternals-collector
 Version: 1
 Type: collector
 Signed: false
 Loaded Time: Mon, 20 Feb 2017 11:17:17 MST
 ```

3. Enjoy!

## Documentation
### Collected Metrics
The plugin collects the following:

Namespace | Description
----------|------------
/intel/sysinternals/handleCount | gets the number of handles
/intel/sysinternals/processCount | gets the number of processes
/intel/sysinternals/threadCount | gets the number of threads

### Examples

Load sysinternals plugin

See available metrics for your system
```
snaptel.exe metric list
```

Run the plugin using a task manifest file (the below example uses the Snap mock file plugin to output the data to a text file)
```
{ 
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "30s"
    },
    "max-failures": 10,
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/sysinternals/handleCount": {},
                "/intel/sysinternals/processCount": {},
                "/intel/sysinternals/threadCount": {}
            },
            "process": [
                {
                    "plugin_name": "passthru-grpc",
                    "process": null,
                    "publish": [
                        {
                            "plugin_name": "mock-file-grpc",
                            "config": {
                                "file": "C:\\SnapLogs\\snap-sysinternals-file.log"
                            }
                        }
                    ]
                
            ]
        }
    }
}
```
Create task:
```
snaptel.exe task create -t sysinternal-file.yaml
```
Stop task (use the task ID that was given the the task was initially created):
```
snaptel.exe task stop 4a156b0f-582f-4a13-8d67-120a2ba72e1d
Task stopped:
ID: 4a156b0f-582f-4a13-8d67-120a2ba72e1d
```

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@meinfield](https://github.com/meinfield/)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
