// +build windows

package main

import (
	"errors"
	"os"
	"os/exec"
)

var supportedCLIs = map[string]string{
	// "bash":       "wsl.exe",
	"powershell": "powershell.exe",
	"cmd":        "cmd.exe",
}

func new(c string) (*cli, error) {
	if ci, ok := supportedCLIs[c]; ok {
		ps, err := exec.LookPath(ci)
		if err != nil {
			return nil, err
		}
		return &cli{
			shell: ps,
			cli:   c,
		}, nil
	}

	return nil, errors.New("unsupported cli")

}

func (p *cli) execute(env map[string]string, cmds ...string) (exitCode int, err error) {
	if p.cli == "powershell" {
		cmds = append([]string{"-NoProfile", "-NonInteractive"}, cmds...)
	}
	if p.cli == "cmd" {
		cmds = append([]string{"/C"}, cmds...)
	}
	if p.cli == "bash" {
		cmds = append([]string{"bash", "-c"}, cmds...)
	}

	cmd := exec.Command(p.shell, cmds...)
	cmd.Env = os.Environ()
	for name, value := range env {
		cmd.Env = append(cmd.Env, name+"="+value)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), err
		}
	}

	return 0, nil
}
