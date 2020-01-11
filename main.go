package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const defaultMakefileName = `makefile.yml`

// BuildVersion - version of the build
var BuildVersion string

// BuildName - name of the build
var BuildName string

type cli struct {
	shell string
}

func main() {
	var filename string = defaultMakefileName
	flag.StringVar(&filename, "file", defaultMakefileName, "the yml file to be loaded")
	flag.StringVar(&filename, "f", defaultMakefileName, "the yml file to be loaded")
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "shows the help message")
	flag.BoolVar(&showHelp, "help", false, "shows the help message")
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "shows the version message")
	flag.BoolVar(&showVersion, "version", false, "shows the version message")
	flag.Usage = help
	flag.Parse()
	if showVersion {
		fmt.Println("Version: ", BuildVersion)
		os.Exit(0)
	}
	if showHelp || len(os.Args) == 1 {
		help()
		os.Exit(0)
	}

	wantedTarget := os.Args[len(os.Args)-1]
	if wantedTarget[0:1] == "-" || len(os.Args) == 1 {
		wantedTarget = ""
	}

	var cinterface *cli
	var vars, env map[string]string
	var targets map[string]interface{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error loading file: %s\n", err)
		os.Exit(1)
	}
	m := make(map[string]interface{})
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		fmt.Printf("Error reading yml data: %s\n", err)
		os.Exit(1)
	}

	if c, ok := m["cli"]; ok {
		cinterface, err = new(c.(string))
		if err != nil {
			fmt.Printf("Error loading provided cli: %s\n", err)
			os.Exit(1)
		}
	} else {
		cinterface, err = new("bash")
		if err != nil {
			fmt.Printf("Error loading bash cli: %s\n", err)
			os.Exit(1)
		}
	}

	if v, ok := m["vars"]; ok {
		vs := v.(map[interface{}]interface{})
		if len(vs) > 0 {
			vars = make(map[string]string)
			for k, v := range vs {
				vars[k.(string)] = v.(string)
			}

		}
	}

	if v, ok := m["env"]; ok {
		vs := v.(map[interface{}]interface{})
		if len(vs) > 0 {
			env = make(map[string]string)
			for k, v := range vs {
				env[k.(string)] = v.(string)
			}

		}
	}

	if _, ok := m["targets"]; !ok || len(m["targets"].(map[interface{}]interface{})) == 0 {
		fmt.Println("No targets defined")
		os.Exit(1)
	}

	if v, ok := m["targets"]; ok {
		vs := v.(map[interface{}]interface{})
		if len(vs) > 0 {
			targets = make(map[string]interface{})
			for k, v := range vs {
				targets[k.(string)] = v
			}
		}
	}

	if t, ok := targets[wantedTarget]; !ok {
		fmt.Printf("no target found for %s\n", wantedTarget)
		os.Exit(1)
	} else {
		target := strings.TrimSpace(t.(string))
		for k, v := range vars {
			target = strings.ReplaceAll(target, "${"+k+"}", v)
			target = strings.ReplaceAll(target, "$("+k+")", v)
		}
		for k, v := range vars {
			target = strings.ReplaceAll(target, "${"+k+"}", v)
			target = strings.ReplaceAll(target, "$("+k+")", v)
		}
		for k, v := range vars {
			target = strings.ReplaceAll(target, "${"+k+"}", v)
			target = strings.ReplaceAll(target, "$("+k+")", v)
		}
		exitCode, err := cinterface.execute(env, target)
		if err != nil {
			fmt.Printf("error executing: %v\n", err)
			os.Exit(1)
		}
		os.Exit(exitCode)
	}

	os.Exit(0)
}

// Help shows the detailed command options
func help() {
	fmt.Print(`Usage: gomake [options] [target] ...

Options:

	-f FILE, --file=path/to/file.yml 	Read pointed file as a the targets file.
	-v, --version 						Prints the version.
	-h, --help                  		Prints this message.
`)
}
