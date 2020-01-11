package main

import (
	"os"
	"os/exec"
)

var supportedCLIs = []string{
	"bash", "powershell", "cmd",
}

func new(c string) (*cli, error) {
	ps, err := exec.LookPath(c + ".exe")
	if err != nil {
		return nil, err
	}
	return &cli{
		shell: ps,
	}, nil
}

func (p *cli) execute(env map[string]string, cmds ...string) (exitCode int, err error) {
	switch p.shell {
	case "powershell":
		cmds = append([]string{"-NoProfile", "-NonInteractive"}, cmds...)
	case "bash":
		cmds = append([]string{"-c"}, cmds...)
	}
	if p.shell == "powershell" {

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
