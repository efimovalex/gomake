package main

import (
	"flag"
	"fmt"
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
	filename                      string = defaultMakefileName
	defaults                      []string
	showHelp, doInit, showVersion bool
	cinterface                    *cli
	vars, env, helpTargets        map[string]string
	targets                       map[string]interface{}
	targetOrder                   []string
)

type content struct {
	Targets yaml.Node `yaml:"targets"`
}

func main() {
	flag.StringVar(&filename, "file", defaultMakefileName, "The yml file to be loaded")
	flag.StringVar(&filename, "f", defaultMakefileName, "The yml file to be loaded")

	flag.BoolVar(&showHelp, "h", false, "Shows the help message")
	flag.BoolVar(&showHelp, "help", false, "Shows the help message")

	flag.BoolVar(&doInit, "i", false, "Creates a makefile.yml")
	flag.BoolVar(&doInit, "init", false, "Creates a makefile.yml")

	flag.BoolVar(&showVersion, "v", false, "Shows the version message")
	flag.BoolVar(&showVersion, "version", false, "Shows the version message")
	flag.Usage = help
	flag.Parse()

	if showVersion {
		fmt.Println("Version: ", BuildVersion)
		exit(0, nil)
	}

	if doInit {
		err := initFileCreate(filename)
		if err != nil {
			exit(1, err)
		}

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

func executeTarget(env map[string]string, wantedTarget string) (int, error) {
	if _, ok := targets[wantedTarget]; !ok {
		return 1, fmt.Errorf("no target found for %s", wantedTarget)
	}

	target := strings.TrimSpace(targets[wantedTarget].(string))
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

	if strings.Contains(target, "gomake "+wantedTarget) {
		fmt.Print("Loop detected: exiting.\n")
		exit(1, nil)
	}

	return cinterface.execute(env, target)
}

func help() {
	helpWithTargets(nil)
}

// Help shows the detailed command options
func helpWithTargets(targets map[string]interface{}) {
	fmt.Print(`Usage: gomake [options] <target>

Options:
	-f FILE, --file=path/to/file.yml    Read pointed file as a the targets file.
	-i, --init                          Creates a makefile.yml to help things get started. Can be combilend with --filename flag.
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
		fmt.Printf(" %s %s %s\n", t, strings.Repeat(" ", pad-len(t)), strings.Title(helpTargets[t]))
	}
}

func exit(code int, err error) {
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}
