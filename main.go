package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultMakefileName = `makefile.yml`

// BuildVersion - version of the build
var BuildVersion string

// BuildName - name of the build
var BuildName string

type cli struct {
	shell string
	cli   string
}

var (
	filename    string = defaultMakefileName
	defaults    []string
	showHelp    bool
	showVersion bool
	cinterface  *cli
	vars, env   map[string]string
	targets     map[string]interface{}
	helpTargets map[string]string
	targetOrder []string
)

type content struct {
	Targets yaml.Node `yaml:"targets"`
}

func main() {
	flag.StringVar(&filename, "file", defaultMakefileName, "the yml file to be loaded")
	flag.StringVar(&filename, "f", defaultMakefileName, "the yml file to be loaded")

	flag.BoolVar(&showHelp, "h", false, "shows the help message")
	flag.BoolVar(&showHelp, "help", false, "shows the help message")

	flag.BoolVar(&showVersion, "v", false, "shows the version message")
	flag.BoolVar(&showVersion, "version", false, "shows the version message")
	flag.Usage = help
	flag.Parse()

	if showVersion {
		fmt.Println("Version: ", BuildVersion)
		exit(0, nil)
	}

	err := parseFile(filename)
	if err != nil {
		exit(1, err)
	}

	if showHelp || len(flag.Args()) == 0 && len(defaults) == 0 {
		helpWithTargets(targets)
		exit(0, nil)
	}

	if len(flag.Args()) == 0 && len(defaults) != 0 {
		for _, t := range defaults {
			exitCode, err := executeTarget(env, t)
			if err != nil {
				exit(exitCode, err)
			}
		}

		exit(0, nil)
	}

	wantedTarget := flag.Args()[len(flag.Args())-1]

	exitCode, err := executeTarget(env, wantedTarget)
	if err != nil {
		exit(exitCode, err)
	}

	exit(0, nil)
}

func executeTarget(env map[string]string, target string) (int, error) {
	if t, ok := targets[target]; !ok {
		return 1, fmt.Errorf("no target found for %s", target)
	} else {
		target = strings.TrimSpace(t.(string))
		for k, v := range vars {
			target = strings.ReplaceAll(target, "${"+k+"}", v)
			target = strings.ReplaceAll(target, "$("+k+")", v)
		}
		// Replace vars that contain vars, since map acces is random
		for k, v := range vars {
			target = strings.ReplaceAll(target, "${"+k+"}", v)
			target = strings.ReplaceAll(target, "$("+k+")", v)
		}

		// vars associated to env variables
		for k := range env {
			for vk, v := range vars {
				env[k] = strings.ReplaceAll(env[k], "${"+vk+"}", v)
				env[k] = strings.ReplaceAll(env[k], "$("+vk+")", v)
			}
		}

		return cinterface.execute(env, target)
	}
}

func parseFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error loading file: %s", err)
	}
	m := make(map[string]interface{})
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("error reading yml data: %s", err)
	}

	if c, ok := m["cli"]; ok {
		cinterface, err = new(c.(string))
		if err != nil {
			return fmt.Errorf("error loading provided cli: %s", err)
		}
	} else {
		cinterface, err = new("bash")
		if err != nil {
			return fmt.Errorf("error loading bash cli: %s", err)
		}
	}

	if v, ok := m["vars"]; ok {
		vs := v.(map[string]interface{})
		if len(vs) > 0 {
			vars = make(map[string]string, len(vs))
			for k, v := range vs {
				vars[k] = v.(string)
			}

		}
	}

	if v, ok := m["env"]; ok {
		vs := v.(map[string]interface{})
		if len(vs) > 0 {
			env = make(map[string]string, len(vs))
			for k, v := range vs {
				env[k] = v.(string)
			}
		}
	}

	if v, ok := m["default"]; ok {
		defaults = strings.Split(strings.TrimSpace(v.(string)), " ")
	}

	if v, ok := m["targets"]; ok {
		vs := v.(map[string]interface{})
		helpTargets = make(map[string]string, len(vs))
		if len(vs) > 0 {
			targets = make(map[string]interface{}, len(vs))
			for k, v := range vs {
				targets[k] = v
			}
		} else {
			return fmt.Errorf("no targets defined")
		}
	} else {
		return fmt.Errorf("no targets defined")
	}

	if len(targets) != 0 {
		c := content{}
		targetOrder = []string{}
		err = yaml.Unmarshal(data, &c)
		for _, v := range c.Targets.Content {
			if v.Style == 0 {
				if _, ok := targets[v.Value]; ok {
					helpTargets[v.Value] = strings.TrimSpace(v.HeadComment)
					targetOrder = append(targetOrder, v.Value)
				}
			}
		}
	}

	return nil
}

func help() {
	helpWithTargets(nil)
}

// Help shows the detailed command options
func helpWithTargets(targets map[string]interface{}) {
	fmt.Print(`Usage: gomake [options] [target] ...

Options:
	-f FILE, --file=path/to/file.yml    Read pointed file as a the targets file.
	-v, --version                       Prints the version.
	-h, --help                          Prints this message.

`)
	if len(targets) > 0 {
		fmt.Println("Targets:")
	}
	pad := 0
	for _, t := range targetOrder {
		if len(t) >= pad {
			pad = len(t) + 4
		}
	}
	for _, t := range targetOrder {
		fmt.Printf(" %s %s %s\n", t, strings.Repeat(" ", pad-len(t)), helpTargets[t])
	}
}

func exit(code int, err error) {
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}
