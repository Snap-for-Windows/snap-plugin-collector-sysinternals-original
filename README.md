#Snap Collector Plugin - SysInternals

##Getting Started
###Installation
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
'''
glide.exe install
'''
5. Build Snap through Go (you cannot be in the scripts folder!)
'''
go build scripts\build_all.sh
'''
6. Download the sysinternals plugin
7. Install dependencies for the plugin using glide
8. Build the project through Go
'''
go build github.com\Snap-for-Windows\snap-plugin-collector-sysinternals
'''
9. Run snapteld.exe with the appropriate flags (example below)
'''
snapteld.exe --plugin-trust 0 --log-level 1
'''
10. Load the plugin through snaptel.exe
'''
snaptel.exe plugin load snap-plugin-collector-sysinternals.exe
'''
11. Run the plugin using a task (the below example uses the Snap mock file plugin to output the data to a text file)
'''
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
                }
            ]
        }
    }
}
'''
'''
snaptel.exe task create -t sysinternal-file.yaml
'''
12. Enjoy!

##Documentation
