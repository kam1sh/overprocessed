package process

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kam1sh/overprocessed/pkg/process/api"
)

type Process struct {
	cmd      *exec.Cmd
	exitCode int
}

func NewProcess(cmd string, args ...string) api.Process {
	c := exec.Command(cmd, args...)
	return &Process{
		cmd:      c,
		exitCode: -1,
	}
}

func NewProcessEnv(cmd string, args []string, env map[string]string, includeSystem bool) api.Process {
	c := exec.Command(cmd, args...)
	if includeSystem {
		for _, v := range os.Environ() {
			name, value, _ := strings.Cut(v, "=")
			if _, exists := env[name]; !exists {
				env[name] = value
			}
		}
	}
	envArr := make([]string, len(env))
	for k, v := range env {
		envArr = append(envArr, fmt.Sprintf("%v=%v", k, v))
	}
	c.Env = envArr
	return &Process{
		cmd:      c,
		exitCode: -1,
	}
}

func (p *Process) Wait() error {
	err := p.cmd.Wait()
	if _, ok := err.(*exec.ExitError); ok {
		err = nil
	}
	p.exitCode = p.cmd.ProcessState.ExitCode()
	return err
}

func (p *Process) ExitCode() int {
	return p.exitCode
}

func (p *Process) Start() error {
	err := p.cmd.Start()
	return err
}

func (p *Process) Cmd() *exec.Cmd {
	return p.cmd
}
