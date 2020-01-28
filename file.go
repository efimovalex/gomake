package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

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

	if v, ok := m["default"]; ok && v != nil {
		defaults = strings.Split(strings.TrimSpace(v.(string)), " ")
	}

	if v, ok := m["targets"]; ok {
		vs := v.(map[string]interface{})
		helpTargets = make(map[string]string, len(vs))
		if len(vs) > 0 {
			targets = make(map[string]interface{}, len(vs))
			for k, v := range vs {
				if strings.Contains(k, " ") {
					for _, tk := range strings.Split(k, " ") {
						targets[tk] = v
					}
				} else {
					targets[k] = v
				}
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

func initFileCreate(filename string) error {
	info, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot stat the file %s because: %s", filename, err)
	}

	if !os.IsNotExist(err) && !info.IsDir() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\nFile %s already exists. Are you sure you want to overwrite it? [y/N]: ", filename)
		approval, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("cannot read response because: %s", err)
		}
		approval = strings.TrimSpace(approval)
		if strings.ToLower(approval) != "y" && strings.ToLower(approval) != "yes" {
			return nil
		}
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}

	_, err = f.WriteString(initFile)
	if err != nil {
		return fmt.Errorf("error writing to file: %s", err)
	}

	return nil
}
